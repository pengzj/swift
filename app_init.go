package swift

import (
	"flag"
	"path/filepath"
	"io/ioutil"
	"log"
	"./master"
	"./server"
	"encoding/json"
)

func (app *Application) init()  {
	app.loadDefaultConfig()
}

func (app *Application) loadDefaultConfig()  {
	app.loadMaster()

	serverId := flag.String("serverId", "", "input correct serverId")
	flag.Parse()


	app.serverId = *serverId

	if app.serverId == "" {
		app.serverType = SERVER_MASTER
	}

	if app.IsMaster() == false {
		app.loadServer()
	}
}

func (app *Application) loadServer()  {
	filePath := filepath.Join(app.configPath, "servers.json")
	in, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	var servers  []server.Server
	err = json.Unmarshal(in, &servers)
	if err != nil {
		log.Fatal(err)
	}

	for _, s := range servers {
		if s.Id == app.serverId {
			app.server = &s
			break
		}
	}
}

func (app *Application) loadMaster()  {
	filePath := filepath.Join(app.configPath, "master.json")
	in, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	m := new(master.Master)
	err = json.Unmarshal(in, m)
	if err != nil {
		log.Fatal(err)
	}
	app.master = m
}