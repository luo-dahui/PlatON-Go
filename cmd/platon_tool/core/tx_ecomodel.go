/*
经济模型rpc接口查詢
platon_call
*/

package core

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"gopkg.in/urfave/cli.v1"
)

var (
	EcoModelTxCmd = cli.Command{
		Name:   "tx_ecomodel",
		Usage:  "经济模型rpc接口查詢",
		Action: tx_ecomodel,
		Flags:  EcoModelCallCmdFlags,
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

func tx_ecomodel(c *cli.Context) error {

	parseConfigJson(c.String(ConfigPathFlag.Name))
	// rpc api
	action := c.String("action")
	action = strings.ToLower(action)

	funcName := c.String("funcName")
	var rlp string
	switch mapNameToFuncType[funcName] {
	case 1100, 1101, 1102, 1200, 1201, 1202, 2102, 2103, 3002:
		rlp = getRlpData(mapNameToFuncType[funcName], nil, nil)
	case 1103:
		{
			delAddress := c.String("address")
			if delAddress == "" {
				delAddress = config.Staking.DelegateAddress.Hex()
			}
			rlp = getRlpData(mapNameToFuncType[funcName], nil, delAddress)
		}
	case 1105:
		{
			nodeId := c.String("nodeId")
			if nodeId == "" {
				nodeId = config.Staking.NodeId.String()
			}
			rlp = getRlpData(mapNameToFuncType[funcName], nil, nodeId)
		}
	case 2100, 2101:
		{
			proposalID := c.String("proposalID")
			if proposalID == "" {
				proposalID = config.Gov.ProposalID.String()
			}
			rlp = getRlpData(mapNameToFuncType[funcName], nil, proposalID)
		}
	case 2104:
		{
			module := c.String("module")
			if module == "" {
				module = config.Gov.Module
			}
			module = strings.ToLower(module)

			name := c.String("name")
			if name == "" {
				name = config.Gov.Name
			}

			var gov Gov
			gov.Module = module
			gov.Name = name

			rlp = getRlpData(mapNameToFuncType[funcName], nil, gov)
		}
	case 2106:
		{
			module := c.String("module")
			if module == "" {
				module = config.Gov.Module
			}
			module = strings.ToLower(module)
			rlp = getRlpData(mapNameToFuncType[funcName], nil, module)
		}
	case 4100:
		{
			// 锁仓金额释放到账账户
			account := c.String("address")
			if account == "" {
				account = config.Restricting.Account.Hex()
			}
			rlp = getRlpData(mapNameToFuncType[funcName], nil, account)
		}
	case 5100:
		{
			// 查询账户在各节点未提取委托奖励
			var reward Reward
			account := c.String("address")
			if account == "" {
				account = config.Reward.Account.Hex()
			}
			reward.Account = common.HexToAddress(account)

			// nodeid
			nodeid := c.String("nodeId")
			if nodeid == "" {
				reward.NodeIds = config.Reward.NodeIds
			} else {
				node := discover.MustHexID(nodeid)
				reward.NodeIds = append(reward.NodeIds, node)
			}

			rlp = getRlpData(mapNameToFuncType[funcName], nil, reward)
			var res RewardResult
			handleCall(rlp, mapNameToAddress[action], &res)
			for i := 0; i < len(res.Ret); i++ {
				fmt.Println("===============================")
				fmt.Printf("nodeid: %s \n", res.Ret[i].NodeId.String())
				fmt.Printf("stakingNumber: %d \n", res.Ret[i].StakingNum)
				bigRes, _ := hexutil.DecodeBig(res.Ret[i].Reward)
				balance := new(big.Int).Div(bigRes, big.NewInt(1e18))
				fmt.Printf("reward: %vLAT\n", balance)
			}

			return nil
		}
	default:
		{
			fmt.Printf("funcName:%s is unknown!!!!", funcName)
			return nil
		}
	}

	result, _ := handleCall(rlp, mapNameToAddress[action], nil)
	fmt.Printf("result:\n %s \n", result)

	return nil
}
