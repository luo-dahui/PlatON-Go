package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"github.com/PlatONnetwork/PlatON-Go/cmd/utils"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/console"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

const (
	deployContract = iota
	invokeContract

	DefaultConfigFilePath = "/config.json"
)

var (
	config = Config{}
)

type JsonParam struct {
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Id      int         `json:"id"`
}

type TxParams struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Gas      string `json:"gas"`
	GasPrice string `json:"gasPrice"`
	Value    string `json:"value"`
	Data     string `json:"data"`
}

type RawTxParams struct {
	TxParams
	Nonce int64 `json:"Nonce"`
}

// 查询经济模型合约接口
type CallEcomodelParams struct {
	To   string `json:"to"`
	Data string `json:"data"`
}

type Staking struct {
	NodeId          discover.NodeID `json:"nodeid"`
	AmountType      uint16          `json:"amountType"`
	Amount          *big.Int        `json:"amount"`
	DelegateAddress common.Address  `json:"delegateAddress"`
}

type Gov struct {
	ProposalID common.Hash `json:"proposalid"`
	Module     string      `json:"module"`
	Name       string      `json:"name"`
}

type Restricting struct {
	Account common.Address `json:"account"`
}

type Reward struct {
	Account common.Address    `json:"account"`
	NodeIds []discover.NodeID `json:"nodeIds"`
}

type DeployParams struct {
	From     string `json:"from"`
	Gas      string `json:"gas"`
	GasPrice string `json:"gasPrice"`
	Data     string `json:"data"`
}

type Tx struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Gas      string `json:"gas"`
	GasPrice string `json:"gasPrice"`
	Value    string `json:"value"`
	Wallet   string `json:"wallet"`
}

type Call struct {
	TxHash string `json:"txhash"`
}

type Config struct {
	ChainID     *big.Int    `json:"chainId"`
	Url         string      `json:"url"`
	Tx          Tx          `json:"tx"`
	Call        Call        `json:"call"`
	Staking     Staking     `json:"staking"`
	Gov         Gov         `json:"gov"`
	Restricting Restricting `json:"restricting"`
	Reward      Reward      `json:"reward"`
}

type FuncDesc struct {
	Name   string `json:"name"`
	Inputs []struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"inputs"`
	Outputs []struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"outputs"`
	Constant string `json:"constant"`
	Type     string `json:"type"`
}

type Response struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  string `json:"result"`
	Id      int    `json:"id"`
	Error   struct {
		Code    int32  `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type Receipt struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  struct {
		BlockHash         string       `json:"blockHash"`
		BlockNumber       string       `json:"blockNumber"`
		ContractAddress   string       `json:"contractAddress"`
		CumulativeGasUsed string       `json:"cumulativeGas_used"`
		From              string       `json:"from"`
		GasUsed           string       `json:"gasUsed"`
		Logs              []*types.Log `json:"logs"              gencodec:"required"`
		Root              string       `json:"root"`
		To                string       `json:"to"`
		TransactionHash   string       `json:"transactionHash"`
		TransactionIndex  string       `json:"transactionIndex"`
	} `json:"result"`
}

func parseConfigJson(configPath string) error {
	if configPath == "" {
		dir, _ := os.Getwd()
		configPath = dir + DefaultConfigFilePath
	}

	if !filepath.IsAbs(configPath) {
		configPath, _ = filepath.Abs(configPath)
	}

	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(fmt.Errorf("parse config file error,%s", err.Error()))
	}

	if err := json.Unmarshal(bytes, &config); err != nil {
		panic(fmt.Errorf("parse config to json error,%s", err.Error()))
	}
	return nil
}

func parseAbiFromJson(fileName string) ([]FuncDesc, error) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("parse abi file error: %s", err.Error())
	}
	var a []FuncDesc
	if err := json.Unmarshal(bytes, &a); err != nil {
		return nil, fmt.Errorf("parse abi to json error: %s", err.Error())
	}
	return a, nil
}

func parseFuncFromAbi(fileName string, funcName string) (*FuncDesc, error) {
	funcs, err := parseAbiFromJson(fileName)
	if err != nil {
		return nil, err
	}

	for _, value := range funcs {
		if value.Name == funcName {
			return &value, nil
		}
	}
	return nil, fmt.Errorf("function %s not found in %s", funcName, fileName)
}

/**
  Find the method called by parsing abi
*/
func GetFuncNameAndParams(f string) (string, []string) {
	funcName := string(f[0:strings.Index(f, "(")])

	paramString := string(f[strings.Index(f, "(")+1 : strings.LastIndex(f, ")")])
	if paramString == "" {
		return funcName, []string{}
	}

	params := strings.Split(paramString, ",")
	for index, param := range params {
		if strings.HasPrefix(param, "\"") {
			params[index] = param[strings.Index(param, "\"")+1 : strings.LastIndex(param, "\"")]
		}
	}
	return funcName, params

}

/**
  Self-test method for encrypting parameters
*/
func encodeParam(abiPath string, funcName string, funcParams string) error {
	// Determine if the method exists
	abiFunc, err := parseFuncFromAbi(abiPath, funcName)
	if err != nil {
		return err
	}

	// Parsing the method of the call
	funcName, inputParams := GetFuncNameAndParams(funcParams)

	// Determine if the parameters are correct
	if len(abiFunc.Inputs) != len(inputParams) {
		return fmt.Errorf("incorrect number of parameters ,request=%d,get=%d\n", len(abiFunc.Inputs), len(inputParams))
	}

	paramArr := [][]byte{
		Int32ToBytes(111),
		[]byte(funcName),
	}

	for i, v := range inputParams {
		input := abiFunc.Inputs[i]
		p, e := StringConverter(v, input.Type)
		if e != nil {
			return err
		}
		paramArr = append(paramArr, p)
	}

	paramBytes, _ := rlp.EncodeToBytes(paramArr)

	fmt.Printf(hexutil.Encode(paramBytes))

	return nil
}

// 输入密码
func promptPassphrase(confirmation bool) string {
	passphrase, err := console.Stdin.PromptPassword("Passphrase: ")
	if err != nil {
		utils.Fatalf("Failed to read passphrase: %v", err)
	}

	if confirmation {
		confirm, err := console.Stdin.PromptPassword("Repeat passphrase: ")
		if err != nil {
			utils.Fatalf("Failed to read passphrase confirmation: %v", err)
		}
		if passphrase != confirm {
			utils.Fatalf("Passphrases do not match")
		}
	}

	return passphrase
}
