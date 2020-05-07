package reward

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

type DelegateRewardInfo struct {
	NodeID     discover.NodeID `json:"nodeID"`
	StakingNum uint64          `json:"stakingNum"`
	Reward     *big.Int        `json:"reward" rlp:"nil"`
}

func TestDecreaseDelegateReward(t *testing.T) {
	var receives []DelegateRewardReceipt
	var list DelegateRewardPerList

	receives = []DelegateRewardReceipt{

		{big.NewInt(200), 2},
		{big.NewInt(550), 3},
		{big.NewInt(400), 4},
		{big.NewInt(400), 5},
		{big.NewInt(800), 6},

		{big.NewInt(600), 7},
	}

	list.Pers = []*DelegateRewardPer{
		&DelegateRewardPer{big.NewInt(300), 1, nil, nil},
		&DelegateRewardPer{big.NewInt(500), 2, nil, nil},
		&DelegateRewardPer{big.NewInt(550), 3, nil, nil},
		&DelegateRewardPer{big.NewInt(800), 4, nil, nil},
		&DelegateRewardPer{big.NewInt(550), 5, nil, nil},
	}
	index := list.DecreaseTotalAmount(receives)
	if index != 4 {
		t.Errorf("receives index is wrong,%v", index)
	}
	if list.Pers[1].Left.Cmp(big.NewInt(300)) != 0 {
		t.Errorf("first Left  is wrong,%v", list.Pers[1].Left)
	}

	if list.Pers[len(list.Pers)-1].Left.Cmp(big.NewInt(150)) != 0 {
		t.Errorf("latest Left  is wrong,%v", list.Pers[1].Left)
	}
}

func TestSize(t *testing.T) {
	delegate := new(big.Int).Mul(new(big.Int).SetInt64(10000000), big.NewInt(params.LAT))
	reward, _ := new(big.Int).SetString("135840374364973262032076", 10)
	per := new(big.Int).Div(reward, delegate)
	key := DelegateRewardPerKey(discover.MustHexID("0aa9805681d8f77c05f317efc141c97d5adb511ffb51f5a251d2d7a4a3a96d9a12adf39f06b702f0ccdff9eddc1790eb272dca31b0c47751d49b5931c58701e7"), 100, 10)

	list := NewDelegateRewardPerList()
	for i := 0; i < DelegateRewardPerLength; i++ {
		list.AppendDelegateRewardPer(NewDelegateRewardPer(uint64(i), per, delegate))
	}
	val, err := rlp.EncodeToBytes(list)
	if err != nil {
		t.Error(err)
		return
	}
	length := len(key) + len(val)

	log.Print("size of per", length*101/(1024*1024))

}

func Test3(t *testing.T) {
	strPubKey := "063effe3e402551a3eca044b6f1da86574a8e6729016ffb900a543825475bc63de056e8c2dd8fab6bfc59a7c9b9af6193f6d633b43c8cddcc8afaa751101ab5c"
	bytesPubKey, _ := hex.DecodeString(strPubKey)
	address := common.BytesToAddress(crypto.Keccak256(bytesPubKey)[12:]).Hex()
	fmt.Printf("strPubKey:%v\n", strPubKey)
	fmt.Printf("address:%v\n", address)
}
