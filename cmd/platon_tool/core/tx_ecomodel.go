/*
经济模型rpc接口查詢
platon_call
*/

package core

import (
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"math/big"
	"strconv"
	"strings"
	"time"

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

// 发送经济模型交易
func sendEcModelTx(fromAddress, toAddress, rlpdata, value string, priv *ecdsa.PrivateKey) (string, error) {

	// 发送交易
	hash, err := SendRawTransaction(fromAddress, toAddress, rlpdata, value, priv)
	if err != nil {
		utils.Fatalf("Send transaction error: %v", err)
	}

	fmt.Printf("tx hash: %s \n", hash)
	fmt.Println("==========get Transaction Receipt, please wait=============")

	// Get transaction receipt according to result
	ch := make(chan string, 1)
	exit := make(chan string, 1)
	go GetTxReceipt(hash, ch, exit, true)

	/*
	  Loop call to get transactionReceipt... until 200s timeout
	*/
	select {
	case retCode := <-ch:
		fmt.Printf("retCode: %s\n", retCode)
		if retCode == "0" {
			fmt.Println("transaction succeed.\n")
		} else {
			fmt.Println("transaction failed.\n")
			GetMsgByErrCode(retCode)
		}
	case <-time.After(time.Second * 200):
		exit <- "exit"
		fmt.Printf("get transaction receipt timeout...more than 200 second.\n")
	}

	return hash, nil
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

// 修改验证人信息
func getRlpDataByInEditCandidate(c *cli.Context, funcType uint16) string {
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

	return getRlpData(funcType, nil, config.Staking)
}

// 增持质押
func getRlpDataByInCreaseStaking(c *cli.Context, funcType uint16) string {
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

	return getRlpData(funcType, nil, config.Staking)
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
	//fmt.Println(from)

	// rpc api
	action := c.String("action")
	action = strings.ToLower(action)

	funcName := c.String("funcName")
	var rlpData string
	switch mapNameToFuncType[funcName] {
	case 1001:
		{
			rlpData = getRlpDataByInEditCandidate(c, mapNameToFuncType[funcName])
		}
	case 1002:
		{
			rlpData = getRlpDataByInCreaseStaking(c, mapNameToFuncType[funcName])
		}
	case 1003:
		{
			nodeId := c.String("nodeId")
			if nodeId == "" {
				nodeId = config.Staking.NodeId.String()
			}
			rlpData = getRlpData(mapNameToFuncType[funcName], nil, nodeId)
		}
	default:
		{
			fmt.Printf("funcName:%s is unknown!!!!", funcName)
			return nil
		}
	}

	// 发送交易
	sendEcModelTx(from, mapNameToAddress[action], rlpData, "", key.PrivateKey)

	return nil
}
