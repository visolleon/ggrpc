package ggrpc

import (
	"context"
	"fmt"
	"net"

	log "github.com/Visolleon/logger"

	"github.com/Visolleon/utils"

	ini "gopkg.in/ini.v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

// Server RPC服务对象
type Server struct {
	RPCInfo    *RPCServerInfo // RPC信息
	GRPCServer *grpc.Server
	Listener   net.Listener
}

// Init 初始化RPC服务
func (s *Server) Init(rpcInfo *RPCServerInfo) {
	s.RPCInfo = rpcInfo
	var opts []grpc.ServerOption

	// TLS认证
	if rpcInfo.IsOpenTLS {
		creds, err := credentials.NewServerTLSFromFile(rpcInfo.Pem, rpcInfo.Key)
		if err != nil {
			log.Errorf("Failed to generate credentials %v", err)
			panic(err)
		}

		opts = append(opts, grpc.Creds(creds))
	}

	// Token认证
	if rpcInfo.IsUseSign {
		// 注册interceptor
		opts = append(opts, grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			var err error
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				err = log.Errorf("can not found token info")
			}
			var (
				appid string // AppID
				guid  string // 请求流水号，每次都生成
				time  string // 时间戳
				sign  string // MD5加密后字符串
			)

			if val, ok := md["appid"]; ok {
				appid = val[0]
			}

			if val, ok := md["guid"]; ok {
				guid = val[0]
			}

			if val, ok := md["time"]; ok {
				time = val[0]
			}

			if val, ok := md["sign"]; ok {
				sign = val[0]
			}

			signStr := fmt.Sprintf("%s|%s|%s", appid, guid, time)
			checkSign := utils.MD5EncodeStr(signStr)

			// log.Printf("Sign: signStr=%s, sign=%s, checkSign=%s \n", signStr, sign, checkSign)
			if appid != rpcInfo.AppID || sign != checkSign {
				err = log.Errorf("Sign invalid: signStr=%s, sign=%s", signStr, sign)
			}
			if err != nil {
				return nil, err
			}
			// 继续处理请求
			return handler(ctx, req)
		}))
	}

	// 实例化grpc Server
	s.GRPCServer = grpc.NewServer(opts...)

	// 注册HelloService
	// pb.RegisterHelloServer(s.GRPCServer, &service.HelloService{})
}

// GetUUID 获取Server的UUID
func (s *Server) GetUUID() string {
	return s.RPCInfo.UUID
}

// RegHandler RPC注册方法
type RegHandler func(*grpc.Server)

// Register 注册函数对象
// @Example:
// 	server.Register(func(grpcServer *grpc.Server) {
// 		pb.RegisterHelloServer(grpcServer, &service.HelloService{})
// 	})
func (s *Server) Register(regFunc RegHandler) {
	regFunc(s.GRPCServer)
}

// Start 启动RPC服务
func (s *Server) Start() {
	if s.GetUUID() == "" {
		log.Printf("Please set server uuid first.")
		panic(0)
	}
	var address = fmt.Sprintf("%s:%d", s.RPCInfo.IP, s.RPCInfo.Port)
	var err error
	s.Listener, err = net.Listen("tcp", address)
	defer func() {
		s.GRPCServer.Stop()
		s.Listener.Close()
	}()

	if err != nil {
		log.Printf("Failed to listen: %v", err)
		panic(err)
	}

	log.Println("RPC service listen on:" + address)
	s.GRPCServer.Serve(s.Listener)
}

// NewGPRCServer 新建RPC服务对象
func NewGPRCServer(rpcInfo *RPCServerInfo) *Server {
	var server = &Server{}
	server.Init(rpcInfo)
	return server
}

// InitServerFromIniConfig 初始化RPC服务器
func InitServerFromIniConfig(cfg *ini.File) *Server {
	// 配置当前RPC服务信息
	RPCServerSetting = LoadRPCServerConfig(cfg)
	return NewGPRCServer(RPCServerSetting)
}
