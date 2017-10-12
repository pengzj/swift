package master

import (
	"../rpc"
)

type Master struct {
	Id string `json:"id"`
	Host string `json:"host"`
	Port string `json:"port"`

	rpcServer *rpc.RpcServer `json:"-"`
}

func (master *Master) Start()  {
	master.startRpcServer()
}

func (master *Master) startRpcServer()  {

}

func (master *Master) Stop()  {
	
}

