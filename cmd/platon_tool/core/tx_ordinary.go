/*
普通rpc交易接口(普通转账)
*/

package core

import (
	"gopkg.in/urfave/cli.v1"
)

var (
	OrdinaryTxCmd = cli.Command{
		Name:   "tx_ordinary",
		Usage:  "普通转账交易接口",
		Action: tx_ordinary,
		Flags:  OrdinaryTxCmdFlags,
	}
)

func tx_ordinary(c *cli.Context) error {
	return nil
}
