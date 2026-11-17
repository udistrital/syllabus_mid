package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/udistrital/syllabus_mid/models"
	"github.com/udistrital/syllabus_mid/utils"
	"github.com/udistrital/utils_oas/requestresponse"
)

// campos que nunca deben reenviarse al CRUD en una edición de contenido
var camposReservadosContenido = []string{
	"_id", "syllabus_code", "version", "syllabus_actual", "estado",
	"creado_por", "flujo", "fecha_creacion", "fecha_modificacion", "documento",
}

func respErr(status int, message string) requestresponse.APIResponse {
	return requestresponse.APIResponseDTO(false, status, nil, message)
}

// apiError traduce un error de dominio (QueryError) al código HTTP asociado.
func apiError(err error) requestresponse.APIResponse {
	return requestresponse.APIResponseDTO(false, statusFromError(err), nil, err.Error())
}

// buildSeguimientoEnRevision devuelve una copia del seguimiento con la fecha de
// revisión de consejo actualizada (se registra al proponer la propuesta).
func buildSeguimientoEnRevision(seguimiento map[string]any, nowISO string) map[string]any {
	copia := map[string]any{}
	for k, v := range seguimiento {
		copia[k] = v
	}
	copia["fechaRevisionConsejo"] = nowISO
	return copia
}

func buildDraftPayload(origen map[string]any, req models.AsignarRequest,
	coord, docente models.Identidad) map[string]any {
	payload := map[string]any{}
	if origen != nil {
		for k, v := range origen {
			payload[k] = v
		}
	}
	for _, k := range append(append([]string{}, camposReservadosContenido...), "seguimiento") {
		delete(payload, k)
	}

	if req.SyllabusCode != "" {
		payload["syllabus_code"] = req.SyllabusCode
	}
	if req.EspacioAcademicoId != 0 {
		payload["espacio_academico_id"] = req.EspacioAcademicoId
	}
	if len(req.ProyectoCurricularIds) > 0 {
		payload["proyecto_curricular_ids"] = req.ProyectoCurricularIds
	}
	if len(req.PlanEstudiosIds) > 0 {
		payload["plan_estudios_ids"] = req.PlanEstudiosIds
	}

	payload["estado"] = models.EstadoAsignado
	payload["tercero_id"] = coord.TerceroId
	payload["creado_por"] = coord.TerceroId
	payload["flujo"] = map[string]any{
		"asignado_a":             docente.TerceroId,
		"asignado_a_documento":   docente.Documento,
		"asignado_por":           coord.TerceroId,
		"asignado_por_documento": coord.Documento,
		"fecha_limite":           req.FechaLimite,
	}
	return payload
}

// AsignarSyllabus crea la propuesta (versión inicial o nueva versión) y la
// asigna a un docente dentro de un plazo.
func AsignarSyllabus(ctx context.Context, body []byte) requestresponse.APIResponse {
	var req models.AsignarRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return respErr(http.StatusBadRequest, "La petición no tiene un cuerpo JSON válido")
	}
	if req.DocumentoDocente == "" {
		return respErr(http.StatusBadRequest, "El documento del docente es obligatorio")
	}
	if req.FechaLimite == "" {
		return respErr(http.StatusBadRequest, "La fecha límite es obligatoria")
	}

	coord, err := utils.ResolverIdentidad(ctx, req.Documento)
	if err != nil {
		return apiError(err)
	}

	docenteTercero, err := utils.GetTerceroByDocumento(ctx, req.DocumentoDocente)
	if err != nil {
		return apiError(err)
	}
	if docenteTercero == nil {
		return respErr(http.StatusNotFound,
			fmt.Sprintf("No se encontró un tercero activo con documento %s", req.DocumentoDocente))
	}
	docente := models.Identidad{
		TerceroId:      docenteTercero.Id,
		Documento:      req.DocumentoDocente,
		NombreCompleto: docenteTercero.NombreCompleto,
	}

	var origen map[string]any
	if req.SyllabusCode != "" {
		activo, err := utils.GetSyllabusActivo(ctx, req.SyllabusCode)
		if err != nil {
			return apiError(err)
		}
		if activo == nil {
			return respErr(http.StatusNotFound, "No existe una versión activa para el syllabus indicado")
		}
		abierto, err := utils.GetBorradorAbierto(ctx, req.SyllabusCode)
		if err != nil {
			return apiError(err)
		}
		if abierto != nil {
			return respErr(http.StatusConflict, "Ya existe una propuesta en curso para este syllabus")
		}
		ok, err := utils.EsCoordinadorDePrograma(ctx, coord.Documento,
			utils.SyllabusIntSlice(activo, "proyecto_curricular_ids"))
		logs.Info("log si es coordinador del syllabus publicado")
		fmt.Println(ok)
		fmt.Println("")
		if err != nil {
			return apiError(err)
		}
		if !ok {
			return respErr(http.StatusForbidden, "No es coordinador de este syllabus")
		}
		origen = activo
	} else {
		if req.EspacioAcademicoId == 0 || len(req.ProyectoCurricularIds) == 0 || len(req.PlanEstudiosIds) == 0 {
			return respErr(http.StatusBadRequest,
				"Faltan datos del espacio académico, proyecto curricular y plan de estudios")
		}
		ok, err := utils.EsCoordinadorDePrograma(ctx, coord.Documento, req.ProyectoCurricularIds)
		if err != nil {
			return apiError(err)
		}
		if !ok {
			return respErr(http.StatusForbidden, "No es coordinador del programa indicado")
		}
	}

	created, err := utils.PostSyllabus(ctx, buildDraftPayload(origen, req, coord, docente))
	if err != nil {
		return apiError(err)
	}
	return requestresponse.APIResponseDTO(true, http.StatusCreated, created, "Syllabus asignado al docente")
}

