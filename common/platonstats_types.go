package common

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strings"

	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/log"
)

//uint32(年度<<16 | 结算周期<<8 | 共识周期)
// 年度：0-中间；1-开始；2-结束
// 结算周期：0-中间；1-开始；2-结束
// 共识周期：0-中间；1-开始；2-选举；3-结束
//type BlockType uint32

type BlockType uint8

type NodeID [512 / 8]byte

type Input []byte

var nodeIdT = reflect.TypeOf(NodeID{})

// MarshalText returns the hex representation of a.
func (a NodeID) MarshalText() ([]byte, error) {
	return hexutil.Bytes(a[:]).MarshalText()
}

// UnmarshalText parses a hash in hex syntax.
func (a *NodeID) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("common.NodeID", input, a[:])
}

// UnmarshalJSON parses a hash in hex syntax.
func (a *NodeID) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(nodeIdT, input, a[:])
}
func (n NodeID) TerminalString() string {
	return hex.EncodeToString(n[:8])
}

var inputT = reflect.TypeOf(Input{})

// MarshalText returns the hex representation of a.
func (a Input) MarshalText() ([]byte, error) {
	return hexutil.Bytes(a[:]).MarshalText()
}

// UnmarshalText parses a hash in hex syntax.
func (a *Input) UnmarshalText(input []byte) error {
	hexBytes, err := hexutil.Decode(string(input[1 : len(input)-1]))
	if err != nil {
		return err
	}
	aa := make(Input, len(hexBytes))

	err = hexutil.UnmarshalFixedText("common.Input", input, aa)
	if err != nil {
		return err
	}
	a = &aa
	return nil
}

// UnmarshalJSON parses a hash in hex syntax.
func (a *Input) UnmarshalJSON(input []byte) error {
	//string(input)="0x0102030405", so, firstly remove the "", and then, to decode it.
	hexBytes, err := hexutil.Decode(string(input[1 : len(input)-1]))
	if err != nil {
		return err
	}
	aa := make(Input, len(hexBytes))
	err = hexutil.UnmarshalFixedJSON(inputT, input, aa)
	if err != nil {
		return err
	}
	a = &aa
	return nil
}

func MustHexID(in string) NodeID {
	id, err := HexID(in)
	if err != nil {
		panic(err)
	}
	return id
}

func HexID(in string) (NodeID, error) {
	var id NodeID
	b, err := hex.DecodeString(strings.TrimPrefix(in, "0x"))
	if err != nil {
		return id, err
	} else if len(b) != len(id) {
		return id, fmt.Errorf("wrong length, want %d hex chars", len(id)*2)
	}
	copy(id[:], b)
	return id, nil
}

const (
	GenesisBlock BlockType = iota
	GeneralBlock
	ConsensusBeginBlock
	ConsensusElectionBlock
	ConsensusEndBlock
	EpochBeginBlock
	EpochEndBlock
	EndOfYear
)

type EmbedTransferTx struct {
	TxHash Hash     `json:"txHash,omitempty"`
	From   Address  `json:"from,omitempty"`
	To     Address  `json:"to,omitempty"`
	Amount *big.Int `json:"amount,omitempty"`
}

type EmbedContractTx struct {
	TxHash          Hash    `json:"txHash,omitempty"`
	From            Address `json:"from,omitempty"`
	ContractAddress Address `json:"contractAddress,omitempty"`
	Input           string  `json:"input,omitempty"` //hex string
}

type EcCommonConfig struct {
	MaxEpochMinutes     uint64 `json:"maxEpochMinutes"`     // 结算周期最大值（分钟）
	NodeBlockTimeWindow uint64 `json:"nodeBlockTimeWindow"` // 出块窗口时间 (秒)
	PerRoundBlocks      uint64 `json:"perRoundBlocks"`      // 一个共识论，每个节点出块数量
	MaxConsensusVals    uint64 `json:"maxConsensusVals"`    // 每个共识论，验证人数量的最大值
	AdditionalCycleTime uint64 `json:"additionalCycleTime"` // 增发周期 (分钟)
}

