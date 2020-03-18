// 生成经济模型rpc接口的rlp data

package core

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/restricting"
)

// createStaking
type Ppos_1000 struct {
	Typ                uint16
	BenefitAddress     common.Address
	NodeId             discover.NodeID
	ExternalId         string
	NodeName           string
	Website            string
	Details            string
	Amount             *big.Int
	ProgramVersion     uint32
	ProgramVersionSign common.VersionSign
	BlsPubKey          bls.PublicKeyHex
	BlsProof           bls.SchnorrProofHex
}

// editorCandidate
type Ppos_1001 struct {
	BenefitAddress common.Address
	NodeId         discover.NodeID
	ExternalId     string
	NodeName       string
	Website        string
	Details        string
}

// increaseStaking
type Ppos_1002 struct {
	NodeId discover.NodeID
	Typ    uint16
	Amount *big.Int
}

// withdrewCandidate
type Ppos_1003 struct {
	NodeId discover.NodeID
}

// delegate
type Ppos_1004 struct {
	Typ    uint16
	NodeId discover.NodeID
	Amount *big.Int
}

// withdrewDelegate
type Ppos_1005 struct {
	StakingBlockNum uint64
	NodeId          discover.NodeID
	Amount          *big.Int
}

// getRelatedListByDelAddr
type Ppos_1103 struct {
	Addr common.Address
}

// getDelegateInfo
type Ppos_1104 struct {
	StakingBlockNum uint64
	DelAddr         common.Address
	NodeId          discover.NodeID
}

// getCandidateInfo
type Ppos_1105 struct {
	NodeId discover.NodeID
}

// submitText
type Ppos_2000 struct {
	Verifier discover.NodeID
	PIPID    string
	//EndVotingRounds uint64
}

// submitVersion
type Ppos_2001 struct {
	Verifier        discover.NodeID
	PIPID           string
	NewVersion      uint32
	EndVotingRounds uint64
}

// submitParam
type Ppos_2002 struct {
	Verifier discover.NodeID
	PIPID    string
	Module   string
	Name     string
	NewValue string
}

// submitCancel
type Ppos_2005 struct {
	Verifier        discover.NodeID
	PIPID           string
	EndVotingRounds uint64
	TobeCanceled    common.Hash
}

// vote
type Ppos_2003 struct {
	Verifier       discover.NodeID
	ProposalID     common.Hash
	Option         uint8
	ProgramVersion uint32
	VersionSign    common.VersionSign
}

type Ppos_20031 struct {
	Verifier       discover.NodeID
	ProposalID     common.Hash
	Option         uint8
	ProgramVersion uint32
	VersionSign    common.VersionSign
}

//declareVersion
type Ppos_2004 struct {
	Verifier       discover.NodeID
	ProgramVersion uint32
	VersionSign    common.VersionSign
}

// getProposal
type Ppos_2100 struct {
	ProposalID common.Hash
}

// getTallyResult
type Ppos_2101 struct {
	ProposalID common.Hash
}

// listProposal
type Ppos_2102 struct {
}

// getActiveVersion
type Ppos_2103 struct {
}

// getGovernParamValue
type Ppos_2104 struct {
	Module string
	Name   string
}

// getAccuVerifiersCount
type Ppos_2105 struct {
	ProposalID common.Hash
	BlockHash  common.Hash
}

// listGovernParam
type Ppos_2106 struct {
	Module string
}

// reportDuplicateSign
type Ppos_3000 struct {
	Data string
}

// checkDuplicateSign
type Ppos_3001 struct {
	Etype       uint32
	Addr        common.Address
	BlockNumber uint64
}

// CreateRestrictingPlan
type Ppos_4000 struct {
	Account common.Address
	Plans   []restricting.RestrictingPlan
}

// GetRestrictingInfo
type Ppos_4100 struct {
	Account common.Address
}

