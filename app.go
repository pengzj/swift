package swift

import (
	"github.com/pengzj/swift/server"
	"github.com/pengzj/swift/connector/option"
	"github.com/pengzj/swift/hub"
	"github.com/pengzj/swift/db"
	"database/sql"
	"github.com/pengzj/swift/rpc"
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

func (app *Application) RegisterBeforeHandler(handler func(*hub.Session, string)error)  {
	hub.RegisterBeforeHandler(handler)
}

func (app *Application) HandleFunc(name string, handler func(*hub.Session, []byte) []byte)  {
	hub.HandleFunc(name, handler)
}

func (app *Application) RegisterHandler(name string)  {
	hub.RegisterHandle(name)
}

func (app *Application) Route(serverType string, handler func(session *hub.Session)string) {
	hub.Route(serverType, handler)
}

func (app *Application) RegisterDB(name, dbType, dsn string)  {
	db.Register(name,dbType,dsn)
}

func (app *Application) RegisterRPC(handler func())  {
	rpc.RegisterRPC(handler)
}


func (app *Application) GetDB(name string) *sql.DB {
	return db.Get(name)
}


var std *Application

func CreateApp() *Application {
	std = new(Application)
	return std
}

func GetApp() *Application {
	return std
}


