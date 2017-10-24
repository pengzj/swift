package hub

import (
	"encoding/json"
	"log"
	"errors"
)

var (
	ERR_HANDLER_EXIST = errors.New("handler exists")
)

type interEntity struct {
	handlerMap map[string]func(*Session, []byte)[]byte
	handlerList []string
	routeMap map[string]func(*Session)string
	beforeHandler func(*Session, string)error
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
	name, err := GetHandlerName(handleId)
	if err != nil {
		return nil, err
	}
	return inter.handlerMap[name], nil
}

func GetHandlerName(handlerId int) (string, error) {
	if len(inter.handlerList) <= handlerId {
		return "", ERR_HANDLER_EXIST
	}
	list := inter.handlerList
	name := list[handlerId]
	return name, nil
}

func RegisterBeforeHandler(handler func(*Session, string)error)  {
	if inter.beforeHandler != nil {
		log.Fatal("beforeHandler register twice")
	}
	inter.beforeHandler = handler
}

func GetBeforeHandler() func(*Session, string)error {
	return inter.beforeHandler
}

func GetHandlerId(name string) (int, error)  {
	var handlerId = 0
	for id, val :=range inter.handlerList {
		if name == val {
			handlerId = id
			break
		}
	}

	return handlerId, nil
}

func Route(serverType string,  handler func(session *Session) string)  {
	inter.routeMap[serverType] = handler
}

func GetRouteHandler(serverType string) func(*Session)string  {
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

