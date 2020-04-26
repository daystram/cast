package v1

import "github.com/astaxie/beego"

// Simple Availability Test
type PingController struct {
	beego.Controller
}

// @Title Ping
// @Success 200 {object} models.Object
// @router / [get]
func (c *PingController) GetAll() {
	c.Data["json"] = "pong"
	c.ServeJSON()
}
