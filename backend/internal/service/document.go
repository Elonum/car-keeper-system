package service

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/carkeeper/backend/internal/apperr"
	"github.com/carkeeper/backend/internal/authz"
	"github.com/carkeeper/backend/internal/model"
	"github.com/carkeeper/backend/internal/repository"
	"github.com/carkeeper/backend/internal/storage"
	"github.com/google/uuid"
)

// Allowed document Content-Types (keep in sync with product policy).
var allowedDocumentMIME = map[string]struct{}{
	"application/pdf": {},
	"image/jpeg":      {},
	"image/png":       {},
	"image/webp":      {},
	"application/msword": {},
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": {},
	"application/vnd.ms-excel": {},
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": {},
}

type DocumentService struct {
	repo    *repository.Repository
	store   storage.FileStorage
	maxSize int64
}

func NewDocumentService(repos *repository.Repository, store storage.FileStorage, maxUploadBytes int64) *DocumentService {
	return &DocumentService{repo: repos, store: store, maxSize: maxUploadBytes}
}

// DocumentCreateInput is validated business input for a new file document.
type DocumentCreateInput struct {
	Reader       io.Reader
	Size         int64
	FileName     string
	MimeType     string
	DocumentType string
	OrderID      *uuid.UUID
	ApptID       *uuid.UUID
	Requester    uuid.UUID
	Role         string
}

func (s *DocumentService) Create(ctx context.Context, in DocumentCreateInput) (*model.Document, error) {
	if in.Size > s.maxSize {
		return nil, fmt.Errorf("file exceeds maximum size of %d bytes", s.maxSize)
	}
	if in.Size <= 0 {
		return nil, fmt.Errorf("empty or invalid file")
	}
	mt := strings.TrimSpace(strings.ToLower(in.MimeType))
	if _, ok := allowedDocumentMIME[mt]; !ok {
		return nil, fmt.Errorf("unsupported file type")
	}
	if !model.ValidDocumentType(in.DocumentType) {
		return nil, fmt.Errorf("invalid document_type")
	}

	var ownerUserID uuid.UUID
	var orderID *uuid.UUID
	var apptID *uuid.UUID

	switch {
	case in.OrderID != nil && in.ApptID == nil:
		order, err := s.repo.Order.GetByID(ctx, *in.OrderID)
		if err != nil {
			return nil, err
		}
		if !authz.CanViewOrder(order.UserID, in.Requester, in.Role) {
			return nil, fmt.Errorf("%w", apperr.ErrForbidden)
		}
		ownerUserID = order.UserID
		orderID = in.OrderID
	case in.ApptID != nil && in.OrderID == nil:
		appt, err := s.repo.ServiceAppointment.GetByID(ctx, *in.ApptID)
		if err != nil {
			return nil, err
		}
		if !authz.IsOwnerOrStaff(appt.OwnerUserID, in.Requester, in.Role) {
			return nil, fmt.Errorf("%w", apperr.ErrForbidden)
		}
		ownerUserID = appt.OwnerUserID
		apptID = in.ApptID
	default:
		return nil, fmt.Errorf("provide exactly one of order_id or service_appointment_id")
	}

	docID := uuid.New()
	key := docID.String()
	if err := s.store.Store(ctx, key, in.Reader, in.Size); err != nil {
		return nil, fmt.Errorf("store file: %w", err)
	}

	fn := sanitizeStoredFileName(in.FileName)
	var sizePtr *int64
	if in.Size > 0 {
		sizePtr = &in.Size
	}
	doc := model.Document{
		DocumentID:           docID,
		UserID:               ownerUserID,
		OrderID:              orderID,
		ServiceAppointmentID: apptID,
		DocumentType:         in.DocumentType,
		FilePath:             key,
		FileName:             stringPtrOrNil(fn),
		FileSize:             sizePtr,
		MimeType:             stringPtr(mt),
	}

	if err := s.repo.Document.Insert(ctx, doc); err != nil {
		_ = s.store.Remove(ctx, key)
		return nil, err
	}

	return s.repo.Document.GetByID(ctx, docID)
}

func (s *DocumentService) List(ctx context.Context, requester uuid.UUID, role string, orderID, apptID *uuid.UUID) ([]model.Document, error) {
	switch {
	case orderID != nil && apptID != nil:
		return nil, fmt.Errorf("use only one filter: order_id or service_appointment_id")
	case orderID != nil:
		order, err := s.repo.Order.GetByID(ctx, *orderID)
		if err != nil {
			return nil, err
		}
		if !authz.CanViewOrder(order.UserID, requester, role) {
			return nil, fmt.Errorf("%w", apperr.ErrForbidden)
		}
		return s.repo.Document.ListByOrderID(ctx, *orderID)
	case apptID != nil:
		appt, err := s.repo.ServiceAppointment.GetByID(ctx, *apptID)
		if err != nil {
			return nil, err
		}
		if !authz.IsOwnerOrStaff(appt.OwnerUserID, requester, role) {
			return nil, fmt.Errorf("%w", apperr.ErrForbidden)
		}
		return s.repo.Document.ListByServiceAppointmentID(ctx, *apptID)
	default:
		return s.repo.Document.ListByUserID(ctx, requester)
	}
}

func (s *DocumentService) Get(ctx context.Context, documentID uuid.UUID, requester uuid.UUID, role string) (*model.Document, error) {
	d, err := s.repo.Document.GetByID(ctx, documentID)
	if err != nil {
		return nil, err
	}
	if !s.canAccessDocument(d, requester, role) {
		return nil, fmt.Errorf("%w", apperr.ErrNotFound)
	}
	return d, nil
}

func (s *DocumentService) OpenFile(ctx context.Context, documentID uuid.UUID, requester uuid.UUID, role string) (io.ReadCloser, *model.Document, error) {
	d, err := s.repo.Document.GetByID(ctx, documentID)
	if err != nil {
		return nil, nil, err
	}
	if !s.canAccessDocument(d, requester, role) {
		return nil, nil, fmt.Errorf("%w", apperr.ErrNotFound)
	}
	rc, err := s.store.Open(ctx, d.FilePath)
	if err != nil {
		return nil, nil, fmt.Errorf("open file: %w", err)
	}
	return rc, d, nil
}

func (s *DocumentService) Delete(ctx context.Context, documentID uuid.UUID, requester uuid.UUID, role string) error {
	d, err := s.repo.Document.GetByID(ctx, documentID)
	if err != nil {
		return err
	}
	if !s.canAccessDocument(d, requester, role) {
		return fmt.Errorf("%w", apperr.ErrForbidden)
	}
	key := d.FilePath
	if err := s.repo.Document.Delete(ctx, documentID); err != nil {
		return err
	}
	_ = s.store.Remove(ctx, key)
	return nil
}

func (s *DocumentService) canAccessDocument(d *model.Document, requester uuid.UUID, role string) bool {
	if authz.IsStaff(role) {
		return true
	}
	if d.UserID == requester {
		return true
	}
	return false
}

func stringPtr(s string) *string {
	return &s
}

func stringPtrOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func sanitizeStoredFileName(name string) string {
	base := filepath.Base(strings.ReplaceAll(name, "\\", "/"))
	var b strings.Builder
	for _, r := range base {
		if r == unicode.ReplacementChar {
			continue
		}
		if r < 32 || r == '"' || r == '|' || r == '>' || r == '<' {
			continue
		}
		b.WriteRune(r)
	}
	out := strings.TrimSpace(b.String())
	if len(out) > 255 {
		out = out[:255]
	}
	if out == "" || out == "." || out == ".." {
		return "document"
	}
	return out
}
