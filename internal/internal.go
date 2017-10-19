package internal

import (
	"google.golang.org/grpc"
	"log"
)

type Server struct {
	Type string
	Id string
	Host string
	Port string
}

type Internal struct {
	rpcClientMap map[string]*grpc.ClientConn
	serverMap map[string]Server
}


var std *Internal



func PutServers(servers []Server)  {
	for _, server :=range servers {
		std.serverMap[server.Id] = server
	}
}


func GetServersByType(serverType string) []Server {
	var servers []Server
	for _, server :=range std.serverMap {
		if server.Type == serverType {
			servers = append(servers, server)
		}
	}
	return servers
}

func GetServerById(serverId string) Server  {
	return std.serverMap[serverId]
}

func SetClientConn(serverId string, clientConn *grpc.ClientConn)  {
	std.rpcClientMap[serverId] = clientConn
}

func loadClientConnByServer(server Server)  {
	conn, err := grpc.Dial(server.Host + ":" +server.Port, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	std.rpcClientMap[server.Id] = conn
}

func GetClientConnByServerId(serverId string) *grpc.ClientConn  {
	conn := std.rpcClientMap[serverId]
	if conn != nil {
		return conn
	}

	server := std.serverMap[serverId]
	loadClientConnByServer(server)

	return std.rpcClientMap[serverId]
}

func init() {
	std = new(Internal)
	std.serverMap = map[string]Server{}
	std.rpcClientMap = map[string]*grpc.ClientConn{}
}