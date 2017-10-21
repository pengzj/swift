# swift
Swift is an easy to use, fast, distributed, supporting multi process  multi process game server framework written by Golang.  

# warning
**This framework is under development, so please don't try it at the moment.**


# how to use
````
import (
	"github.com/pengzj/swift"
	"./app/servers/connector/handler"
)

func main()  {
	app := swift.CreateApp()
	app.SetConfigPath("/Users/francis.peng/Work/PfsGame/config")
	app.SetLogPath("/Users/francis.peng/Project/logs/pfs")

	app.HandleFunc("user.login", handler.UserLogin)
	app.RegisterHandler("onChat")

	app.Run()
}

````

# server config directory
**include files**:  
|-- master.json  
|-- servers.json  
|-- cert.pem  
|-- key.pem  

master.json
```
{
  "id": "master-server-1",
  "host": "127.0.0.1",
  "port": "3015"
}
```
servers.json  
````
[
  {"type": "connector", "id": "connector-server-1","clientHost": "127.0.0.1", "clientPort": "3301", "host": "127.0.0.1", "port": "3401", "frontend": true, "connType": "tcp"},
  {"type": "connector", "id": "connector-server-2","clientHost": "127.0.0.1", "clientPort": "3302", "host": "127.0.0.1", "port": "3402", "frontend": true, "connType": "tcp"},
  {"type": "rank", "id": "rank-server-1","host": "127.0.0.1", "port": "3403", "frontend": false}
]
````


# to do list
 
data  async  to db
optimize internal interface  
log  
manage tool  
cronjob  
service    
test case  
benchmark 

