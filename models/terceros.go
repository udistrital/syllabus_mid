package models

// Tercero es la representación mínima de un tercero en terceros_crud.
type Tercero struct {
	Id             int    `json:"Id"`
	NombreCompleto string `json:"NombreCompleto"`
	UsuarioWSO2    string `json:"UsuarioWSO2"`
}

// DatosIdentificacion es un registro de datos_identificacion de terceros_crud.
type DatosIdentificacion struct {
	Numero    string  `json:"Numero"`
	Activo    bool    `json:"Activo"`
	TerceroId Tercero `json:"TerceroId"`
}
