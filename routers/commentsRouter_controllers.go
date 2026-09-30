package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

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
