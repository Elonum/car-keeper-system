package upload

import (
	"net/http"
	"strings"
)

var allowedImageMIME = map[string]struct{}{
	"image/jpeg": {},
	"image/png":  {},
	"image/webp": {},
}

// ResolveImageMIME picks a trusted image MIME from file magic bytes only.
func ResolveImageMIME(head []byte, hint, fileName string) string {
	_ = hint
	_ = fileName
	if len(head) == 0 {
		return ""
	}
	detected := strings.TrimSpace(strings.ToLower(http.DetectContentType(head)))
	if i := strings.Index(detected, ";"); i >= 0 {
		detected = strings.TrimSpace(detected[:i])
	}
	if isAllowedImageMIME(detected) {
		return detected
	}
	return ""
}

func isAllowedImageMIME(m string) bool {
	if m == "" {
		return false
	}
	_, ok := allowedImageMIME[m]
	return ok
}

// AllowedModelImageExtensions lists catalog model image extensions.
func AllowedModelImageExtensions() []string {
	return []string{".jpg", ".jpeg", ".png", ".webp"}
}
