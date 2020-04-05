go_platon_tool工具

是一款go语言版的platon工具，主要用于发送platon经济模型相关的交易和查询命令


## 环境准备
### windows系统

- 安装choco

```
Set-ExecutionPolicy RemoteSigned \
iwr https://chocolatey.org/install.ps1 -UseBasicParsing | iex
```

- 安装golang

```
choco install golang
```

- 安装mingw（包括gcc,g++等）

```
choco install mingw
```

- 配置环境变量（GOROOT,GOPATH以及bls库）

```powershell
D:\core\PlatonGo\src\github.com\PlatONnetwork\PlatON-Go\crypto\bls\bls_win\lib
```

### ubuntu



### 配置文件config.json

```json
{
    "chainId":101,
    "url":"http://192.168.112.33:6666",
    "tx":{
        "from":"0x914d53aad47dbe7d0186a608ef5c3538306a6f22",
        "to":"0x1708aebb40044977e7aae5472e7a35dae59510fc",
        "gas":"2100000",
        "gasPrice":"1000000000000",
        "value":"1000000000000000000",
        "wallet":"D://reward.json"
    },
    "call":{
        "txhash":"0xb119d5a0f47965d2289fbc91366739df99a45cdc6b1805f7e60bd94257c836d"
    },
    "staking":{
        "nodeid":"0xe2181d8dc731b14117ba6d982ce163fc7b9b14bbbaf9cb3c343ef72c24cf3ed568cac6ecbc30fddf9012320fab99f6be6ab37132d083cb514100bdb4b90fff5e",
        "delegateAddress":"0x914d53aad47dbe7d0186a608ef5c3538306a6f22"
    },
    "gov":{
        "proposalid":"0x44c2b07551e3195acfc6ef674d78992bfeb445c7804f198c964ae6113af5a0e0",
        "module":"staking",
        "name":"stakeThreshold"
    },
    "restricting":{
        "account":"0x914d53aad47dbe7d0186a608ef5c3538306a6f22"
    },
    "reward":{
        "account":"0x914d53aad47dbe7d0186a608ef5c3538306a6f22",
        "nodeIds":[]
    }
}
```

配置文件说明：

>- url: 连接节点的ip和rpc端口
>- gas:
>- gasPrice:
>- from:
>- staking:



----



## 交易相关命令

主要有两种签名交易的方式：

- **本地钱包文件签名方式：钱包文件放到本地进行解锁**

- **服务器代理签名方式：钱包文件放到节点上**

  

### 普通交易(Tx_Ordinary)

- 代理签名方式

```shell
./platon_tool.exe sendTransaction --from "0x914d53aad47dbe7d0186a608ef5c3538306a6f22" --to "0x1708aebb40044977e7aae5472e7a35dae59510fc" --value "1000000000000000000" --config D://config.json
```

> 此种方式要求钱包在节点上的keystore目录下。并提示输入from地址对应的钱包密码。

- 私钥签名方式

```shell
./platon_tool.exe sendRawTransaction --wallet D://reward.json --to "0x1708aebb40044977e7aae5472e7a35dae59510fc" --value "1000000000000000000" --config D://config.json
```

>此种方式需要将钱包保存在本地，不需要将钱包放到节点上。
>
>- wallet: 发交易的钱包文件路径。（不输入时，从config.json配置文件里面的tx.wallet参数中读取）
>- to: 接收转账地址。（不输入时，从config.json配置文件里面的tx.to参数中读取）
>- value: 转账金额（单位：von, 不输入时，从config.json配置文件里面的tx.value参数中读取）
>- config: 参数配置文件，保存用于签名的chainId。

- 获取交易回执

```shell
./platon_tool.exe getTxReceipt --hash "0x8403252dadc7abc1bf73b859566a304bb55bc6eac7f69d537bc64a9cf7a37b03" --config D://config.json
```

> - `--hash` 为交易hash，不输入时，从config.json配置文件里面的call.txhash参数中读取



--------------

### 经济模型交易(Tx_EcoModel)

