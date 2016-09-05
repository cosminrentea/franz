package controllers

import "github.com/astaxie/beego"

type HealthController struct {
	beego.Controller
}

func (this *HealthController) Get() {
	this.Ctx.Output.SetStatus(200)
	this.Ctx.Output.Body([]byte("{}"))
}
