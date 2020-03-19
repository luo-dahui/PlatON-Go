package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/PlatONnetwork/PlatON-Go/cmd/platon_tool/core"
	"github.com/PlatONnetwork/PlatON-Go/cmd/utils"
	"gopkg.in/urfave/cli.v1"
)

var (
	app = utils.NewApp("", "the wasm command line interface")
)

func init() {

	// Initialize the CLI app
	app.Commands = []cli.Command{
		core.OrdinaryCallCmd,       // 普通查询
		core.EcoModelCallCmd,       // 经济模型合约查询
		core.OrdinaryTxCmd,         // 普通转账交易（无）
		core.SendTransactionCmd,    // 代理签名交易
		core.SendRawTransactionCmd, // 私钥签名交易
		core.GetTxReceiptCmd,       // 获取交易回执
		core.TestCmd,
	}
	sort.Sort(cli.CommandsByName(app.Commands))
	app.After = func(ctx *cli.Context) error {
		return nil
	}
}

func main() {

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