[接口说明文档](http://192.168.9.66/Juzix-Platon-Doc/Dark/blob/develop/03-%E7%B3%BB%E7%BB%9F%E8%AE%BE%E8%AE%A1/01-%E6%9E%B6%E6%9E%84%E8%AE%BE%E8%AE%A1/PlatON%E5%BA%95%E5%B1%82/PlatON%E5%86%85%E7%BD%AE%E5%90%88%E7%BA%A6%E5%8F%8ARPC%E6%8E%A5%E5%8F%A3%E8%AF%B4%E6%98%8E.md#withdrawDelegateReward )

- **使用本地钱包文件签名方式**

#### 质押

合约地址：0x1000000000000000000000000000000000000002

- 发起质押 (funcType:1000)

```shell
./platon_tool.exe tx_ecomodel --action staking --funcName createStaking --config D://config.json 
```

- 修改质押信息 (funcType:1001)

```shell
./platon_tool.exe tx_ecomodel --action staking --funcName editCandidate --config D://config.json 
```

- 增持质押 (funcType:1002)

```shell
./platon_tool.exe tx_ecomodel --action staking --funcName increaseStaking --nodeId "" --amountType 0 --amount 100000000000000000 --config D://config.json 
```

> 参数说明：
>
> - `--nodeId`为质押节点Id，不输入时，从config.json配置文件里面的staking.nodeId参数中读取。
> - `--amount`为增持质押金额，不输入时，从config.json配置文件里面的staking.Amount参数中读取。
> - `--amountType`为增持质押金额类型，不输入时，从config.json配置文件里面的staking.AmountType参数中读取。

- 撤销质押(一次性发起全部撤销，多次到账) (funcType:1003)

```shell
./platon_tool tx_ecomodel --action staking --funcName withdrewStaking --nodeId "" --config D://config.json 
```

> 参数说明：
>
> - `--nodeId`为质押节点Id，不输入时，从config.json配置文件里面的staking.nodeId参数中读取。





#### 治理

合约地址： 0x1000000000000000000000000000000000000005

#### 举报惩罚

合约地址： 0x1000000000000000000000000000000000000004

#### 锁仓计划

合约地址： 0x1000000000000000000000000000000000000001

#### 奖励

合约地址： 0x1000000000000000000000000000000000000006

---------



### wasm合约(Wasm_Tx)





---



## 查询相关命令

### 普通查询(Call_Ordinary)

- 无参数

  传入查询接口名，如：platon_blockNumber, admin_nodeInfo

```shell
查询块高：
curl -X POST --data '{"jsonrpc":"2.0","method":"platon_blockNumber","params":[],"id":73}' 192.168.112.33:6789 -H "Content-Type: application/json"

查询节点信息：
curl -X POST --data '{"jsonrpc":"2.0","method":"net_peerCount","params":[],"id":73}' 192.168.112.33:6789 -H "Content-Type: application/json"
```

- 有参数

```SHELL
查询金额：
curl -X POST --data '{"jsonrpc":"2.0","method":"platon_getBalance","params":["0x2ad92510527a4b97ffbe9a390207d42d305bedb6", "latest"],"id":73}' 192.168.112.33:6789 -H "Content-Type: application/json"
```

```shell
test --config E://code//PlatON//src//github.com//PlatONnetwork//PlatON-Go//cmd//platon_tool//config.json --action platon_getBalance

// 获取节点信息
./platon_tool.exe call_ordinary --config E://code//PlatON//src//github.com//PlatONnetwork//PlatON-Go//cmd//platon_tool//config.json --action NodeInfo

// 获取金额
./platon_tool.exe call_ordinary --config E://code//PlatON//src//github.com//PlatONnetwork//PlatON-Go//cmd//platon_tool//config.json --action getbalance --address 0x2ad92510527a4b97ffbe9a390207d42d305bedb6

// 获取交易回执（合约交易需要解析logs[data]）
./platon_tool.exe call_ordinary --config E://code//PlatON//src//github.com//PlatONnetwork//PlatON-Go//cmd//platon_tool//config.json --action getTxReceipt --txhash 0xbe0e22ebf8d5eda9a2155e9115dfbd97f557ae622c2df5e5b42d9220647dfa98
```

通过rpc接口直接查询，工具不提供此类查询，意义不大。

- 获取交易回执

```shell
./platon_tool.exe getTxReceipt --hash "0x8403252dadc7abc1bf73b859566a304bb55bc6eac7f69d537bc64a9cf7a37b03" --config D://config.json
```

> - `--hash` 为交易hash，不输入时，从config.json配置文件里面的call.txhash参数中读取



---



### 经济模型合约(Call_EcoModel)

[接口说明文档](http://192.168.9.66/Juzix-Platon-Doc/Dark/blob/develop/03-%E7%B3%BB%E7%BB%9F%E8%AE%BE%E8%AE%A1/01-%E6%9E%B6%E6%9E%84%E8%AE%BE%E8%AE%A1/PlatON%E5%BA%95%E5%B1%82/PlatON%E5%86%85%E7%BD%AE%E5%90%88%E7%BA%A6%E5%8F%8ARPC%E6%8E%A5%E5%8F%A3%E8%AF%B4%E6%98%8E.md#withdrawDelegateReward )

#### 质押

合约地址：0x1000000000000000000000000000000000000002

- 查询当前结算周期的验证人队列 (funcType:1100)

```shell
./platon_tool.exe call_ecomodel --action staking --funcName getVerifierList --config E://code//PlatON//src//github.com//PlatONnetwork//PlatON-Go//cmd//platon_tool//config.json 
```

-  查询当前共识周期的验证人列表(funcType:1101)

```shell
./platon_tool.exe call_ecomodel --action staking --funcName getValidatorList --config D://config.json 
```

- 查询所有实时的候选人列表(funcType:1102)

```shell
./platon_tool.exe call_ecomodel --action staking --funcName getCandidateList --config D://config.json 
```

- 查询当前账户地址所委托的节点的NodeID和质押Id(funcType:1103)

```shell
./platon_tool call_ecomodel --action staking --funcName getRelatedListByDelAddr --address "0x914d53aad47dbe7d0186a608ef5c3538306a6f22" --config D://config.json
```

> `--address`为委托地址，不输入时，从config.json配置文件里面的staking.delegateAddress参数中读取。

- 查询当前单个委托信息(funcType:1104)

  **待补充**

- 查询当前节点的质押信息(funcType:1105)

```shell
./platon_tool.exe call_ecomodel --action staking --funcName getCandidateInfo --nodeId "0x003b9cebca9e0b031be9107c736e7393c217d5066b5a5473e3d03aab35bc7b3d1eadca6c69fcd94f7c266057af87e2f3dfc746d660d656bf703427302e1e8cd0" --config D://config.json
```

> `--nodeId`为节点，不输入时，从config.json配置文件里面的staking.nodeid参数中读取。

#### 治理

合约地址： 0x1000000000000000000000000000000000000005

- 查询提案(funcType:2100)

```shell
./platon_tool.exe call_ecomodel --action gov --funcName getProposal --proposalID "0x44c2b07551e3195acfc6ef674d78992bfeb445c7804f198c964ae6113af5a0e0" --config D://config.json
```

> `--proposalID`为提案，不输入时，从config.json配置文件里面的gov.proposalID参数中读取。

- 查询提案结果(funcType:2101)

```shell
./platon_tool.exe call_ecomodel --action gov --funcName getTallyResult --proposalID "0x44c2b07551e3195acfc6ef674d78992bfeb445c7804f198c964ae6113af5a0e0" --config D://config.json
```

> `--proposalID`为提案，不输入时，从config.json配置文件里面的gov.proposalID参数中读取。

- 查询提案列表(funcType:2102)

```shell
./platon_tool.exe call_ecomodel --action gov --funcName listProposal --config D://config.json
```

- 查询节点的链生效版本(funcType:2103)

```shell
./platon_tool.exe call_ecomodel --action gov --funcName getActiveVersion --config D://config.json
```

- 查询当前块高的治理参数值(funcType:2104)

```shell
./platon_tool.exe call_ecomodel --action gov --funcName getGovernParamValue --module staking --name stakeThreshold --config D://config.json
```

>`--module`为模块名，不输入时，从config.json配置文件里面的gov.module参数中读取；
>
>`--name`为参数名，不输入时，从config.json配置文件里面的gov.name参数中读取；

- 查询治理参数列表(funcType:2106)

```shell
./platon_tool.exe call_ecomodel --action gov --funcName listGovernParam --module staking --config D://config.json
```

> `--module`为模块名，不输入时，从config.json配置文件里面的gov.module参数中读取，如果gov.module为""，表示查询所有治理参数。
>
> 模块：
>
> - staking: 质押模块
> - slashing: 惩罚模块
> - block: 区块相关

#### 举报惩罚

合约地址： 0x1000000000000000000000000000000000000004

- 查询零出块的节点列表(funcType:3002)

```shell
./platon_tool.exe call_ecomodel --action slashing --funcName ZeroProduceNodeList --config D://config.json
```

#### 锁仓计划

合约地址： 0x1000000000000000000000000000000000000001

- 获取锁仓信息(funcType:4100)

```shell
./platon_tool.exe call_ecomodel --action restricting --funcName GetRestrictingInfo --address "e697317a9f490f3f4847b995d819cf79a3b47a91" --config D://config.json
```

> `--address`为锁仓释放到账账户地址，不输入时，从config.json配置文件里面的restricting.Address参数中读取.

#### 奖励

合约地址： 0x1000000000000000000000000000000000000006

- 查询账户在各节点未提取委托奖励(funcType:5100)

```shell
./platon_tool.exe call_ecomodel --action reward --funcName getDelegateReward --address "0x914d53aad47dbe7d0186a608ef5c3538306a6f22" --nodeId "e2181d8dc731b14117ba6d982ce163fc7b9b14bbbaf9cb3c343ef72c24cf3ed568cac6ecbc30fddf9012320fab99f6be6ab37132d083cb514100bdb4b90fff5e" --config D://config.json
```

> - --address`为委托账户地址，不输入时，从config.json配置文件里面的reward.Address参数中读取，
> - --nodeid为委托的节点id(单个)，不输入时，委托的节点id列表从config.json配置文件里面的reward.nodeIds参数中读取，nodeIds配置为空时，表示查询账户委托的所有节点。



-----



### wams合约查询(Wasm_Call)