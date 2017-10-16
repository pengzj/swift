package swift

import (
	"./server"
	"./connector/option"
	"./internal"
	"./hub"
)

type EnumState uint8

const (
	STATE_INITED EnumState = iota  //app has inited
	STATE_START  // app start
	STATE_STARTED // app has started
	START_STOPED // app has stoped
)

type Application struct {
	state EnumState
	configPath string
	logPath string
	connectorOptions map[string]*option.ConnectorOption
	server *server.Server
}



func (app *Application) SetConfigPath(path string)  {
	app.configPath = path
}

func (app *Application) SetLogPath(path string)  {
	app.logPath = path
}


func (app * Application) IsMaster() bool {
	return  app.server.Type == SERVER_MASTER;
}

func (app *Application) getServerType() string {
	return app.server.Type;
}

func (app *Application) Run()  {
	app.init()
	app.startServers()
}

func (app *Application) HandleFunc(name string, handler func(interface{})(result []byte))  {
	internal.HandleFunc(name, handler())
}

func (app *Application) Route(serverType string, handler func(session *hub.Session)string) {
	hub.GetHub().Route(serverType, handler)
}


var std *Application

func CreateApp() *Application {
	std = new(Application)
	return std
}

func GetApp() *Application {
	return std
}


