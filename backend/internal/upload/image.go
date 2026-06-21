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
	if mimeType := sniffImageMIME(head); mimeType != "" {
		return mimeType
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

func sniffImageMIME(head []byte) string {
	if len(head) >= 12 && string(head[0:4]) == "RIFF" && string(head[8:12]) == "WEBP" {
		return "image/webp"
	}
	if len(head) >= 3 && head[0] == 0xFF && head[1] == 0xD8 && head[2] == 0xFF {
		return "image/jpeg"
	}
	if len(head) >= 8 && head[0] == 0x89 && head[1] == 'P' && head[2] == 'N' && head[3] == 'G' {
		return "image/png"
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
