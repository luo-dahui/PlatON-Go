/*
经济模型rpc接口查詢
platon_call
*/

package core

import (
	"encoding/hex"
	"fmt"
	"strings"

	"gopkg.in/urfave/cli.v1"
)

var (
	EcoModelCallCmd = cli.Command{
		Name:   "call_ecomodel",
		Usage:  "经济模型rpc接口查詢",
		Action: call_ecomodel,
		Flags:  EcoModelCallCmdFlags,
	}
)

// 合约名称--->合约地址
var mapNameToAddress = map[string]string{
	"staking":     "0x1000000000000000000000000000000000000002",
	"gov":         "0x1000000000000000000000000000000000000005",
	"slashing":    "0x1000000000000000000000000000000000000004",
	"restricting": "0x1000000000000000000000000000000000000001",
	"reward":      "0x1000000000000000000000000000000000000006",
}

// 接口名--->接口号
var mapNameToFuncType = map[string]uint16{
	"getVerifierList":         1100,
	"getValidatorList":        1101,
	"getCandidateList":        1102,
	"getRelatedListByDelAddr": 1103,
	"getCandidateInfo":        1105,
	"getPackageReward":        1200,
	"getStakingReward":        1201,
	"getAvgPackTime":          1202,

	"getProposal":         2100,
	"getTallyResult":      2101,
	"listProposal":        2102,
	"getActiveVersion":    2103,
	"getGovernParamValue": 2104,
	"listGovernParam":     2106,

	"ZeroProduceNodeList": 3002,
}

func handleCall(rlpdata, toAddress string) (string, error) {

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

	return string(res_data), nil
}

func call_ecomodel(c *cli.Context) error {

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
	default:
		{
			fmt.Printf("funcName:%s is unknown!!!!", funcName)
			return nil
		}
	}

	result, _ := handleCall(rlp, mapNameToAddress[action])
	fmt.Printf("result:\n %s \n", result)

	return nil
}
