package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
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

const initProgramVersion = uint32(VersionMajor<<16 | VersionMinor<<8 | VersionPatch)

type keyInfo struct {
	nodePublicKey  string
	nodePrivateKey string
	versionSign    string
	blsPublicKey   string
	blsPrivateKey  string
	blsProof       string
}

func GetCurrentPath() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func main() {

	nodeCountFlag := flag.Int("count", 0, "Number of nodes")
	flag.Parse()

	if 0 == *nodeCountFlag {
		utils.Fatalf("指定生成节点信息个数为0!!!")
	}

	sheetName := "keyInfos"
	strVersion := fmt.Sprintf("%d", initProgramVersion)
	f := excelize.NewFile()
	// 创建一个工作表
	index := f.NewSheet(sheetName)
	// 设置工作簿的默认工作表
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")
	f.SetCellValue(sheetName, "A1", "nodeid")
	f.SetCellValue(sheetName, "B1", "nodePrivateKey")
	f.SetCellValue(sheetName, "C1", "version")
	f.SetCellValue(sheetName, "D1", "version_sign")
	f.SetCellValue(sheetName, "E1", "blsPublicKey")
	f.SetCellValue(sheetName, "F1", "blsPrivateKey")
	f.SetCellValue(sheetName, "G1", "blsProof")

	nIndex := 2
	for i := 0; i < *nodeCountFlag; i++ {
		keyInfo := keyInfo{}
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
		keyInfo.nodePrivateKey = hex.EncodeToString(crypto.FromECDSA(nodePrivateKey))
		keyInfo.nodePublicKey = hex.EncodeToString(crypto.FromECDSAPub(&nodePrivateKey.PublicKey)[1:])

		// fmt.Printf("nodeid: %v\n", keyInfo.nodePublicKey)
		// fmt.Printf("node prikey: %v\n", keyInfo.nodePrivateKey)

		// 生成版本签名
		node.GetCryptoHandler().SetPrivateKey(nodePrivateKey)
		/*
			node.GetCryptoHandler().SetPrivateKey(
				crypto.HexMustToECDSA("40a2d01c7b10d19dbdd8b0c04be82d368b3d65a0a3f35e5c9c99eb81229298f7"))
		*/

		versionSign := node.GetCryptoHandler().MustSign(initProgramVersion)
		keyInfo.versionSign = hex.EncodeToString(versionSign)
		// fmt.Printf("version: %v\n", initProgramVersion)
		// fmt.Printf("sign: %v\n", keyInfo.versionSign)

		// 生成blspubkey和blskey
		err = bls.Init(int(bls.BLS12_381))
		if err != nil {
			utils.Fatalf("Failed to generate random bls private key: %v", err)
		}
		var blsPrivateKey bls.SecretKey
		blsPrivateKey.SetByCSPRNG()
		blsPubKey := blsPrivateKey.GetPublicKey()
		keyInfo.blsPrivateKey = hex.EncodeToString(blsPrivateKey.GetLittleEndian())
		keyInfo.blsPublicKey = hex.EncodeToString(blsPubKey.Serialize())

		// fmt.Printf("bls public key: %v\n", keyInfo.blsPublicKey)
		// fmt.Printf("bls prikey: %v\n", keyInfo.blsPrivateKey)

		// 生成零知识证明
		proof, _ := blsPrivateKey.MakeSchnorrNIZKP()
		proofByte, _ := proof.MarshalText()
		var proofHex bls.SchnorrProofHex
		proofHex.UnmarshalText(proofByte)
		keyInfo.blsProof = proofHex.String()
		// fmt.Printf("string proof: %v\n", keyInfo.blsProof)

		/*
			err = proof.VerifySchnorrNIZK(*blsPubKey)
			if err != nil {
				utils.Fatalf("verify bls private key: %v", err)
			}
		*/

		// fmt.Print("=======================\n")

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", nIndex), keyInfo.nodePublicKey)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", nIndex), keyInfo.nodePrivateKey)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", nIndex), strVersion)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", nIndex), keyInfo.versionSign)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", nIndex), keyInfo.blsPublicKey)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", nIndex), keyInfo.blsPrivateKey)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", nIndex), keyInfo.blsProof)

		nIndex++
	}
	strPath := GetCurrentPath()
	strPath = filepath.Join(strPath, "keyInfos.xlsx")
	if err := f.SaveAs(strPath); err != nil {
		fmt.Println(err)
	}
	// 保存xlsx文件
	fmt.Printf("生成虚拟节点信息成功, 生成文件路径: %v\n", strPath)

}
