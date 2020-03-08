/*
普通rpc接口查詢
*/

package core

import (
	"encoding/json"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/rpc"
	"gopkg.in/urfave/cli.v1"
	"math/big"
	"strings"
)

var (
	OrdinaryCallCmd = cli.Command{
		Name:   "call_ordinary",
		Usage:  "普通rpc接口查詢",
		Action: call_ordinary,
		Flags:  OrdinaryCallCmdFlags,
	}
)

// 获取返回值为string类型的接口(无入参)
func showStringValue(client *rpc.Client, action string, toBig bool) error {
	var result string
	err := client.Call(&result, action)

	if err != nil {
		fmt.Println("client.Call err", err)
		return err
	}
	if toBig == true {
		bigRes, _ := hexutil.DecodeBig(result)
		fmt.Printf("\n%v: %v\n", action, bigRes)
	} else {
		fmt.Printf("\n%v: %v\n", action, result)
	}
	return nil
}

// 获取返回值为string类型的接口(有入参)
func showStringValueByParams(action string, toBig bool, params interface{}, unit string) (string, error) {
	r, _ := Send(params, action)
	resp := parseResponse(r)

	if toBig == true {
		bigRes, _ := hexutil.DecodeBig(resp.Result)
		fmt.Printf("\n%v: %v %s\n", action, bigRes, unit)
	} else {
		fmt.Printf("\n%v: %v %s\n", action, resp.Result, unit)
	}

	return resp.Result, nil
}

func getTransactionReceipt(txHash string) {
	var receipt = Receipt{}
	res, e := Send([]string{txHash}, "platon_getTransactionReceipt")
	if e != nil {
		panic(fmt.Sprintf("send http post to get transaction receipt error！\n %s", e.Error()))
	}

	e = json.Unmarshal([]byte(res), &receipt)
	if e != nil {
		panic(fmt.Sprintf("parse get receipt result error ! \n %s", e.Error()))
	}

	fmt.Printf("\nresult:%v\n", res)
}

type Platon struct {
	Network uint64              `json:"network"` // PlatON network ID (1=Frontier, 2=Morden, Ropsten=3, Rinkeby=4)
	Genesis common.Hash         `json:"genesis"` // SHA3 hash of the host's genesis block
	Config  *params.ChainConfig `json:"config"`  // Chain configuration for the fork rules
	Head    common.Hash         `json:"head"`    // SHA3 hash of the host's best owned block
}

type NodeInfo struct {
	ID     string `json:"id"`        // Unique node identifier (also the encryption key)
	Name   string `json:"name"`      // Name of the node, including client type, version, OS, custom data
	BlsPub string `json:"blsPubKey"` // BLS public key
	Enode  string `json:"enode"`     // Enode URL for adding this peer from remote peers
	IP     string `json:"ip"`        // IP address of the node
	Ports  struct {
		Discovery int `json:"discovery"` // UDP listening port for discovery protocol
		Listener  int `json:"listener"`  // TCP listening port for RLPx
	} `json:"ports"`
	ListenAddr string            `json:"listenAddr"`
	Protocols  map[string]Platon `json:"protocols"`

	//	Protocols  map[string]interface{} `json:"protocols"`
}

func getNodeInfo(client *rpc.Client, action string) error {
	var result NodeInfo
	err := client.Call(&result, action)
	if err != nil {
		fmt.Println("client.Call err", err)
		return err
	}

	platon := result.Protocols["platon"]

	fmt.Printf("\n%s: %d\n", "ChainId", platon.Config.ChainID)
	fmt.Printf("\n%s: %d\n", "NetWord", platon.Network)
	fmt.Printf("\n%s: %v\n", "NodeId", result.ID)
	fmt.Printf("\n%s: %v\n", "BlsPubKey", result.BlsPub)
	fmt.Printf("\n%s: %v\n", "Genesis Hash", platon.Genesis.Hex())

	return nil
}

func showInterfaceValue(client *rpc.Client, keyname, action string, result interface{},
	args ...interface{}) error {

	err := client.Call(&result, action)
	if err != nil {
		fmt.Println("client.Call err", err)
		return err
	}
	/*
		for i := 0; i < argsLen; i++ {
			key := args[i]
			value := result[key]
			fmt.Printf("\n%v", value)
		}
		params, err := json.Marshal(args)
		if err != nil {
			fmt.Printf("\n%v", params)
			return err
		}*/
	fmt.Printf("\n%v: %v\n", keyname, result)
	return nil
}

func call_ordinary(c *cli.Context) error {
	parseConfigJson(c.String(ConfigPathFlag.Name))
	client, err := rpc.Dial(config.Url)
	if err != nil {
		fmt.Println("rpc.Dial err", err)
		return err
	}

	// rpc api
	action := c.String("action")
	action = strings.ToLower(action)
	switch action {
	case "blocknumber":
		{
			showStringValue(client, "platon_blockNumber", true)
		}
	case "datadir":
		{
			showStringValue(client, "admin_datadir", false)
		}
	case "nodeinfo":
		{
			getNodeInfo(client, "admin_nodeInfo")
		}
	case "getbalance":
		{
			address := c.String("address")
			if address == "" {
				fmt.Print("Incorrect Usage: flag needs an argument: -address")
			} else {
				params := make([]interface{}, 2)
				params[0] = address
				params[1] = "latest"
				res, _ := showStringValueByParams("platon_getBalance", true, params, "Von")
				bigRes, _ := hexutil.DecodeBig(res)
				balance := new(big.Int).Div(bigRes, big.NewInt(1e18))
				fmt.Printf("balance: %v LAT\n", balance)
			}
		}
	case "gettxreceipt":
		{
			txHash := c.String("txhash")
			if txHash == "" {
				fmt.Print("Incorrect Usage: flag needs an argument: -txhash")
			} else {
				// Get transaction receipt according to result
				getTransactionReceipt(txHash)
			}
		}
	default:
		{
			fmt.Printf("not found this rpc api: %s , please check it!!!!\n", action)
			return nil
		}
	}
	return nil
}
