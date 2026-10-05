package mocks

import "github.com/udistrital/syllabus_mid/models"

// GetProgramasAcademicosByIdentificacion es un mock temporal de la respuesta de
// academica_core mientras se define su endpoint real. Devuelve los programas
// asociados a la identificación (hoy solo se conoce la data de ejemplo).
func GetProgramasAcademicosByIdentificacion(identificacion string) []models.ProgramaAcademicoCore {
	if identificacion != "45512868" {
		return []models.ProgramaAcademicoCore{}
	}

	opcionAdmision := "1"

	return []models.ProgramaAcademicoCore{
		{
			Codigo:              54,
			Nombre:              "MAESTRÍA EN ESTUDIOS EDUCATIVOS AFROCOLOMBIANOS Y AFROLATINOAMERICANOS (DIST)",
			DependenciaCodigo:   24,
			Estado:              "A",
			CodigoIcfes:         nil,
			EmpleadoCodigo:      3531,
			TipoCarrera:         12,
			ResolucionSup:       "RESOLUCIÓN 007 DE 26 SEPTIEMBRE 2024",
			FechaAprobIcfes:     nil,
			FechaUltRenov:       nil,
			Abreviatura:         "MAE. EST. AFROCOL. AFROLAT. D.",
			Jornada:             "DIURNA",
			ProgramaCoordinador: nil,
			IdentificacionCoord: "45512868",
			SeOfrece:            "N",
			CodigoSnies:         "119092",
			JustificacionCodigo: nil,
			IndicadorCiclo:      "N",
			NotaAprobacion:      35,
			Email:               "maeafro@udistrital.edu.co",
			OpcionAdmision:      &opcionAdmision,
			Alias:               "PROYECTO CURRICULAR",
		},
	}
}
