package models

import "time"

// ProgramaAcademicoCore representa un registro del reporte de academica_core.
// Los tipos corresponden a los valores reales que expone el servicio para
// MNTAC.ACCRA; los campos nullable se modelan con punteros.
type ProgramaAcademicoCore struct {
	Codigo              int        `json:"CRA_COD"`
	Nombre              string     `json:"CRA_NOMBRE"`
	DependenciaCodigo   int        `json:"CRA_DEP_COD"`
	Estado              string     `json:"CRA_ESTADO"`
	CodigoIcfes         *string    `json:"CRA_COD_ICFES"`
	EmpleadoCodigo      int        `json:"CRA_EMP_COD"`
	TipoCarrera         int        `json:"CRA_TIP_CRA"`
	ResolucionSup       string     `json:"CRA_RESOL_SUP"`
	FechaAprobIcfes     *time.Time `json:"CRA_FECHA_APROB_ICFES"`
	FechaUltRenov       *time.Time `json:"CRA_FECHA_ULT_RENOV"`
	Abreviatura         string     `json:"CRA_ABREV"`
	Jornada             string     `json:"CRA_JORNADA"`
	ProgramaCoordinador *string    `json:"CRA_PRO_COOR"`
	IdentificacionCoord string     `json:"CRA_EMP_NRO_IDEN"`
	SeOfrece            string     `json:"CRA_SE_OFRECE"`
	CodigoSnies         string     `json:"CRA_COD_SNIES"`
	JustificacionCodigo *string    `json:"CRA_JUST_COD"`
	IndicadorCiclo      string     `json:"CRA_IND_CICLO"`
	NotaAprobacion      int        `json:"CRA_NOTA_APROB"`
	Email               string     `json:"CRA_EMAIL"`
	OpcionAdmision      *string    `json:"CRA_OPCION_ADM"`
	Alias               string     `json:"CRA_ALIAS"`
}
