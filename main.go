package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
	"github.com/cosminrentea/franz/controllers"
)

func main() {
	setupLogging()
	beego.Router("/", &controllers.MainController{})
	beego.Router("/message/", &controllers.MessageController{}, "get:ListMessages;post:NewMessage")
	beego.Router("/message/:id:int", &controllers.MessageController{}, "get:GetMessage;put:UpdateMessage")
	beego.Router("/admin/healthcheck", &controllers.HealthController{})
	log.Debug("starting")
	beego.Run()
}

func setupLogging() {
	logs.Register("logruslogstash", NewLogrusLogstash)
	nl := logs.NewLogger()
	nl.SetLogger("logruslogstash",
		fmt.Sprintf(`{"Level": "%v", "Env": "%v", "ServiceName": "%v"}`,
			config.ExpandValueEnv("${FRANZ_LOG||error}"),
			config.ExpandValueEnv("${FRANZ_ENV||dev}"),
			beego.AppConfig.String("AppName")))
}
