// @APIVersion 1.0.0
// @Title SGA MID - Syllabus
// @Description Microservicio del SGA MID, complementa los endpoints del syllabus, permitiendo consultar información del syllabus y consumir el endpoint generador plantilla del syllabus.
package routers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/syllabus_mid/controllers"
	"github.com/udistrital/utils_oas/errorhandler"
)

func init() {

	beego.ErrorController(&errorhandler.ErrorHandlerController{})

	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/syllabus",
			beego.NSInclude(
				&controllers.SyllabusController{},
			),
		),
		beego.NSNamespace("/syllabus",
			beego.NSInclude(
				&controllers.SyllabusLegacyController{},
			),
		),
		beego.NSNamespace("/vinculacion",
			beego.NSInclude(
				&controllers.VinculacionController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
