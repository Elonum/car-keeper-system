package service

import (
	"bytes"
	"context"
	"errors"
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
	"github.com/carkeeper/backend/internal/upload"
	"github.com/google/uuid"
)

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
	Size         int64 // optional hint from multipart; actual size is measured from bytes read
	FileName     string
	MimeType     string
	DocumentType string
	OrderID      *uuid.UUID
	ApptID       *uuid.UUID
	Requester    uuid.UUID
	Role         string
}

func (s *DocumentService) maxUploadSize() int64 {
	if s.maxSize > 0 {
		return s.maxSize
	}
	return 15 << 20
}

func (s *DocumentService) Create(ctx context.Context, in DocumentCreateInput) (*model.Document, error) {
	max := s.maxUploadSize()
	if in.Size > max {
		return nil, apperr.PayloadTooLarge(fmt.Sprintf("file exceeds maximum size of %d bytes", max))
	}

	buf := bytes.NewBuffer(make([]byte, 0, 4096))
	limited := io.LimitReader(in.Reader, max+1)
	written, err := io.Copy(buf, limited)
	if err != nil {
		return nil, apperr.BadRequest("failed to read uploaded file")
	}
	if written > max {
		return nil, apperr.PayloadTooLarge(fmt.Sprintf("file exceeds maximum size of %d bytes", max))
	}
	if written == 0 {
		return nil, apperr.BadRequest("file is empty")
	}

	data := buf.Bytes()
	mt := upload.ResolveDocumentMIME(data, in.MimeType, in.FileName)
	if mt == "" {
		return nil, apperr.BadRequest("unsupported file type; allowed: PDF, JPEG, PNG, WebP, Word, Excel")
	}
	if !model.ValidDocumentType(in.DocumentType) {
		return nil, apperr.BadRequest("invalid document_type")
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
		if !authz.IsOwnerOrHasPermission(appt.OwnerUserID, in.Requester, in.Role, authz.PermAppointmentsViewAny) {
			return nil, fmt.Errorf("%w", apperr.ErrForbidden)
		}
		ownerUserID = appt.OwnerUserID
		apptID = in.ApptID
	default:
		return nil, apperr.BadRequest("provide exactly one of order_id or service_appointment_id")
	}

	docID := uuid.New()
	key := docID.String()
	payload := bytes.NewReader(data)
	if err := s.store.Store(ctx, key, payload, written); err != nil {
		return nil, apperr.Internal(fmt.Errorf("store file: %w", err))
	}

	fn := sanitizeStoredFileName(in.FileName)
	size := written
	doc := model.Document{
		DocumentID:           docID,
		UserID:               ownerUserID,
		OrderID:              orderID,
		ServiceAppointmentID: apptID,
		DocumentType:         in.DocumentType,
		FilePath:             key,
		FileName:             stringPtrOrNil(fn),
		FileSize:             &size,
		MimeType:             stringPtr(mt),
		FileAvailable:        true,
	}

	if err := s.repo.Document.Insert(ctx, doc); err != nil {
		_ = s.store.Remove(ctx, key)
		return nil, err
	}

	out, err := s.repo.Document.GetByID(ctx, docID)
	if err != nil {
		return nil, err
	}
	s.enrichFileAvailable(ctx, out)
	s.redactDocumentContext(out, in.Role)
	return out, nil
}

func (s *DocumentService) List(ctx context.Context, requester uuid.UUID, role string, orderID, apptID *uuid.UUID) ([]model.Document, error) {
	var list []model.Document
	var err error

	switch {
	case orderID != nil && apptID != nil:
		return nil, apperr.BadRequest("use only one filter: order_id or service_appointment_id")
	case orderID != nil:
		order, err := s.repo.Order.GetByID(ctx, *orderID)
		if err != nil {
			return nil, err
		}
		if !authz.CanViewOrder(order.UserID, requester, role) {
			return nil, fmt.Errorf("%w", apperr.ErrForbidden)
		}
		list, err = s.repo.Document.ListByOrderID(ctx, *orderID)
	case apptID != nil:
		appt, err := s.repo.ServiceAppointment.GetByID(ctx, *apptID)
		if err != nil {
			return nil, err
		}
		if !authz.IsOwnerOrHasPermission(appt.OwnerUserID, requester, role, authz.PermAppointmentsViewAny) {
			return nil, fmt.Errorf("%w", apperr.ErrForbidden)
		}
		list, err = s.repo.Document.ListByServiceAppointmentID(ctx, *apptID)
	default:
		if authz.HasPermission(role, authz.PermDocumentsViewAny) {
			list, err = s.repo.Document.ListAll(ctx)
		} else {
			list, err = s.repo.Document.ListByUserID(ctx, requester)
		}
	}
	if err != nil {
		return nil, err
	}
	for i := range list {
		s.enrichFileAvailable(ctx, &list[i])
		s.redactDocumentContext(&list[i], role)
	}
	return list, nil
}

func (s *DocumentService) Get(ctx context.Context, documentID uuid.UUID, requester uuid.UUID, role string) (*model.Document, error) {
	d, err := s.repo.Document.GetByID(ctx, documentID)
	if err != nil {
		return nil, err
	}
	if !s.canAccessDocument(d, requester, role) {
		return nil, fmt.Errorf("%w", apperr.ErrNotFound)
	}
	s.enrichFileAvailable(ctx, d)
	s.redactDocumentContext(d, role)
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
		if errors.Is(err, storage.ErrNotFound) {
			return nil, nil, apperr.NotFoundErr("file is not available on the server")
		}
		return nil, nil, apperr.Internal(fmt.Errorf("open file: %w", err))
	}
	d.FileAvailable = true
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

func (s *DocumentService) enrichFileAvailable(ctx context.Context, d *model.Document) {
	if d == nil || d.FilePath == "" {
		d.FileAvailable = false
		return
	}
	exists, err := s.store.Exists(ctx, d.FilePath)
	if err != nil {
		d.FileAvailable = false
		return
	}
	d.FileAvailable = exists
}

func (s *DocumentService) redactDocumentContext(d *model.Document, role string) {
	if d == nil {
		return
	}
	if !authz.HasPermission(role, authz.PermDocumentsViewAny) {
		d.OwnerEmail = nil
	}
}

func (s *DocumentService) canAccessDocument(d *model.Document, requester uuid.UUID, role string) bool {
	if authz.HasPermission(role, authz.PermDocumentsViewAny) {
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
