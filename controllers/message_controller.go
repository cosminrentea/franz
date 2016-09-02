package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/cosminrentea/franz/models"
)

type MessageController struct {
	beego.Controller
}

// Example:
//
//   req: GET /message/
//   res: 200 {"Messages": [
//          {"ID": 1, "Title": "Learn Go", "Done": false},
//          {"ID": 2, "Title": "Buy bread", "Done": true}
//        ]}
func (this *MessageController) ListMessages() {
	res := struct{ Messages []*models.Message }{models.DefaultMessageList.All()}
	this.Data["json"] = res
	this.ServeJSON()
}

// Examples:
//
//   req: POST /message/ {"Title": ""}
//   res: 400 empty title
//
//   req: POST /message/ {"Title": "Buy bread"}
//   res: 200
func (this *MessageController) NewMessage() {
	req := struct{ Title string }{}
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &req); err != nil {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte("empty title"))
		return
	}
	t, err := models.NewMessage(req.Title)
	if err != nil {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte(err.Error()))
		return
	}
	err = models.DefaultMessageList.Send(t)
	if err != nil {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte(err.Error()))
		return
	}
	models.DefaultMessageList.Save(t)
}

// Examples:
//
//   req: GET /message/1
//   res: 200 {"ID": 1, "Title": "Buy bread", "Done": true}
//
//   req: GET /message/42
//   res: 404 message not found
func (this *MessageController) GetMessage() {
	id := this.Ctx.Input.Param(":id")
	beego.Info("Message is ", id)
	intid, _ := strconv.ParseInt(id, 10, 64)
	t, ok := models.DefaultMessageList.Find(intid)
	beego.Info("Found", ok)
	if !ok {
		this.Ctx.Output.SetStatus(404)
		this.Ctx.Output.Body([]byte("message not found"))
		return
	}
	this.Data["json"] = t
	this.ServeJSON()
}

// Example:
//
//   req: PUT /message/1 {"ID": 1, "Title": "Learn Go", "Done": true}
//   res: 200
//
//   req: PUT /message/2 {"ID": 2, "Title": "Learn Go", "Done": true}
//   res: 400 inconsistent message IDs
func (this *MessageController) UpdateMessage() {
	id := this.Ctx.Input.Param(":id")
	beego.Info("Message is ", id)
	intid, _ := strconv.ParseInt(id, 10, 64)
	var t models.Message
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &t); err != nil {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte(err.Error()))
		return
	}
	if t.ID != intid {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte("inconsistent message IDs"))
		return
	}
	if _, ok := models.DefaultMessageList.Find(intid); !ok {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte("message not found"))
		return
	}
	models.DefaultMessageList.Save(&t)
}
