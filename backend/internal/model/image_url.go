package model

import "github.com/google/uuid"

// BuildModelImageURL returns a cache-busted public path for a model image.
// The v query param changes whenever image_key changes so browsers fetch the new file.
func BuildModelImageURL(modelID uuid.UUID, imageKey string) string {
	if imageKey == "" {
		return ""
	}
	return "/api/catalog/models/" + modelID.String() + "/image?v=" + imageKey
}
