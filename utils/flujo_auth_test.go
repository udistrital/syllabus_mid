package utils

import (
	"testing"

	"github.com/udistrital/syllabus_mid/models"
)

func TestPlazoVencido(t *testing.T) {
	pasado := map[string]any{"flujo": map[string]any{"fecha_limite": "2000-01-01"}}
	if !PlazoVencido(pasado) {
		t.Fatalf("se esperaba plazo vencido para una fecha pasada")
	}

	futuro := map[string]any{"flujo": map[string]any{"fecha_limite": "2999-12-31"}}
	if PlazoVencido(futuro) {
		t.Fatalf("no se esperaba plazo vencido para una fecha futura")
	}

	sinPlazo := map[string]any{"flujo": map[string]any{}}
	if PlazoVencido(sinPlazo) {
		t.Fatalf("sin fecha límite no debe considerarse vencido")
	}
}

func TestEsDocenteAsignado(t *testing.T) {
	draft := map[string]any{
		"flujo": map[string]any{
			"asignado_a":           float64(1042),
			"asignado_a_documento": "80111222",
		},
	}

	porId := models.Identidad{TerceroId: 1042}
	if !EsDocenteAsignado(draft, porId) {
		t.Fatalf("debía reconocer al docente por tercero_id")
	}

	porDocumento := models.Identidad{TerceroId: 0, Documento: "80111222"}
	if !EsDocenteAsignado(draft, porDocumento) {
		t.Fatalf("debía reconocer al docente por documento")
	}

	otro := models.Identidad{TerceroId: 999, Documento: "70000000"}
	if EsDocenteAsignado(draft, otro) {
		t.Fatalf("no debía reconocer a un tercero distinto")
	}
}

func TestSyllabusIntSlice(t *testing.T) {
	s := map[string]any{"proyecto_curricular_ids": []any{float64(375), 12.0}}
	got := SyllabusIntSlice(s, "proyecto_curricular_ids")
	if len(got) != 2 || got[0] != 375 || got[1] != 12 {
		t.Fatalf("conversión inesperada: %v", got)
	}
}

func TestIntersectaProgramas(t *testing.T) {
	codigos := map[int]bool{375: true, 412: true}

	if !intersectaProgramas(codigos, []int{412, 999}) {
		t.Fatalf("debía encontrar intersección en 412")
	}
	if intersectaProgramas(codigos, []int{999, 1000}) {
		t.Fatalf("no debía encontrar intersección")
	}
	if intersectaProgramas(codigos, nil) {
		t.Fatalf("sin proyectos no hay intersección")
	}
	if intersectaProgramas(map[int]bool{}, []int{375}) {
		t.Fatalf("sin programas del coordinador no hay intersección")
	}
}