type EcStakingConfig struct {
	StakeThreshold          *big.Int `json:"stakeThreshold"`          // 质押门槛（LAT)
	OperatingThreshold      *big.Int `json:"operatingThreshold"`      // //增加/减少委托，或者增加质押时，允许的最小数量；当减少委托时，如果剩余委托数量小于此值，这此次减少委托操作，将导致撤销所有委托。
	MaxValidators           uint64   `json:"maxValidators"`           // 每个结算周期最多的备选节点数量（101个）
	UnStakeFreezeDuration   uint64   `json:"unStakeFreezeDuration"`   // 解除质押时，资金将被冻结的结算周期数量（注意节点如参与治理投票的情况）
	RewardPerMaxChangeRange uint16   `json:"rewardPerMaxChangeRange"` // 修改质押信息时，可以修改委托分红比例。但是和原比例的偏差有个允许范围。
	RewardPerChangeInterval uint16   `json:"rewardPerChangeInterval"` // 修改质押信息时，可以修改委托分红比例。但是需要和上次修改分红比例间隔一段时间(epoch)
}

type EcRewardConfig struct {
	NewBlockRate          uint64 `json:"newBlockRate"`          // 在计算下颚结算周期的奖励时，出块奖励占总奖励的百分比。（剩下的就是质押奖励）
	PlatONFoundationYear  uint32 `json:"platonFoundationYear"`  // 基金会从第几次参与增发（为了鼓励开发者，PlatON开始几次的增发金额都转入开发者基金，后续的增发，才有部分转入基金会）。链的创世块中已做了一次（增）发行。所以配置项值需从2开始，当配置成2时，链第一年末的增发就需分配资金到基金会。实际上，配置成1，和配置成2的效果一样。
	IncreaseIssuanceRatio uint16 `json:"increaseIssuanceRatio"` // 增发比例，基数是目前的发行总量（也即上次增发后的总量）
}

type EcSlashingConfig struct {
	SlashFractionDuplicateSign uint32 `json:"slashFractionDuplicateSign"` // 节点双签时的处罚比例（1%%,基数是有效质押）
	DuplicateSignReportReward  uint32 `json:"duplicateSignReportReward"`  // 节点举报其它节点双签的奖励百分比（1%，基数是双签时的处罚金)
	MaxEvidenceAge             uint32 `json:"maxEvidenceAge"`             // 法定证据追诉期（epoch数量），超过这个期限，将不再惩罚。
	SlashBlocksReward          uint32 `json:"slashBlocksReward"`          // 0出块惩罚时，惩罚多少个区块的出块奖励
	ZeroProduceCumulativeTime  uint16 `json:"zeroProduceCumulativeTime"`  // 0出块统计的时间范围（共识论）
	ZeroProduceNumberThreshold uint16 `json:"zeroProduceNumberThreshold"` // 在0出块统计时间内，节点0出块次数达到此值，将被处罚。
	ZeroProduceFreezeDuration  uint64 `json:"zeroProduceFreezeDuration"`  // 节点在0出块惩罚后，如果剩余质押金足够（大于质押金最低要求），节点将被冻结指定的时间（epoch），冻结期满后，质押状态将重新恢复到正常状态。如果节点先被0出块处罚，接着被双签举报，那冻结器满后，会被解质押。
}

type GenesisData struct {
	ChainID                   *big.Int           `json:"chainID,omitempty"`           //链ID
	PlatONFundAccount         Address            `json:"platONFundAccount,omitempty"` //PlatON基金会地址
	CDFAccount                Address            `json:"cDFAccount,omitempty"`        //开发者基金地址
	IssueAmount               *big.Int           `json:"issueAmount,omitempty"`       //发行金额
	EcCommonConfig            *EcCommonConfig    `json:"ecCommonConfig,omitempty"`    //配置项
	EcStakingConfig           *EcStakingConfig   `json:"ecStakingConfig,omitempty"`   //配置项
	EcRewardConfig            *EcRewardConfig    `json:"ecStakingConfig,omitempty"`   //配置项
	EcSlashingConfig          *EcSlashingConfig  `json:"ecSlashingConfig,omitempty"`  //配置项
	AllocItemList             []*AllocItem       `json:"allocItemList,omitempty"`
	StakingItemList           []*StakingItem     `json:"stakingItemList,omitempty"`
	RestrictingCreateItemList []*RestrictingItem `json:"restrictingCreateItemList,omitempty"`
	InitFundItemList          []*InitFundItem    `json:"initFundItemList,omitempty"`
}
type AllocItem struct {
	Address Address  `json:"address,omitempty"`
	Amount  *big.Int `json:"amount,omitempty"`
}

