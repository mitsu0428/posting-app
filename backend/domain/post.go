package domain

import "time"

type Post struct {
	ID           int       `json:"id" db:"id"`
	Title        string    `json:"title" db:"title"`
	Content      string    `json:"content" db:"content"`
	ThumbnailURL *string   `json:"thumbnail_url,omitempty" db:"thumbnail_url"`
	Status       string    `json:"status" db:"status"`
	UserID       int       `json:"user_id" db:"user_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type Reply struct {
	ID          int       `json:"id" db:"id"`
	PostID      int       `json:"post_id" db:"post_id"`
	Content     string    `json:"content" db:"content"`
	UserID      int       `json:"user_id" db:"user_id"`
	IsAnonymous bool      `json:"is_anonymous" db:"is_anonymous"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type PostRepository interface {
	Create(post *Post) error
	GetByID(id int) (*Post, error)
	List(status string, page, limit int) ([]*Post, int, error)
	GetByUserID(userID int) ([]*Post, error)
	Update(post *Post) error
	Delete(id int) error
	Approve(id int) error
	Reject(id int) error
}

type ReplyRepository interface {
	Create(reply *Reply) error
	GetByPostID(postID int) ([]*Reply, error)
}

type PostUsecase interface {
	CreatePost(userID int, title, content string, thumbnailURL *string) (*Post, error)
	GetPost(id int) (*Post, error)
	ListPosts(status string, page, limit int) ([]*Post, int, error)
	GetUserPosts(userID int) ([]*Post, error)
	CreateReply(postID, userID int, content string, isAnonymous bool) (*Reply, error)
	GetReplies(postID int) ([]*Reply, error)
}

type AdminUsecase interface {
	ListAllPosts(status string) ([]*Post, error)
	ApprovePost(id int) error
	RejectPost(id int) error
	DeletePost(id int) error
	ListUsers() ([]*User, error)
	DeactivateUser(id int) error
}