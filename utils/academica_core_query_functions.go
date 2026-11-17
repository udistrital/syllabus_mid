package utils

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/syllabus_mid/models"
	"github.com/udistrital/utils_oas/request"
)

type docenteActivoResponse struct {
	Success bool                     `json:"Success"`
	Status  int                      `json:"Status"`
	Message string                   `json:"Message"`
	Data    models.DocenteActivoCore `json:"Data"`
}

// GetProgramasAcademicosByCoordinador consulta academica_core para obtener los
// programas académicos vinculados a la identificación del coordinador.
func GetProgramasAcademicosByCoordinador(ctx context.Context, identificacion string) ([]models.ProgramaAcademicoCore, error) {
	var programas []models.ProgramaAcademicoCore

	_, err := request.GetWithContext(ctx,
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

// GetDocenteActivoByDocumento consulta academica_core para verificar si el
// documento corresponde a un docente con carga académica activa. Retorna
// (nil, nil) cuando el servicio confirma que no es docente activo.
func GetDocenteActivoByDocumento(ctx context.Context, documento string) (*models.DocenteActivoCore, error) {
	endpoint := beego.AppConfig.String("AcademicaCore") +
		"docentes-activos/" + url.PathEscape(documento)

	var resp docenteActivoResponse
	status, err := request.GetWithContext(ctx, endpoint, &resp)
	if err != nil {
		if status == http.StatusNotFound {
			return nil, nil
		}
		return nil, &QueryError{
			Status:  http.StatusBadGateway,
			Message: fmt.Sprintf("Error al consultar el docente activo en academica_core para la identificación %s: %v", documento, err),
		}
	}
	if !resp.Success || resp.Status != http.StatusOK {
		return nil, &QueryError{
			Status:  http.StatusBadGateway,
			Message: fmt.Sprintf("academica_core respondió de forma inconsistente al consultar el docente activo para la identificación %s", documento),
		}
	}

	docente := resp.Data
	docente.Cedula = strings.TrimSpace(docente.Cedula)
	docente.Nombre = strings.TrimSpace(docente.Nombre)
	if docente.Cedula != documento || docente.Nombre == "" {
		return nil, &QueryError{
			Status:  http.StatusBadGateway,
			Message: fmt.Sprintf("academica_core devolvió datos inconsistentes para el docente con identificación %s", documento),
		}
	}
	return &docente, nil
}
