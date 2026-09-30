package models

// VinculacionPrograma es el objeto de respuesta por cada programa asociado al
// coordinador consultado.
type VinculacionPrograma struct {
	IdOikos           int    `json:"IdOikos"`
	Codigo            int    `json:"Codigo"`
	Nombre            string `json:"Nombre"`
	CorreoElectronico string `json:"CorreoElectronico"`
	Coordinador       string `json:"Coordinador"`
	Identificacion    string `json:"Identificacion"`
	PadreOikos        int    `json:"PadreOikos"`
}
