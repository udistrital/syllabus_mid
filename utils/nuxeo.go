package utils

import (
	"context"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
)

const maxNuxeoImageBytes = 2 * 1024 * 1024

var nuxeoUUIDPattern = regexp.MustCompile(`^[a-zA-Z0-9-]{8,64}$`)

type nuxeoDocumentResponse struct {
	File    string `json:"file"`
	Content struct {
		MimeType string `json:"mime-type"`
	} `json:"file:content"`
}

// GetNuxeoImageDataURI consulta NUXEO por el documento y devuelve su contenido
// como data URI. El tipo real se detecta por los bytes, porque NUXEO reporta el
// mime segun el tipo de documento y no siempre coincide con el archivo.
func GetNuxeoImageDataURI(ctx context.Context, uuid string) (string, error) {
	uuid = strings.TrimSpace(uuid)
	if !nuxeoUUIDPattern.MatchString(uuid) {
		return "", fmt.Errorf("identificador de imagen invalido")
	}

	baseURL := strings.TrimRight(beego.AppConfig.String("NuxeoService"), "/") + "/"
	var response nuxeoDocumentResponse
	if _, err := request.GetWithContext(ctx, baseURL+"document/"+uuid, &response); err != nil {
		return "", fmt.Errorf("consulta a NUXEO: %w", err)
	}

	file := strings.ReplaceAll(strings.ReplaceAll(response.File, "\n", ""), "\r", "")
	file = strings.TrimSpace(file)
	if file == "" {
		return "", fmt.Errorf("NUXEO no devolvio el contenido de la imagen")
	}

	mime := detectImageMime(file)
	if mime == "" {
		mime = normalizeReportedMime(response.Content.MimeType)
	}
	if mime == "" {
		return "", fmt.Errorf("tipo de imagen no permitido: %q", response.Content.MimeType)
	}

	if base64.StdEncoding.DecodedLen(len(file)) > maxNuxeoImageBytes {
		return "", fmt.Errorf("la imagen supera el limite de 2 MB")
	}

	return "data:" + mime + ";base64," + file, nil
}

func detectImageMime(file string) string {
	switch {
	case strings.HasPrefix(file, "iVBORw0KGgo"):
		return "image/png"
	case strings.HasPrefix(file, "/9j/"):
		return "image/jpeg"
	case strings.HasPrefix(file, "UklGR"):
		return "image/webp"
	default:
		return ""
	}
}

func normalizeReportedMime(value string) string {
	mime := strings.TrimSpace(strings.Split(value, ";")[0])
	if mime == "image/jpg" {
		mime = "image/jpeg"
	}
	switch mime {
	case "image/png", "image/jpeg", "image/webp":
		return mime
	default:
		return ""
	}
}
