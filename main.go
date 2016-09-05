package main

import (
	"github.com/astaxie/beego"
	"github.com/cosminrentea/franz/controllers"
	"github.com/astaxie/beego/logs"
	"fmt"
	"os"
)

func main() {
	setupLogging()
	beego.Router("/", &controllers.MainController{})
	beego.Router("/message/", &controllers.MessageController{}, "get:ListMessages;post:NewMessage")
	beego.Router("/message/:id:int", &controllers.MessageController{}, "get:GetMessage;put:UpdateMessage")
	beego.Run()
}

func setupLogging() {
	logs.Register("logruslogstash", NewLogrusLogstash)
	nl := logs.NewLogger()
	env := os.Getenv("FRANZ_ENV")
	if env == "" {
		env = "dev"
	}
	nl.SetLogger("logruslogstash",
		fmt.Sprintf(`{"Level":"%v", "Env": "%v", "ServiceName": "%v"}`,
			os.Getenv("FRANZ_LOG"), env, beego.AppConfig.String("AppName")))
}
