package swift

import (
	"./connector/option"
	"google.golang.org/grpc"
)

func (app *Application) SetConnectionOption(serverType string, option *option.ConnectorOption)  {
	app.connectorOptions[serverType] = option
}

func (app *Application) GetConnectionOption(serverType string) *option.ConnectorOption {
	return app.connectorOptions[serverType]
}

func (app *Application) startServers()  {
	option := app.GetConnectionOption(app.server.Type)
	app.server.Start(option)
}


func (app *Application) GetRpcClient(serverType string) *grpc.ClientConn  {
	// todo route and return conn
	return nil
}

func (app *Application) Broadcast(serverType, route string, data []byte)  {

}