// GuardarBorrador permite al docente asignado guardar el contenido de la propuesta.
func GuardarBorrador(ctx context.Context, code string, body []byte) requestresponse.APIResponse {
	var identReq models.TransicionRequest
	_ = json.Unmarshal(body, &identReq)

	ident, err := utils.ResolverIdentidad(ctx, identReq.Documento)
	if err != nil {
		return apiError(err)
	}
	draft, err := utils.GetBorradorAbierto(ctx, code)
	if err != nil {
		return apiError(err)
	}
	if draft == nil {
		return respErr(http.StatusNotFound, "No hay una propuesta en curso para este syllabus")
	}
	if !utils.EsDocenteAsignado(draft, ident) {
		return respErr(http.StatusForbidden, "Solo el docente asignado puede editar la propuesta")
	}
	estado := utils.SyllabusString(draft, "estado")
	if estado != models.EstadoAsignado && estado != models.EstadoRechazado {
		return respErr(http.StatusConflict, "La propuesta no es editable en estado "+estado)
	}
	if utils.PlazoVencido(draft) {
		id := utils.SyllabusString(draft, "_id")
		_, _ = utils.PutEstado(ctx, id, map[string]any{
			"estado":            models.EstadoExpirado,
			"tercero_id":        ident.TerceroId,
			"tercero_documento": ident.Documento,
			"observacion":       "Plazo de edición vencido",
		})
		return respErr(http.StatusConflict, "El plazo de edición venció")
	}

	// Actualización parcial: se reenvía únicamente lo enviado por el docente
	// (sin campos de identidad/flujo/versionado, que el CRUD protege).
	contenido := map[string]any{}
	if err := json.Unmarshal(body, &contenido); err != nil {
		return respErr(http.StatusBadRequest, "El contenido de la propuesta no es válido")
	}
	for _, k := range camposReservadosContenido {
		delete(contenido, k)
	}
	contenido["tercero_id"] = ident.TerceroId

	updated, err := utils.PutSyllabusContenido(ctx, utils.SyllabusString(draft, "_id"), contenido)
	if err != nil {
		return apiError(err)
	}
	return requestresponse.APIResponseDTO(true, http.StatusOK, updated, "Borrador guardado")
}

