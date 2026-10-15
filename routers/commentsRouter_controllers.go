package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/udistrital/syllabus_mid/controllers:FlujoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/syllabus_mid/controllers:FlujoController"],
		beego.ControllerComments{
			Method:           "Aprobar",
			Router:           "/:code/aprobar",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/syllabus_mid/controllers:FlujoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/syllabus_mid/controllers:FlujoController"],
		beego.ControllerComments{
			Method:           "GetBorrador",
			Router:           "/:code/borrador",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/syllabus_mid/controllers:FlujoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/syllabus_mid/controllers:FlujoController"],
		beego.ControllerComments{
			Method:           "GuardarBorrador",
			Router:           "/:code/borrador",
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/syllabus_mid/controllers:FlujoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/syllabus_mid/controllers:FlujoController"],
		beego.ControllerComments{
			Method:           "Cancelar",
			Router:           "/:code/cancelar",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/syllabus_mid/controllers:FlujoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/syllabus_mid/controllers:FlujoController"],
		beego.ControllerComments{
			Method:           "GetEstado",
			Router:           "/:code/estado",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/syllabus_mid/controllers:FlujoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/syllabus_mid/controllers:FlujoController"],
		beego.ControllerComments{
			Method:           "Proponer",
			Router:           "/:code/proponer",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/syllabus_mid/controllers:FlujoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/syllabus_mid/controllers:FlujoController"],
		beego.ControllerComments{
			Method:           "Reabrir",
			Router:           "/:code/reabrir",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/syllabus_mid/controllers:FlujoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/syllabus_mid/controllers:FlujoController"],
		beego.ControllerComments{
			Method:           "Reasignar",
			Router:           "/:code/reasignar",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/syllabus_mid/controllers:FlujoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/syllabus_mid/controllers:FlujoController"],
		beego.ControllerComments{
			Method:           "Rechazar",
			Router:           "/:code/rechazar",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/syllabus_mid/controllers:FlujoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/syllabus_mid/controllers:FlujoController"],
		beego.ControllerComments{
			Method:           "Asignar",
			Router:           "/asignar",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/syllabus_mid/controllers:FlujoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/syllabus_mid/controllers:FlujoController"],
		beego.ControllerComments{
			Method:           "BuscarDocente",
			Router:           "/docente",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/syllabus_mid/controllers:FlujoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/syllabus_mid/controllers:FlujoController"],
		beego.ControllerComments{
			Method:           "MisAsignaciones",
			Router:           "/mis-asignaciones",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/syllabus_mid/controllers:FlujoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/syllabus_mid/controllers:FlujoController"],
		beego.ControllerComments{
			Method:           "Pendientes",
			Router:           "/pendientes",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/syllabus_mid/controllers:SyllabusController"] = append(beego.GlobalControllerRouter["github.com/udistrital/syllabus_mid/controllers:SyllabusController"],
		beego.ControllerComments{
			Method:           "PostSyllabusTemplate",
			Router:           "/syllabus_template",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/syllabus_mid/controllers:SyllabusLegacyController"] = append(beego.GlobalControllerRouter["github.com/udistrital/syllabus_mid/controllers:SyllabusLegacyController"],
		beego.ControllerComments{
			Method:           "GetSyllabusLegacy",
			Router:           "/:qp_syllabus",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/syllabus_mid/controllers:VinculacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/syllabus_mid/controllers:VinculacionController"],
		beego.ControllerComments{
			Method:           "PostVinculacion",
			Router:           "/",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
