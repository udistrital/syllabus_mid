package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/syllabus_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
)

// SyllabusController operations for Espacios_academicos
type SyllabusController struct {
	beego.Controller
}

// URLMapping ...
func (c *SyllabusController) URLMapping() {
	c.Mapping("PostSyllabusTemplate", c.PostSyllabusTemplate)
}

// PostSyllabusTemplate ...
// @Title PostSyllabusTemplate
// @Description post Syllabus template
// @Param   body        body    {}  true        "body generar plantilla del syllabus"
// @Success 200 {}
// @Failure 403 :body {}
// @router /syllabus_template [post]
func (c *SyllabusController) PostSyllabusTemplate() {
	defer errorhandler.HandlePanic(&c.Controller)
	bodyData := c.Ctx.Input.RequestBody
	respuesta := services.PostSyllabusTemplate(c.Ctx.Request.Context(), bodyData)
	c.Data["json"] = respuesta
	c.Ctx.Output.SetStatus(respuesta.Status)
	c.ServeJSON()
}
