package main

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/cmd/platon_tool/core"
	"github.com/PlatONnetwork/PlatON-Go/cmd/utils"
	"gopkg.in/urfave/cli.v1"
	"os"
	"sort"
)

var (
	app = utils.NewApp("", "the wasm command line interface")
)

func init() {

	// Initialize the CLI app
	app.Commands = []cli.Command{
		core.OrdinaryCallCmd,
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
	/*
		client, err := rpc.Dial("http://192.168.112.33:6789")
		if err != nil {
			fmt.Println("rpc.Dial err", err)
			return
		}

		var account []string
		err = client.Call(&account, "platon_accounts")
		var result string
		//var result hexutil.Big
		err = client.Call(&result, "platon_getBalance", account[0], "latest")
		//err = ec.c.CallContext(ctx, &result, "eth_getBalance", account, "latest")

		if err != nil {
			fmt.Println("client.Call err", err)
			return
		}

		var result2 NodeInfo
		err = client.Call(&result2, "admin_nodeInfo")

		if err != nil {
			fmt.Println("client.Call err", err)
			return
		}

		fmt.Printf("node info: %s\n", result2)*/
}
