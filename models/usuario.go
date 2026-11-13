package models

// Usuario representa la respuesta del endpoint oauth2/userinfo de WSO2.
type Usuario struct {
	Documento          string `json:"documento"`
	DocumentoCompuesto string `json:"documento_compuesto"`
	Email              string `json:"email"`
	Role               string `json:"role"`
	Sub                string `json:"sub"`
}
