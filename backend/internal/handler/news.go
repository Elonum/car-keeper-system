package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/carkeeper/backend/internal/authz"
	"github.com/carkeeper/backend/internal/middleware"
	"github.com/carkeeper/backend/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) GetNews(w http.ResponseWriter, r *http.Request) {
	role, _ := middleware.GetUserRole(r.Context())
	scope := r.URL.Query().Get("scope")

	var filter *bool
	if authz.IsStaff(role) {
		switch scope {
		case "all":
			filter = nil
		case "unpublished":
			f := false
			filter = &f
		default:
			t := true
			filter = &t
		}
	} else {
		if scope != "" {
			Forbidden(w, "Invalid query parameter: scope is restricted to staff")
			return
		}
		t := true
		filter = &t
	}

	news, err := h.services.News.GetNews(r.Context(), filter)
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

	item, err := h.services.News.GetNewsByID(r.Context(), newsID)
	if err != nil {
		NotFound(w, err.Error())
		return
	}

	role, _ := middleware.GetUserRole(r.Context())
	if !item.IsPublished && !authz.IsStaff(role) {
		NotFound(w, "News not found")
		return
	}

	Success(w, item)
}

func (h *Handler) CreateNews(w http.ResponseWriter, r *http.Request) {
	userID, ok := RequireStaff(w, r)
	if !ok {
		return
	}

	var create model.NewsCreate
	if err := json.NewDecoder(r.Body).Decode(&create); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}

	create.Title = strings.TrimSpace(create.Title)
	create.Content = strings.TrimSpace(create.Content)
	if create.Title == "" || utf8.RuneCountInString(create.Title) > 255 {
		BadRequest(w, "title is required (max 255 characters)")
		return
	}
	if create.Content == "" {
		BadRequest(w, "content is required")
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
	if _, ok := RequireStaff(w, r); !ok {
		return
	}

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
	if _, ok := RequireStaff(w, r); !ok {
		return
	}

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
	if _, ok := RequireStaff(w, r); !ok {
		return
	}

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
	if _, ok := RequireStaff(w, r); !ok {
		return
	}

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