type StakingItem struct {
	NodeID         NodeID   `json:"nodeID,omitempty"`
	StakingAddress Address  `json:"stakingAddress,omitempty"`
	BenefitAddress Address  `json:"benefitAddress,omitempty"`
	NodeName       string   `json:"nodeName,omitempty"`
	Amount         *big.Int `json:"amount,omitempty"`
}

type RestrictingItem struct {
	From        Address    `json:"from,omitempty"`
	DestAddress Address    `json:"destAddress,omitempty"`
	Plans       []*big.Int `json:"plans,omitempty"`
}

type InitFundItem struct {
	From   Address  `json:"from,omitempty"`
	To     Address  `json:"to,omitempty"`
	Amount *big.Int `json:"amount,omitempty"`
}

func (g *GenesisData) AddEconomicConfig(common *EcCommonConfig, staking *EcStakingConfig,
	reward *EcRewardConfig, slashing *EcSlashingConfig) {
	g.EcCommonConfig = common
	g.EcStakingConfig = staking
	g.EcRewardConfig = reward
	g.EcSlashingConfig = slashing
}

func (g *GenesisData) AddAllocItem(address Address, amount *big.Int) {
	g.AllocItemList = append(g.AllocItemList, &AllocItem{Address: address, Amount: amount})
}
func (g *GenesisData) AddRestrictingCreateItem(from, to Address, plans []*big.Int) {
	g.RestrictingCreateItemList = append(g.RestrictingCreateItemList, &RestrictingItem{From: from, DestAddress: to, Plans: plans})
}

func (g *GenesisData) AddInitFundItem(from, to Address, initAmount *big.Int) {
	g.InitFundItemList = append(g.InitFundItemList, &InitFundItem{From: from, To: to, Amount: initAmount})
}

func (g *GenesisData) AddStakingItem(nodeID NodeID, nodeName string, stakingAddress, benefitAddress Address, amount *big.Int) {
	g.StakingItemList = append(g.StakingItemList, &StakingItem{NodeID: nodeID, NodeName: nodeName, StakingAddress: stakingAddress, BenefitAddress: benefitAddress, Amount: amount})
}

type AdditionalIssuanceData struct {
	AdditionalNo     uint32          `json:"additionalNo,omitempty"`     //增发周期
	AdditionalBase   *big.Int        `json:"additionalBase,omitempty"`   //增发基数
	AdditionalRate   uint16          `json:"additionalRate,omitempty"`   //增发比例 单位：万分之一
	AdditionalAmount *big.Int        `json:"additionalAmount,omitempty"` //增发金额
	IssuanceItemList []*IssuanceItem `json:"issuanceItemList,omitempty"` //增发分配
}

type IssuanceItem struct {
	Address Address  `json:"address,omitempty"` //增发金额分配地址
	Amount  *big.Int `json:"amount,omitempty"`  //增发金额
}

func (d *AdditionalIssuanceData) AddIssuanceItem(address Address, amount *big.Int) {
	//todo: test
	d.IssuanceItemList = append(d.IssuanceItemList, &IssuanceItem{Address: address, Amount: amount})
}

// 分配奖励，包括出块奖励，质押奖励
// 注意：委托人不一定每次都能参与到出块奖励的分配中（共识论跨结算周期时会出现，此时节点虽然还在出块，但是可能已经不在当前结算周期的101备选人列表里了，那这个出块节点的委托人在当前结算周期，就不参与这个块的出块奖励分配）
type RewardData struct {
	BlockRewardAmount   *big.Int         `json:"blockRewardAmount,omitempty"`   //出块奖励
	DelegatorReward     bool             `json:"delegatorReward"`               //出块奖励中，分配给委托人的奖励
	StakingRewardAmount *big.Int         `json:"stakingRewardAmount,omitempty"` //一结算周期内所有101节点的质押奖励
	CandidateInfoList   []*CandidateInfo `json:"candidateInfoList,omitempty"`   //备选节点信息
}

type CandidateInfo struct {
	NodeID       NodeID  `json:"nodeId,omitempty"`       //备选节点ID
	MinerAddress Address `json:"minerAddress,omitempty"` //备选节点的矿工地址（收益地址）
}

type ZeroSlashingItem struct {
	NodeID         NodeID   `json:"nodeId,omitempty"`         //备选节点ID
	SlashingAmount *big.Int `json:"slashingAmount,omitempty"` //0出块处罚金(从质押金扣)
}

