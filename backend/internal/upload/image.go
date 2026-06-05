package upload

import (
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

var allowedImageMIME = map[string]struct{}{
	"image/jpeg": {},
	"image/png":  {},
	"image/webp": {},
}

// ResolveImageMIME picks a trusted image MIME from magic bytes, client hint, and file name.
func ResolveImageMIME(head []byte, hint, fileName string) string {
	if len(head) > 0 {
		detected := strings.TrimSpace(strings.ToLower(http.DetectContentType(head)))
		if i := strings.Index(detected, ";"); i >= 0 {
			detected = strings.TrimSpace(detected[:i])
		}
		if isAllowedImageMIME(detected) {
			return detected
		}
	}
	candidate := strings.TrimSpace(strings.ToLower(hint))
	if i := strings.Index(candidate, ";"); i >= 0 {
		candidate = strings.TrimSpace(candidate[:i])
	}
	if isAllowedImageMIME(candidate) {
		return candidate
	}
	if ext := strings.ToLower(filepath.Ext(fileName)); ext != "" {
		byExt := mime.TypeByExtension(ext)
		byExt = strings.TrimSpace(strings.ToLower(strings.Split(byExt, ";")[0]))
		if isAllowedImageMIME(byExt) {
			return byExt
		}
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
