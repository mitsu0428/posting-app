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
	IsDeleted    bool       `json:"is_deleted" db:"is_deleted"`
	GroupID      *int       `json:"group_id" db:"group_id"`
	Group        *Group     `json:"group,omitempty"`
	Categories   []Category `json:"categories,omitempty"`
	LikesCount   int        `json:"likes_count" db:"likes_count"`
	IsLiked      bool       `json:"is_liked" db:"is_liked"`
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

type Category struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Color       string    `json:"color" db:"color"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type Like struct {
	ID        int       `json:"id" db:"id"`
	PostID    int       `json:"post_id" db:"post_id"`
	UserID    int       `json:"user_id" db:"user_id"`
	User      *User     `json:"user,omitempty"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Group struct {
	ID          int           `json:"id" db:"id"`
	Name        string        `json:"name" db:"name"`
	Description string        `json:"description" db:"description"`
	OwnerID     int           `json:"owner_id" db:"owner_id"`
	Owner       *User         `json:"owner,omitempty"`
	IsActive    bool          `json:"is_active" db:"is_active"`
	Members     []GroupMember `json:"members,omitempty"`
	MemberCount int           `json:"member_count" db:"member_count"`
	CreatedAt   time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at" db:"updated_at"`
}

type GroupMember struct {
	ID       int       `json:"id" db:"id"`
	GroupID  int       `json:"group_id" db:"group_id"`
	UserID   int       `json:"user_id" db:"user_id"`
	User     *User     `json:"user,omitempty"`
	Role     string    `json:"role" db:"role"`
	JoinedAt time.Time `json:"joined_at" db:"joined_at"`
}

type PostStatus string

const (
	PostStatusPending  PostStatus = "pending"
	PostStatusApproved PostStatus = "approved"
	PostStatusRejected PostStatus = "rejected"
)
