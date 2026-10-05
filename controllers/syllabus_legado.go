package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/syllabus_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
)

type SyllabusLegacyController struct {
	beego.Controller
}

// URLMapping ...
func (c *SyllabusLegacyController) URLMapping() {
	c.Mapping("GetSyllabusLegacy", c.GetSyllabusLegacy)
}

// GetSyllabusLegacy ...
// @Title GetSyllabusLegacy
// @Description get syllabus
// @Param	qp_syllabus		path	string	true	"Parámetros de consulta cifrados (id proyecto curricular, plan de estudio, id espacio académico) syllabus"
// @Success 200 {}
// @Failure 404 not found resource
// @router /:qp_syllabus [get]
func (c *SyllabusLegacyController) GetSyllabusLegacy() {
	defer errorhandler.HandlePanic(&c.Controller)
	encodedParamsPlan := c.Ctx.Input.Param(":qp_syllabus")
	resultado := services.GetSyllabusLegacy(encodedParamsPlan)
	c.Data["json"] = resultado
	c.Ctx.Output.SetStatus(resultado.Status)
	c.ServeJSON()
}
