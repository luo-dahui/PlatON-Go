/*
经济模型rpc接口查詢
platon_call
*/

package core

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"strconv"
	"strings"

	"github.com/PlatONnetwork/PlatON-Go/accounts/keystore"
	"github.com/PlatONnetwork/PlatON-Go/cmd/utils"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"gopkg.in/urfave/cli.v1"
)

var (
	EcoModelTxCmd = cli.Command{
		Name:   "tx_ecomodel",
		Usage:  "经济模型rpc接口查詢",
		Action: tx_ecomodel,
		Flags:  EcoModelTxCmdFlags,
	}
)

func handleTx(rlpdata, toAddress string, v interface{}) (string, error) {

	callEcomodelParams := CallEcomodelParams{
		To:   toAddress,
		Data: rlpdata,
	}

	callParams := make([]interface{}, 2)
	callParams[0] = callEcomodelParams
	callParams[1] = "latest"

	r, err := Send(callParams, "platon_call")
	if err != nil {
		return "", fmt.Errorf("send http post to invokeContract contract error")
	}
	resp := parseResponse(r)

	if len(resp.Result) > 1 {
		if resp.Result[0:2] == "0x" || resp.Result[0:2] == "0X" {
			resp.Result = resp.Result[2:]
		}
	}
	res_data, _ := hex.DecodeString(resp.Result)

	json.Unmarshal(res_data, v)

	return string(res_data), nil
}

func getBigValueByString(value string) *big.Int {
	var ok bool
	var bigValue *big.Int
	if len(value) > 2 && value[0:2] == "0x" {
		bigValue, ok = new(big.Int).SetString(value[2:], 16)
		if ok == false || config.Staking.Amount == nil {
			panic(fmt.Errorf("Amount is error"))
		}
	} else {
		bigValue, _ = new(big.Int).SetString(value, 10)
	}

	return bigValue
}

func tx_ecomodel(c *cli.Context) error {
	parseConfigJson(c.String(ConfigPathFlag.Name))

	// 加载钱包文件
	walletFile := c.String(WalletFilePathFlag.Name)
	if walletFile == "" {
		walletFile = config.Tx.Wallet
	}

	keyjson, err := ioutil.ReadFile(walletFile)
	if err != nil {
		utils.Fatalf("Failed to read the wallet file at '%s': %v", walletFile, err)
	}

	// Decrypt key with passphrase.
	passphrase := promptPassphrase(false)
	key, err := keystore.DecryptKey(keyjson, passphrase)
	if err != nil {
		utils.Fatalf("the wallet password is error: %v", err)
	}

	// privateKey := hex.EncodeToString(crypto.FromECDSA(key.PrivateKey))
	from := key.Address.Hex()
	fmt.Println(from)

	// rpc api
	action := c.String("action")
	action = strings.ToLower(action)

	funcName := c.String("funcName")
	var rlp string
	switch mapNameToFuncType[funcName] {
	case 1002:
		{
			nodeId := c.String("nodeId")
			if nodeId != "" {
				config.Staking.NodeId, _ = discover.HexID(nodeId)
			}

			amountType := c.String("amountType")
			if amountType != "" {
				typ, _ := strconv.Atoi(amountType)
				config.Staking.AmountType = uint16(typ)
			}

			amount := c.String("amount")
			if amount != "" {
				config.Staking.Amount = getBigValueByString(amount)
			}

			rlp = getRlpData(mapNameToFuncType[funcName], nil, config.Staking)
		}
	default:
		{
			fmt.Printf("funcName:%s is unknown!!!!", funcName)
			return nil
		}
	}

	result, _ := handleTx(rlp, mapNameToAddress[action], nil)
	fmt.Printf("result:\n %s \n", result)

	return nil
}
