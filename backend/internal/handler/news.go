package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/carkeeper/backend/internal/middleware"
	"github.com/carkeeper/backend/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) GetNews(w http.ResponseWriter, r *http.Request) {
	var isPublished *bool
	if isPublishedStr := r.URL.Query().Get("is_published"); isPublishedStr != "" {
		if val, err := strconv.ParseBool(isPublishedStr); err == nil {
			isPublished = &val
		}
	}

	news, err := h.services.News.GetNews(r.Context(), isPublished)
	if err != nil {
		InternalServerError(w, err.Error())
		return
	}
	Success(w, news)
}

func (h *Handler) GetNewsByID(w http.ResponseWriter, r *http.Request) {
	newsIDStr := chi.URLParam(r, "id")
	newsID, err := uuid.Parse(newsIDStr)
	if err != nil {
		BadRequest(w, "Invalid news ID")
		return
	}

	news, err := h.services.News.GetNewsByID(r.Context(), newsID)
	if err != nil {
		NotFound(w, err.Error())
		return
	}
	Success(w, news)
}

func (h *Handler) CreateNews(w http.ResponseWriter, r *http.Request) {
	userIDVal, ok := middleware.GetUserIDUUID(r.Context())
	if !ok {
		Unauthorized(w, "User not authenticated")
		return
	}

	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		userIDStr, ok := userIDVal.(string)
		if !ok {
			Unauthorized(w, "Invalid user ID")
			return
		}
		var err error
		userID, err = uuid.Parse(userIDStr)
		if err != nil {
			Unauthorized(w, "Invalid user ID format")
			return
		}

		// Check if user is admin or manager
		role, _ := middleware.GetUserRole(r.Context())
		if role != "admin" && role != "manager" {
			Unauthorized(w, "Only admins and managers can create news")
			return
		}
	}

	var create model.NewsCreate
	if err := json.NewDecoder(r.Body).Decode(&create); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}

	news, err := h.services.News.CreateNews(r.Context(), userID, create)
	if err != nil {
		BadRequest(w, err.Error())
		return
	}

	Success(w, news)
}

func (h *Handler) UpdateNews(w http.ResponseWriter, r *http.Request) {
	newsIDStr := chi.URLParam(r, "id")
	newsID, err := uuid.Parse(newsIDStr)
	if err != nil {
		BadRequest(w, "Invalid news ID")
		return
	}

	var update struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}

	if err := h.services.News.UpdateNews(r.Context(), newsID, update.Title, update.Content); err != nil {
		BadRequest(w, err.Error())
		return
	}

	news, err := h.services.News.GetNewsByID(r.Context(), newsID)
	if err != nil {
		InternalServerError(w, err.Error())
		return
	}
	Success(w, news)
}

func (h *Handler) PublishNews(w http.ResponseWriter, r *http.Request) {
	newsIDStr := chi.URLParam(r, "id")
	newsID, err := uuid.Parse(newsIDStr)
	if err != nil {
		BadRequest(w, "Invalid news ID")
		return
	}

	if err := h.services.News.PublishNews(r.Context(), newsID); err != nil {
		BadRequest(w, err.Error())
		return
	}

	news, err := h.services.News.GetNewsByID(r.Context(), newsID)
	if err != nil {
		InternalServerError(w, err.Error())
		return
	}
	Success(w, news)
}

func (h *Handler) UnpublishNews(w http.ResponseWriter, r *http.Request) {
	newsIDStr := chi.URLParam(r, "id")
	newsID, err := uuid.Parse(newsIDStr)
	if err != nil {
		BadRequest(w, "Invalid news ID")
		return
	}

	if err := h.services.News.UnpublishNews(r.Context(), newsID); err != nil {
		BadRequest(w, err.Error())
		return
	}

	news, err := h.services.News.GetNewsByID(r.Context(), newsID)
	if err != nil {
		InternalServerError(w, err.Error())
		return
	}
	Success(w, news)
}

func (h *Handler) DeleteNews(w http.ResponseWriter, r *http.Request) {
	newsIDStr := chi.URLParam(r, "id")
	newsID, err := uuid.Parse(newsIDStr)
	if err != nil {
		BadRequest(w, "Invalid news ID")
		return
	}

	if err := h.services.News.DeleteNews(r.Context(), newsID); err != nil {
		InternalServerError(w, err.Error())
		return
	}

	Success(w, map[string]string{"message": "News deleted"})
}