type DuplicatedSignSlashingSetting struct {
	PenaltyRatioByValidStakings uint32 `json:"penaltyRatioByValidStakings,omitempty"` //unit:1%%		//罚金 = 有效质押 * PenaltyRatioByValidStakings / 10000
	RewardRatioByPenalties      uint32 `json:"rewardRatioByPenalties,omitempty"`      //unit:1%		//给举报人的赏金=罚金 * RewardRatioByPenalties / 100
}

type StakingSetting struct {
	OperatingThreshold *big.Int `json:"operatingThreshold,omitempty"` //质押，委托操作，要求的最小数量；当某次操作后，剩余数量小于此值时，这剩余数量将随此次操作一次处理完。
}

type StakingFrozenItem struct {
	NodeID        NodeID  `json:"nodeId,omitempty"`        //备选节点ID
	NodeAddress   Address `json:"nodeAddress,omitempty"`   //备选节点地址
	FrozenEpochNo uint64  `json:"frozenEpochNo,omitempty"` //质押资金，被解冻的结算周期（此周期最后一个块的endBlocker里）
	Recovery      bool    `json:"recovery"`                //Recover=true；表示冻结期结束后，质押将变成有效质押；Recover=false, 表示冻结期结束后，质押将原来退回质押钱包（或者和锁仓合约）
}

type RestrictingReleaseItem struct {
	DestAddress   Address  `json:"destAddress,omitempty,omitempty"` //释放地址
	ReleaseAmount *big.Int `json:"releaseAmount,omitempty"`         //释放金额
	LackingAmount *big.Int `json:"lackingAmount,omitempty"`         //欠释放金额
}

//todo:改名
//撤消委托后领取的奖励（全部减持）
type WithdrawDelegation struct {
	TxHash          Hash     `json:"txHash,omitempty"`                    //委托用户撤销节点的全部委托的交易HASH
	DelegateAddress Address  `json:"delegateAddress,omitempty,omitempty"` //委托用户地址
	NodeID          NodeID   `json:"nodeId,omitempty"`                    //委托用户委托的节点ID
	RewardAmount    *big.Int `json:"rewardAmount,omitempty"`              //委托用户从此节点获取的全部委托奖励
}

//处理委托
type FixDelegation struct {
	NodeID                              NodeID   `json:"nodeId,omitempty"`                              //委托用户委托的节点ID
	StakingBlockNumber                  uint64   `json:"stakingBlockNumber,omitempty"`                  //委托用户委托的节点ID
	ImproperValidRestrictingAmount      *big.Int `json:"improperValidRestrictingAmount,omitempty"`      //需要退回的生效期的，挪用的锁仓金额
	ImproperHesitatingRestrictingAmount *big.Int `json:"improperHesitatingRestrictingAmount,omitempty"` //需要退回的犹豫期的，挪用的锁仓金额
	Withdraw                            bool     `json:"withdraw"`                                      //true/false
	RewardAmount                        *big.Int `json:"rewardAmount,omitempty"`                        //Withdraw=true时，委托用户从此节点获取的全部委托奖励
}

//处理质押
type FixStaking struct {
	NodeID                              NodeID   `json:"nodeId,omitempty"`                              //节点ID
	StakingBlockNumber                  uint64   `json:"stakingBlockNumber,omitempty"`                  //节点首次质押的快高
	ImproperValidRestrictingAmount      *big.Int `json:"improperValidRestrictingAmount,omitempty"`      //需要退回的生效期的，挪用的锁仓金额
	ImproperHesitatingRestrictingAmount *big.Int `json:"improperHesitatingRestrictingAmount,omitempty"` //需要退回的犹豫期的，挪用的锁仓金额
	FurtherOperation                    string   `json:"furtherOperation,omitempty"`                    //在退回挪用的锁仓金额后的进一步操作，NOP：已经撤消并在冻结期,则不需要额外做什么；WITHDRAW：需要撤消并冻结1个结算周期，REDUCE：继续质押
}

//处理Issue1625
type FixIssue1625 struct {
	FixDelegationList []*FixDelegation `json:"fixDelegationList,omitempty"` //节点ID
	FixStakingList    []*FixStaking    `json:"fixStakingList,omitempty"`    //节点首次质押的快高
}

var ExeBlockDataCollector = make(map[uint64]*ExeBlockData)

