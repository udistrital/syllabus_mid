package utils

import (
	"fmt"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/udistrital/syllabus_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// GetProgramasAcademicosByCoordinador consulta academica_core para obtener los
// programas académicos vinculados a la identificación del coordinador.
func GetProgramasAcademicosByCoordinador(identificacion string) ([]models.ProgramaAcademicoCore, error) {
	var programas []models.ProgramaAcademicoCore

	err := request.GetJson(
		beego.AppConfig.String("AcademicaCore")+
			fmt.Sprintf("programas-academicos/coordinador/%s", identificacion), &programas)
	if err != nil {
		return nil, &QueryError{
			Status:  http.StatusBadGateway,
			Message: fmt.Sprintf("Error al consultar los programas académicos en academica_core para la identificación %s: %v", identificacion, err),
		}
	}
	return programas, nil
}