// ProponerSyllabus marca la propuesta como enviada a revisión.
func ProponerSyllabus(ctx context.Context, code string, body []byte) requestresponse.APIResponse {
	var req models.TransicionRequest
	_ = json.Unmarshal(body, &req)

	ident, err := utils.ResolverIdentidad(ctx, req.Documento)
	if err != nil {
		return apiError(err)
	}
	draft, err := utils.GetBorradorAbierto(ctx, code)
	if err != nil {
		return apiError(err)
	}
	if draft == nil {
		return respErr(http.StatusNotFound, "No hay una propuesta en curso para este syllabus")
	}
	if !utils.EsDocenteAsignado(draft, ident) {
		return respErr(http.StatusForbidden, "Solo el docente asignado puede proponer la propuesta")
	}
	estado := utils.SyllabusString(draft, "estado")
	if estado != models.EstadoAsignado && estado != models.EstadoRechazado {
		return respErr(http.StatusConflict, "La propuesta ya fue enviada o no admite envío en estado "+estado)
	}
	if utils.PlazoVencido(draft) {
		id := utils.SyllabusString(draft, "_id")
		_, _ = utils.PutEstado(ctx, id, map[string]any{
			"estado":            models.EstadoExpirado,
			"tercero_id":        ident.TerceroId,
			"tercero_documento": ident.Documento,
			"observacion":       "Plazo de edición vencido",
		})
		return respErr(http.StatusConflict, "El plazo de edición venció")
	}

	// Al proponer se registra la fecha de revisión de consejo (el revisor/aprobador
	// se registra al publicar).
	id := utils.SyllabusString(draft, "_id")
	nowISO := time.Now().UTC().Format(time.RFC3339)
	seguimiento := buildSeguimientoEnRevision(utils.SeguimientoMap(draft), nowISO)
	if _, err := utils.PutSyllabusContenido(ctx, id, map[string]any{"seguimiento": seguimiento}); err != nil {
		return apiError(err)
	}

	updated, err := utils.PutEstado(ctx, id, map[string]any{
		"estado":            models.EstadoEnRevision,
		"tercero_id":        ident.TerceroId,
		"tercero_documento": ident.Documento,
		"observacion":       req.Observacion,
	})
	if err != nil {
		return apiError(err)
	}
	return requestresponse.APIResponseDTO(true, http.StatusOK, updated, "Propuesta enviada a revisión")
}

// AprobarSyllabus publica la propuesta como nueva versión activa.
func AprobarSyllabus(ctx context.Context, code string, body []byte) requestresponse.APIResponse {
	var req models.TransicionRequest
	_ = json.Unmarshal(body, &req)

	ident, err := utils.ResolverIdentidad(ctx, req.Documento)
	if err != nil {
		return apiError(err)
	}
	draft, err := utils.GetBorradorAbierto(ctx, code)
	if err != nil {
		return apiError(err)
	}
	if draft == nil {
		return respErr(http.StatusNotFound, "No hay una propuesta en curso para este syllabus")
	}
	if utils.SyllabusString(draft, "estado") != models.EstadoEnRevision {
		return respErr(http.StatusConflict, "La propuesta no está en revisión")
	}
	ok, err := utils.EsCoordinadorDePrograma(ctx, ident.Documento,
		utils.SyllabusIntSlice(draft, "proyecto_curricular_ids"))
	if err != nil {
		return apiError(err)
	}
	if !ok {
		return respErr(http.StatusForbidden, "Solo el coordinador del programa puede aprobar la propuesta")
	}
	if req.NumeroActa == "" || req.Archivo == "" {
		return respErr(http.StatusBadRequest,
			"El número de acta y el documento del acta son obligatorios para publicar")
	}

	published, err := utils.PutPublicar(ctx, utils.SyllabusString(draft, "_id"), map[string]any{
		"tercero_id":        ident.TerceroId,
		"tercero_documento": ident.Documento,
		"nombre_aprobador":  ident.NombreCompleto,
		"numeroActa":        req.NumeroActa,
		"archivo":           req.Archivo,
		"observacion":       req.Observacion,
	})
	if err != nil {
		return apiError(err)
	}
	return requestresponse.APIResponseDTO(true, http.StatusOK, published, "Propuesta aprobada y publicada")
}

// RechazarSyllabus rechaza la propuesta y, opcionalmente, la reabre con nuevo plazo.
func RechazarSyllabus(ctx context.Context, code string, body []byte) requestresponse.APIResponse {
	var req models.TransicionRequest
	_ = json.Unmarshal(body, &req)
	if req.MotivoRechazo == "" {
		return respErr(http.StatusBadRequest, "El motivo del rechazo es obligatorio")
	}

	ident, err := utils.ResolverIdentidad(ctx, req.Documento)
	if err != nil {
		return apiError(err)
	}
	draft, err := utils.GetBorradorAbierto(ctx, code)
	if err != nil {
		return apiError(err)
	}
	if draft == nil {
		return respErr(http.StatusNotFound, "No hay una propuesta en curso para este syllabus")
	}
	if utils.SyllabusString(draft, "estado") != models.EstadoEnRevision {
		return respErr(http.StatusConflict, "Solo una propuesta en revisión puede rechazarse")
	}
	ok, err := utils.EsCoordinadorDePrograma(ctx, ident.Documento,
		utils.SyllabusIntSlice(draft, "proyecto_curricular_ids"))
	if err != nil {
		return apiError(err)
	}
	if !ok {
		return respErr(http.StatusForbidden, "Solo el coordinador del programa puede rechazar la propuesta")
	}

	id := utils.SyllabusString(draft, "_id")
	updated, err := utils.PutEstado(ctx, id, map[string]any{
		"estado":            models.EstadoRechazado,
		"motivo_rechazo":    req.MotivoRechazo,
		"tercero_id":        ident.TerceroId,
		"tercero_documento": ident.Documento,
		"observacion":       req.Observacion,
	})
	if err != nil {
		return apiError(err)
	}

	if req.FechaLimite != "" {
		reopened, err := utils.PutEstado(ctx, id, map[string]any{
			"estado":            models.EstadoAsignado,
			"fecha_limite":      req.FechaLimite,
			"tercero_id":        ident.TerceroId,
			"tercero_documento": ident.Documento,
			"observacion":       "Reapertura con nuevo plazo",
		})
		if err != nil {
			return apiError(err)
		}
		updated = reopened
	}
	return requestresponse.APIResponseDTO(true, http.StatusOK, updated, "Propuesta rechazada")
}

