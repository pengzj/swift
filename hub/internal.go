package hub

import (
	"errors"
	"encoding/json"
	"log"
)

type interEntity struct {
	handlerMap map[string]func(*Session, []byte)[]byte
	handlerList []string
	routeMap map[string]func(*Session)string
}

var inter *interEntity

func HandleFunc(name string, handler func(*Session, []byte) []byte)  {
	if inter.handlerMap[name] != nil {
		panic("func " + name + " register twice")
	}
	inter.handlerMap[name] = handler
	inter.handlerList = append(inter.handlerList, name)
}

func RegisterHandle(name string)  {
	inter.handlerList = append(inter.handlerList, name)
}

func GetHandler(handleId int) (func(*Session, []byte) []byte, error) {
	if len(inter.handlerList) <= handleId {
		return nil, errors.New("handler exist")
	}
	list := inter.handlerList
	name := list[handleId]
	return inter.handlerMap[name], nil
}

func GetHandlerId(name string) (int, error)  {
	var handlerId = 0
	for id, val :=range inter.handlerList {
		if name == val {
			handlerId = id
			break
		}
	}
	if handlerId == 0 {
	//	return  0, errors.New("no found handle")
	}
	return handlerId, nil
}

func Route(serverType string,  handler func(session *Session) string)  {
	inter.routeMap[serverType] = handler
}

func GetRouteHandle(serverType string) func(*Session)string  {
	return inter.routeMap[serverType]
}

func GetHandlers() []byte {
	var data []byte
	type route struct{
		Id int
		Name string
	}
	var routes []route
	for id, name :=range inter.handlerList {
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

func NotFound() []byte {
	return inter.handlerMap["notFound"](nil, []byte{})
}

