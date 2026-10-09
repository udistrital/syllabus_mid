package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/syllabus_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
)

// VinculacionController operations for Vinculacion
type VinculacionController struct {
	beego.Controller
}

// URLMapping ...
func (c *VinculacionController) URLMapping() {
	c.Mapping("PostVinculacion", c.PostVinculacion)
}

// PostVinculacion ...
// @Title PostVinculacion
// @Description consulta la vinculación del usuario y resuelve su programa académico
// @Param   body        body    {}  true        "body datos de identidad"
// @Success 200 {}
// @Failure 403 :body {}
// @router / [post]
func (c *VinculacionController) PostVinculacion() {
	defer errorhandler.HandlePanic(&c.Controller)
	bodyData := c.Ctx.Input.RequestBody
	respuesta := services.PostVinculacion(bodyData)
	c.Data["json"] = respuesta
	c.Ctx.Output.SetStatus(respuesta.Status)
	c.ServeJSON()
}
