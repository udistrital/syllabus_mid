package utils

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/syllabus_mid/helpers"
	"github.com/udistrital/utils_oas/request"
)

func GetSyllabusTemplateData(spaceData, syllabusData, facultyData, projectData map[string]interface{}, languages string) map[string]interface{} {
	var propositos []interface{}
	var contenidoTematicoDescripcion string
	var contenidoTematicoDetalle []interface{}
	var evaluacionDescripcion string
	var evaluacionDetalle []interface{}
	var idiomas string
	var bibliografia map[string]interface{}
	var seguimiento map[string]interface{}
	var objetivosEspecificos []string
	var versionSyllabus string

	if syllabusData["objetivos_especificos"] != nil {
		objetivos := syllabusData["objetivos_especificos"].([]any)
		for _, objetivo := range objetivos {
			objetivoStr := fmt.Sprintf("%v", objetivo.(map[string]interface{})["objetivo"])
			objetivosEspecificos = append(objetivosEspecificos, objetivoStr)
		}
	} else {
		objetivosEspecificos = []string{}
	}

	contenido := syllabusData["contenido"]
	if contenido != nil {
		contenidoTematicoDescripcion = fmt.Sprintf(
			"%v",
			helpers.DefaultToMapString(contenido.(map[string]interface{}),
				"descripcion", ""))

		if contenido.(map[string]interface{})["temas"] == nil {
			contenidoTematicoDetalle = []interface{}{}
		} else {
			contenidoTematicoDetalle = contenido.(map[string]interface{})["temas"].([]interface{})
		}
	}

	evaluacion := syllabusData["evaluacion"]
	if evaluacion != nil {
		evaluacionDescripcion = fmt.Sprintf(
			"%v",
			helpers.DefaultToMapString(evaluacion.(map[string]interface{}), "descripcion", ""))

		if evaluacion.(map[string]interface{})["evaluaciones"] == nil {
			evaluacionDetalle = []any{}
		} else {
			evaluacionDetalle = evaluacion.(map[string]interface{})["evaluaciones"].([]interface{})
		}
	}

	if syllabusData["idioma_espacio_id"] != nil {
		idiomas = languages
	}

	if syllabusData["bibliografia"] != nil {
		bibliografia = syllabusData["bibliografia"].(map[string]interface{})
	}

	if syllabusData["seguimiento"] != nil {
		seguimiento = syllabusData["seguimiento"].(map[string]interface{})
	} else {
		seguimiento = map[string]interface{}{}
	}

	if syllabusData["resultados_aprendizaje"] != nil {
		propositos = syllabusData["resultados_aprendizaje"].([]interface{})
	} else {
		propositos = []interface{}{}
	}

	fechaRevConsejo := strings.Split(
		helpers.DefaultToMapString(seguimiento, "fechaRevisionConsejo", "").(string),
		"T")[0]
	fechaAprobConsejo := strings.Split(
		helpers.DefaultToMapString(seguimiento, "fechaAprobacionConsejo", "").(string),
		"T")[0]
	numActa := helpers.DefaultToMapString(seguimiento, "numeroActa", "").(string)

	if versionSyll := helpers.DefaultToMapString(syllabusData, "version", 0); versionSyll.(float64) > 0 {
		versionSyllabus = fmt.Sprintf("%v", versionSyll)
	} else {
		versionSyllabus = ""
	}

	syllabusTemplateData := map[string]interface{}{
		"nombre_facultad":                helpers.DefaultToMapString(facultyData, "Nombre", ""),
		"nombre_proyecto_curricular":     helpers.DefaultToMapString(projectData, "proyecto_curricular_nombre", ""),
		"cod_plan_estudio":               helpers.DefaultToMapString(syllabusData, "plan_estudios_id", ""),
		"nombre_espacio_academico":       helpers.DefaultToMapString(spaceData, "nombre_espacio_academico", ""),
		"cod_espacio_academico":          helpers.DefaultToMapString(spaceData, "cod_espacio_academico", ""),
		"num_creditos":                   helpers.DefaultToMapString(spaceData, "num_creditos", ""),
		"htd":                            helpers.DefaultToMapString(spaceData, "htd", ""),
		"htc":                            helpers.DefaultToMapString(spaceData, "htc", ""),
		"hta":                            helpers.DefaultToMapString(spaceData, "hta", ""),
		"es_asignatura":                  helpers.DefaultTo(spaceData["es_asignatura"], false),
		"es_catedra":                     helpers.DefaultTo(spaceData["es_catedra"], false),
		"es_obligatorio_basico":          helpers.DefaultTo(spaceData["es_obligatorio_basico"], false),
		"es_obligatorio_comp":            helpers.DefaultTo(spaceData["es_obligatorio_comp"], false),
		"es_electivo_int":                helpers.DefaultTo(spaceData["es_electivo_int"], false),
		"es_electivo_ext":                helpers.DefaultTo(spaceData["es_electivo_ext"], false),
		"es_electivo":                    helpers.DefaultTo(spaceData["es_electivo"], false),
		"es_teorico":                     false,
		"es_practico":                    false,
		"es_teorico_practico":            false,
		"es_presencial":                  false,
		"es_presencial_tic":              false,
		"es_virtual":                     false,
		"otra_modalidad":                 false,
		"cual_otra_modalidad":            "",
		"idiomas":                        helpers.DefaultTo(idiomas, ""),
		"sugerencias":                    helpers.DefaultToMapString(syllabusData, "sugerencias", ""),
		"justificacion":                  helpers.DefaultToMapString(syllabusData, "justificacion", ""),
		"objetivo_general":               helpers.DefaultToMapString(syllabusData, "objetivo_general", ""),
		"objetivos_especificos":          objetivosEspecificos,
		"propositos":                     propositos,
		"contenido_tematico_descripcion": helpers.DefaultTo(contenidoTematicoDescripcion, ""),
		"contenido_tematico_detalle":     contenidoTematicoDetalle,
		"estrategias_ensenanza":          syllabusData["estrategias"],
		"evaluacion_descripcion":         helpers.DefaultTo(evaluacionDescripcion, ""),
		"evaluacion_detalle":             evaluacionDetalle,
		"medios_recursos":                helpers.DefaultToMapString(syllabusData, "recursos_educativos", ""),
		"practicas_salidas":              helpers.DefaultToMapString(syllabusData, "practicas_academicas", ""),
		"bibliografia_basica":            bibliografia["basicas"],
		"bibliografia_complementaria":    bibliografia["complementarias"],
		"bibliografia_paginas":           bibliografia["paginasWeb"],
		"fecha_rev_consejo":              fechaRevConsejo,
		"fecha_aprob_consejo":            fechaAprobConsejo,
		"num_acta":                       numActa,
		"version_syllabus":               versionSyllabus}

	return syllabusTemplateData
}

