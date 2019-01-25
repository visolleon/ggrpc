package ggrpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Visolleon/utils"
	"github.com/pborman/uuid"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// customCredential 自定义认证
type customCredential struct {
	appID     string
	guid      string
	time      string
	sign      string
	IsUseSign bool
}

// GetRequestMetadata 实现自定义认证接口
func (c customCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	c.guid = uuid.New()
	c.time = fmt.Sprintf("%d", time.Now().UnixNano())
	signStr := fmt.Sprintf("%s|%s|%s", c.appID, c.guid, c.time)
	c.sign = utils.MD5EncodeStr(signStr)

	return map[string]string{
		"appid": c.appID,
		"guid":  c.guid,
		"time":  c.time,
		"sign":  c.sign,
	}, nil
}

// RequireTransportSecurity 自定义认证是否开启TLS
func (c customCredential) RequireTransportSecurity() bool {
	return c.IsUseSign
}

// Client RPC客户端
type Client struct {
	RPCInfo *RPCServerInfo // RPC信息
	address string
	opts    []grpc.DialOption
}

// Init 初始化
func (c *Client) Init(rpcInfo *RPCServerInfo) {
	c.RPCInfo = rpcInfo
	c.address = fmt.Sprintf("%s:%d", rpcInfo.IP, rpcInfo.Port)
	var opts []grpc.DialOption
	if rpcInfo.IsOpenTLS {
		// TLS连接
		creds, err := credentials.NewClientTLSFromFile(rpcInfo.Pem, rpcInfo.TLSName)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	// 指定自定义认证
	cCred := &customCredential{
		appID:     rpcInfo.AppID,
		IsUseSign: rpcInfo.IsUseSign,
	}

	opts = append(opts, grpc.WithPerRPCCredentials(cCred))
	c.opts = opts
}

// RegInterceptor 注册拦截器
func (c *Client) RegInterceptor(hander grpc.UnaryClientInterceptor) {
	c.opts = append(c.opts, grpc.WithUnaryInterceptor(hander))
}

// InvokeHandler 调用方法
type InvokeHandler func(*grpc.ClientConn) error

// Invoke 调用
func (c *Client) Invoke(hander InvokeHandler) error {
	conn, err := grpc.Dial(c.address, c.opts...)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	return hander(conn)
}

// NewClient 新建Client对象
func NewClient(rpcInfo *RPCServerInfo) *Client {
	var c = &Client{}
	c.Init(rpcInfo)
	return c
}
