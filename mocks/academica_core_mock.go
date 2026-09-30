package mocks

import "github.com/udistrital/syllabus_mid/models"

// GetProgramasAcademicosByIdentificacion es un mock temporal de la respuesta de
// academica_core mientras se define su endpoint real. Devuelve los programas
// asociados a la identificación (hoy solo se conoce la data de ejemplo).
func GetProgramasAcademicosByIdentificacion(identificacion string) []models.ProgramaAcademicoCore {
	if identificacion != "45512868" {
		return []models.ProgramaAcademicoCore{}
	}

	return []models.ProgramaAcademicoCore{
		{
			Codigo:              54,
			Nombre:              "MAESTRÍA EN ESTUDIOS EDUCATIVOS AFROCOLOMBIANOS Y AFROLATINOAMERICANOS (DIST)",
			DependenciaCodigo:   24,
			Estado:              "A",
			CodigoIcfes:         "",
			EmpleadoCodigo:      3531,
			TipoCarrera:         12,
			ResolucionSup:       "RESOLUCIÓN 007 DE 26 SEPTIEMBRE 2024",
			FechaAprobIcfes:     "",
			FechaUltRenov:       "",
			Abreviatura:         "MAE. EST. AFROCOL. AFROLAT. D.",
			Jornada:             "DIURNA",
			ProgramaCoordinador: "",
			IdentificacionCoord: "45512868",
			SeOfrece:            "N",
			CodigoSnies:         "119092",
			JustificacionCodigo: "",
			IndicadorCiclo:      "N",
			NotaAprobacion:      35,
			Email:               "maeafro@udistrital.edu.co",
			OpcionAdmision:      "1",
			Alias:               "PROYECTO CURRICULAR",
		},
	}
}
