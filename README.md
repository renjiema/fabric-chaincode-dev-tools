# fabric-chaincode-dev-tools

>本项目是基于[官方文档](https://hyperledger-fabric.readthedocs.io/en/latest/peer-chaincode-devmode.html) 集成的docker版智能合约调试工具

## 使用说明
合约路径在.env文件中进行配置

- 服务启动
```shell
make start
```
默认情况下不启动chaincode合约链码，需要本地启动指定链码服务。

* 服务启动(带链码容器)

```shell
make start-chaincode
```

启动合约容器需在.env文件中配置合约路径

- 服务停用
```shell
make stop
```
- 合约重启(适用于链码容器方式启动服务)
```shell
make chaincode-reload
```
- 合约调用
```shell
# 进入cli容器
docker exec -it cli bash
# 使用命令行进行调用，例如
peer channel list
```
* http接口调用合约

提交交易接口：

```shell
curl -H "Content-Type: application/json" -X POST -d '{"funcName":"name","args":["arg1","arg2"]}' http://localhost:8080/invoke
```

查询交易接口：

```shell
curl -H "Content-Type: application/json" -X POST -d '{"funcName": "name","args": ["arg1","arg2"]}' http://localhost:8080/query
```

* IDEA/GoLand中Debug调试链码

在Debug Configuration中添加如下环境变量：

```
CORE_CHAINCODE_ID_NAME=mycc:1.0
CORE_CHAINCODE_LOGLEVEL=debug
CORE_PEER_TLS_ENABLED=false
```

并在Program arguments中添加peer地址参数：

```
-peer.address 127.0.0.1:7052
```

vscode在launch.json中配置以上参数。

### 通道和合约说明

> 默认通道名：ch1，默认合约名：mycc
### 常见错误

- chaincode definition for 'mycc' exists, but chaincode is not installed
```
解决方案：更新fabric的docker镜像，最新的2.2.x镜像可不安装链码
```
