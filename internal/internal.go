package internal

import (
	"google.golang.org/grpc"
	"context"
	"github.com/pengzj/swift/logger"
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

	secretKey string
}


var std *Internal


type customCredential struct {

}

func (c customCredential) GetRequestMetadata(ctx context.Context, uri ...string)  (map[string]string, error) {
	return map[string]string {
		"token": std.secretKey,
	}, nil
}

func (c customCredential) RequireTransportSecurity() bool {
	return false
}

func PutServers(servers []Server)  {
	for _, server :=range servers {
		std.serverMap[server.Id] = server
	}
}


func GetServersByType(serverType string) []Server {
	var servers []Server = make([]Server,0)
	for _, server :=range std.serverMap {
		if server.Type == serverType {
			servers = append(servers, server)
		}
	}
	return servers
}

func GetServers() []Server  {
	var servers []Server = make([]Server,0)
	for _, server :=range std.serverMap {
		servers = append(servers, server)
	}
	return servers
}

func GetServerById(serverId string) Server  {
	return std.serverMap[serverId]
}

func SetClientConn(serverId string, clientConn *grpc.ClientConn)  {
	std.rpcClientMap[serverId] = clientConn
}

func SetSecretKey(key string)  {
	std.secretKey = key
}


func GetSecretKey() string {
	return std.secretKey
}

func loadClientConnByServer(server Server)  {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithPerRPCCredentials(new(customCredential)))
	conn, err := grpc.Dial(server.Host + ":" +server.Port, opts...)
	if err != nil {
		logger.Fatal("client conn error ",err)
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