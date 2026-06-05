package service

import (
	"context"

	"github.com/carkeeper/backend/internal/apperr"
	"github.com/carkeeper/backend/internal/model"
	"github.com/carkeeper/backend/internal/repository"
	"github.com/carkeeper/backend/internal/validate"
	"github.com/google/uuid"
)

type NewsService struct {
	repo *repository.Repository
}

func NewNewsService(repos *repository.Repository) *NewsService {
	return &NewsService{repo: repos}
}

func (s *NewsService) GetNews(ctx context.Context, isPublished *bool) ([]model.NewsWithAuthor, error) {
	return s.repo.News.GetAll(ctx, isPublished)
}

func (s *NewsService) GetNewsByID(ctx context.Context, newsID uuid.UUID) (*model.NewsWithAuthor, error) {
	return s.repo.News.GetByID(ctx, newsID)
}

func (s *NewsService) CreateNews(ctx context.Context, authorID uuid.UUID, create model.NewsCreate) (*model.News, error) {
	title, content, msg := validate.NewsFields(create.Title, create.Content)
	if msg != "" {
		return nil, apperr.BadRequest(msg)
	}
	create.Title = title
	create.Content = content
	return s.repo.News.Create(ctx, authorID, create)
}

func (s *NewsService) UpdateNews(ctx context.Context, newsID uuid.UUID, title, content string) error {
	title, content, msg := validate.NewsFields(title, content)
	if msg != "" {
		return apperr.BadRequest(msg)
	}
	return s.repo.News.Update(ctx, newsID, title, content)
}

func (s *NewsService) PublishNews(ctx context.Context, newsID uuid.UUID) error {
	return s.repo.News.Publish(ctx, newsID)
}

func (s *NewsService) UnpublishNews(ctx context.Context, newsID uuid.UUID) error {
	return s.repo.News.Unpublish(ctx, newsID)
}

func (s *NewsService) DeleteNews(ctx context.Context, newsID uuid.UUID) error {
	return s.repo.News.Delete(ctx, newsID)
}

