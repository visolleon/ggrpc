package ggrpc

import ini "gopkg.in/ini.v1"

var (
	// RPCServerSetting RPC Server Setting
	RPCServerSetting *RPCServerInfo
)

// LoadRPCServerConfig 设置RPC服务配置信息
func LoadRPCServerConfig(cfg *ini.File) *RPCServerInfo {
	_rpcServerSetting := &RPCServerInfo{}
	_rpcServerSetting.UUID = cfg.Section("rpc").Key("UUID").MustString("")
	_rpcServerSetting.AppID = cfg.Section("rpc").Key("APPID").MustString("101010")
	_rpcServerSetting.AppKey = cfg.Section("rpc").Key("APPKEY").MustString("...")
	_rpcServerSetting.IsUseSign = cfg.Section("rpc").Key("IsUseSign").MustBool(false)
	_rpcServerSetting.Key = cfg.Section("rpc").Key("KEY").MustString("keys/server.key")
	_rpcServerSetting.Pem = cfg.Section("rpc").Key("PEM").MustString("keys/server.pem")
	_rpcServerSetting.IsOpenTLS = cfg.Section("rpc").Key("IsOpenTLS").MustBool(false)
	_rpcServerSetting.Port = cfg.Section("rpc").Key("PORT").MustInt(50052)
	return _rpcServerSetting
}
