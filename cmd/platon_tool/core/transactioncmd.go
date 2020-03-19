package core

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"strconv"
	"strings"

	"github.com/PlatONnetwork/PlatON-Go/accounts/keystore"
	"github.com/PlatONnetwork/PlatON-Go/cmd/utils"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"gopkg.in/urfave/cli.v1"
)

var (
	SendTransactionCmd = cli.Command{
		Name:   "sendTransaction",
		Usage:  "send a transaction",
		Action: sendTransactionCmd,
		Flags:  sendTransactionCmdFlags,
	}
	SendRawTransactionCmd = cli.Command{
		Name:   "sendRawTransaction",
		Usage:  "send a raw transaction",
		Action: sendRawTransactionCmd,
		Flags:  sendRawTransactionCmdFlags,
	}
	GetTxReceiptCmd = cli.Command{
		Name:   "getTxReceipt",
		Usage:  "get transaction receipt by hash",
		Action: getTxReceiptCmd,
		Flags:  getTxReceiptCmdFlags,
	}
)

func getTxReceiptCmd(c *cli.Context) {
	parseConfigJson(c.String(ConfigPathFlag.Name))
	hash := c.String(TransactionHashFlag.Name)
	if hash == "" {
		hash = config.Call.TxHash
	}

	GetTxReceipt(hash)
}

func GetTxReceipt(txHash string) (Receipt, error) {
	var receipt = Receipt{}
	res, _ := Send([]string{txHash}, "platon_getTransactionReceipt")
	e := json.Unmarshal([]byte(res), &receipt)
	if e != nil {
		panic(fmt.Sprintf("parse get receipt result error ! \n %s", e.Error()))
	}

	if receipt.Result.BlockHash == "" {
		panic("no receipt found")
	}
	out, _ := json.MarshalIndent(receipt, "", "  ")
	fmt.Println(string(out))
	return receipt, nil
}

func sendTransactionCmd(c *cli.Context) error {

	parseConfigJson(c.String(ConfigPathFlag.Name))
	from := c.String(TxFromFlag.Name)
	if from == "" {
		from = config.Tx.From
	}
	to := c.String(TxToFlag.Name)
	if to == "" {
		to = config.Tx.To
	}
	value := c.String(TransferValueFlag.Name)
	if value == "" {
		value = config.Tx.Value
	}

	hash, err := SendTransaction(from, to, value)
	if err != nil {
		utils.Fatalf("Send transaction error: %v", err)
	}

	fmt.Printf("tx hash: %s", hash)
	return nil
}

func sendRawTransactionCmd(c *cli.Context) error {

	parseConfigJson(c.String(ConfigPathFlag.Name))
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

	to := c.String(TxToFlag.Name)
	if to == "" {
		to = config.Tx.To
	}
	value := c.String(TransferValueFlag.Name)
	if value == "" {
		value = config.Tx.Value
	}

	hash, err := SendRawTransaction(from, to, value, key.PrivateKey)
	if err != nil {
		utils.Fatalf("Send transaction error: %v", err)
	}

	fmt.Printf("tx hash: %s", hash)
	return nil
}

// 此种方式需要将钱包放到节点的keystore目录下（不安全）
func SendTransaction(from, to, value string) (string, error) {
	var tx TxParams
	if from == "" {
		from = config.Tx.From
	}
	tx.From = from
	tx.To = to
	tx.Gas = config.Tx.Gas
	tx.GasPrice = config.Tx.GasPrice

	if !strings.HasPrefix(value, "0x") {
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			panic(fmt.Sprintf("transfer value to int error.%s", err))
		}
		value = hexutil.EncodeBig(big.NewInt(intValue))
	}
	tx.Value = value

	// 输入钱包密码
	passphrase := promptPassphrase(false)
	//password := "88888888"
	// params := make([]TxParams, 2)
	params := make([]interface{}, 2)
	params[0] = tx
	params[1] = passphrase

	res, _ := Send(params, "personal_sendTransaction")
	response := parseResponse(res)

	return response.Result, nil
}

func SendRawTransaction(from, to, value string, priv *ecdsa.PrivateKey) (string, error) {

	var v int64
	var err error
	if strings.HasPrefix(value, "0x") {
		bigValue, _ := hexutil.DecodeBig(value)
		v = bigValue.Int64()
	} else {
		v, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			panic(fmt.Sprintf("transfer value to int error.%s", err))
		}
	}

	nonce := getNonce(from)
	//nonce++

	newTx := getSignedTransaction(from, to, v, priv, nonce)

	hash, err := sendRawTransaction(newTx)
	if err != nil {
		panic(err)
	}
	return hash, nil
}

func sendRawTransaction(transaction *types.Transaction) (string, error) {
	bytes, _ := rlp.EncodeToBytes(transaction)
	res, err := Send([]string{hexutil.Encode(bytes)}, "platon_sendRawTransaction")
	if err != nil {
		panic(err)
	}
	response := parseResponse(res)

	return response.Result, nil
}

func getSignedTransaction(from, to string, value int64, priv *ecdsa.PrivateKey,
	nonce uint64) *types.Transaction {

	gas, _ := strconv.Atoi(config.Tx.Gas)
	gasPrice, _ := new(big.Int).SetString(config.Tx.GasPrice, 10)
	newTx, err := types.SignTx(types.NewTransaction(nonce, common.HexToAddress(to),
		big.NewInt(value), uint64(gas), gasPrice, []byte{}),
		types.NewEIP155Signer(config.ChainID), priv)

	if err != nil {
		panic(fmt.Errorf("sign error,%s", err.Error()))
	}
	return newTx
}

func getNonce(addr string) uint64 {
	res, _ := Send([]string{addr, "latest"}, "platon_getTransactionCount")
	response := parseResponse(res)
	nonce, _ := hexutil.DecodeBig(response.Result)
	//fmt.Println(addr, nonce)
	fmt.Printf("address:%v, nonce:%v \n", addr, nonce)
	return nonce.Uint64()
}
