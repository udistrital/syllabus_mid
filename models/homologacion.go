package models

// HomologacionResponse representa la respuesta de homologación por SNIES.
// Varios campos pueden venir null o no estar definidos, por eso se modelan como any.
type HomologacionResponse struct {
	Homologacion HomologacionProyecto `json:"homologacion"`
}

type HomologacionProyecto struct {
	ProyectoSnies any `json:"proyecto_snies"`
	IdOikos       any `json:"id_oikos"`
	IdArgo        any `json:"id_argo"`
	IdSnies       any `json:"id_snies"`
}
