package validate

import (
	"strings"
)

const (
	NewsTitleMax   = 255
	NewsContentMax = 50_000
)

// NewsFields validates title and content for create/update.
func NewsFields(title, content string) (t, c string, errMsg string) {
	t = strings.TrimSpace(title)
	c = strings.TrimSpace(content)
	if t == "" || runeLen(t) > NewsTitleMax {
		return "", "", "title is required (max 255 characters)"
	}
	if c == "" {
		return "", "", "content is required"
	}
	if runeLen(c) > NewsContentMax {
		return "", "", "content is too long"
	}
	if hasDisallowedControlRunes(t, false) {
		return "", "", "title contains invalid characters"
	}
	if hasDisallowedControlRunes(c, true) {
		return "", "", "content contains invalid characters"
	}
	return t, c, ""
}