func GetSyllabusTemplate(syllabusTemplateData map[string]interface{}, syllabusTemplateResponse *map[string]interface{}, format string) {
	var url string
	if strings.ToLower(format) == "pdf" {
		url = beego.AppConfig.String("SyllabusService") + "v2/syllabus/template"
	} else {
		url = beego.AppConfig.String("SyllabusService") + "syllabus/template/spreadsheet"
	}
	if err := helpers.SendJson(
		url,
		"POST",
		&syllabusTemplateResponse,
		syllabusTemplateData); err != nil {
		panic(map[string]interface{}{
			"funcion": "GenerarTemplate",
			"err":     "Error al generar el documento del syllabus ",
			"status":  "400",
			"log":     err})
	}
}

func GetAcademicSpaceData(pensumId, carreraCod, asignaturaCod int) (map[string]any, error) {
	var spaceResponse map[string]interface{}

	spaceErr := request.GetJsonWSO2(
		beego.AppConfig.String("AcademicaEspacioAcademicoService")+
			fmt.Sprintf("detalle_espacio_academico/%v/%v/%v", pensumId, carreraCod, asignaturaCod),
		&spaceResponse)

	if spaceErr == nil && fmt.Sprintf("%v", spaceResponse) != "map[espacios_academicos:map[]]" && fmt.Sprintf("%v", spaceResponse) != "map[]]" {
		spaces := spaceResponse["espacios_academicos"].(map[string]interface{})["espacio_academico"].([]interface{})
		if len(spaces) > 0 {
			space := spaces[0].(map[string]interface{})

			esAsignatura := strings.ToLower(fmt.Sprintf("%v", space["tipo"])) == "asignatura"
			spaceType := strings.ToLower(fmt.Sprintf("%v", space["cea_abr"]))
			spaceData := map[string]interface{}{
				"nombre_espacio_academico": fmt.Sprintf("%v", helpers.DefaultToMapString(space, "asi_nombre", "")),
				"cod_espacio_academico":    fmt.Sprintf("%v", helpers.DefaultToMapString(space, "asi_cod", "")),
				"num_creditos":             fmt.Sprintf("%v", helpers.DefaultToMapString(space, "pen_cre", "")),
				"htd":                      fmt.Sprintf("%v", helpers.DefaultToMapString(space, "pen_nro_ht", "")),
				"htc":                      fmt.Sprintf("%v", helpers.DefaultToMapString(space, "pen_nro_hp", "")),
				"hta":                      fmt.Sprintf("%v", helpers.DefaultToMapString(space, "pen_nro_aut", "")),
				"es_asignatura":            esAsignatura,
				"es_catedra":               !esAsignatura,
				"es_obligatorio_basico":    spaceType == "ob",
				"es_obligatorio_comp":      spaceType == "oc",
				"es_electivo_int":          spaceType == "ei",
				"es_electivo_ext":          spaceType == "ee",
				"es_electivo":              spaceType == "e",
			}
			return spaceData, nil
		} else {
			return nil, fmt.Errorf("Espacio académico no encontrado")
		}
	} else {
		return nil, fmt.Errorf("Espacio académico no encontrado")
	}
}

func GetIdiomas(idiomaIds []interface{}) (string, error) {
	var idiomaResponse []map[string]interface{}
	idiomasStr := ""

	idiomaErr := request.GetJson(
		beego.AppConfig.String("IdiomaService")+"idioma",
		&idiomaResponse)

	if idiomaErr == nil {
		for i, id := range idiomaIds {
			for _, idioma := range idiomaResponse {
				if idioma["Id"] == id {
					if i == len(idiomaIds)-1 {
						idiomasStr += idioma["Nombre"].(string)
					} else {
						idiomasStr += idioma["Nombre"].(string) + ", "
					}
					break
				}
			}
		}
		return idiomasStr, nil
	} else {
		return "", fmt.Errorf("Idiomas no encontrados")
	}
}
