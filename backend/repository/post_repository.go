package repository

import (
	"database/sql"
	"posting-app/domain"
)

type postRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) domain.PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(post *domain.Post) error {
	query := `
		INSERT INTO posts (title, content, thumbnail_url, user_id, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`
	
	return r.db.QueryRow(query, post.Title, post.Content, post.ThumbnailURL,
		post.UserID, post.Status).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
}

func (r *postRepository) GetByID(id int) (*domain.Post, error) {
	post := &domain.Post{}
	query := `
		SELECT id, title, content, thumbnail_url, status, user_id, created_at, updated_at
		FROM posts WHERE id = $1`
	
	err := r.db.QueryRow(query, id).Scan(
		&post.ID, &post.Title, &post.Content, &post.ThumbnailURL,
		&post.Status, &post.UserID, &post.CreatedAt, &post.UpdatedAt)
	
	if err != nil {
		return nil, err
	}
	
	return post, nil
}

func (r *postRepository) List(status string, page, limit int) ([]*domain.Post, int, error) {
	offset := (page - 1) * limit
	
	var query string
	var args []interface{}
	
	if status != "" {
		query = `
			SELECT id, title, content, thumbnail_url, status, user_id, created_at, updated_at
			FROM posts WHERE status = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
		args = []interface{}{status, limit, offset}
	} else {
		query = `
			SELECT id, title, content, thumbnail_url, status, user_id, created_at, updated_at
			FROM posts ORDER BY created_at DESC LIMIT $1 OFFSET $2`
		args = []interface{}{limit, offset}
	}
	
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var posts []*domain.Post
	for rows.Next() {
		post := &domain.Post{}
		err := rows.Scan(
			&post.ID, &post.Title, &post.Content, &post.ThumbnailURL,
			&post.Status, &post.UserID, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		posts = append(posts, post)
	}

	var countQuery string
	var countArgs []interface{}
	if status != "" {
		countQuery = "SELECT COUNT(*) FROM posts WHERE status = $1"
		countArgs = []interface{}{status}
	} else {
		countQuery = "SELECT COUNT(*) FROM posts"
		countArgs = []interface{}{}
	}
	
	var total int
	err = r.db.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (r *postRepository) GetByUserID(userID int) ([]*domain.Post, error) {
	query := `
		SELECT id, title, content, thumbnail_url, status, user_id, created_at, updated_at
		FROM posts WHERE user_id = $1 ORDER BY created_at DESC`
	
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*domain.Post
	for rows.Next() {
		post := &domain.Post{}
		err := rows.Scan(
			&post.ID, &post.Title, &post.Content, &post.ThumbnailURL,
			&post.Status, &post.UserID, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *postRepository) Update(post *domain.Post) error {
	query := `
		UPDATE posts 
		SET title = $2, content = $3, thumbnail_url = $4, status = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1`
	
	_, err := r.db.Exec(query, post.ID, post.Title, post.Content,
		post.ThumbnailURL, post.Status)
	return err
}

func (r *postRepository) Delete(id int) error {
	query := `DELETE FROM posts WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *postRepository) Approve(id int) error {
	query := `UPDATE posts SET status = 'approved' WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *postRepository) Reject(id int) error {
	query := `UPDATE posts SET status = 'rejected' WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

type replyRepository struct {
	db *sql.DB
}

func NewReplyRepository(db *sql.DB) domain.ReplyRepository {
	return &replyRepository{db: db}
}

func (r *replyRepository) Create(reply *domain.Reply) error {
	query := `
		INSERT INTO replies (post_id, content, user_id, is_anonymous)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`
	
	return r.db.QueryRow(query, reply.PostID, reply.Content, reply.UserID,
		reply.IsAnonymous).Scan(&reply.ID, &reply.CreatedAt)
}

func (r *replyRepository) GetByPostID(postID int) ([]*domain.Reply, error) {
	query := `
		SELECT id, post_id, content, user_id, is_anonymous, created_at
		FROM replies WHERE post_id = $1 ORDER BY created_at ASC`
	
	rows, err := r.db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var replies []*domain.Reply
	for rows.Next() {
		reply := &domain.Reply{}
		err := rows.Scan(
			&reply.ID, &reply.PostID, &reply.Content, &reply.UserID,
			&reply.IsAnonymous, &reply.CreatedAt)
		if err != nil {
			return nil, err
		}
		replies = append(replies, reply)
	}

	return replies, nil
}