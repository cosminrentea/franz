package main

import (
	"github.com/astaxie/beego"
	"github.com/cosminrentea/franz/controllers"
)

func main() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/task/", &controllers.TaskController{}, "get:ListTasks;post:NewTask")
	beego.Router("/task/:id:int", &controllers.TaskController{}, "get:GetTask;put:UpdateTask")
	beego.Run()
}
