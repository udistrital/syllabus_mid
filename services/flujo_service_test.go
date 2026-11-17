package services

import (
	"testing"

	"github.com/udistrital/syllabus_mid/models"
)

func TestBuildDraftPayloadDesdeActivo(t *testing.T) {
	origen := map[string]any{
		"_id":                     "abc",
		"version":                 float64(3),
		"syllabus_actual":         true,
		"estado":                  "activo",
		"syllabus_code":           "uuid-1",
		"espacio_academico_id":    float64(3401),
		"proyecto_curricular_ids": []any{float64(375)},
		"plan_estudios_ids":       []any{float64(12)},
		"justificacion":           "contenido previo",
		"seguimiento":             map[string]any{"numeroActa": "123"},
	}
	req := models.AsignarRequest{
		SyllabusCode: "uuid-1",
		FechaLimite:  "2026-11-30",
	}
	coord := models.Identidad{TerceroId: 887, Documento: "79456789", NombreCompleto: "Coord X"}
	docente := models.Identidad{TerceroId: 1042, Documento: "80111222", NombreCompleto: "Docente Y"}

	payload := buildDraftPayload(origen, req, coord, docente)

	if payload["estado"] != models.EstadoAsignado {
		t.Fatalf("estado esperado asignado, got %v", payload["estado"])
	}
	if _, ok := payload["version"]; ok {
		t.Fatalf("no debe conservar version")
	}
	if _, ok := payload["syllabus_actual"]; ok {
		t.Fatalf("no debe conservar syllabus_actual")
	}
	if _, ok := payload["_id"]; ok {
		t.Fatalf("no debe conservar _id")
	}
	if _, ok := payload["seguimiento"]; ok {
		t.Fatalf("no debe conservar seguimiento")
	}
	if payload["justificacion"] != "contenido previo" {
		t.Fatalf("debe conservar el contenido previo")
	}

	flujo, ok := payload["flujo"].(map[string]any)
	if !ok {
		t.Fatalf("debe incluir flujo")
	}
	if flujo["asignado_a"] != 1042 || flujo["fecha_limite"] != "2026-11-30" {
		t.Fatalf("flujo inesperado: %v", flujo)
	}
}

func TestBuildSeguimientoEnRevision(t *testing.T) {
	original := map[string]any{
		"elaboro":       "Docente Y",
		"fecha_elaboro": "2026-10-01T10:00:00Z",
		"numeroActa":    "",
	}
	resultado := buildSeguimientoEnRevision(original, "2026-10-13T12:00:00Z")

	if resultado["fechaRevisionConsejo"] != "2026-10-13T12:00:00Z" {
		t.Fatalf("debe fijar fechaRevisionConsejo")
	}
	if resultado["elaboro"] != "Docente Y" || resultado["fecha_elaboro"] != "2026-10-01T10:00:00Z" {
		t.Fatalf("debe preservar los campos previos del seguimiento")
	}
	if _, mutado := original["fechaRevisionConsejo"]; mutado {
		t.Fatalf("no debe mutar el mapa original")
	}
}

func TestMarcaAguaTexto(t *testing.T) {
	if got := marcaAguaTexto(models.EstadoEnRevision); got != "EN REVISIÓN" {
		t.Fatalf("en_revision inesperado: %q", got)
	}
	for _, estado := range []string{models.EstadoAsignado, models.EstadoRechazado, models.EstadoExpirado} {
		if got := marcaAguaTexto(estado); got != "BORRADOR" {
			t.Fatalf("%s inesperado: %q", estado, got)
		}
	}
}
