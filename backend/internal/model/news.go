package model

import (
	"time"

	"github.com/google/uuid"
)

type News struct {
	NewsID      uuid.UUID  `db:"news_id" json:"news_id"`
	Title       string     `db:"title" json:"title"`
	Content     string     `db:"content" json:"content"`
	AuthorID    *uuid.UUID `db:"author_id" json:"author_id,omitempty"`
	PublishedAt *time.Time `db:"published_at" json:"published_at,omitempty"`
	IsPublished bool       `db:"is_published" json:"is_published"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
}

type NewsCreate struct {
	Title   string `json:"title" validate:"required,min=1,max=255"`
	Content string `json:"content" validate:"required"`
}

type NewsWithAuthor struct {
	News
	AuthorName *string `db:"author_name" json:"author_name,omitempty"`
}