// ReabrirSyllabus reabre una propuesta rechazada o expirada con un nuevo plazo.
func ReabrirSyllabus(ctx context.Context, code string, body []byte) requestresponse.APIResponse {
	var req models.TransicionRequest
	_ = json.Unmarshal(body, &req)
	if req.FechaLimite == "" {
		return respErr(http.StatusBadRequest, "La nueva fecha límite es obligatoria")
	}

	ident, err := utils.ResolverIdentidad(ctx, req.Documento)
	if err != nil {
		return apiError(err)
	}
	draft, err := utils.GetBorradorAbierto(ctx, code)
	if err != nil {
		return apiError(err)
	}
	if draft == nil {
		return respErr(http.StatusNotFound, "No hay una propuesta en curso para este syllabus")
	}
	estado := utils.SyllabusString(draft, "estado")
	if estado != models.EstadoRechazado && estado != models.EstadoExpirado {
		return respErr(http.StatusConflict, "Solo una propuesta rechazada o expirada puede reabrirse")
	}
	ok, err := utils.EsCoordinadorDePrograma(ctx, ident.Documento,
		utils.SyllabusIntSlice(draft, "proyecto_curricular_ids"))
	if err != nil {
		return apiError(err)
	}
	if !ok {
		return respErr(http.StatusForbidden, "Solo el coordinador del programa puede reabrir la propuesta")
	}

	updated, err := utils.PutEstado(ctx, utils.SyllabusString(draft, "_id"), map[string]any{
		"estado":            models.EstadoAsignado,
		"fecha_limite":      req.FechaLimite,
		"tercero_id":        ident.TerceroId,
		"tercero_documento": ident.Documento,
		"observacion":       req.Observacion,
	})
	if err != nil {
		return apiError(err)
	}
	return requestresponse.APIResponseDTO(true, http.StatusOK, updated, "Propuesta reabierta")
}

