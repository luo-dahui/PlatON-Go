package core

import "gopkg.in/urfave/cli.v1"

var (
	ActionFlag = cli.StringFlag{
		Name:  "action",
		Usage: "rpc api name",
	}

	FuncNameFlag = cli.StringFlag{
		Name:  "funcName",
		Usage: "economic model function name",
	}

	AmountTypeFlag = cli.StringFlag{
		Name:  "amountType",
		Usage: "Amount Type(0:Amount of free,1:Restricting Amount)",
	}

	AmountFlag = cli.StringFlag{
		Name:  "amount",
		Usage: "Amount(unit:Von)",
	}

	ConfigPathFlag = cli.StringFlag{
		Name:  "config",
		Usage: "config path",
	}
	WalletFilePathFlag = cli.StringFlag{
		Name:  "wallet",
		Value: "",
		Usage: "wallet file path",
	}
	PKFilePathFlag = cli.StringFlag{
		Name:  "pkfile",
		Value: "",
		Usage: "private key file path",
	}
	StabExecTimesFlag = cli.IntFlag{
		Name:  "times",
		Value: 1000,
		Usage: "execute times",
	}
	SendTxIntervalFlag = cli.IntFlag{
		Name:  "interval",
		Value: 10,
		Usage: "Time interval for sending transactions",
	}
	AccountSizeFlag = cli.IntFlag{
		Name:  "size",
		Value: 10,
		Usage: "account size",
	}
	TxJsonDataFlag = cli.StringFlag{
		Name:  "data",
		Usage: "transaction data",
	}
	ContractWasmFilePathFlag = cli.StringFlag{
		Name:  "code",
		Usage: "wasm file path",
	}
	ContractAddrFlag = cli.StringFlag{
		Name: "addr",

		Usage: "the contract address",
	}
	ContractFuncNameFlag = cli.StringFlag{
		Name:  "func",
		Usage: "function and param ,eg :set(1,\"a\")",
	}
	TransactionTypeFlag = cli.IntFlag{
		Name:  "type",
		Value: 2,
		Usage: "tx type ,default 2",
	}
	ContractAbiFilePathFlag = cli.StringFlag{
		Name:  "abi",
		Usage: "abi file path",
	}
	TransactionHashFlag = cli.StringFlag{
		Name:  "hash",
		Usage: "tx hash",
	}
	TxFromFlag = cli.StringFlag{
		Name:  "from",
		Usage: "transaction sender addr",
	}
	TxToFlag = cli.StringFlag{
		Name:  "to",
		Usage: "transaction acceptor addr",
	}
	TransferValueFlag = cli.StringFlag{
		Name:  "value",
		Value: "0xDE0B6B3A7640000", //1LAT
		Usage: "transfer value",
	}

	testCmdFlags = []cli.Flag{
		ConfigPathFlag,
		ActionFlag,
	}
	// 普通查詢相關
	AddressFlag = cli.StringFlag{
		Name:  "address",
		Usage: "the Account address",
	}
	TxHashFlag = cli.StringFlag{
		Name:  "txhash",
		Usage: "the transaction hash",
	}

	OrdinaryCallCmdFlags = []cli.Flag{
		ConfigPathFlag,
		ActionFlag,
		AddressFlag, // 查询余额等
		TxHashFlag,  // 交易hash
	}

	// 普通交易
	OrdinaryTxCmdFlags = []cli.Flag{
		ConfigPathFlag,
		ActionFlag,
		TxFromFlag,
		TxToFlag,
		TransferValueFlag,
	}

	// 经济模型查询相关
	NodeIdFlag = cli.StringFlag{
		Name:  "nodeId",
		Usage: "the Node ID",
	}
	ProposalIDFlag = cli.StringFlag{
		Name:  "proposalID",
		Usage: "the Proposal ID",
	}

	ModuleFlag = cli.StringFlag{
		Name:  "module",
		Usage: "the module",
	}

	NameFlag = cli.StringFlag{
		Name:  "name",
		Usage: "the param's name of module",
	}

	EcoModelCallCmdFlags = []cli.Flag{
		ActionFlag,
		FuncNameFlag,
		AddressFlag,
		NodeIdFlag,
		ProposalIDFlag,
		ModuleFlag,
		NameFlag,
		ConfigPathFlag,
	}

	EcoModelTxCmdFlags = []cli.Flag{
		ActionFlag,
		FuncNameFlag,
		AmountTypeFlag,
		AmountFlag,
		WalletFilePathFlag,
		AddressFlag,
		NodeIdFlag,
		ProposalIDFlag,
		ModuleFlag,
		NameFlag,
		ConfigPathFlag,
	}

	// 合約相關
	deployCmdFlags = []cli.Flag{
		ContractWasmFilePathFlag,
		ContractAbiFilePathFlag,
		ConfigPathFlag,
	}
	invokeCmdFlags = []cli.Flag{
		ContractFuncNameFlag,
		ContractAbiFilePathFlag,
		ContractAddrFlag,
		ConfigPathFlag,
		TransactionTypeFlag,
	}

	sendTransactionCmdFlags = []cli.Flag{
		TxFromFlag,
		TxToFlag,
		TransferValueFlag,
		ConfigPathFlag,
	}
	sendRawTransactionCmdFlags = []cli.Flag{
		PKFilePathFlag,
		WalletFilePathFlag,
		TxFromFlag,
		TxToFlag,
		TransferValueFlag,
		ConfigPathFlag,
	}
	getTxReceiptCmdFlags = []cli.Flag{
		TransactionHashFlag,
		ConfigPathFlag,
	}

	stabilityCmdFlags = []cli.Flag{
		PKFilePathFlag,
		StabExecTimesFlag,
		SendTxIntervalFlag,
		ConfigPathFlag,
	}
	stabPrepareCmdFlags = []cli.Flag{
		PKFilePathFlag,
		AccountSizeFlag,
		TransferValueFlag,
		ConfigPathFlag,
	}
)
