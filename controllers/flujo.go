package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/syllabus_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
)

// FlujoController expone el flujo de aprobación de syllabus.
type FlujoController struct {
	beego.Controller
}

// Asignar ...
// @Title Asignar
// @Description Asigna la creación/edición de un syllabus a un docente
// @Param body body {} true "datos de asignación"
// @Success 201 {}
// @Failure 400
// @router /asignar [post]
func (c *FlujoController) Asignar() {
	defer errorhandler.HandlePanic(&c.Controller)
	respuesta := services.AsignarSyllabus(c.Ctx.Request.Context(), c.Ctx.Input.RequestBody)
	c.Data["json"] = respuesta
	c.Ctx.Output.SetStatus(respuesta.Status)
	c.ServeJSON()
}

// GetBorrador ...
// @router /:code/borrador [get]
func (c *FlujoController) GetBorrador() {
	defer errorhandler.HandlePanic(&c.Controller)
	respuesta := services.GetBorrador(c.Ctx.Request.Context(),
		c.Ctx.Input.Param(":code"), c.GetString("documento"))
	c.Data["json"] = respuesta
	c.Ctx.Output.SetStatus(respuesta.Status)
	c.ServeJSON()
}

// GuardarBorrador ...
// @router /:code/borrador [put]
func (c *FlujoController) GuardarBorrador() {
	defer errorhandler.HandlePanic(&c.Controller)
	respuesta := services.GuardarBorrador(c.Ctx.Request.Context(),
		c.Ctx.Input.Param(":code"), c.Ctx.Input.RequestBody)
	c.Data["json"] = respuesta
	c.Ctx.Output.SetStatus(respuesta.Status)
	c.ServeJSON()
}

// Proponer ...
// @router /:code/proponer [post]
func (c *FlujoController) Proponer() {
	defer errorhandler.HandlePanic(&c.Controller)
	respuesta := services.ProponerSyllabus(c.Ctx.Request.Context(),
		c.Ctx.Input.Param(":code"), c.Ctx.Input.RequestBody)
	c.Data["json"] = respuesta
	c.Ctx.Output.SetStatus(respuesta.Status)
	c.ServeJSON()
}

// Aprobar ...
// @router /:code/aprobar [post]
func (c *FlujoController) Aprobar() {
	defer errorhandler.HandlePanic(&c.Controller)
	respuesta := services.AprobarSyllabus(c.Ctx.Request.Context(),
		c.Ctx.Input.Param(":code"), c.Ctx.Input.RequestBody)
	c.Data["json"] = respuesta
	c.Ctx.Output.SetStatus(respuesta.Status)
	c.ServeJSON()
}

// Rechazar ...
// @router /:code/rechazar [post]
func (c *FlujoController) Rechazar() {
	defer errorhandler.HandlePanic(&c.Controller)
	respuesta := services.RechazarSyllabus(c.Ctx.Request.Context(),
		c.Ctx.Input.Param(":code"), c.Ctx.Input.RequestBody)
	c.Data["json"] = respuesta
	c.Ctx.Output.SetStatus(respuesta.Status)
	c.ServeJSON()
}

// Reabrir ...
// @router /:code/reabrir [post]
func (c *FlujoController) Reabrir() {
	defer errorhandler.HandlePanic(&c.Controller)
	respuesta := services.ReabrirSyllabus(c.Ctx.Request.Context(),
		c.Ctx.Input.Param(":code"), c.Ctx.Input.RequestBody)
	c.Data["json"] = respuesta
	c.Ctx.Output.SetStatus(respuesta.Status)
	c.ServeJSON()
}

// Reasignar ...
// @router /:code/reasignar [post]
func (c *FlujoController) Reasignar() {
	defer errorhandler.HandlePanic(&c.Controller)
	respuesta := services.ReasignarSyllabus(c.Ctx.Request.Context(),
		c.Ctx.Input.Param(":code"), c.Ctx.Input.RequestBody)
	c.Data["json"] = respuesta
	c.Ctx.Output.SetStatus(respuesta.Status)
	c.ServeJSON()
}

// Cancelar ...
// @router /:code/cancelar [post]
func (c *FlujoController) Cancelar() {
	defer errorhandler.HandlePanic(&c.Controller)
	respuesta := services.CancelarSyllabus(c.Ctx.Request.Context(),
		c.Ctx.Input.Param(":code"), c.Ctx.Input.RequestBody)
	c.Data["json"] = respuesta
	c.Ctx.Output.SetStatus(respuesta.Status)
	c.ServeJSON()
}

// GetEstado ...
// @router /:code/estado [get]
func (c *FlujoController) GetEstado() {
	defer errorhandler.HandlePanic(&c.Controller)
	respuesta := services.GetEstado(c.Ctx.Request.Context(),
		c.Ctx.Input.Param(":code"), c.GetString("documento"))
	c.Data["json"] = respuesta
	c.Ctx.Output.SetStatus(respuesta.Status)
	c.ServeJSON()
}

// Pendientes ...
// @router /pendientes [get]
func (c *FlujoController) Pendientes() {
	defer errorhandler.HandlePanic(&c.Controller)
	respuesta := services.PendientesSyllabus(c.Ctx.Request.Context(), c.GetString("documento"))
	c.Data["json"] = respuesta
	c.Ctx.Output.SetStatus(respuesta.Status)
	c.ServeJSON()
}

// MisAsignaciones ...
// @router /mis-asignaciones [get]
func (c *FlujoController) MisAsignaciones() {
	defer errorhandler.HandlePanic(&c.Controller)
	respuesta := services.MisAsignacionesSyllabus(c.Ctx.Request.Context(), c.GetString("documento"))
	c.Data["json"] = respuesta
	c.Ctx.Output.SetStatus(respuesta.Status)
	c.ServeJSON()
}

// BuscarDocente ...
// @router /docente [get]
func (c *FlujoController) BuscarDocente() {
	defer errorhandler.HandlePanic(&c.Controller)
	respuesta := services.BuscarDocente(c.Ctx.Request.Context(), c.GetString("documento"))
	c.Data["json"] = respuesta
	c.Ctx.Output.SetStatus(respuesta.Status)
	c.ServeJSON()
}
