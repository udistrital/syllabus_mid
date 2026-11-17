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
