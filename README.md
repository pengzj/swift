# swift
Swift is an easy to use, fast, distributed, supporting multi process  multi process game server framework written by Golang.  

# warning
**This framework is in development, so please don't try it at the moment.**


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

