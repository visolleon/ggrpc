package ggrpc

// RPCServerInfo RPC服务配置实体类
type RPCServerInfo struct {
	UUID string `json:"uuid"` // RPC服务UUID，用于服务器判别

	IP   string `json:"ip"`   // PRC IP地址
	Port int    `json:"port"` // RPC端口号

	IsOpenTLS bool   `json:"isOpenTLS"` // 是否启用TLS认证，如果启用，则必须要指明公钥Key和私钥Pem的文件地址
	Key       string `json:"key"`       // 公钥，供服务端使用
	Pem       string `json:"pem"`       // 私钥
	TLSName   string `json:"tlsName"`   // 加密名，供客户端使用

	IsUseSign bool   `json:"isUseSign"` // 是否启用MD5签名认证，如果启用，则必须要设置AppID和AppKey
	AppID     string `json:"appId"`     // 自定义验证APPID，暂时固定值，以后扩展使用
	AppKey    string `json:"appKey"`    // 自定义验证APPKEY
}
