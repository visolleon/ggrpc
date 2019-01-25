package ggrpc

import (
	"errors"
	"log"
	"math/rand"

	"github.com/Visolleon/utils"

	"google.golang.org/grpc"
)

// ServerGroup RPC服务器群管理信息
type ServerGroup struct {
	InfoList   []*RPCServerInfo // RPC信息库列表
	r          *rand.Rand       // 随机种子
	ServerList []*Client        // RPC Server列表
	configFile string           // 配置文件
}

// Init 初始化配置
func (group *ServerGroup) Init(configFile string) error {
	var err error
	group.r = utils.NewRand()
	group.InfoList = make([]*RPCServerInfo, 0)
	err = utils.FromJSONFile(configFile, &group.InfoList)
	if err != nil || len(group.InfoList) == 0 {
		log.Panicf("Load RPC server list failed, %v", err)
	} else {
		group.configFile = configFile
		group.ServerList = make([]*Client, 0)
		// 初始化Client Server对象
		for _, info := range group.InfoList {
			if info != nil {
				c := &Client{}
				c.Init(info)
				group.ServerList = append(group.ServerList, c)
			}
		}
	}
	log.Printf("group.ServerList[0]: %+v\n", group.ServerList[0])
	return err
}

// GoWatchFile 监控配置文件
func (group *ServerGroup) GoWatchFile() {
	log.Printf("### watch Config file: %s \n", group.configFile)
	go utils.Watcher(group.configFile, func(fileName string) {
		group.Init(group.configFile)
	})
}

// RandOne 随机选取一个服务
func (group *ServerGroup) RandOne() (*Client, error) {
	var (
		err     error
		rpcInfo *Client
	)
	if len(group.ServerList) > 0 {
		index := group.r.Intn(len(group.ServerList))
		rpcInfo = group.ServerList[index]
	} else {
		err = errors.New("Not rpc server")
	}
	if rpcInfo == nil {
		err = errors.New("Get rpc server error")
	}
	return rpcInfo, err
}

// GetByUUID 根据UIID来获取RPC服务器信息
func (group *ServerGroup) GetByUUID(uuid string) (*Client, error) {
	var (
		err     error
		rpcInfo *Client
	)
	for i := 0; i < len(group.ServerList); i++ {
		if group.ServerList[i].RPCInfo.UUID == uuid {
			rpcInfo = group.ServerList[i]
			break
		}
	}
	if rpcInfo == nil {
		err = errors.New("Not found rpc server")
	}
	return rpcInfo, err
}

// RegInterceptor 全体注册拦截器
func (group *ServerGroup) RegInterceptor(hander grpc.UnaryClientInterceptor) {
	for _, c := range group.ServerList {
		c.RegInterceptor(hander)
	}
}

// InvokeRand 随机调用
func (group *ServerGroup) InvokeRand(hander InvokeHandler) error {
	c, err := group.RandOne()
	if err == nil {
		err = c.Invoke(hander)
	}
	return err
}

// InvokeByUUID 调用
func (group *ServerGroup) InvokeByUUID(uuid string, hander InvokeHandler) error {
	c, err := group.GetByUUID(uuid)
	if err == nil {
		err = c.Invoke(hander)
	}
	return err
}

// InvokeAll RPC集群群体调用
func (group *ServerGroup) InvokeAll(hander InvokeHandler) error {
	var err error
	for _, c := range group.ServerList {
		// TODO: 考虑协程处理...
		err = c.Invoke(hander)
	}
	return err
}