// ReasignarSyllabus cambia el docente de una propuesta abierta conservando su
// contenido (mismo documento/versión), con plazo nuevo. Disponible para
// asignado/rechazado/expirado; una propuesta en revisión debe rechazarse antes.
func ReasignarSyllabus(ctx context.Context, code string, body []byte) requestresponse.APIResponse {
	var req models.ReasignarRequest
	_ = json.Unmarshal(body, &req)
	if req.DocumentoDocente == "" {
		return respErr(http.StatusBadRequest, "El documento del nuevo docente es obligatorio")
	}
	if req.FechaLimite == "" {
		return respErr(http.StatusBadRequest, "La nueva fecha límite es obligatoria")
	}

	ident, err := utils.ResolverIdentidad(ctx, req.Documento)
	if err != nil {
		return apiError(err)
	}
	draft, err := utils.GetBorradorAbierto(ctx, code)
	if err != nil {
		return apiError(err)
	}
	if draft == nil {
		return respErr(http.StatusNotFound, "No hay una propuesta en curso para este syllabus")
	}
	estado := utils.SyllabusString(draft, "estado")
	if estado != models.EstadoAsignado && estado != models.EstadoRechazado && estado != models.EstadoExpirado {
		return respErr(http.StatusConflict,
			"No se puede reasignar en estado "+estado+"; rechace la propuesta primero")
	}
	ok, err := utils.EsCoordinadorDePrograma(ctx, ident.Documento,
		utils.SyllabusIntSlice(draft, "proyecto_curricular_ids"))
	if err != nil {
		return apiError(err)
	}
	if !ok {
		return respErr(http.StatusForbidden, "Solo el coordinador del programa puede reasignar la propuesta")
	}

	flujo := utils.FlujoMap(draft)
	if utils.FlujoString(flujo, "asignado_a_documento") == req.DocumentoDocente {
		return respErr(http.StatusBadRequest, "El docente indicado ya tiene asignada esta propuesta")
	}

	nuevoDocente, err := utils.GetTerceroByDocumento(ctx, req.DocumentoDocente)
	if err != nil {
		return apiError(err)
	}
	if nuevoDocente == nil {
		return respErr(http.StatusNotFound,
			fmt.Sprintf("No se encontró un tercero activo con documento %s", req.DocumentoDocente))
	}

	observacion := req.Observacion
	if observacion == "" {
		observacion = "Reasignación de docente"
	}
	updated, err := utils.PutEstado(ctx, utils.SyllabusString(draft, "_id"), map[string]any{
		"estado":               models.EstadoAsignado,
		"asignado_a":           nuevoDocente.Id,
		"asignado_a_documento": req.DocumentoDocente,
		"fecha_limite":         req.FechaLimite,
		"tercero_id":           ident.TerceroId,
		"tercero_documento":    ident.Documento,
		"observacion":          observacion,
	})
	if err != nil {
		return apiError(err)
	}
	return requestresponse.APIResponseDTO(true, http.StatusOK, updated, "Propuesta reasignada")
}

// CancelarSyllabus descarta una propuesta en curso.
func CancelarSyllabus(ctx context.Context, code string, body []byte) requestresponse.APIResponse {
	var req models.TransicionRequest
	_ = json.Unmarshal(body, &req)

	ident, err := utils.ResolverIdentidad(ctx, req.Documento)
	if err != nil {
		return apiError(err)
	}
	draft, err := utils.GetBorradorAbierto(ctx, code)
	if err != nil {
		return apiError(err)
	}
	if draft == nil {
		return respErr(http.StatusNotFound, "No hay una propuesta en curso para este syllabus")
	}
	ok, err := utils.EsCoordinadorDePrograma(ctx, ident.Documento,
		utils.SyllabusIntSlice(draft, "proyecto_curricular_ids"))
	if err != nil {
		return apiError(err)
	}
	if !ok {
		return respErr(http.StatusForbidden, "Solo el coordinador del programa puede cancelar la propuesta")
	}

	updated, err := utils.PutEstado(ctx, utils.SyllabusString(draft, "_id"), map[string]any{
		"estado":            models.EstadoCancelado,
		"tercero_id":        ident.TerceroId,
		"tercero_documento": ident.Documento,
		"observacion":       req.Observacion,
	})
	if err != nil {
		return apiError(err)
	}
	return requestresponse.APIResponseDTO(true, http.StatusOK, updated, "Propuesta cancelada")
}

// GetBorrador devuelve la propuesta en curso al docente asignado o al coordinador.
func GetBorrador(ctx context.Context, code string, documento string) requestresponse.APIResponse {
	ident, err := utils.ResolverIdentidad(ctx, documento)
	if err != nil {
		return apiError(err)
	}
	draft, err := utils.GetBorradorAbierto(ctx, code)
	if err != nil {
		return apiError(err)
	}
	if draft == nil {
		return respErr(http.StatusNotFound, "No hay una propuesta en curso para este syllabus")
	}
	if autorizado, err := autorizadoSobreDraft(ctx, draft, ident); err != nil {
		return apiError(err)
	} else if !autorizado {
		return respErr(http.StatusForbidden, "No tiene acceso a esta propuesta")
	}
	return requestresponse.APIResponseDTO(true, http.StatusOK, draft, "Propuesta obtenida")
}

