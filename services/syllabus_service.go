package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/syllabus_mid/models"
	"github.com/udistrital/syllabus_mid/utils"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
)

// esEstadoBorrador indica si el estado corresponde a una propuesta no publicada.
func esEstadoBorrador(estado string) bool {
	for _, e := range models.EstadosBorrador {
		if e == estado {
			return true
		}
	}
	return false
}

// marcaAguaTexto define la marca de agua para las propuestas no publicadas
// (todo lo que no es la versión final/publicada).
func marcaAguaTexto(estado string) string {
	if estado == models.EstadoEnRevision {
		return "EN REVISIÓN"
	}
	return "BORRADOR"
}

func PostSyllabusTemplate(ctx context.Context, data []byte) requestresponse.APIResponse {
	var syllabusRequest map[string]interface{}
	var syllabusResponse map[string]interface{}
	var syllabusTemplateResponse map[string]interface{}
	var syllabusTemplateData map[string]interface{}
	var syllabusData map[string]interface{}

	if err := json.Unmarshal(data, &syllabusRequest); err == nil {
		syllabusCode := syllabusRequest["syllabusCode"]

		estadoStr := ""
		if estado, hasEstado := syllabusRequest["estado"]; hasEstado && estado != nil {
			estadoStr = fmt.Sprintf("%v", estado)
		}

		if esEstadoBorrador(estadoStr) {
			// Borrador/propuesta: se trae la versión no publicada por su estado
			// (el estado explícito evita el filtro por defecto que oculta borradores).
			_, syllabusErr := request.GetWithContext(ctx, beego.AppConfig.String("SyllabusService")+
				fmt.Sprintf("syllabus?query=syllabus_code:%v,estado:%v&limit=1&offset=0", syllabusCode, estadoStr), &syllabusResponse)
			if syllabusErr != nil || syllabusResponse["Success"] == false {
				if syllabusErr == nil {
					syllabusErr = fmt.Errorf("SyllabusService: %v", syllabusResponse["Message"])
				}
				logs.Error(syllabusErr.Error())
				return requestresponse.APIResponseDTO(false, 404, nil, syllabusErr.Error())
			}
			syllabusList, ok := syllabusResponse["Data"].([]interface{})
			if !ok || len(syllabusList) < 1 {
				err := fmt.Errorf("SyllabusService: propuesta no encontrada para el syllabusCode %v y estado %v", syllabusCode, estadoStr)
				logs.Error(err.Error())
				return requestresponse.APIResponseDTO(false, 404, nil, err.Error())
			}
			syllabusData = syllabusList[0].(map[string]interface{})
		} else if syllabusVersion, hasVersion := syllabusRequest["version"]; hasVersion {
			_, syllabusErr := request.GetWithContext(ctx, beego.AppConfig.String("SyllabusService")+
				fmt.Sprintf("syllabus?query=syllabus_code:%v,version:%v&limit=1&offset=0", syllabusCode, syllabusVersion), &syllabusResponse)
			if syllabusErr != nil || syllabusResponse["Success"] == false {
				if syllabusErr == nil {
					syllabusErr = fmt.Errorf("SyllabusService: %v", syllabusResponse["Message"])
				}
				logs.Error(syllabusErr.Error())
				return requestresponse.APIResponseDTO(false, 404, nil, syllabusErr.Error())
			}
			syllabusList := syllabusResponse["Data"].([]interface{})
			if len(syllabusList) < 1 {
				err := fmt.Errorf("SyllabusService: syllabus not found by syllabusCode and version")
				logs.Error(err.Error())
				return requestresponse.APIResponseDTO(false, 404, nil, err.Error())
			} else {
				syllabusData = syllabusList[0].(map[string]interface{})
			}
		} else {
			_, syllabusErr := request.GetWithContext(ctx, beego.AppConfig.String("SyllabusService")+
				fmt.Sprintf("syllabus/%v", syllabusCode), &syllabusResponse)
			if syllabusErr != nil || syllabusResponse["Success"] == false {
				if syllabusErr == nil {
					syllabusErr = fmt.Errorf("SyllabusService: %v", syllabusResponse["Message"])
				}
				logs.Error(syllabusErr.Error())
				return requestresponse.APIResponseDTO(false, 404, nil, syllabusErr.Error())
			}
			syllabusData = syllabusResponse["Data"].(map[string]interface{})
		}

		// Se valida que se envíen los datos planId y proyectoId  en el payload de la request
		if syllabusRequest["planId"] == nil || syllabusRequest["proyectoId"] == nil {
			err := fmt.Errorf("SyllabusTemplateService: Incomplete data to generate the document. Plan de estudios y/o Proyecto Curricular")
			logs.Error(err.Error())
			return requestresponse.APIResponseDTO(false, 404, nil, err.Error())
		}
		// Se asigna proyecto_curricular_id y plan_estudios_id al syllabusData
		syllabusData["proyecto_curricular_id"] = syllabusRequest["proyectoId"]
		syllabusData["plan_estudios_id"] = syllabusRequest["planId"]

		spaceData, spaceErr := utils.GetAcademicSpaceData(ctx,
			int(syllabusData["plan_estudios_id"].(float64)),
			int(syllabusData["proyecto_curricular_id"].(float64)),
			int(syllabusData["espacio_academico_id"].(float64)))

		projectData, projectErr := utils.GetProyectoCurricular(ctx, int(syllabusData["proyecto_curricular_id"].(float64)))

		if spaceErr == nil && projectErr == nil {
			facultyData, facultyErr := utils.GetFacultadDelProyectoC(ctx, projectData["id_oikos"].(string))
			idiomas := ""

			if syllabusData["idioma_espacio_id"] != nil {
				idiomasStr, idiomaErr := utils.GetIdiomas(ctx, syllabusData["idioma_espacio_id"].([]interface{}))
				if idiomaErr == nil {
					idiomas = idiomasStr
				}
			}

			if facultyErr == nil {
				syllabusTemplateData = utils.GetSyllabusTemplateData(
					spaceData, syllabusData,
					facultyData, projectData, idiomas)

				// Las propuestas no publicadas llevan marca de agua en el PDF.
				if esEstadoBorrador(estadoStr) {
					syllabusTemplateData["marca_agua"] = marcaAguaTexto(estadoStr)
				}

				utils.GetSyllabusTemplate(ctx, syllabusTemplateData, &syllabusTemplateResponse)
				document := syllabusTemplateResponse["body"].(map[string]interface{})["Data"]
				return requestresponse.APIResponseDTO(true,
					201, map[string]interface{}{"document": document},
					"Generated Syllabus Template OK")
			} else {
				err := fmt.Errorf(
					"SyllabusTemplateService: Incomplete data to generate the document. Facultad y/o Idioma")
				logs.Error(err.Error())
				return requestresponse.APIResponseDTO(false, 404, nil, err.Error())
			}
		} else {
			err := fmt.Errorf(
				"SyllabusTemplateService: Incomplete data to generate the document. Espacio Académico y/o Proyecto Curricular")
			logs.Error(err.Error())
			return requestresponse.APIResponseDTO(false, 404, nil, err.Error())
		}
	}
	return requestresponse.APIResponseDTO(false, 400, nil, "Invalid request body")
}
