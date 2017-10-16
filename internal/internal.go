package internal

import (
	"google.golang.org/grpc"
	"encoding/json"
	"log"
)

type Internal struct {
	handlerMap map[string]func([]byte)[]byte
	handlerList []string
	rpcClientMap map[string]*grpc.ClientConn

	serverMap map[string]Server
}

type Server struct {
	Type string
	Id string
	ClientHost string
	ClientPort string
	Host string
	Port string
	Frontend bool
}


var std *Internal


func HandleFunc(name string, handler func(interface{}) []byte)  {
	if std.handlerMap[name] != nil {
		panic("func " + name + " register twice")
	}
	std.handlerMap[name] = handler
	std.handlerList = append(std.handlerList, name)
}

func GetHandler(handleId int) (func(interface{}) []byte, error) {
	if len(std.handlerList) <= handleId {
		return nil, error.Error("handler exist")
	}
	name := std.handlerList[handleId]
	return std.handlerMap[name], nil
}

func GetRoutes() []byte {
	var data []byte
	type route struct{
		Id int
		Name string
	}
	var routes []route
	for id, name :=range std.handlerList {
		routes = append(routes, route{
			Id:id,
			Name:name,
		})
	}
	data, err := json.Marshal(routes)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func NotFound()  []byte {
	return std.handlerMap["notFound"]()
}

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

func GetClientConn(serverId string) *grpc.ClientConn  {
	return std.rpcClientMap[serverId]
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

	std.handlerMap["notFound"] = func() []byte {
		var body = struct {
			Code int
			Message string
		}{
			Code: 404,
			Message:"method not exists",
		}
		data, _ := json.Marshal(body)
		return data
	}
	std.handlerList[0] = "notFound"
}