// GetEstado devuelve el resumen del flujo de un syllabus.
func GetEstado(ctx context.Context, code string, documento string) requestresponse.APIResponse {
	ident, err := utils.ResolverIdentidad(ctx, documento)
	if err != nil {
		return apiError(err)
	}
	draft, err := utils.GetBorradorAbierto(ctx, code)
	if err != nil {
		return apiError(err)
	}
	if draft == nil {
		activo, err := utils.GetSyllabusActivo(ctx, code)
		if err != nil {
			return apiError(err)
		}
		if activo == nil {
			return respErr(http.StatusNotFound, "Syllabus no encontrado")
		}
		return requestresponse.APIResponseDTO(true, http.StatusOK, map[string]any{
			"estado": models.EstadoActivo,
		}, "Syllabus activo")
	}
	if autorizado, err := autorizadoSobreDraft(ctx, draft, ident); err != nil {
		return apiError(err)
	} else if !autorizado {
		return respErr(http.StatusForbidden, "No tiene acceso al estado de esta propuesta")
	}
	version, _ := utils.SyllabusInt(draft, "version")
	return requestresponse.APIResponseDTO(true, http.StatusOK, map[string]any{
		"estado":  utils.SyllabusString(draft, "estado"),
		"flujo":   utils.FlujoMap(draft),
		"version": version,
		"_id":     utils.SyllabusString(draft, "_id"),
	}, "Estado del flujo obtenido")
}

// PendientesSyllabus lista las propuestas abiertas (en revisión y en curso:
// asignado/rechazado/expirado) de los programas del coordinador. La
// autorización se recalcula por programa vigente: al cambiar el coordinador,
// su reemplazo hereda la gestión de las asignaciones (incluidas las del
// coordinador anterior), y quien deja de serlo pierde el acceso.
func PendientesSyllabus(ctx context.Context, documento string) requestresponse.APIResponse {
	ident, err := utils.ResolverIdentidad(ctx, documento)
	if err != nil {
		return apiError(err)
	}
	permitidos, err := utils.GetCodigosProgramaDeCoordinador(ctx, ident.Documento)
	if err != nil {
		return apiError(err)
	}
	drafts, err := utils.GetBorradoresAbiertos(ctx)
	if err != nil {
		return apiError(err)
	}
	resultado := make([]map[string]any, 0)
	for _, draft := range drafts {
		for _, proyectoId := range utils.SyllabusIntSlice(draft, "proyecto_curricular_ids") {
			if permitidos[proyectoId] {
				resultado = append(resultado, draft)
				break
			}
		}
	}
	return requestresponse.APIResponseDTO(true, http.StatusOK, resultado, "Propuestas abiertas del coordinador")
}

// MisAsignacionesSyllabus lista las propuestas asignadas a un docente.
func MisAsignacionesSyllabus(ctx context.Context, documento string) requestresponse.APIResponse {
	ident, err := utils.ResolverIdentidad(ctx, documento)
	if err != nil {
		return apiError(err)
	}
	drafts, err := utils.GetBorradoresPorDocente(ctx, ident.Documento, ident.TerceroId)
	if err != nil {
		return apiError(err)
	}
	return requestresponse.APIResponseDTO(true, http.StatusOK, drafts, "Asignaciones del docente")
}

// BuscarDocente resuelve un tercero a partir del número de identificación.
func BuscarDocente(ctx context.Context, documento string) requestresponse.APIResponse {
	documento = strings.TrimSpace(documento)
	if documento == "" {
		return respErr(http.StatusBadRequest, "El documento es obligatorio")
	}
	tercero, err := utils.GetTerceroByDocumento(ctx, documento)
	if err != nil {
		return apiError(err)
	}
	if tercero == nil {
		return respErr(http.StatusNotFound,
			fmt.Sprintf("No se encontró un tercero activo con documento %s", documento))
	}
	docente, err := utils.GetDocenteActivoByDocumento(ctx, documento)
	if err != nil {
		return apiError(err)
	}
	if docente == nil {
		return respErr(http.StatusNotFound,
			fmt.Sprintf("El tercero con documento %s existe, pero no está registrado como docente activo", documento))
	}
	return requestresponse.APIResponseDTO(true, http.StatusOK, models.DocenteBusqueda{
		Id:             tercero.Id,
		NombreCompleto: tercero.NombreCompleto,
		UsuarioWSO2:    tercero.UsuarioWSO2,
		Cedula:         docente.Cedula,
		Nombre:         docente.Nombre,
	}, "Docente encontrado")
}

// autorizadoSobreDraft indica si la identidad es el docente asignado o un
// coordinador del programa del syllabus.
func autorizadoSobreDraft(ctx context.Context, draft map[string]any, ident models.Identidad) (bool, error) {
	if utils.EsDocenteAsignado(draft, ident) {
		return true, nil
	}
	ok, err := utils.EsCoordinadorDePrograma(ctx, ident.Documento,
		utils.SyllabusIntSlice(draft, "proyecto_curricular_ids"))
	if err != nil {
		logs.Warn("error validando coordinador: %v", err)
		return false, err
	}
	return ok, nil
}
