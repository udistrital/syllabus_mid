package test

import (
	"testing"

	"github.com/udistrital/syllabus_mid/utils"
)

func plantillaDe(t *testing.T, syllabusData map[string]interface{}) map[string]interface{} {
	t.Helper()
	return utils.GetSyllabusTemplateData(
		map[string]interface{}{},
		syllabusData,
		map[string]interface{}{},
		map[string]interface{}{},
		"")
}

func TestPropositosMatrizDeRelacion(t *testing.T) {
	html := "<table><thead><tr><th>Propósitos de Formación</th><th>Relación</th></tr></thead>" +
		"<tbody><tr><td>Propósito 1</td><td>Relación 1</td></tr></tbody></table>"

	plantilla := plantillaDe(t, map[string]interface{}{
		"resultados_aprendizaje": []interface{}{
			map[string]interface{}{"pfa_matriz": html},
		},
	})

	propositos, ok := plantilla["propositos"].([]interface{})
	if !ok || len(propositos) != 1 {
		t.Fatalf("se esperaba un propósito en la matriz, se obtuvieron %v", plantilla["propositos"])
	}

	proposito, ok := propositos[0].(map[string]interface{})
	if !ok {
		t.Fatalf("el propósito no es un mapa: %v", propositos[0])
	}

	if proposito["pfa_matriz"] != html {
		t.Fatalf("no se conservó el HTML de la matriz: %v", proposito["pfa_matriz"])
	}
}

func TestPropositosListado(t *testing.T) {
	plantilla := plantillaDe(t, map[string]interface{}{
		"resultados_aprendizaje": []interface{}{
			map[string]interface{}{"pfa_programa": "PFA del programa"},
		},
	})

	propositos, ok := plantilla["propositos"].([]interface{})
	if !ok || len(propositos) != 1 {
		t.Fatalf("se esperaba un propósito en el listado, se obtuvieron %v", plantilla["propositos"])
	}

	proposito := propositos[0].(map[string]interface{})
	if proposito["pfa_programa"] != "PFA del programa" {
		t.Fatalf("no se conservó el propósito del listado: %v", proposito["pfa_programa"])
	}
}

func TestPropositosFormatoJerarquico(t *testing.T) {
	plantilla := plantillaDe(t, map[string]interface{}{
		"resultados_aprendizaje": []interface{}{
			map[string]interface{}{
				"competencia": "Competencia 1",
				"resultados": []interface{}{
					map[string]interface{}{
						"id":                  "01",
						"dominio":             "Cognitivo - Conocer",
						"resultado_detallado": "Resultado 1",
					},
				},
			},
		},
	})

	propositos, ok := plantilla["propositos"].([]interface{})
	if !ok || len(propositos) != 1 {
		t.Fatalf("se esperaba un propósito jerárquico, se obtuvieron %v", plantilla["propositos"])
	}

	proposito := propositos[0].(map[string]interface{})
	if proposito["competencia"] != "Competencia 1" || proposito["resultado_detallado"] != "Resultado 1" {
		t.Fatalf("no se conservó el formato jerárquico: %v", proposito)
	}
}

func TestPropositosSinResultados(t *testing.T) {
	plantilla := plantillaDe(t, map[string]interface{}{})

	propositos, ok := plantilla["propositos"].([]interface{})
	if !ok || len(propositos) != 0 {
		t.Fatalf("se esperaba la lista de propósitos vacía, se obtuvieron %v", plantilla["propositos"])
	}
}

func TestEsPropuesta(t *testing.T) {
	if !utils.EsPropuesta(map[string]interface{}{"es_propuesta": true}) {
		t.Fatalf("se esperaba que el documento con es_propuesta=true fuera propuesta")
	}

	for _, datos := range []map[string]interface{}{
		{},
		{"es_propuesta": false},
		{"es_propuesta": "true"},
	} {
		if utils.EsPropuesta(datos) {
			t.Fatalf("se esperaba que el documento %v no fuera propuesta", datos)
		}
	}
}

func TestBuildSpaceDataDesdePropuesta(t *testing.T) {
	spaceData := utils.BuildSpaceDataDesdePropuesta(map[string]interface{}{
		"es_propuesta":       true,
		"propuesta_nombre":   "Cálculo Diferencial",
		"propuesta_creditos": 4,
	})

	if spaceData["nombre_espacio_academico"] != "Cálculo Diferencial" {
		t.Fatalf("no se conservó el nombre de la asignatura: %v", spaceData["nombre_espacio_academico"])
	}

	if spaceData["num_creditos"] != "4" {
		t.Fatalf("no se conservaron los créditos de la propuesta: %v", spaceData["num_creditos"])
	}

	for _, campo := range []string{"cod_espacio_academico", "htd", "htc", "hta", "naturaleza"} {
		if spaceData[campo] != "" {
			t.Fatalf("el campo %v depende del espacio real y debe quedar vacío: %v", campo, spaceData[campo])
		}
	}
}

func TestPlantillaMarcaSinCabeceraParaPropuesta(t *testing.T) {
	propuesta := plantillaDe(t, map[string]interface{}{
		"es_propuesta":         true,
		"propuesta_nombre":     "Cálculo Diferencial",
		"espacio_academico_id": 0,
	})

	if propuesta["sin_cabecera"] != true {
		t.Fatalf("la plantilla de la propuesta debe pedir omitir la cabecera: %v", propuesta["sin_cabecera"])
	}

	normal := plantillaDe(t, map[string]interface{}{})
	if normal["sin_cabecera"] != false {
		t.Fatalf("la plantilla de un syllabus normal debe conservar la cabecera: %v", normal["sin_cabecera"])
	}
}
