package handler

import (
	"errors"
	"io"
	"mime"
	"net/http"
	"strings"

	"github.com/carkeeper/backend/internal/apperr"
	"github.com/carkeeper/backend/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

const multipartMemoryLimit = 8 << 20

func (h *Handler) CreateDocument(w http.ResponseWriter, r *http.Request) {
	requester, role, ok := RequesterAndRole(w, r)
	if !ok {
		return
	}

	maxB := h.cfg.Storage.MaxUploadBytes
	if maxB <= 0 {
		maxB = 15 << 20
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxB+multipartMemoryLimit)

	if err := r.ParseMultipartForm(multipartMemoryLimit); err != nil {
		BadRequest(w, "invalid multipart form")
		return
	}
	defer func() { _ = r.MultipartForm.RemoveAll() }()

	file, header, err := r.FormFile("file")
	if err != nil {
		BadRequest(w, "file field is required")
		return
	}
	defer file.Close()

	docType := strings.TrimSpace(r.FormValue("document_type"))
	orderID, errOrder := parseOptionalUUID(r.FormValue("order_id"))
	apptID, errAppt := parseOptionalUUID(r.FormValue("service_appointment_id"))
	if errOrder != nil || errAppt != nil {
		BadRequest(w, "invalid order_id or service_appointment_id")
		return
	}

	size := header.Size
	if size <= 0 {
		BadRequest(w, "could not determine file size; upload a file with known length")
		return
	}
	if size > maxB {
		BadRequest(w, "file too large")
		return
	}

	mt := header.Header.Get("Content-Type")
	if i := strings.Index(mt, ";"); i >= 0 {
		mt = strings.TrimSpace(mt[:i])
	}

	doc, err := h.services.Document.Create(r.Context(), service.DocumentCreateInput{
		Reader:       file,
		Size:         size,
		FileName:     header.Filename,
		MimeType:     mt,
		DocumentType: docType,
		OrderID:      orderID,
		ApptID:       apptID,
		Requester:    requester,
		Role:         role,
	})
	if err != nil {
		switch {
		case errors.Is(err, apperr.ErrForbidden):
			Forbidden(w, "not allowed to attach document to this resource")
		default:
			BadRequest(w, err.Error())
		}
		return
	}

	Success(w, doc)
}

func parseOptionalUUID(s string) (*uuid.UUID, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	id, err := uuid.Parse(s)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (h *Handler) ListDocuments(w http.ResponseWriter, r *http.Request) {
	requester, role, ok := RequesterAndRole(w, r)
	if !ok {
		return
	}
	orderID, errO := parseOptionalUUID(r.URL.Query().Get("order_id"))
	apptID, errA := parseOptionalUUID(r.URL.Query().Get("service_appointment_id"))
	if errO != nil || errA != nil {
		BadRequest(w, "invalid query id")
		return
	}

	list, err := h.services.Document.List(r.Context(), requester, role, orderID, apptID)
	if err != nil {
		if errors.Is(err, apperr.ErrForbidden) {
			Forbidden(w, "not allowed to list these documents")
			return
		}
		BadRequest(w, err.Error())
		return
	}
	Success(w, list)
}

func (h *Handler) GetDocument(w http.ResponseWriter, r *http.Request) {
	requester, role, ok := RequesterAndRole(w, r)
	if !ok {
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "documentID"))
	if err != nil {
		BadRequest(w, "invalid document id")
		return
	}
	doc, err := h.services.Document.Get(r.Context(), id, requester, role)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			NotFound(w, "document not found")
			return
		}
		NotFound(w, err.Error())
		return
	}
	Success(w, doc)
}

func (h *Handler) DownloadDocument(w http.ResponseWriter, r *http.Request) {
	requester, role, ok := RequesterAndRole(w, r)
	if !ok {
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "documentID"))
	if err != nil {
		BadRequest(w, "invalid document id")
		return
	}
	rc, doc, err := h.services.Document.OpenFile(r.Context(), id, requester, role)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			NotFound(w, "document not found")
			return
		}
		NotFound(w, err.Error())
		return
	}
	defer rc.Close()

	ct := "application/octet-stream"
	if doc.MimeType != nil && *doc.MimeType != "" {
		ct = *doc.MimeType
	}
	w.Header().Set("Content-Type", ct)

	fn := "document"
	if doc.FileName != nil && *doc.FileName != "" {
		fn = *doc.FileName
	}
	cd := mime.FormatMediaType("attachment", map[string]string{"filename": fn})
	if cd == "" {
		cd = `attachment; filename="document"`
	}
	w.Header().Set("Content-Disposition", cd)

	w.WriteHeader(http.StatusOK)
	_, _ = io.Copy(w, rc)
}

func (h *Handler) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	requester, role, ok := RequesterAndRole(w, r)
	if !ok {
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "documentID"))
	if err != nil {
		BadRequest(w, "invalid document id")
		return
	}
	if err := h.services.Document.Delete(r.Context(), id, requester, role); err != nil {
		if errors.Is(err, apperr.ErrForbidden) {
			Forbidden(w, "not allowed to delete this document")
			return
		}
		BadRequest(w, err.Error())
		return
	}
	Success(w, map[string]string{"message": "document deleted"})
}
