package models

// Estados del flujo de aprobación (deben coincidir con el CRUD).
const (
	EstadoActivo     = "activo"
	EstadoAsignado   = "asignado"
	EstadoEnRevision = "en_revision"
	EstadoRechazado  = "rechazado"
	EstadoExpirado   = "expirado"
	EstadoCancelado  = "cancelado"
	EstadoAprobado   = "aprobado"
)

// EstadosBorrador agrupa los estados de una propuesta no publicada.
var EstadosBorrador = []string{EstadoAsignado, EstadoEnRevision, EstadoRechazado, EstadoExpirado}

// Identidad representa al tercero autenticado que ejecuta una acción del flujo.
type Identidad struct {
	TerceroId      int
	Documento      string
	NombreCompleto string
}

// HistorialEstado es una entrada de auditoría del flujo.
type HistorialEstado struct {
	Estado           string `json:"estado"`
	Fecha            string `json:"fecha,omitempty"`
	TerceroId        int    `json:"tercero_id,omitempty"`
	TerceroDocumento string `json:"tercero_documento,omitempty"`
	Observacion      string `json:"observacion,omitempty"`
}

// Flujo es la metadata de la asignación/aprobación de un syllabus.
type Flujo struct {
	AsignadoA            int               `json:"asignado_a,omitempty"`
	AsignadoADocumento   string            `json:"asignado_a_documento,omitempty"`
	AsignadoPor          int               `json:"asignado_por,omitempty"`
	AsignadoPorDocumento string            `json:"asignado_por_documento,omitempty"`
	FechaAsignacion      string            `json:"fecha_asignacion,omitempty"`
	FechaLimite          string            `json:"fecha_limite,omitempty"`
	FechaEnvio           string            `json:"fecha_envio,omitempty"`
	FechaResolucion      string            `json:"fecha_resolucion,omitempty"`
	ResueltoPor          int               `json:"resuelto_por,omitempty"`
	MotivoRechazo        string            `json:"motivo_rechazo,omitempty"`
	Historial            []HistorialEstado `json:"historial,omitempty"`
}

// AsignarRequest es el cuerpo para asignar la creación/edición a un docente.
type AsignarRequest struct {
	Documento             string `json:"documento"`
	DocumentoDocente      string `json:"documento_docente"`
	FechaLimite           string `json:"fecha_limite"`
	SyllabusCode          string `json:"syllabus_code"`
	EspacioAcademicoId    int    `json:"espacio_academico_id"`
	ProyectoCurricularIds []int  `json:"proyecto_curricular_ids"`
	PlanEstudiosIds       []int  `json:"plan_estudios_ids"`
	Observacion           string `json:"observacion"`
}

// TransicionRequest es el cuerpo común de proponer/aprobar/rechazar/reabrir/cancelar.
type TransicionRequest struct {
	Documento     string `json:"documento"`
	Observacion   string `json:"observacion"`
	FechaLimite   string `json:"fecha_limite"`
	MotivoRechazo string `json:"motivo_rechazo"`
	NumeroActa    string `json:"numeroActa"`
	Archivo       string `json:"archivo"`
}
