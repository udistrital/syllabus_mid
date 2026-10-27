package utils

import (
	"context"
	"net/http"
	"time"

	"github.com/udistrital/syllabus_mid/models"
)

// GetCodigosProgramaDeCoordinador devuelve el conjunto de códigos de programa
// (CRA_COD, "id oracle") de los programas académicos vinculados al coordinador.
// Se usan los ids de academica_core porque son los mismos que el syllabus
// almacena en proyecto_curricular_ids; no se debe convertir a id_oikos.
func GetCodigosProgramaDeCoordinador(ctx context.Context, documento string) (map[int]bool, error) {
	if documento == "" {
		return nil, &QueryError{Status: http.StatusBadRequest, Message: "El documento del coordinador es obligatorio"}
	}
	programas, err := GetProgramasAcademicosByCoordinador(ctx, documento)
	if err != nil {
		return nil, err
	}
	codigos := map[int]bool{}
	for _, programa := range programas {
		codigos[programa.Codigo] = true
	}
	return codigos, nil
}

// intersectaProgramas indica si alguno de los proyectos curriculares está en el
// conjunto de códigos de programa del coordinador.
func intersectaProgramas(codigos map[int]bool, proyectoIds []int) bool {
	for _, proyectoId := range proyectoIds {
		if codigos[proyectoId] {
			return true
		}
	}
	return false
}

// EsCoordinadorDePrograma valida que el coordinador tenga alguno de los
// proyectos curriculares indicados dentro de sus programas vinculados.
func EsCoordinadorDePrograma(ctx context.Context, documento string, proyectoIds []int) (bool, error) {
	codigos, err := GetCodigosProgramaDeCoordinador(ctx, documento)
	if err != nil {
		return false, err
	}
	return intersectaProgramas(codigos, proyectoIds), nil
}

// EsDocenteAsignado valida que la identidad corresponda al docente asignado.
func EsDocenteAsignado(syllabus map[string]any, identidad models.Identidad) bool {
	flujo := FlujoMap(syllabus)
	if id, ok := FlujoInt(flujo, "asignado_a"); ok && identidad.TerceroId != 0 && id == identidad.TerceroId {
		return true
	}
	if doc := FlujoString(flujo, "asignado_a_documento"); doc != "" && doc == identidad.Documento {
		return true
	}
	return false
}

// PlazoVencido indica si la fecha_limite del flujo ya pasó (se cuenta hasta el
// final del día indicado).
func PlazoVencido(syllabus map[string]any) bool {
	flujo := FlujoMap(syllabus)
	limite := FlujoString(flujo, "fecha_limite")
	if limite == "" {
		return false
	}
	fecha, err := time.Parse("2006-01-02", limite)
	if err != nil {
		if fechaRFC, errRFC := time.Parse(time.RFC3339, limite); errRFC == nil {
			fecha = fechaRFC
		} else {
			return false
		}
	}
	fin := time.Date(fecha.Year(), fecha.Month(), fecha.Day(), 23, 59, 59, 0, fecha.Location())
	return time.Now().After(fin)
}
