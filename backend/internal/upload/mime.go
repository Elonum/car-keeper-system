package upload

import (
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

var allowedDocumentMIME = map[string]struct{}{
	"application/pdf":    {},
	"image/jpeg":         {},
	"image/png":          {},
	"image/webp":         {},
	"application/msword": {},
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": {},
	"application/vnd.ms-excel": {},
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": {},
}

// ResolveDocumentMIME picks a trusted MIME from magic bytes, client hint, and file name.
func ResolveDocumentMIME(head []byte, hint, fileName string) string {
	if len(head) > 0 {
		detected := strings.TrimSpace(strings.ToLower(http.DetectContentType(head)))
		if i := strings.Index(detected, ";"); i >= 0 {
			detected = strings.TrimSpace(detected[:i])
		}
		if isAllowedDocumentMIME(detected) {
			return detected
		}
	}
	candidate := strings.TrimSpace(strings.ToLower(hint))
	if i := strings.Index(candidate, ";"); i >= 0 {
		candidate = strings.TrimSpace(candidate[:i])
	}
	if isAllowedDocumentMIME(candidate) {
		return candidate
	}
	if ext := strings.ToLower(filepath.Ext(fileName)); ext != "" {
		byExt := mime.TypeByExtension(ext)
		byExt = strings.TrimSpace(strings.ToLower(strings.Split(byExt, ";")[0]))
		if isAllowedDocumentMIME(byExt) {
			return byExt
		}
	}
	return ""
}

func isAllowedDocumentMIME(m string) bool {
	if m == "" {
		return false
	}
	_, ok := allowedDocumentMIME[m]
	return ok
}

// AllowedDocumentExtensions lists client-facing file extensions (lowercase, with dot).
func AllowedDocumentExtensions() []string {
	return []string{
		".pdf",
		".jpg", ".jpeg", ".png", ".webp",
		".doc", ".docx",
		".xls", ".xlsx",
	}
}
