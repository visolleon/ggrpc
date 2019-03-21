grpc封装库
=================
一个简单对GPRC做封装的库

安装

```sh
go get github.com/Visolleon/ggrpc
```

### 生成proto

```pb
syntax = "proto3"; // 指定proto版本
package proto;     // 指定包名

option go_package = "proto";

// 玩家信息结构
message PlayerStruct {
    // 用户ID
    int64 Uid = 1;
    // 用户Token
    string Token = 2;
    // 用户昵称
    string Name = 3;
}
```

```sh
protoc --go_out=plugins=grpc:. *.proto
```

### 使用`protoc`命令编译`.proto`文件:

* -I 参数：指定import路径，可以指定多个 -I参数，按顺序查找，默认只查找当前目录
* --go_out ：golang编译支持，支持以下参数:
    * plugins=plugin1+plugin2 - 指定插件，目前只有grpc，即：plugins=grpc
    * M 参数 - 指定导入的.proto文件路径编译后对应的golang包名(不指定本参数默认就是.proto文件中import语句的路径)
    * import_prefix=xxx - 为所有import路径添加前缀，主要用于编译子目录内的多个proto文件，这个参数按理说很有用，尤其适用替代hello_http.proto编译时的M参数，但是实际使用时有个蛋疼的问题，自己尝试看看吧
    * import_path=foo/bar - 用于指定未声明package或go_package的文件的包名，最右面的斜线前的字符会被忽略
    * :编译文件路径  .proto文件路径(支持通配符)
    * 同一个包内包含多个.proto文件时使用通配符同时编译所有文件，单独编译每个文件会存在变量命名冲突，例如编译hello_http.proto那里所示

> 完整示例：

```protobuf
protoc --go_out=plugins=grpc,Mfoo/bar.proto=bar,import_prefix=,import_path=foo/bar:. ./*.proto
```

### 新建RPC服务器

> ini配置文件设置
```ini
[rpc]
; rpc服务的唯一标识
UUID=Service_001
; rpc服务端口号
PORT = 51001
; 启用TLS认证，详细可查询grpc的OpenTLS
IsOpenTLS = true
KEY = keys/server.key
PEM = keys/server.pem
; 启用SIGN认证
IsUseSign = true
APPID = b659e4f4cea1ad1f011e48355c14667d
APPKEY = 37e267fe-f348-3cc1-85ff-ee111cf6e2b7
```
```go
// 新建rpc Server
var server = ggrpc.InitServerFromIniConfig(config.Cfg)
// 注册RPC服务提供调用的类和方法
proto.RegisterBattleServer(server.GRPCServer, &service.UserService{})
// 启动RPC服务
server.Start()
```

### 调用RPC服务
RPC调用配置文件，支持集群，或可以根据uuid有目标性的调用
```json
[
    {
        "uuid": "PushServer_001",
        "ip": "127.0.0.1",
        "port": 40041,
        "isOpenTLS": true,
        "key": "keys/server.key",
        "pem": "keys/server.pem",
        "tlsName": "MLDD.RPC",
        "isUseSign": true,
        "appId": "b659e4f4cea1ad1f011e48355c14667d",
        "appKey": "37e267fe-f348-3cc1-85ff-ee111cf6e2b7",
        "websocketURL": "127.0.0.1:8002"
    },
    ...
]
```
```go
// 新建RPC集群调用（集群用于负载或者有目的性的调用）
var PushServerGroup = &ggrpc.ServerGroup{}
// 读取配置文件
err := PushServerGroup.Init(config.PushServerConfigFile)
if err != nil {
    log.Panicf("Load push server error, %v", err)
}
// 监控配置文件，当配置文件发生变化后，程序自动更新
PushServerGroup.GoWatchFile()
```

## LICENSE
[MIT](LICENSE)