package service

import (
	"context"

	"github.com/carkeeper/backend/internal/model"
	"github.com/carkeeper/backend/internal/repository"
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
	return s.repo.News.Create(ctx, authorID, create)
}

func (s *NewsService) UpdateNews(ctx context.Context, newsID uuid.UUID, title, content string) error {
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