type decDataConfig struct {
	P1000  Ppos_1000
	P1001  Ppos_1001
	P1002  Ppos_1002
	P1003  Ppos_1003
	P1004  Ppos_1004
	P1005  Ppos_1005
	P1103  Ppos_1103
	P1104  Ppos_1104
	P1105  Ppos_1105
	P2000  Ppos_2000
	P2001  Ppos_2001
	P2002  Ppos_2002
	P2005  Ppos_2005
	P2003  Ppos_2003
	P20031 []Ppos_20031
	P2004  Ppos_2004
	P2100  Ppos_2100
	P2101  Ppos_2101
	P2102  Ppos_2102
	P2103  Ppos_2103
	P2104  Ppos_2104
	P2105  Ppos_2105
	P2106  Ppos_2106
	P3000  Ppos_3000
	P3001  Ppos_3001
	P4000  Ppos_4000
	P4100  Ppos_4100
}

func getRlpData(funcType uint16, cfg *decDataConfig, selParams interface{}) string {
	rlpData := ""

	var params [][]byte
	params = make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(funcType)
	params = append(params, fnType)

	switch funcType {
	case 1000:
		{
			typ, _ := rlp.EncodeToBytes(cfg.P1000.Typ)
			benefitAddress, _ := rlp.EncodeToBytes(cfg.P1000.BenefitAddress.Bytes())
			nodeId, _ := rlp.EncodeToBytes(cfg.P1000.NodeId)
			externalId, _ := rlp.EncodeToBytes(cfg.P1000.ExternalId)
			nodeName, _ := rlp.EncodeToBytes(cfg.P1000.NodeName)
			website, _ := rlp.EncodeToBytes(cfg.P1000.Website)
			details, _ := rlp.EncodeToBytes(cfg.P1000.Details)
			amount, _ := rlp.EncodeToBytes(cfg.P1000.Amount)
			programVersion, _ := rlp.EncodeToBytes(cfg.P1000.ProgramVersion)
			programVersionSign, _ := rlp.EncodeToBytes(cfg.P1000.ProgramVersionSign)
			blsPubKey, _ := rlp.EncodeToBytes(cfg.P1000.BlsPubKey)
			blsProof, _ := rlp.EncodeToBytes(cfg.P1000.BlsProof)

			params = append(params, typ)
			params = append(params, benefitAddress)
			params = append(params, nodeId)
			params = append(params, externalId)
			params = append(params, nodeName)
			params = append(params, website)
			params = append(params, details)
			params = append(params, amount)
			params = append(params, programVersion)
			params = append(params, programVersionSign)
			params = append(params, blsPubKey)
			params = append(params, blsProof)
		}
	case 1001:
		{
			benefitAddress, _ := rlp.EncodeToBytes(cfg.P1001.BenefitAddress.Bytes())
			nodeId, _ := rlp.EncodeToBytes(cfg.P1001.NodeId)
			externalId, _ := rlp.EncodeToBytes(cfg.P1001.ExternalId)
			nodeName, _ := rlp.EncodeToBytes(cfg.P1001.NodeName)
			website, _ := rlp.EncodeToBytes(cfg.P1001.Website)
			details, _ := rlp.EncodeToBytes(cfg.P1001.Details)

			params = append(params, benefitAddress)
			params = append(params, nodeId)
			params = append(params, externalId)
			params = append(params, nodeName)
			params = append(params, website)
			params = append(params, details)
		}
	case 1002:
		{
			nodeId, _ := rlp.EncodeToBytes(cfg.P1002.NodeId)
			typ, _ := rlp.EncodeToBytes(cfg.P1002.Typ)
			amount, _ := rlp.EncodeToBytes(cfg.P1002.Amount)

			params = append(params, nodeId)
			params = append(params, typ)
			params = append(params, amount)
		}
	case 1003:
		{
			nodeId, _ := rlp.EncodeToBytes(cfg.P1003.NodeId)
			params = append(params, nodeId)
		}
	case 1004:
		{
			typ, _ := rlp.EncodeToBytes(cfg.P1004.Typ)
			nodeId, _ := rlp.EncodeToBytes(cfg.P1004.NodeId)
			amount, _ := rlp.EncodeToBytes(cfg.P1004.Amount)

			params = append(params, typ)
			params = append(params, nodeId)
			params = append(params, amount)
		}
	case 1005:
		{
			stakingBlockNum, _ := rlp.EncodeToBytes(cfg.P1005.StakingBlockNum)
			nodeId, _ := rlp.EncodeToBytes(cfg.P1005.NodeId)
			amount, _ := rlp.EncodeToBytes(cfg.P1005.Amount)

			params = append(params, stakingBlockNum)
			params = append(params, nodeId)
			params = append(params, amount)
		}
	case 1100:
	case 1101:
	case 1102:
	case 1200:
	case 1201:
	case 1202:
	case 5000:
	case 1103, 4100:
		{
			// 委托地址
			address := common.HexToAddress(selParams.(string))
			rlpAddress, _ := rlp.EncodeToBytes(address.Bytes())

			params = append(params, rlpAddress)
		}
	case 1104:
		{
			stakingBlockNum, _ := rlp.EncodeToBytes(cfg.P1104.StakingBlockNum)
			delAddr, _ := rlp.EncodeToBytes(cfg.P1104.DelAddr.Bytes())
			nodeId, _ := rlp.EncodeToBytes(cfg.P1104.NodeId)

			params = append(params, stakingBlockNum)
			params = append(params, delAddr)
			params = append(params, nodeId)
		}
	case 1105:
		{
			// 节点id
			nodeId, _ := discover.HexID(selParams.(string))
			rlpBytes, _ := rlp.EncodeToBytes(nodeId.Bytes())

			params = append(params, rlpBytes)
		}
	case 2000:
		{
			verifier, _ := rlp.EncodeToBytes(cfg.P2000.Verifier)
			pipID, _ := rlp.EncodeToBytes(cfg.P2000.PIPID)
			//endVotingRounds, _ := rlp.EncodeToBytes(cfg.P2000.EndVotingRounds)
			params = append(params, verifier)
			params = append(params, pipID)
			//params = append(params, endVotingRounds)
		}
	case 2001:
		{
			verifier, _ := rlp.EncodeToBytes(cfg.P2001.Verifier)
			pipID, _ := rlp.EncodeToBytes(cfg.P2001.PIPID)
			newVersion, _ := rlp.EncodeToBytes(cfg.P2001.NewVersion)
			endVotingRounds, _ := rlp.EncodeToBytes(cfg.P2001.EndVotingRounds)
			params = append(params, verifier)
			params = append(params, pipID)
			params = append(params, newVersion)
			params = append(params, endVotingRounds)
		}
	case 2002:
		{
			verifier, _ := rlp.EncodeToBytes(cfg.P2002.Verifier)
			pipID, _ := rlp.EncodeToBytes(cfg.P2002.PIPID)
			module, _ := rlp.EncodeToBytes(cfg.P2002.Module)
			name, _ := rlp.EncodeToBytes(cfg.P2002.Name)
			newValue, _ := rlp.EncodeToBytes(cfg.P2002.NewValue)

			params = append(params, verifier)
			params = append(params, pipID)
			params = append(params, module)
			params = append(params, name)
			params = append(params, newValue)
		}
	case 2005:
		{
			verifier, _ := rlp.EncodeToBytes(cfg.P2005.Verifier)
			pipID, _ := rlp.EncodeToBytes(cfg.P2005.PIPID)
			endVotingRounds, _ := rlp.EncodeToBytes(cfg.P2005.EndVotingRounds)
			tobeCanceled, _ := rlp.EncodeToBytes(cfg.P2005.TobeCanceled)
			params = append(params, verifier)
			params = append(params, pipID)
			params = append(params, endVotingRounds)
			params = append(params, tobeCanceled)
		}
	case 2003:
		{
			verifier, _ := rlp.EncodeToBytes(cfg.P2003.Verifier)
			proposalID, _ := rlp.EncodeToBytes(cfg.P2003.ProposalID.Bytes())
			op, _ := rlp.EncodeToBytes(cfg.P2003.Option)
			programVersion, _ := rlp.EncodeToBytes(cfg.P2003.ProgramVersion)
			versionSign, _ := rlp.EncodeToBytes(cfg.P2003.VersionSign)
			params = append(params, verifier)
			params = append(params, proposalID)
			params = append(params, op)
			params = append(params, programVersion)
			params = append(params, versionSign)
		}
	case 20031:
		{
			for i := 0; i < len(cfg.P20031); i++ {
				fnType, _ := rlp.EncodeToBytes(uint16(2003))
				verifier, _ := rlp.EncodeToBytes(cfg.P20031[i].Verifier)
				proposalID, _ := rlp.EncodeToBytes(cfg.P20031[i].ProposalID.Bytes())
				op, _ := rlp.EncodeToBytes(cfg.P20031[i].Option)
				programVersion, _ := rlp.EncodeToBytes(cfg.P20031[i].ProgramVersion)
				versionSign, _ := rlp.EncodeToBytes(cfg.P20031[i].VersionSign)
				params = make([][]byte, 0)
				params = append(params, fnType)
				params = append(params, verifier)
				params = append(params, proposalID)
				params = append(params, op)
				params = append(params, programVersion)
				params = append(params, versionSign)

				buf := new(bytes.Buffer)
				err := rlp.Encode(buf, params)
				if err != nil {
					panic(fmt.Errorf("%d encode rlp data fail: %v", funcType, err))
				} else {
					rlpData = hexutil.Encode(buf.Bytes())
					fmt.Printf("RLP= %s\n", rlpData)
				}
			}
			return ""
		}

	case 2004:
		{
			verifier, _ := rlp.EncodeToBytes(cfg.P2004.Verifier)
			programVersion, _ := rlp.EncodeToBytes(cfg.P2004.ProgramVersion)
			versionSign, _ := rlp.EncodeToBytes(cfg.P2004.VersionSign)
			params = append(params, verifier)
			params = append(params, programVersion)
			params = append(params, versionSign)
		}
	case 2100, 2101:
		{ // ProposalID
			strParam := common.HexToHash(selParams.(string))
			rlpBytes, _ := rlp.EncodeToBytes(strParam.Bytes())

			params = append(params, rlpBytes)
		}
	case 2106:
		{
			module, _ := rlp.EncodeToBytes(selParams.(string))
			params = append(params, module)
		}
	case 2102:
	case 2103:
	case 2104:
		{
			gov := selParams.(Gov)
			module, _ := rlp.EncodeToBytes(gov.Module)
			name, _ := rlp.EncodeToBytes(gov.Name)
			params = append(params, module)
			params = append(params, name)
		}
	case 2105:
		{
			proposalID, _ := rlp.EncodeToBytes(cfg.P2105.ProposalID.Bytes())
			blockHash, _ := rlp.EncodeToBytes(cfg.P2105.BlockHash.Bytes())
			params = append(params, proposalID)
			params = append(params, blockHash)
		}
	case 3000:
		{
			data, _ := rlp.EncodeToBytes(cfg.P3000.Data)
			params = append(params, data)
		}
	case 3001:
		{
			etype, _ := rlp.EncodeToBytes(cfg.P3001.Etype)
			addr, _ := rlp.EncodeToBytes(cfg.P3001.Addr.Bytes())
			blockNumber, _ := rlp.EncodeToBytes(cfg.P3001.BlockNumber)
			params = append(params, etype)
			params = append(params, addr)
			params = append(params, blockNumber)
		}
	case 3002:
	case 4000:
		{
			account, _ := rlp.EncodeToBytes(cfg.P4000.Account.Bytes())
			plans, _ := rlp.EncodeToBytes(cfg.P4000.Plans)
			params = append(params, account)
			params = append(params, plans)
		}
	case 5100:
		{
			reward := selParams.(Reward)
			account, _ := rlp.EncodeToBytes(reward.Account)
			nodeIds, _ := rlp.EncodeToBytes(reward.NodeIds)
			params = append(params, account)
			params = append(params, nodeIds)
		}
	default:
		{
			//	panic(fmt.Errorf("funcType:%d is unknown!!!!", funcType))
			fmt.Printf("funcType:%d is unknown!!!!\n", funcType)
			return ""
		}
	}

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		panic(fmt.Errorf("%d encode rlp data fail: %v", funcType, err))
	} else {
		rlpData = hexutil.Encode(buf.Bytes())
		fmt.Printf("funcType:%d rlp_data = %s\n", funcType, rlpData)
	}

	return rlpData
}

func main() {
	var cfg decDataConfig
	getRlpData(uint16(1100), &cfg, nil)
}