func PopExeBlockData(blockNumber uint64) *ExeBlockData {
	if exeBlockData, ok := ExeBlockDataCollector[blockNumber]; ok && exeBlockData != nil {
		delete(ExeBlockDataCollector, blockNumber)
		return exeBlockData
	}
	return nil
}

func InitExeBlockData(blockNumber uint64) {
	exeBlockData := &ExeBlockData{
		ZeroSlashingItemList:       make([]*ZeroSlashingItem, 0),
		StakingFrozenItemList:      make([]*StakingFrozenItem, 0),
		RestrictingReleaseItemList: make([]*RestrictingReleaseItem, 0),
		EmbedTransferTxList:        make([]*EmbedTransferTx, 0),
		EmbedContractTxList:        make([]*EmbedContractTx, 0),
		FixIssue1625Map:            make(map[Address]*FixIssue1625),
	}

	ExeBlockDataCollector[blockNumber] = exeBlockData
}

func GetExeBlockData(blockNumber uint64) *ExeBlockData {
	return ExeBlockDataCollector[blockNumber]
}

type ExeBlockData struct {
	ActiveVersion                 string                         `json:"activeVersion,omitempty"` //如果当前块有升级提案生效，则填写新版本,0.14.0
	AdditionalIssuanceData        *AdditionalIssuanceData        `json:"additionalIssuanceData,omitempty"`
	RewardData                    *RewardData                    `json:"rewardData,omitempty"`
	ZeroSlashingItemList          []*ZeroSlashingItem            `json:"zeroSlashingItemList,omitempty"`
	DuplicatedSignSlashingSetting *DuplicatedSignSlashingSetting `json:"duplicatedSignSlashingSetting,omitempty"`
	StakingSetting                *StakingSetting                `json:"stakingSetting,omitempty"`
	StakingFrozenItemList         []*StakingFrozenItem           `json:"stakingFrozenItemList,omitempty"`
	RestrictingReleaseItemList    []*RestrictingReleaseItem      `json:"restrictingReleaseItemList,omitempty"`
	EmbedTransferTxList           []*EmbedTransferTx             `json:"embedTransferTxList,omitempty"`    //一个显式交易引起的内置转账交易：一般有两种情况：1是部署，或者调用合约时，带上了value，则这个value会转账给合约地址；2是调用合约，合约内部调用transfer()函数完成转账
	EmbedContractTxList           []*EmbedContractTx             `json:"embedContractTxList,omitempty"`    //一个显式交易引起的内置合约交易。这个显式交易显然也是个合约交易，在这个合约里，又调用了其他合约（包括内置合约）
	WithdrawDelegationList        []*WithdrawDelegation          `json:"withdrawDelegationList,omitempty"` //当委托用户撤回节点的全部委托时，需要的统计信息（由于Alaya在运行中，只能兼容Alaya的bug）
	FixIssue1625Map               map[Address]*FixIssue1625      `json:"fixIssue1625Map,omitempty"`
}

func CollectAdditionalIssuance(blockNumber uint64, additionalIssuanceData *AdditionalIssuanceData) {
	if exeBlockData, ok := ExeBlockDataCollector[blockNumber]; ok && exeBlockData != nil {
		json, _ := json.Marshal(additionalIssuanceData)
		log.Debug("CollectAdditionalIssuance", "blockNumber", blockNumber, "additionalIssuanceData", string(json))
		exeBlockData.AdditionalIssuanceData = additionalIssuanceData
	}
}

func CollectStakingFrozenItem(blockNumber uint64, nodeId NodeID, nodeAddress NodeAddress, frozenEpochNo uint64, recovery bool) {
	if exeBlockData, ok := ExeBlockDataCollector[blockNumber]; ok && exeBlockData != nil {
		log.Debug("CollectStakingFrozenItem", "blockNumber", blockNumber, "nodeId", Bytes2Hex(nodeId[:]), "nodeAddress", nodeAddress.Hex(), "frozenEpochNo", frozenEpochNo, "recovery", recovery)
		exeBlockData.StakingFrozenItemList = append(exeBlockData.StakingFrozenItemList, &StakingFrozenItem{NodeID: nodeId, NodeAddress: Address(nodeAddress), FrozenEpochNo: frozenEpochNo, Recovery: recovery})
	}
}

