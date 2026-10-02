package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	// "github.com/udistrital/syllabus_mid/mocks"

	"github.com/udistrital/syllabus_mid/models"
	"github.com/udistrital/syllabus_mid/utils"
	"github.com/udistrital/utils_oas/requestresponse"
)

func statusFromError(err error) int {
	var queryError *utils.QueryError
	if errors.As(err, &queryError) {
		return queryError.Status
	}
	return http.StatusInternalServerError
}

func PostVinculacion(data []byte) requestresponse.APIResponse {
	var vinculacionRequest map[string]interface{}

	if err := json.Unmarshal(data, &vinculacionRequest); err != nil {
		return requestresponse.APIResponseDTO(false, http.StatusBadRequest, nil, "La petición no tiene un cuerpo JSON válido")
	}

	identificacion, hasIdentificacion := vinculacionRequest["Identificacion"].(string)
	if !hasIdentificacion || identificacion == "" {
		return requestresponse.APIResponseDTO(false, http.StatusBadRequest, nil, "La petición no contiene el campo Identificacion")
	}

	// El nombre del coordinador es opcional en la petición.
	nombreCompleto, _ := vinculacionRequest["NombreCompleto"].(string)

	// programas := mocks.GetProgramasAcademicosByIdentificacion(identificacion)
	programas, err := utils.GetProgramasAcademicosByCoordinador(identificacion)
	if err != nil {
		return requestresponse.APIResponseDTO(false, statusFromError(err), nil, err.Error())
	}
	if len(programas) == 0 {
		return requestresponse.APIResponseDTO(false, http.StatusNotFound, nil,
			fmt.Sprintf("No se encontraron programas vinculados para la identificación %s", identificacion))
	}

	resultado := make([]models.VinculacionPrograma, 0, len(programas))
	for _, programa := range programas {
		idOikos, err := utils.GetIdOikosBySnies(programa.CodigoSnies)
		if err != nil {
			return requestresponse.APIResponseDTO(false, statusFromError(err), nil,
				fmt.Sprintf("%v (programa: %s)", err, programa.Nombre))
		}

		idOikosInt, err := strconv.Atoi(idOikos)
		if err != nil {
			return requestresponse.APIResponseDTO(false, http.StatusBadGateway, nil,
				fmt.Sprintf("El id_oikos recibido de homologación no es numérico: %s (programa: %s)", idOikos, programa.Nombre))
		}

		padreOikos, err := utils.GetPadreDependenciaOikos(idOikos)
		if err != nil {
			return requestresponse.APIResponseDTO(false, statusFromError(err), nil,
				fmt.Sprintf("%v (programa: %s)", err, programa.Nombre))
		}

		resultado = append(resultado, models.VinculacionPrograma{
			IdOikos:           idOikosInt,
			Codigo:            programa.Codigo,
			Nombre:            programa.Nombre,
			CorreoElectronico: programa.Email,
			Coordinador:       nombreCompleto,
			Identificacion:    programa.IdentificacionCoord,
			PadreOikos:        padreOikos,
		})
	}

	return requestresponse.APIResponseDTO(true, http.StatusOK, resultado, "Vinculación consultada")
}
