go_platon_tool工具

是一款go语言版的platon工具，主要用于发送platon经济模型相关的交易和查询命令





go_platon_tool命令

## 交易相关

### 普通交易(Ordinary_Tx)

### 经济模型合约(EcoModel_Tx)

- 质押
- 治理
- 举报惩罚
- 锁仓计划
- 奖励

### 

### wasm合约(Wasm_Tx)



## 查询相关

### 经济模型合约(Call_EcoModel)

[接口说明文档](http://192.168.9.66/Juzix-Platon-Doc/Dark/blob/develop/03-%E7%B3%BB%E7%BB%9F%E8%AE%BE%E8%AE%A1/01-%E6%9E%B6%E6%9E%84%E8%AE%BE%E8%AE%A1/PlatON%E5%BA%95%E5%B1%82/PlatON%E5%86%85%E7%BD%AE%E5%90%88%E7%BA%A6%E5%8F%8ARPC%E6%8E%A5%E5%8F%A3%E8%AF%B4%E6%98%8E.md#withdrawDelegateReward )

#### 质押

- 

- 治理
- 举报惩罚
- 锁仓计划
- 奖励

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

- wams合约查询(Wasm_Call)