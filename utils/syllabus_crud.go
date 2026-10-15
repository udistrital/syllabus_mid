package utils

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/syllabus_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// crudResponse replica el sobre de respuesta del CRUD de syllabus.
type crudResponse struct {
	Success bool `json:"Success"`
	Status  int  `json:"Status"`
	Message any  `json:"Message"`
	Data    any  `json:"Data"`
}

func crudBase() string {
	return beego.AppConfig.String("SyllabusService")
}

func doGetSyllabusList(ctx context.Context, query string, limit int) ([]map[string]any, error) {
	endpoint := crudBase() + "syllabus?limit=" + strconv.Itoa(limit) + "&offset=0&query=" +
		url.QueryEscape(query)

	var resp crudResponse
	status, err := request.GetWithContext(ctx, endpoint, &resp)
	if err != nil {
		return nil, &QueryError{
			Status:  status,
			Message: fmt.Sprintf("Error consultando syllabus: %v", err),
		}
	}
	list, _ := resp.Data.([]any)
	out := make([]map[string]any, 0, len(list))
	for _, item := range list {
		if m, ok := item.(map[string]any); ok {
			out = append(out, m)
		}
	}
	return out, nil
}

// GetSyllabusActivo obtiene la versión publicada de un syllabus_code.
func GetSyllabusActivo(ctx context.Context, code string) (map[string]any, error) {
	list, err := doGetSyllabusList(ctx, fmt.Sprintf("syllabus_code:%s,syllabus_actual:true", code), 1)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return list[0], nil
}

// GetBorradorAbierto obtiene la propuesta en curso de un syllabus_code (si existe).
func GetBorradorAbierto(ctx context.Context, code string) (map[string]any, error) {
	query := fmt.Sprintf("syllabus_code:%s,estado:%s", code, strings.Join(models.EstadosBorrador, "|"))
	list, err := doGetSyllabusList(ctx, query, 1)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return list[0], nil
}

// GetBorradoresPorEstado devuelve todas las propuestas en un estado.
func GetBorradoresPorEstado(ctx context.Context, estado string) ([]map[string]any, error) {
	return doGetSyllabusList(ctx, fmt.Sprintf("estado:%s", estado), 200)
}

// GetBorradoresAbiertos devuelve todas las propuestas no publicadas (en
// cualquier estado borrador: asignado, en_revision, rechazado, expirado).
func GetBorradoresAbiertos(ctx context.Context) ([]map[string]any, error) {
	query := fmt.Sprintf("estado:%s", strings.Join(models.EstadosBorrador, "|"))
	return doGetSyllabusList(ctx, query, 200)
}

// GetBorradoresPorDocente devuelve las asignaciones de un docente por documento
// o por tercero_id.
func GetBorradoresPorDocente(ctx context.Context, documento string, terceroId int) ([]map[string]any, error) {
	estados := strings.Join(models.EstadosBorrador, "|")
	resultado := make([]map[string]any, 0)
	vistos := map[string]bool{}

	agregar := func(items []map[string]any) {
		for _, item := range items {
			id := asString(item["_id"])
			if id != "" && vistos[id] {
				continue
			}
			if id != "" {
				vistos[id] = true
			}
			resultado = append(resultado, item)
		}
	}

	if documento != "" {
		items, err := doGetSyllabusList(ctx,
			fmt.Sprintf("flujo.asignado_a_documento:%s,estado:%s", documento, estados), 200)
		if err != nil {
			return nil, err
		}
		agregar(items)
	}
	if terceroId > 0 {
		items, err := doGetSyllabusList(ctx,
			fmt.Sprintf("flujo.asignado_a:%d,estado:%s", terceroId, estados), 200)
		if err != nil {
			return nil, err
		}
		agregar(items)
	}
	return resultado, nil
}

// PostSyllabus crea una nueva versión (borrador o publicación) del syllabus.
func PostSyllabus(ctx context.Context, payload any) (map[string]any, error) {
	var resp crudResponse
	status, err := request.PostWithContext(ctx, crudBase()+"syllabus", payload, &resp)
	if err != nil {
		return nil, &QueryError{Status: status, Message: fmt.Sprintf("Error creando el borrador del syllabus: %v", err)}
	}
	created, _ := resp.Data.(map[string]any)
	return created, nil
}

// PutSyllabusContenido actualiza el contenido de un borrador por su _id.
func PutSyllabusContenido(ctx context.Context, id string, payload any) (map[string]any, error) {
	var resp crudResponse
	status, err := request.PutWithContext(ctx, crudBase()+"syllabus/"+id, payload, &resp)
	if err != nil {
		return nil, &QueryError{Status: status, Message: fmt.Sprintf("Error guardando el contenido del syllabus: %v", err)}
	}
	updated, _ := resp.Data.(map[string]any)
	return updated, nil
}

// PutEstado aplica una transición de estado validada por el CRUD.
func PutEstado(ctx context.Context, id string, payload any) (map[string]any, error) {
	var resp crudResponse
	status, err := request.PutWithContext(ctx, crudBase()+"syllabus/"+id+"/estado", payload, &resp)
	if err != nil {
		return nil, &QueryError{Status: status, Message: fmt.Sprintf("Error aplicando la transición de estado: %v", err)}
	}
	updated, _ := resp.Data.(map[string]any)
	return updated, nil
}

// PutPublicar aprueba y publica la propuesta como versión actual.
func PutPublicar(ctx context.Context, id string, payload any) (map[string]any, error) {
	var resp crudResponse
	status, err := request.PutWithContext(ctx, crudBase()+"syllabus/"+id+"/publicar", payload, &resp)
	if err != nil {
		return nil, &QueryError{Status: status, Message: fmt.Sprintf("Error publicando la propuesta: %v", err)}
	}
	published, _ := resp.Data.(map[string]any)
	return published, nil
}

// --- Lectores de mapas (JSON dinámico del CRUD) ---

func toInt(v any) (int, bool) {
	switch n := v.(type) {
	case float64:
		return int(n), true
	case int:
		return n, true
	case int64:
		return int(n), true
	case string:
		if i, err := strconv.Atoi(n); err == nil {
			return i, true
		}
	}
	return 0, false
}

func SyllabusString(s map[string]any, key string) string {
	return asString(s[key])
}

func SyllabusInt(s map[string]any, key string) (int, bool) {
	return toInt(s[key])
}

func SyllabusIntSlice(s map[string]any, key string) []int {
	out := make([]int, 0)
	for _, v := range asSlice(s[key]) {
		if i, ok := toInt(v); ok {
			out = append(out, i)
		}
	}
	logs.Info("listado de id de programas en syllabus")
	fmt.Println(out)
	fmt.Println("")
	return out
}

func FlujoMap(s map[string]any) map[string]any {
	return asMap(s["flujo"])
}

func SeguimientoMap(s map[string]any) map[string]any {
	return asMap(s["seguimiento"])
}

func FlujoInt(flujo map[string]any, key string) (int, bool) {
	return toInt(flujo[key])
}

func FlujoString(flujo map[string]any, key string) string {
	return asString(flujo[key])
}

func BadRequest(msg string) *QueryError {
	return &QueryError{Status: http.StatusBadRequest, Message: msg}
}
