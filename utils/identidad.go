package utils

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/astaxie/beego"
	"github.com/udistrital/syllabus_mid/models"
	"github.com/udistrital/utils_oas/request"
)

const authorizationKey = "Authorization"

// userInfoURL es el endpoint de WSO2 con el que se resuelve la identidad a partir
// del token Bearer (el gateway ya validó la firma).
const userInfoURL = "https://autenticacion.portaloas.udistrital.edu.co/oauth2/userinfo"

// GetUsuarioActual obtiene el usuario asociado al token de la petición.
// Retorna (nil, nil) si no hay token (por ejemplo, en local).
func GetUsuarioActual(ctx context.Context) (*models.Usuario, error) {
	token, _ := ctx.Value(authorizationKey).(string)
	if token == "" {
		return nil, nil
	}
	var usuario models.Usuario
	status, err := request.GetWithContext(ctx, userInfoURL, &usuario)
	if err != nil {
		return nil, &QueryError{
			Status:  status,
			Message: fmt.Sprintf("No se pudo resolver la identidad del usuario: %v", err),
		}
	}
	return &usuario, nil
}

// GetTerceroByDocumento resuelve el tercero a partir del número de identificación.
func GetTerceroByDocumento(ctx context.Context, documento string) (*models.Tercero, error) {
	if documento == "" {
		return nil, &QueryError{Status: http.StatusBadRequest, Message: "El documento es obligatorio"}
	}
	endpoint := beego.AppConfig.String("TercerosService") +
		"datos_identificacion?query=Activo:true,Numero:" + url.QueryEscape(documento) +
		"&sortby=FechaCreacion&order=desc"

	var datos []models.DatosIdentificacion
	status, err := request.GetWithContext(ctx, endpoint, &datos)
	if err != nil {
		return nil, &QueryError{
			Status:  status,
			Message: fmt.Sprintf("Error consultando terceros para el documento %s: %v", documento, err),
		}
	}
	if len(datos) == 0 {
		return nil, nil
	}
	tercero := datos[0].TerceroId
	return &tercero, nil
}

// ResolverIdentidad determina la identidad del actor: prioriza el token y usa
// documentoBody como respaldo (entornos locales sin token).
func ResolverIdentidad(ctx context.Context, documentoBody string) (models.Identidad, error) {
	documento := ""
	if usuario, err := GetUsuarioActual(ctx); err != nil {
		return models.Identidad{}, err
	} else if usuario != nil && usuario.Documento != "" {
		documento = usuario.Documento
	}
	if documento == "" {
		documento = documentoBody
	}
	if documento == "" {
		return models.Identidad{}, &QueryError{
			Status:  http.StatusUnauthorized,
			Message: "No se pudo determinar la identidad del usuario (token o documento)",
		}
	}

	tercero, err := GetTerceroByDocumento(ctx, documento)
	if err != nil {
		return models.Identidad{}, err
	}
	if tercero == nil {
		return models.Identidad{}, &QueryError{
			Status:  http.StatusNotFound,
			Message: fmt.Sprintf("No se encontró un tercero activo con documento %s", documento),
		}
	}
	return models.Identidad{
		TerceroId:      tercero.Id,
		Documento:      documento,
		NombreCompleto: tercero.NombreCompleto,
	}, nil
}
