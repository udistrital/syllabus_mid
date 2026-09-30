package models

// ProgramaAcademicoCore representa un registro del reporte de academica_core.
// Los campos any se dejan flexibles porque no se conoce su tipo/valor real
// (pueden venir como int o string en el JSON).
type ProgramaAcademicoCore struct {
	Codigo              int    `json:"CRA_COD"`
	Nombre              string `json:"CRA_NOMBRE"`
	DependenciaCodigo   int    `json:"CRA_DEP_COD"`
	Estado              string `json:"CRA_ESTADO"`
	CodigoIcfes         any    `json:"CRA_COD_ICFES"`
	EmpleadoCodigo      int    `json:"CRA_EMP_COD"`
	TipoCarrera         int    `json:"CRA_TIP_CRA"`
	ResolucionSup       string `json:"CRA_RESOL_SUP"`
	FechaAprobIcfes     any    `json:"CRA_FECHA_APROB_ICFES"`
	FechaUltRenov       any    `json:"CRA_FECHA_ULT_RENOV"`
	Abreviatura         string `json:"CRA_ABREV"`
	Jornada             string `json:"CRA_JORNADA"`
	ProgramaCoordinador any    `json:"CRA_PRO_COOR"`
	IdentificacionCoord string `json:"CRA_EMP_NRO_IDEN"`
	SeOfrece            string `json:"CRA_SE_OFRECE"`
	CodigoSnies         string `json:"CRA_COD_SNIES"`
	JustificacionCodigo any    `json:"CRA_JUST_COD"`
	IndicadorCiclo      string `json:"CRA_IND_CICLO"`
	NotaAprobacion      int    `json:"CRA_NOTA_APROB"`
	Email               string `json:"CRA_EMAIL"`
	OpcionAdmision      any    `json:"CRA_OPCION_ADM"`
	Alias               string `json:"CRA_ALIAS"`
}
