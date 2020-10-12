// Copyright 2017 The go-ethereum Authors
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

package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"flag"
	"fmt"

	"github.com/PlatONnetwork/PlatON-Go/cmd/utils"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"
	"github.com/PlatONnetwork/PlatON-Go/node"
)

const (
	//These versions are meaning the current code version.
	VersionMajor = 0  // Major version component of the current release
	VersionMinor = 13 // Minor version component of the current release
	VersionPatch = 2  // Patch version component of the current release
)

type nodekeypair struct {
	PrivateKey string
	PublicKey  string
}

type blskeypair struct {
	PrivateKey string
	PublicKey  string
}

func main() {
	nodeCountFlag := flag.Int("count", 1, "Number of nodes")
	flag.Parse()

	for i := 0; i < *nodeCountFlag; i++ {
		// 生成nodeid和nodekey
		// Check if keyfile path given and make sure it doesn't already exist.
		var nodePrivateKey *ecdsa.PrivateKey
		var err error
		// generate random.
		nodePrivateKey, err = crypto.GenerateKey()
		if err != nil {
			utils.Fatalf("Failed to generate random node private key: %v", err)
		}

		// Output some information.
		out := nodekeypair{
			PublicKey:  hex.EncodeToString(crypto.FromECDSAPub(&nodePrivateKey.PublicKey)[1:]),
			PrivateKey: hex.EncodeToString(crypto.FromECDSA(nodePrivateKey)),
		}
		fmt.Printf("nodeid: %v\n", out.PublicKey)
		fmt.Printf("node prikey: %v\n", out.PrivateKey)

		// 生成版本签名
		node.GetCryptoHandler().SetPrivateKey(nodePrivateKey)
		/*
			node.GetCryptoHandler().SetPrivateKey(
				crypto.HexMustToECDSA("40a2d01c7b10d19dbdd8b0c04be82d368b3d65a0a3f35e5c9c99eb81229298f7"))
		*/

		initProgramVersion := uint32(VersionMajor<<16 | VersionMinor<<8 | VersionPatch)
		versionSign := node.GetCryptoHandler().MustSign(initProgramVersion)
		fmt.Printf("version: %v\n", initProgramVersion)
		fmt.Printf("sign: %v\n", hex.EncodeToString(versionSign))

		// 生成blspubkey和blskey
		err = bls.Init(int(bls.BLS12_381))
		if err != nil {
			utils.Fatalf("Failed to generate random bls private key: %v", err)
		}
		var blsPrivateKey bls.SecretKey
		blsPrivateKey.SetByCSPRNG()
		blsPubKey := blsPrivateKey.GetPublicKey()
		blsInfo := blskeypair{
			PrivateKey: hex.EncodeToString(blsPrivateKey.GetLittleEndian()),
			PublicKey:  hex.EncodeToString(blsPubKey.Serialize()),
		}
		fmt.Printf("bls public key: %v\n", blsInfo.PublicKey)
		fmt.Printf("bls prikey: %v\n", blsInfo.PrivateKey)

		// 生成零知识证明
		proof, _ := blsPrivateKey.MakeSchnorrNIZKP()
		proofByte, _ := proof.MarshalText()
		var proofHex bls.SchnorrProofHex
		proofHex.UnmarshalText(proofByte)
		fmt.Printf("proof: %v\n", proofHex)

		/*
			err = proof.VerifySchnorrNIZK(*blsPubKey)
			if err != nil {
				utils.Fatalf("verify bls private key: %v", err)
			}
		*/

		// 保存xlsx文件
	}
	fmt.Printf("生成虚拟节点信息成功!\n")
}
