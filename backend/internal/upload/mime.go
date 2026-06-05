package upload

import (
	"net/http"
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

// ResolveDocumentMIME picks a trusted MIME from file magic bytes only.
func ResolveDocumentMIME(head []byte, hint, fileName string) string {
	_ = hint
	_ = fileName
	if len(head) == 0 {
		return ""
	}
	detected := strings.TrimSpace(strings.ToLower(http.DetectContentType(head)))
	if i := strings.Index(detected, ";"); i >= 0 {
		detected = strings.TrimSpace(detected[:i])
	}
	if isAllowedDocumentMIME(detected) {
		return detected
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
