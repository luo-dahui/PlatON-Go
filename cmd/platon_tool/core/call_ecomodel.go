/*
经济模型rpc接口查詢
platon_call
*/

package core

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/rpc"
	"gopkg.in/urfave/cli.v1"
	"math/big"
	"strings"
)

var (
	EcoModelCallCmd = cli.Command{
		Name:   "call_ecomodel",
		Usage:  "普通rpc接口查詢",
		Action: call_ordinary,
		Flags:  EcoModelCallCmdFlags,
	}
)

func call_ecomodel(c *cli.Context) error {
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
