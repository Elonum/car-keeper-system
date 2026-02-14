package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/carkeeper/backend/database"
	"github.com/carkeeper/backend/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type NewsRepository struct {
	db *database.DB
}

func NewNewsRepository(db *database.DB) *NewsRepository {
	return &NewsRepository{db: db}
}

func (r *NewsRepository) GetAll(ctx context.Context, isPublished *bool) ([]model.NewsWithAuthor, error) {
	query := `
		SELECT 
			n.news_id, n.title, n.content, n.author_id, n.published_at, n.is_published,
			n.created_at, n.updated_at,
			u.first_name || ' ' || u.last_name as author_name
		FROM news n
		LEFT JOIN users u ON n.author_id = u.user_id
		WHERE 1=1
	`
	var args []interface{}
	argPos := 1

	if isPublished != nil {
		query += fmt.Sprintf(` AND n.is_published = $%d`, argPos)
		args = append(args, *isPublished)
		argPos++
	}

	query += ` ORDER BY n.published_at DESC NULLS LAST, n.created_at DESC`

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get news: %w", err)
	}
	defer rows.Close()

	var newsList []model.NewsWithAuthor
	for rows.Next() {
		var news model.NewsWithAuthor
		if err := rows.Scan(
			&news.NewsID, &news.Title, &news.Content, &news.AuthorID,
			&news.PublishedAt, &news.IsPublished, &news.CreatedAt, &news.UpdatedAt,
			&news.AuthorName,
		); err != nil {
			return nil, fmt.Errorf("failed to scan news: %w", err)
		}
		newsList = append(newsList, news)
	}

	return newsList, nil
}

func (r *NewsRepository) GetByID(ctx context.Context, newsID uuid.UUID) (*model.NewsWithAuthor, error) {
	var news model.NewsWithAuthor
	query := `
		SELECT 
			n.news_id, n.title, n.content, n.author_id, n.published_at, n.is_published,
			n.created_at, n.updated_at,
			u.first_name || ' ' || u.last_name as author_name
		FROM news n
		LEFT JOIN users u ON n.author_id = u.user_id
		WHERE n.news_id = $1
	`

	err := r.db.Pool.QueryRow(ctx, query, newsID).Scan(
		&news.NewsID, &news.Title, &news.Content, &news.AuthorID,
		&news.PublishedAt, &news.IsPublished, &news.CreatedAt, &news.UpdatedAt,
		&news.AuthorName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("news not found")
		}
		return nil, fmt.Errorf("failed to get news: %w", err)
	}

	return &news, nil
}

func (r *NewsRepository) Create(ctx context.Context, authorID uuid.UUID, create model.NewsCreate) (*model.News, error) {
	var news model.News
	query := `
		INSERT INTO news (title, content, author_id, is_published)
		VALUES ($1, $2, $3, false)
		RETURNING news_id, title, content, author_id, published_at, is_published, created_at, updated_at
	`

	err := r.db.Pool.QueryRow(ctx, query, create.Title, create.Content, authorID).Scan(
		&news.NewsID, &news.Title, &news.Content, &news.AuthorID,
		&news.PublishedAt, &news.IsPublished, &news.CreatedAt, &news.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create news: %w", err)
	}

	return &news, nil
}

func (r *NewsRepository) Update(ctx context.Context, newsID uuid.UUID, title, content string) error {
	query := `UPDATE news SET title = $1, content = $2 WHERE news_id = $3`
	_, err := r.db.Pool.Exec(ctx, query, title, content, newsID)
	if err != nil {
		return fmt.Errorf("failed to update news: %w", err)
	}
	return nil
}

func (r *NewsRepository) Publish(ctx context.Context, newsID uuid.UUID) error {
	query := `UPDATE news SET is_published = true, published_at = now() WHERE news_id = $1`
	_, err := r.db.Pool.Exec(ctx, query, newsID)
	if err != nil {
		return fmt.Errorf("failed to publish news: %w", err)
	}
	return nil
}

func (r *NewsRepository) Unpublish(ctx context.Context, newsID uuid.UUID) error {
	query := `UPDATE news SET is_published = false, published_at = NULL WHERE news_id = $1`
	_, err := r.db.Pool.Exec(ctx, query, newsID)
	if err != nil {
		return fmt.Errorf("failed to unpublish news: %w", err)
	}
	return nil
}

func (r *NewsRepository) Delete(ctx context.Context, newsID uuid.UUID) error {
	query := `DELETE FROM news WHERE news_id = $1`
	_, err := r.db.Pool.Exec(ctx, query, newsID)
	if err != nil {
		return fmt.Errorf("failed to delete news: %w", err)
	}
	return nil
}

