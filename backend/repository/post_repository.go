package repository

import (
	"database/sql"
	"fmt"
	"strconv"

	"posting-app/domain"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(post *domain.Post) error {
	query := `
		INSERT INTO posts (title, content, thumbnail_url, author_id, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(
		query,
		post.Title,
		post.Content,
		post.ThumbnailURL,
		post.AuthorID,
		post.Status,
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)

	return err
}

func (r *PostRepository) GetByID(id int) (*domain.Post, error) {
	post := &domain.Post{}
	query := `
		SELECT p.id, p.title, p.content, p.thumbnail_url, p.author_id, p.status, p.created_at, p.updated_at,
			   u.id, u.email, u.display_name, u.bio, u.role, u.subscription_status, u.is_active, u.created_at, u.updated_at
		FROM posts p
		JOIN users u ON p.author_id = u.id
		WHERE p.id = $1`

	var author domain.User
	err := r.db.QueryRow(query, id).Scan(
		&post.ID, &post.Title, &post.Content, &post.ThumbnailURL, &post.AuthorID, &post.Status, &post.CreatedAt, &post.UpdatedAt,
		&author.ID, &author.Email, &author.DisplayName, &author.Bio, &author.Role, &author.SubscriptionStatus, &author.IsActive, &author.CreatedAt, &author.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	post.Author = &author

	// Get replies
	replies, err := r.getRepliesByPostID(id)
	if err != nil {
		return nil, err
	}
	post.Replies = replies

	return post, nil
}

func (r *PostRepository) GetApproved(page, limit int) ([]*domain.Post, int, error) {
	offset := (page - 1) * limit

	// Get total count
	var total int
	err := r.db.QueryRow("SELECT COUNT(*) FROM posts p JOIN users u ON p.author_id = u.id WHERE p.status = $1 AND u.is_active = true", domain.PostStatusApproved).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get posts
	query := `
		SELECT p.id, p.title, p.content, p.thumbnail_url, p.author_id, p.status, p.created_at, p.updated_at,
			   u.id, u.email, u.display_name, u.bio, u.role, u.subscription_status, u.is_active, u.created_at, u.updated_at
		FROM posts p
		JOIN users u ON p.author_id = u.id
		WHERE p.status = $1 AND u.is_active = true
		ORDER BY p.created_at DESC 
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, domain.PostStatusApproved, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var posts []*domain.Post
	for rows.Next() {
		post := &domain.Post{}
		author := &domain.User{}
		err := rows.Scan(
			&post.ID, &post.Title, &post.Content, &post.ThumbnailURL, &post.AuthorID, &post.Status, &post.CreatedAt, &post.UpdatedAt,
			&author.ID, &author.Email, &author.DisplayName, &author.Bio, &author.Role, &author.SubscriptionStatus, &author.IsActive, &author.CreatedAt, &author.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		post.Author = author
		posts = append(posts, post)
	}

	return posts, total, nil
}

func (r *PostRepository) GetByUserID(userID, page, limit int) ([]*domain.Post, int, error) {
	offset := (page - 1) * limit

	// Get total count
	var total int
	err := r.db.QueryRow("SELECT COUNT(*) FROM posts WHERE author_id = $1", userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get posts
	query := `
		SELECT p.id, p.title, p.content, p.thumbnail_url, p.author_id, p.status, p.created_at, p.updated_at,
			   u.id, u.email, u.display_name, u.bio, u.role, u.subscription_status, u.is_active, u.created_at, u.updated_at
		FROM posts p
		JOIN users u ON p.author_id = u.id
		WHERE p.author_id = $1
		ORDER BY p.created_at DESC 
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var posts []*domain.Post
	for rows.Next() {
		post := &domain.Post{}
		author := &domain.User{}
		err := rows.Scan(
			&post.ID, &post.Title, &post.Content, &post.ThumbnailURL, &post.AuthorID, &post.Status, &post.CreatedAt, &post.UpdatedAt,
			&author.ID, &author.Email, &author.DisplayName, &author.Bio, &author.Role, &author.SubscriptionStatus, &author.IsActive, &author.CreatedAt, &author.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		post.Author = author
		posts = append(posts, post)
	}

	return posts, total, nil
}

func (r *PostRepository) GetForAdmin(page, limit int, status *domain.PostStatus) ([]*domain.Post, int, error) {
	offset := (page - 1) * limit

	whereClause := "u.is_active = true"
	args := []interface{}{limit, offset}

	if status != nil {
		whereClause += " AND p.status = $3"
		args = append(args, *status)
	}

	// Get total count
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM posts p JOIN users u ON p.author_id = u.id WHERE %s", whereClause)
	var countArgs []interface{}
	if status != nil {
		countArgs = []interface{}{*status}
	}
	err := r.db.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get posts
	query := fmt.Sprintf(`
		SELECT p.id, p.title, p.content, p.thumbnail_url, p.author_id, p.status, p.created_at, p.updated_at,
			   u.id, u.email, u.display_name, u.bio, u.role, u.subscription_status, u.is_active, u.created_at, u.updated_at
		FROM posts p
		JOIN users u ON p.author_id = u.id
		WHERE %s
		ORDER BY p.created_at DESC 
		LIMIT $1 OFFSET $2`, whereClause)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var posts []*domain.Post
	for rows.Next() {
		post := &domain.Post{}
		author := &domain.User{}
		err := rows.Scan(
			&post.ID, &post.Title, &post.Content, &post.ThumbnailURL, &post.AuthorID, &post.Status, &post.CreatedAt, &post.UpdatedAt,
			&author.ID, &author.Email, &author.DisplayName, &author.Bio, &author.Role, &author.SubscriptionStatus, &author.IsActive, &author.CreatedAt, &author.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		post.Author = author
		posts = append(posts, post)
	}

	return posts, total, nil
}

func (r *PostRepository) Update(post *domain.Post) error {
	query := `
		UPDATE posts 
		SET title = $1, content = $2, thumbnail_url = $3, status = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5`

	_, err := r.db.Exec(query, post.Title, post.Content, post.ThumbnailURL, post.Status, post.ID)
	return err
}

func (r *PostRepository) Delete(id int) error {
	query := `DELETE FROM posts WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *PostRepository) UpdateStatus(id int, status domain.PostStatus) error {
	query := `UPDATE posts SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.Exec(query, status, id)
	return err
}

func (r *PostRepository) CreateReply(reply *domain.Reply) error {
	query := `
		INSERT INTO replies (content, post_id, author_id, is_anonymous)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	err := r.db.QueryRow(
		query,
		reply.Content,
		reply.PostID,
		reply.AuthorID,
		reply.IsAnonymous,
	).Scan(&reply.ID, &reply.CreatedAt)

	return err
}

func (r *PostRepository) getRepliesByPostID(postID int) ([]domain.Reply, error) {
	query := `
		SELECT r.id, r.content, r.post_id, r.author_id, r.is_anonymous, r.created_at,
			   u.id, u.email, u.display_name, u.bio, u.role, u.subscription_status, u.is_active, u.created_at, u.updated_at
		FROM replies r
		LEFT JOIN users u ON r.author_id = u.id AND u.is_active = true
		WHERE r.post_id = $1
		ORDER BY r.created_at ASC`

	rows, err := r.db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var replies []domain.Reply
	for rows.Next() {
		reply := domain.Reply{}
		var author *domain.User

		var userID, userEmail, userDisplayName, userBio, userRole, userSubscriptionStatus sql.NullString
		var userIsActive sql.NullBool
		var userCreatedAt, userUpdatedAt sql.NullTime

		err := rows.Scan(
			&reply.ID, &reply.Content, &reply.PostID, &reply.AuthorID, &reply.IsAnonymous, &reply.CreatedAt,
			&userID, &userEmail, &userDisplayName, &userBio, &userRole, &userSubscriptionStatus, &userIsActive, &userCreatedAt, &userUpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if !reply.IsAnonymous && userID.Valid {
			id, _ := strconv.Atoi(userID.String)
			author = &domain.User{
				ID:                 id,
				Email:              userEmail.String,
				DisplayName:        userDisplayName.String,
				Role:               userRole.String,
				SubscriptionStatus: domain.UserSubscriptionStatus(userSubscriptionStatus.String),
				IsActive:           userIsActive.Bool,
				CreatedAt:          userCreatedAt.Time,
				UpdatedAt:          userUpdatedAt.Time,
			}
			if userBio.Valid {
				author.Bio = &userBio.String
			}
		}

		reply.Author = author
		replies = append(replies, reply)
	}

	return replies, nil
}