func CollectRestrictingReleaseItem(blockNumber uint64, destAddress Address, releaseAmount *big.Int, lackingAmount *big.Int) {
	if exeBlockData, ok := ExeBlockDataCollector[blockNumber]; ok && exeBlockData != nil {
		log.Debug("CollectRestrictingReleaseItem", "blockNumber", blockNumber, "destAddress", destAddress, "releaseAmount", releaseAmount, "lackingAmount", lackingAmount)
		exeBlockData.RestrictingReleaseItemList = append(exeBlockData.RestrictingReleaseItemList, &RestrictingReleaseItem{DestAddress: destAddress, ReleaseAmount: releaseAmount, LackingAmount: lackingAmount})
	}
}

func CollectBlockRewardData(blockNumber uint64, blockRewardAmount *big.Int, delegatorReward bool) {
	if exeBlockData, ok := ExeBlockDataCollector[blockNumber]; ok && exeBlockData != nil {
		log.Debug("CollectBlockRewardData", "blockNumber", blockNumber, "blockRewardAmount", blockRewardAmount, "delegatorReward", delegatorReward)
		if exeBlockData.RewardData == nil {
			exeBlockData.RewardData = new(RewardData)
		}
		exeBlockData.RewardData.BlockRewardAmount = blockRewardAmount
		exeBlockData.RewardData.DelegatorReward = delegatorReward
	}
}

func CollectStakingRewardData(blockNumber uint64, stakingRewardAmount *big.Int, candidateInfoList []*CandidateInfo) {
	if exeBlockData, ok := ExeBlockDataCollector[blockNumber]; ok && exeBlockData != nil {
		log.Debug("CollectStakingRewardData", "blockNumber", blockNumber, "stakingRewardAmount", stakingRewardAmount)
		for _, candidateInfo := range candidateInfoList {
			log.Debug("nodeID:" + Bytes2Hex(candidateInfo.NodeID[:]))
		}
		if exeBlockData.RewardData == nil {
			exeBlockData.RewardData = new(RewardData)
		}
		exeBlockData.RewardData.StakingRewardAmount = stakingRewardAmount
		exeBlockData.RewardData.CandidateInfoList = candidateInfoList
	}
}

func CollectDuplicatedSignSlashingSetting(blockNumber uint64, penaltyRatioByValidStakings, rewardRatioByPenalties uint32) {
	if exeBlockData, ok := ExeBlockDataCollector[blockNumber]; ok && exeBlockData != nil {
		log.Debug("CollectDuplicatedSignSlashingSetting", "blockNumber", blockNumber, "penaltyRatioByValidStakings", penaltyRatioByValidStakings, "rewardRatioByPenalties", rewardRatioByPenalties)
		if exeBlockData.DuplicatedSignSlashingSetting == nil {
			//在同一个区块中，只要设置一次即可
			exeBlockData.DuplicatedSignSlashingSetting = &DuplicatedSignSlashingSetting{PenaltyRatioByValidStakings: penaltyRatioByValidStakings, RewardRatioByPenalties: rewardRatioByPenalties}
		}
	}
}

func CollectStakingSetting(blockNumber uint64, operatingThreshold *big.Int) {
	if exeBlockData, ok := ExeBlockDataCollector[blockNumber]; ok && exeBlockData != nil {
		log.Debug("CollectStakingSetting", "blockNumber", blockNumber, "operatingThreshold", operatingThreshold)
		if exeBlockData.StakingSetting == nil {
			exeBlockData.StakingSetting = &StakingSetting{OperatingThreshold: operatingThreshold}
		}
	}
}

func CollectZeroSlashingItem(blockNumber uint64, nodeId NodeID, slashingAmount *big.Int) {
	if exeBlockData, ok := ExeBlockDataCollector[blockNumber]; ok && exeBlockData != nil {
		log.Debug("CollectZeroSlashingItem", "blockNumber", blockNumber, "nodeId", Bytes2Hex(nodeId[:]), "slashingAmount", slashingAmount)
		exeBlockData.ZeroSlashingItemList = append(exeBlockData.ZeroSlashingItemList, &ZeroSlashingItem{NodeID: nodeId, SlashingAmount: slashingAmount})
	}
}

/*func CollectZeroSlashingItem(blockNumber uint64, zeroSlashingItemList []*ZeroSlashingItem) {
	if exeBlockData, ok := ExeBlockDataCollector[blockNumber]; ok && exeBlockData != nil {
		json, _ := json.Marshal(zeroSlashingItemList)
		log.Debug("CollectZeroSlashingItem", "blockNumber", blockNumber, "zeroSlashingItemList", string(json))
		exeBlockData.ZeroSlashingItemList = zeroSlashingItemList
	}
}*/

