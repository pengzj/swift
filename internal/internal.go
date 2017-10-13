package internal

import (
	"google.golang.org/grpc"
	"encoding/json"
	"log"
)

type Internal struct {
	handlerMap map[string]func([]byte)[]byte
	routeList []string
	rpcMap map[string]*grpc.ClientConn
}

var std *Internal



func HandleFunc(name string, handler func(interface{}) []byte)  {
	if std.handlerMap[name] != nil {
		panic("func " + name + " register twice")
	}
	std.handlerMap[name] = handler
	std.routeList = append(std.routeList, name)
}

func GetHandler(handleId int) (func(interface{}) []byte, error) {
	if len(std.routeList) <= handleId {
		return nil, error.Error("handler exist")
	}
	name := std.routeList[handleId]
	return std.handlerMap[name], nil
}

func GetRoutes() []byte {
	var data []byte
	type route struct{
		Id int
		Name string
	}
	var routes []route
	for id, name :=range std.routeList {
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
	std.routeList[0] = "notFound"
}