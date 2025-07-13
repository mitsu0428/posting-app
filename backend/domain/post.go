package domain

import (
	"time"
)

type Post struct {
	ID           int        `json:"id" db:"id"`
	Title        string     `json:"title" db:"title"`
	Content      string     `json:"content" db:"content"`
	ThumbnailURL *string    `json:"thumbnail_url" db:"thumbnail_url"`
	AuthorID     int        `json:"-" db:"author_id"`
	Author       *User      `json:"author,omitempty"`
	Status       PostStatus `json:"status" db:"status"`
	Replies      []Reply    `json:"replies,omitempty"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

type Reply struct {
	ID          int       `json:"id" db:"id"`
	Content     string    `json:"content" db:"content"`
	PostID      int       `json:"post_id" db:"post_id"`
	AuthorID    *int      `json:"-" db:"author_id"`
	Author      *User     `json:"author" db:"-"`
	IsAnonymous bool      `json:"is_anonymous" db:"is_anonymous"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type PostStatus string

const (
	PostStatusPending  PostStatus = "pending"
	PostStatusApproved PostStatus = "approved"
	PostStatusRejected PostStatus = "rejected"
)