func CollectEmbedTransferTx(blockNumber uint64, txHash Hash, from, to Address, amount *big.Int) {
	if exeBlockData, ok := ExeBlockDataCollector[blockNumber]; ok && exeBlockData != nil {
		log.Debug("CollectEmbedTransferTx", "blockNumber", blockNumber, "txHash", txHash.Hex(), "from", from.Bech32(), "to", to.Bech32(), "amount", amount)
		amt := new(big.Int).Set(amount)
		exeBlockData.EmbedTransferTxList = append(exeBlockData.EmbedTransferTxList, &EmbedTransferTx{TxHash: txHash, From: from, To: to, Amount: amt})
	}
}

func CollectEmbedContractTx(blockNumber uint64, txHash Hash, from, contractAddress Address, input []byte) {
	if exeBlockData, ok := ExeBlockDataCollector[blockNumber]; ok && exeBlockData != nil {
		log.Debug("CollectEmbedContractTx", "blockNumber", blockNumber, "txHash", txHash.Hex(), "contractAddress", from.Bech32(), "input", Bytes2Hex(input))
		exeBlockData.EmbedContractTxList = append(exeBlockData.EmbedContractTxList, &EmbedContractTx{TxHash: txHash, From: from, ContractAddress: contractAddress, Input: Bytes2Hex(input)})
	}
}

//撤消委托时，才需要收集委托奖励总金额
func CollectWithdrawDelegation(blockNumber uint64, txHash Hash, delegateAddress Address, nodeId NodeID, delegationRewardAmount *big.Int) {
	if exeBlockData, ok := ExeBlockDataCollector[blockNumber]; ok && exeBlockData != nil {
		log.Debug("CollectWithdrawDelegation", "blockNumber", blockNumber, "txHash", txHash.Hex(), "delegateAddress", delegateAddress.Bech32(), "nodeId", Bytes2Hex(nodeId[:]), "delegationRewardAmount", delegationRewardAmount)
		amt := new(big.Int).Set(delegationRewardAmount)
		exeBlockData.WithdrawDelegationList = append(exeBlockData.WithdrawDelegationList, &WithdrawDelegation{TxHash: txHash, DelegateAddress: delegateAddress, NodeID: nodeId, RewardAmount: amt})
	}
}

func CollectActiveVersion(blockNumber uint64, newVersion uint32) {
	if exeBlockData, ok := ExeBlockDataCollector[blockNumber]; ok && exeBlockData != nil {
		log.Debug("CollectActiveVersion", "blockNumber", blockNumber, "newVersion", newVersion)
		exeBlockData.ActiveVersion = FormatVersion(newVersion)
	}
}

func CollectFixDelegation(blockNumber uint64, account Address, fixDelegation *FixDelegation) {
	if exeBlockData, ok := ExeBlockDataCollector[blockNumber]; ok && exeBlockData != nil {
		log.Debug("CollectFixDelegation", "blockNumber", blockNumber, "account", account.Bech32(), "nodeId", fixDelegation.NodeID.TerminalString(), "fixDelegation", fixDelegation)
		if _, ok := exeBlockData.FixIssue1625Map[account]; ok {
			exeBlockData.FixIssue1625Map[account].FixDelegationList = append(exeBlockData.FixIssue1625Map[account].FixDelegationList, fixDelegation)
		} else {
			//不存在
			exeBlockData.FixIssue1625Map[account] = &FixIssue1625{FixDelegationList: []*FixDelegation{fixDelegation}}
		}
	}
}

func CollectFixStaking(blockNumber uint64, account Address, fixStaking *FixStaking) {
	if exeBlockData, ok := ExeBlockDataCollector[blockNumber]; ok && exeBlockData != nil {
		log.Debug("CollectFixStaking", "blockNumber", blockNumber, "account", account.Bech32(), "nodeId", fixStaking.NodeID.TerminalString(), "fixStaking", fixStaking)
		if _, ok := exeBlockData.FixIssue1625Map[account]; ok {
			exeBlockData.FixIssue1625Map[account].FixStakingList = append(exeBlockData.FixIssue1625Map[account].FixStakingList, fixStaking)
		} else {
			//不存在
			exeBlockData.FixIssue1625Map[account] = &FixIssue1625{FixStakingList: []*FixStaking{fixStaking}}
		}
	}
}