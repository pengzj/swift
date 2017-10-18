package swift

import (
	"flag"
	"path/filepath"
	"io/ioutil"
	"log"
	"./server"
	"encoding/json"
	"./internal"
)

func (app *Application) init()  {
	app.loadDefaultConfig()
}

func (app *Application) loadDefaultConfig()  {
	var serverId string
	flag.StringVar(&serverId, "serverId", "", "input correct serverId")
	flag.Parse()

	var serverType = "";

	if serverId  == "" {
		serverType = SERVER_MASTER
	}

	var internalServers []internal.Server



	var filePath string

	var s *server.Server

	filePath = filepath.Join(app.configPath, "master.json")
	in, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(in, s)
	if err != nil {
		log.Fatal(err)
	}
	s.Type = SERVER_MASTER
	s.IsMaster = true
	s.Frontend = false


	internalServers = append(internalServers, internal.Server{
		Type:s.Type,
		Id:s.Id,
		ClientHost:s.ClientHost,
		ClientPort:s.ClientPort,
		Host:s.Host,
		Port:s.Port,
		Frontend:s.Frontend,
	})


	if serverType == SERVER_MASTER {
		app.server = s
	}

	filePath = filepath.Join(app.configPath, "servers.json")
	var servers []*server.Server
	in, err = ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(in, servers)
	if err != nil {
		log.Fatal(err)
	}


	for _, s :=range servers {
		if s.Id == serverId {
			app.server = s
		}
		internalServers = append(internalServers, internal.Server{
			Type:s.Type,
			Id:s.Id,
			ClientHost:s.ClientHost,
			ClientPort:s.ClientPort,
			Host:s.Host,
			Port:s.Port,
			Frontend:s.Frontend,
		})
	}


	internal.PutServers(internalServers)
}