package swift

import (
	"./master"
	"./server"
	"./connector/option"
	"./internal"
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
	serverId string
	serverType string
	configPath string
	logPath string
	connectorOptions map[string]*option.ConnectorOption
	master *master.Master
	server *server.Server
}



func (app *Application) SetConfigPath(path string)  {
	app.configPath = path
}

func (app *Application) SetLogPath(path string)  {
	app.logPath = path
}


func (app * Application) IsMaster() bool {
	return  app.serverType == SERVER_MASTER;
}

func (app *Application) getServerType() string {
	return app.serverType;
}

func (app *Application) getMaster() *master.Master {
	return app.master
}

func (app *Application) Run()  {
	app.init()
	app.startServers()
}

func (app *Application) HandleFunc(name string, handler func(interface{})(result []byte))  {
	internal.HandleFunc(name, handler())
}


var std *Application

func CreateApp() *Application {
	std = new(Application)
	return std
}

func GetApp() *Application {
	return std
}


