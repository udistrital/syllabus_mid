package models

// DependenciasPadresResponse representa la respuesta de OIKOS para
// dependencia/get_dependencias_padres_by_id/{id}.
type DependenciasPadresResponse struct {
	Body []DependenciaOikos `json:"Body"`
	Type string             `json:"Type"`
}

type DependenciaOikos struct {
	Id       int    `json:"Id"`
	Nombre   string `json:"Nombre"`
	Padre    int    `json:"Padre"`
	Hija     int    `json:"Hija"`
	Opciones any    `json:"Opciones"`
}
