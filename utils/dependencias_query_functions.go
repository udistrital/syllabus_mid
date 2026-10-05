package utils

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/syllabus_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// QueryError es un error de consulta que incluye el código HTTP asociado.
type QueryError struct {
	Status  int
	Message string
}

func (e *QueryError) Error() string {
	return e.Message
}

func GetFacultadDelProyectoC(projectIdOikos string) (map[string]any, error) {
	var facultadResponse map[string]interface{}

	facultadErr := request.GetJson(
		beego.AppConfig.String("OikosService")+
			fmt.Sprintf("dependencia/get_dependencias_padres_by_id/%v", projectIdOikos),
		&facultadResponse)
	if facultadErr == nil && facultadResponse["Type"] == "success" {
		dependencias := facultadResponse["Body"].([]interface{})
		if len(dependencias) > 0 {
			var dependenciaPadre int
			for i := len(dependencias) - 1; i >= 0; i-- {
				if fmt.Sprintf("%v", dependencias[i].(map[string]interface{})["Id"]) == projectIdOikos {
					dependenciaPadre = int(dependencias[i].(map[string]interface{})["Padre"].(float64))
				}

				if dependenciaPadre > 0 && int(dependencias[i].(map[string]interface{})["Id"].(float64)) == dependenciaPadre {
					return dependencias[i].(map[string]interface{}), nil
				}
			}
		}
		return nil, fmt.Errorf("Facultad no encontrada")
	} else {
		return nil, fmt.Errorf("Facultad no encontrada")
	}
}

// GetPadreDependenciaOikos retorna el id de la dependencia padre (campo Padre)
// del registro cuyo Id coincide con idOikos. Un Padre 0 indica la máxima jerarquía.
func GetPadreDependenciaOikos(idOikos string) (int, error) {
	var resp models.DependenciasPadresResponse

	err := request.GetJson(
		beego.AppConfig.String("OikosService")+
			fmt.Sprintf("dependencia/get_dependencias_padres_by_id/%s", idOikos), &resp)
	if err != nil {
		return 0, &QueryError{
			Status:  http.StatusBadGateway,
			Message: fmt.Sprintf("Error al consultar la dependencia padre en OIKOS para el id_oikos %s: %v", idOikos, err),
		}
	}
	if resp.Type != "success" {
		return 0, &QueryError{
			Status:  http.StatusBadGateway,
			Message: fmt.Sprintf("OIKOS no respondió correctamente al consultar la dependencia padre del id_oikos %s", idOikos),
		}
	}

	idOikosInt, err := strconv.Atoi(idOikos)
	if err != nil {
		return 0, &QueryError{
			Status:  http.StatusBadGateway,
			Message: fmt.Sprintf("El id_oikos recibido de homologación no es numérico: %s", idOikos),
		}
	}
	for _, dependencia := range resp.Body {
		if dependencia.Id == idOikosInt {
			return dependencia.Padre, nil
		}
	}
	return 0, &QueryError{
		Status:  http.StatusNotFound,
		Message: fmt.Sprintf("No se encontró la dependencia %s en OIKOS", idOikos),
	}
}

func GetProyectoCurricular(proyectoId int) (map[string]any, error) {
	var proyectoResponse map[string]interface{}

	proyectoErr := request.GetJsonWSO2(
		beego.AppConfig.String("HomologacionDependenciaService")+
			fmt.Sprintf("proyecto_curricular_cod_proyecto/%v", proyectoId),
		&proyectoResponse)
	if proyectoErr == nil && fmt.Sprintf("%v", proyectoResponse) != "map[homologacion:map[]]" && fmt.Sprintf("%v", proyectoResponse) != "map[]]" {
		homologacionData := proyectoResponse["homologacion"].(map[string]interface{})
		proyectoData := map[string]interface{}{
			"proyecto_curricular_nombre": fmt.Sprintf("%v", homologacionData["proyecto_snies"]),
			"id_oikos":                   homologacionData["id_oikos"],
			"id_snies":                   homologacionData["id_snies"],
			"id_argo":                    homologacionData["id_argo"],
		}
		return proyectoData, nil
	} else {
		return nil, fmt.Errorf("Proyecto curricular homologación no encontrado")
	}
}

// GetIdOikosBySnies consulta homologación por código SNIES y retorna el id_oikos.
func GetIdOikosBySnies(snies string) (string, error) {
	var resp models.HomologacionResponse

	err := request.GetJsonWSO2(
		beego.AppConfig.String("HomologacionDependenciaService")+
			fmt.Sprintf("proyecto_curricular_snies/%s", snies), &resp)
	if err != nil {
		return "", &QueryError{
			Status:  http.StatusBadGateway,
			Message: fmt.Sprintf("Error al consultar homologación para el SNIES %s: %v", snies, err),
		}
	}

	if resp.Homologacion.IdOikos == nil {
		return "", &QueryError{
			Status:  http.StatusNotFound,
			Message: fmt.Sprintf("No se encontró id_oikos en homologación para el SNIES %s", snies),
		}
	}

	idOikos := fmt.Sprintf("%v", resp.Homologacion.IdOikos)
	if idOikos == "" || idOikos == "<nil>" {
		return "", &QueryError{
			Status:  http.StatusNotFound,
			Message: fmt.Sprintf("No se encontró id_oikos en homologación para el SNIES %s", snies),
		}
	}

	return idOikos, nil
}
