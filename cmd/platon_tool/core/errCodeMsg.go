// Copyright 2015 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

// rlpdump is a pretty-printer for RLP data.
package core

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/PlatONnetwork/PlatON-Go/x/gov"

	"github.com/PlatONnetwork/PlatON-Go/x/restricting"

	"github.com/PlatONnetwork/PlatON-Go/x/slashing"

	"github.com/PlatONnetwork/PlatON-Go/x/staking"

	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

func GetMsgByErrCode(errCode string) string {

	// parse platon ppos tx receipt errCode
	if errCode != "" {
		data, err := hex.DecodeString(strings.TrimPrefix(errCode, "0x"))
		if err != nil {
			die(err)
		}

		var args [][]byte
		if err := rlp.Decode(bytes.NewReader(data), &args); nil != err {
			die(err)
		}

		stakingErrCode := map[uint32]string{
			staking.ErrWrongBlsPubKey.Code:            staking.ErrWrongBlsPubKey.Msg,
			staking.ErrWrongBlsPubKeyProof.Code:       staking.ErrWrongBlsPubKeyProof.Msg,
			staking.ErrDescriptionLen.Code:            staking.ErrDescriptionLen.Msg,
			staking.ErrWrongProgramVersionSign.Code:   staking.ErrWrongProgramVersionSign.Msg,
			staking.ErrProgramVersionTooLow.Code:      staking.ErrProgramVersionTooLow.Msg,
			staking.ErrDeclVsFialedCreateCan.Code:     staking.ErrDeclVsFialedCreateCan.Msg,
			staking.ErrNoSameStakingAddr.Code:         staking.ErrNoSameStakingAddr.Msg,
			staking.ErrStakeVonTooLow.Code:            staking.ErrStakeVonTooLow.Msg,
			staking.ErrCanAlreadyExist.Code:           staking.ErrCanAlreadyExist.Msg,
			staking.ErrCanNoExist.Code:                staking.ErrCanNoExist.Msg,
			staking.ErrCanStatusInvalid.Code:          staking.ErrCanStatusInvalid.Msg,
			staking.ErrIncreaseStakeVonTooLow.Code:    staking.ErrIncreaseStakeVonTooLow.Msg,
			staking.ErrDelegateVonTooLow.Code:         staking.ErrDelegateVonTooLow.Msg,
			staking.ErrAccountNoAllowToDelegate.Code:  staking.ErrAccountNoAllowToDelegate.Msg,
			staking.ErrCanNoAllowDelegate.Code:        staking.ErrCanNoAllowDelegate.Msg,
			staking.ErrWithdrewDelegateVonTooLow.Code: staking.ErrWithdrewDelegateVonTooLow.Msg,
			staking.ErrDelegateNoExist.Code:           staking.ErrDelegateNoExist.Msg,
			staking.ErrWrongVonOptType.Code:           staking.ErrWrongVonOptType.Msg,
			staking.ErrAccountVonNoEnough.Code:        staking.ErrAccountVonNoEnough.Msg,
			staking.ErrBlockNumberDisordered.Code:     staking.ErrBlockNumberDisordered.Msg,
			staking.ErrDelegateVonNoEnough.Code:       staking.ErrDelegateVonNoEnough.Msg,
			staking.ErrWrongWithdrewDelVonCalc.Code:   staking.ErrWrongWithdrewDelVonCalc.Msg,
			staking.ErrValidatorNoExist.Code:          staking.ErrValidatorNoExist.Msg,
			staking.ErrWrongFuncParams.Code:           staking.ErrWrongFuncParams.Msg,
			staking.ErrWrongSlashType.Code:            staking.ErrWrongSlashType.Msg,
			staking.ErrSlashVonOverflow.Code:          staking.ErrSlashVonOverflow.Msg,
			staking.ErrWrongSlashVonCalc.Code:         staking.ErrWrongSlashVonCalc.Msg,
			staking.ErrGetVerifierList.Code:           staking.ErrGetVerifierList.Msg,
			staking.ErrGetValidatorList.Code:          staking.ErrGetValidatorList.Msg,
			staking.ErrGetCandidateList.Code:          staking.ErrGetCandidateList.Msg,
			staking.ErrGetDelegateRelated.Code:        staking.ErrGetDelegateRelated.Msg,
			staking.ErrQueryCandidateInfo.Code:        staking.ErrQueryCandidateInfo.Msg,
			staking.ErrQueryDelegateInfo.Code:         staking.ErrQueryDelegateInfo.Msg,
		}

		slashingErrCode := map[uint32]string{
			slashing.ErrDuplicateSignVerify.Code: slashing.ErrDuplicateSignVerify.Msg,
			slashing.ErrSlashingExist.Code:       slashing.ErrSlashingExist.Msg,
			slashing.ErrBlockNumberTooHigh.Code:  slashing.ErrBlockNumberTooHigh.Msg,
			slashing.ErrIntervalTooLong.Code:     slashing.ErrIntervalTooLong.Msg,
			slashing.ErrGetCandidate.Code:        slashing.ErrGetCandidate.Msg,
			slashing.ErrAddrMismatch.Code:        slashing.ErrAddrMismatch.Msg,
			slashing.ErrNodeIdMismatch.Code:      slashing.ErrNodeIdMismatch.Msg,
			slashing.ErrBlsPubKeyMismatch.Code:   slashing.ErrBlsPubKeyMismatch.Msg,
			slashing.ErrSlashingFail.Code:        slashing.ErrSlashingFail.Msg,
			slashing.ErrNotValidator.Code:        slashing.ErrNotValidator.Msg,
			slashing.ErrSameAddr.Code:            slashing.ErrSameAddr.Msg,
		}

		restrictingErrCode := map[uint32]string{
			restricting.ErrParamEpochInvalid.Code:                 restricting.ErrParamEpochInvalid.Msg,
			restricting.ErrCountRestrictPlansInvalid.Code:         restricting.ErrCountRestrictPlansInvalid.Msg,
			restricting.ErrLockedAmountTooLess.Code:               restricting.ErrLockedAmountTooLess.Msg,
			restricting.ErrBalanceNotEnough.Code:                  restricting.ErrBalanceNotEnough.Msg,
			restricting.ErrAccountNotFound.Code:                   restricting.ErrAccountNotFound.Msg,
			restricting.ErrSlashingTooMuch.Code:                   restricting.ErrSlashingTooMuch.Msg,
			restricting.ErrStakingAmountEmpty.Code:                restricting.ErrStakingAmountEmpty.Msg,
			restricting.ErrPledgeLockFundsAmountLessThanZero.Code: restricting.ErrPledgeLockFundsAmountLessThanZero.Msg,
			restricting.ErrReturnLockFundsAmountLessThanZero.Code: restricting.ErrReturnLockFundsAmountLessThanZero.Msg,
			restricting.ErrSlashingAmountLessThanZero.Code:        restricting.ErrSlashingAmountLessThanZero.Msg,
			restricting.ErrCreatePlanAmountLessThanZero.Code:      restricting.ErrCreatePlanAmountLessThanZero.Msg,
			restricting.ErrStakingAmountInvalid.Code:              restricting.ErrStakingAmountInvalid.Msg,
			restricting.ErrRestrictBalanceNotEnough.Code:          restricting.ErrRestrictBalanceNotEnough.Msg,
		}

		govErrCode := map[uint32]string{
			gov.ActiveVersionError.Code:                gov.ActiveVersionError.Msg,
			gov.VoteOptionError.Code:                   gov.VoteOptionError.Msg,
			gov.ProposalTypeError.Code:                 gov.ProposalTypeError.Msg,
			gov.ProposalIDEmpty.Code:                   gov.ProposalIDEmpty.Msg,
			gov.ProposalIDExist.Code:                   gov.ProposalIDExist.Msg,
			gov.ProposalNotFound.Code:                  gov.ProposalNotFound.Msg,
			gov.PIPIDEmpty.Code:                        gov.PIPIDEmpty.Msg,
			gov.PIPIDExist.Code:                        gov.PIPIDExist.Msg,
			gov.EndVotingRoundsTooSmall.Code:           gov.EndVotingRoundsTooSmall.Msg,
			gov.EndVotingRoundsTooLarge.Code:           gov.EndVotingRoundsTooLarge.Msg,
			gov.NewVersionError.Code:                   gov.NewVersionError.Msg,
			gov.VotingVersionProposalExist.Code:        gov.VotingVersionProposalExist.Msg,
			gov.PreActiveVersionProposalExist.Code:     gov.PreActiveVersionProposalExist.Msg,
			gov.VotingCancelProposalExist.Code:         gov.VotingCancelProposalExist.Msg,
			gov.TobeCanceledProposalNotFound.Code:      gov.TobeCanceledProposalNotFound.Msg,
			gov.TobeCanceledProposalTypeError.Code:     gov.TobeCanceledProposalTypeError.Msg,
			gov.TobeCanceledProposalNotAtVoting.Code:   gov.TobeCanceledProposalNotAtVoting.Msg,
			gov.ProposerEmpty.Code:                     gov.ProposerEmpty.Msg,
			gov.VerifierInfoNotFound.Code:              gov.VerifierInfoNotFound.Msg,
			gov.VerifierStatusInvalid.Code:             gov.VerifierStatusInvalid.Msg,
			gov.TxSenderDifferFromStaking.Code:         gov.TxSenderDifferFromStaking.Msg,
			gov.TxSenderIsNotVerifier.Code:             gov.TxSenderIsNotVerifier.Msg,
			gov.TxSenderIsNotCandidate.Code:            gov.TxSenderIsNotCandidate.Msg,
			gov.VersionSignError.Code:                  gov.VersionSignError.Msg,
			gov.VerifierNotUpgraded.Code:               gov.VerifierNotUpgraded.Msg,
			gov.ProposalNotAtVoting.Code:               gov.ProposalNotAtVoting.Msg,
			gov.VoteDuplicated.Code:                    gov.VoteDuplicated.Msg,
			gov.DeclareVersionError.Code:               gov.DeclareVersionError.Msg,
			gov.NotifyStakingDeclaredVersionError.Code: gov.NotifyStakingDeclaredVersionError.Msg,
			gov.TallyResultNotFound.Code:               gov.TallyResultNotFound.Msg,
			gov.UnsupportedGovernParam.Code:            gov.UnsupportedGovernParam.Msg,
			gov.VotingParamProposalExist.Code:          gov.VotingParamProposalExist.Msg,
			gov.GovernParamValueError.Code:             gov.GovernParamValueError.Msg,
			gov.ParamProposalIsSameValue.Code:          gov.ParamProposalIsSameValue.Msg,
		}

		codeStr := string(args[0])
		code, err := strconv.Atoi(codeStr)
		if nil != err {
			die(err)
		}

		var Msg string

		if msg, ok := stakingErrCode[uint32(code)]; ok {
			Msg = msg
		}
		if msg, ok := slashingErrCode[uint32(code)]; ok {
			Msg = msg
		}
		if msg, ok := restrictingErrCode[uint32(code)]; ok {
			Msg = msg
		}
		if msg, ok := govErrCode[uint32(code)]; ok {
			Msg = msg
		}

		fmt.Printf("error message: %s !!!\n", Msg)
		fmt.Println()

		return Msg
	}

	return ""
}

func die(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
	os.Exit(1)
}
