package main

import (
	"github.com/astaxie/beego"
	"github.com/cosminrentea/franz/controllers"
)

func main() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/message/", &controllers.MessageController{}, "get:ListMessages;post:NewMessage")
	beego.Router("/message/:id:int", &controllers.MessageController{}, "get:GetMessage;put:UpdateMessage")
	beego.Run()
}
