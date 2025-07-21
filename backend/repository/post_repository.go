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
		INSERT INTO posts (title, content, thumbnail_url, author_id, status, group_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(
		query,
		post.Title,
		post.Content,
		post.ThumbnailURL,
		post.AuthorID,
		post.Status,
		post.GroupID,
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)

	return err
}

func (r *PostRepository) GetByID(id int) (*domain.Post, error) {
	post := &domain.Post{}
	query := `
		SELECT p.id, p.title, p.content, p.thumbnail_url, p.author_id, p.status, p.is_deleted, p.group_id, p.created_at, p.updated_at,
			   u.id, u.email, u.display_name, u.bio, u.role, u.subscription_status, u.is_active, u.created_at, u.updated_at,
			   COALESCE(likes_count.count, 0) as likes_count
		FROM posts p
		JOIN users u ON p.author_id = u.id
		LEFT JOIN (SELECT post_id, COUNT(*) as count FROM likes GROUP BY post_id) likes_count ON p.id = likes_count.post_id
		WHERE p.id = $1 AND p.is_deleted = false`

	var author domain.User
	err := r.db.QueryRow(query, id).Scan(
		&post.ID, &post.Title, &post.Content, &post.ThumbnailURL, &post.AuthorID, &post.Status, &post.IsDeleted, &post.GroupID, &post.CreatedAt, &post.UpdatedAt,
		&author.ID, &author.Email, &author.DisplayName, &author.Bio, &author.Role, &author.SubscriptionStatus, &author.IsActive, &author.CreatedAt, &author.UpdatedAt,
		&post.LikesCount,
	)

	if err != nil {
		return nil, err
	}

	post.Author = &author

	// Get categories
	categories, err := r.GetPostCategories(id)
	if err != nil {
		return nil, err
	}
	post.Categories = categories

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
	err := r.db.QueryRow("SELECT COUNT(*) FROM posts p JOIN users u ON p.author_id = u.id WHERE p.status = $1 AND u.is_active = true AND p.is_deleted = false AND p.group_id IS NULL", domain.PostStatusApproved).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get posts
	query := `
		SELECT p.id, p.title, p.content, p.thumbnail_url, p.author_id, p.status, p.is_deleted, p.group_id, p.created_at, p.updated_at,
			   u.id, u.email, u.display_name, u.bio, u.role, u.subscription_status, u.is_active, u.created_at, u.updated_at,
			   COALESCE(likes_count.count, 0) as likes_count
		FROM posts p
		JOIN users u ON p.author_id = u.id
		LEFT JOIN (SELECT post_id, COUNT(*) as count FROM likes GROUP BY post_id) likes_count ON p.id = likes_count.post_id
		WHERE p.status = $1 AND u.is_active = true AND p.is_deleted = false AND p.group_id IS NULL
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
			&post.ID, &post.Title, &post.Content, &post.ThumbnailURL, &post.AuthorID, &post.Status, &post.IsDeleted, &post.GroupID, &post.CreatedAt, &post.UpdatedAt,
			&author.ID, &author.Email, &author.DisplayName, &author.Bio, &author.Role, &author.SubscriptionStatus, &author.IsActive, &author.CreatedAt, &author.UpdatedAt,
			&post.LikesCount,
		)
		if err != nil {
			return nil, 0, err
		}
		post.Author = author
		
		// Get categories for each post
		categories, err := r.GetPostCategories(post.ID)
		if err != nil {
			return nil, 0, err
		}
		post.Categories = categories
		
		posts = append(posts, post)
	}

	return posts, total, nil
}

func (r *PostRepository) GetByUserID(userID, page, limit int) ([]*domain.Post, int, error) {
	offset := (page - 1) * limit

	// Get total count
	var total int
	err := r.db.QueryRow("SELECT COUNT(*) FROM posts WHERE author_id = $1 AND is_deleted = false", userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get posts
	query := `
		SELECT p.id, p.title, p.content, p.thumbnail_url, p.author_id, p.status, p.is_deleted, p.group_id, p.created_at, p.updated_at,
			   u.id, u.email, u.display_name, u.bio, u.role, u.subscription_status, u.is_active, u.created_at, u.updated_at,
			   COALESCE(likes_count.count, 0) as likes_count
		FROM posts p
		JOIN users u ON p.author_id = u.id
		LEFT JOIN (SELECT post_id, COUNT(*) as count FROM likes GROUP BY post_id) likes_count ON p.id = likes_count.post_id
		WHERE p.author_id = $1 AND p.is_deleted = false
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
			&post.ID, &post.Title, &post.Content, &post.ThumbnailURL, &post.AuthorID, &post.Status, &post.IsDeleted, &post.GroupID, &post.CreatedAt, &post.UpdatedAt,
			&author.ID, &author.Email, &author.DisplayName, &author.Bio, &author.Role, &author.SubscriptionStatus, &author.IsActive, &author.CreatedAt, &author.UpdatedAt,
			&post.LikesCount,
		)
		if err != nil {
			return nil, 0, err
		}
		post.Author = author
		
		// Get categories for each post
		categories, err := r.GetPostCategories(post.ID)
		if err != nil {
			return nil, 0, err
		}
		post.Categories = categories
		
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
	countWhereClause := "u.is_active = true"
	var countArgs []interface{}
	if status != nil {
		countWhereClause += " AND p.status = $1"
		countArgs = []interface{}{*status}
	}
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM posts p JOIN users u ON p.author_id = u.id WHERE %s", countWhereClause)
	err := r.db.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get posts
	query := fmt.Sprintf(`
		SELECT p.id, p.title, p.content, p.thumbnail_url, p.author_id, p.status, p.is_deleted, p.group_id, p.created_at, p.updated_at,
			   u.id, u.email, u.display_name, u.bio, u.role, u.subscription_status, u.is_active, u.created_at, u.updated_at,
			   COALESCE(likes_count.count, 0) as likes_count
		FROM posts p
		JOIN users u ON p.author_id = u.id
		LEFT JOIN (SELECT post_id, COUNT(*) as count FROM likes GROUP BY post_id) likes_count ON p.id = likes_count.post_id
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
			&post.ID, &post.Title, &post.Content, &post.ThumbnailURL, &post.AuthorID, &post.Status, &post.IsDeleted, &post.GroupID, &post.CreatedAt, &post.UpdatedAt,
			&author.ID, &author.Email, &author.DisplayName, &author.Bio, &author.Role, &author.SubscriptionStatus, &author.IsActive, &author.CreatedAt, &author.UpdatedAt,
			&post.LikesCount,
		)
		if err != nil {
			return nil, 0, err
		}
		post.Author = author
		
		// Get categories for each post
		categories, err := r.GetPostCategories(post.ID)
		if err != nil {
			return nil, 0, err
		}
		post.Categories = categories
		
		posts = append(posts, post)
	}

	return posts, total, nil
}

func (r *PostRepository) Update(post *domain.Post) error {
	query := `
		UPDATE posts 
		SET title = $1, content = $2, thumbnail_url = $3, status = $4, group_id = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6`

	_, err := r.db.Exec(query, post.Title, post.Content, post.ThumbnailURL, post.Status, post.GroupID, post.ID)
	return err
}

func (r *PostRepository) Delete(id int) error {
	query := `UPDATE posts SET is_deleted = true, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
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
				Role:               domain.UserRole(userRole.String),
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

// Category related methods
func (r *PostRepository) GetAllCategories() ([]domain.Category, error) {
	query := `SELECT id, name, description, color, created_at, updated_at FROM categories ORDER BY name`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var category domain.Category
		err := rows.Scan(&category.ID, &category.Name, &category.Description, &category.Color, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (r *PostRepository) CreateCategory(category *domain.Category) error {
	query := `INSERT INTO categories (name, description, color) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(query, category.Name, category.Description, category.Color).Scan(&category.ID, &category.CreatedAt, &category.UpdatedAt)
	return err
}

func (r *PostRepository) AddPostCategories(postID int, categoryIDs []int) error {
	if len(categoryIDs) == 0 {
		return nil
	}
	
	// Delete existing categories for this post
	_, err := r.db.Exec("DELETE FROM post_categories WHERE post_id = $1", postID)
	if err != nil {
		return err
	}
	
	// Insert new categories
	for _, categoryID := range categoryIDs {
		_, err := r.db.Exec("INSERT INTO post_categories (post_id, category_id) VALUES ($1, $2)", postID, categoryID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *PostRepository) GetPostCategories(postID int) ([]domain.Category, error) {
	query := `
		SELECT c.id, c.name, c.description, c.color, c.created_at, c.updated_at
		FROM categories c
		JOIN post_categories pc ON c.id = pc.category_id
		WHERE pc.post_id = $1
		ORDER BY c.name`
	
	rows, err := r.db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var category domain.Category
		err := rows.Scan(&category.ID, &category.Name, &category.Description, &category.Color, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

// Like related methods
func (r *PostRepository) ToggleLike(postID, userID int) error {
	// Check if like exists
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM likes WHERE post_id = $1 AND user_id = $2)", postID, userID).Scan(&exists)
	if err != nil {
		return err
	}
	
	if exists {
		// Remove like
		_, err = r.db.Exec("DELETE FROM likes WHERE post_id = $1 AND user_id = $2", postID, userID)
	} else {
		// Add like
		_, err = r.db.Exec("INSERT INTO likes (post_id, user_id) VALUES ($1, $2)", postID, userID)
	}
	return err
}

func (r *PostRepository) GetLikesCount(postID int) (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = $1", postID).Scan(&count)
	return count, err
}

func (r *PostRepository) IsLikedByUser(postID, userID int) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM likes WHERE post_id = $1 AND user_id = $2)", postID, userID).Scan(&exists)
	return exists, err
}

// Group related methods
func (r *PostRepository) CreateGroup(group *domain.Group) error {
	query := `INSERT INTO groups (name, description, owner_id) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(query, group.Name, group.Description, group.OwnerID).Scan(&group.ID, &group.CreatedAt, &group.UpdatedAt)
	if err != nil {
		return err
	}
	
	// Add owner as member
	_, err = r.db.Exec("INSERT INTO group_members (group_id, user_id, role) VALUES ($1, $2, 'owner')", group.ID, group.OwnerID)
	return err
}

func (r *PostRepository) GetUserGroups(userID int) ([]domain.Group, error) {
	query := `
		SELECT g.id, g.name, g.description, g.owner_id, g.is_active, g.created_at, g.updated_at,
			   COUNT(gm.user_id) as member_count
		FROM groups g
		LEFT JOIN group_members gm ON g.id = gm.group_id
		WHERE g.owner_id = $1 OR g.id IN (SELECT group_id FROM group_members WHERE user_id = $1)
		GROUP BY g.id, g.name, g.description, g.owner_id, g.is_active, g.created_at, g.updated_at
		ORDER BY g.name`
		
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []domain.Group
	for rows.Next() {
		var group domain.Group
		err := rows.Scan(&group.ID, &group.Name, &group.Description, &group.OwnerID, &group.IsActive, &group.CreatedAt, &group.UpdatedAt, &group.MemberCount)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	return groups, nil
}

func (r *PostRepository) AddGroupMember(groupID, userID int) error {
	_, err := r.db.Exec("INSERT INTO group_members (group_id, user_id) VALUES ($1, $2)", groupID, userID)
	return err
}

func (r *PostRepository) IsGroupMember(groupID, userID int) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = $1 AND user_id = $2)", groupID, userID).Scan(&exists)
	return exists, err
}

func (r *PostRepository) GetGroupPosts(groupID, userID, page, limit int) ([]*domain.Post, int, error) {
	// Check if user is member of the group
	isMember, err := r.IsGroupMember(groupID, userID)
	if err != nil {
		return nil, 0, err
	}
	if !isMember {
		return nil, 0, fmt.Errorf("user is not a member of this group")
	}

	offset := (page - 1) * limit

	// Get total count
	var total int
	err = r.db.QueryRow("SELECT COUNT(*) FROM posts WHERE group_id = $1 AND is_deleted = false", groupID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get posts
	query := `
		SELECT p.id, p.title, p.content, p.thumbnail_url, p.author_id, p.status, p.is_deleted, p.group_id, p.created_at, p.updated_at,
			   u.id, u.email, u.display_name, u.bio, u.role, u.subscription_status, u.is_active, u.created_at, u.updated_at
		FROM posts p
		JOIN users u ON p.author_id = u.id
		WHERE p.group_id = $1 AND p.is_deleted = false
		ORDER BY p.created_at DESC 
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, groupID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var posts []*domain.Post
	for rows.Next() {
		post := &domain.Post{}
		author := &domain.User{}
		err := rows.Scan(
			&post.ID, &post.Title, &post.Content, &post.ThumbnailURL, &post.AuthorID, &post.Status, &post.IsDeleted, &post.GroupID, &post.CreatedAt, &post.UpdatedAt,
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

func (r *PostRepository) GetUserGroupCount(userID int) (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM groups WHERE owner_id = $1 AND is_active = true", userID).Scan(&count)
	return count, err
}

func (r *PostRepository) UpdateGroup(group *domain.Group) error {
	query := `UPDATE groups SET name = $1, description = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3 AND owner_id = $4`
	_, err := r.db.Exec(query, group.Name, group.Description, group.ID, group.OwnerID)
	return err
}

func (r *PostRepository) DeleteGroup(groupID int) error {
	// Start transaction for cascading delete
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	// Delete group members first
	_, err = tx.Exec("DELETE FROM group_members WHERE group_id = $1", groupID)
	if err != nil {
		return err
	}
	
	// Update posts to remove group association (set group_id to NULL)
	_, err = tx.Exec("UPDATE posts SET group_id = NULL WHERE group_id = $1", groupID)
	if err != nil {
		return err
	}
	
	// Delete the group
	_, err = tx.Exec("DELETE FROM groups WHERE id = $1", groupID)
	if err != nil {
		return err
	}
	
	return tx.Commit()
}

func (r *PostRepository) GetGroupMembers(groupID int) ([]domain.User, error) {
	query := `
		SELECT u.id, u.email, u.password_hash, u.display_name, u.bio, u.role, u.subscription_status, 
		       u.stripe_customer_id, u.is_active, u.email_verified, u.created_at, u.updated_at
		FROM users u
		INNER JOIN group_members gm ON u.id = gm.user_id
		WHERE gm.group_id = $1 AND u.is_active = true
		ORDER BY u.display_name`
	
	rows, err := r.db.Query(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var members []domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.ID, &user.Email, &user.PasswordHash, &user.DisplayName, &user.Bio,
			&user.Role, &user.SubscriptionStatus, &user.StripeCustomerID, 
			&user.IsActive, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		members = append(members, user)
	}
	
	return members, nil
}

func (r *PostRepository) RemoveGroupMember(groupID, userID int) error {
	_, err := r.db.Exec("DELETE FROM group_members WHERE group_id = $1 AND user_id = $2", groupID, userID)
	return err
}
