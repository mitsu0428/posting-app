package usecase

import (
	"fmt"
	"posting-app/domain"
)

type postUsecase struct {
	postRepo  domain.PostRepository
	replyRepo domain.ReplyRepository
}

func NewPostUsecase(postRepo domain.PostRepository, replyRepo domain.ReplyRepository) domain.PostUsecase {
	return &postUsecase{
		postRepo:  postRepo,
		replyRepo: replyRepo,
	}
}

func (u *postUsecase) CreatePost(userID int, title, content string, thumbnailURL *string) (*domain.Post, error) {
	post := &domain.Post{
		Title:        title,
		Content:      content,
		ThumbnailURL: thumbnailURL,
		UserID:       userID,
		Status:       "pending",
	}

	if err := u.postRepo.Create(post); err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return post, nil
}

func (u *postUsecase) GetPost(id int) (*domain.Post, error) {
	post, err := u.postRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("post not found: %w", err)
	}

	if post.Status != "approved" {
		return nil, fmt.Errorf("post not available")
	}

	return post, nil
}

func (u *postUsecase) ListPosts(status string, page, limit int) ([]*domain.Post, int, error) {
	if status == "" {
		status = "approved"
	}

	posts, total, err := u.postRepo.List(status, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list posts: %w", err)
	}

	return posts, total, nil
}

func (u *postUsecase) GetUserPosts(userID int) ([]*domain.Post, error) {
	posts, err := u.postRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user posts: %w", err)
	}

	return posts, nil
}

func (u *postUsecase) CreateReply(postID, userID int, content string, isAnonymous bool) (*domain.Reply, error) {
	post, err := u.postRepo.GetByID(postID)
	if err != nil {
		return nil, fmt.Errorf("post not found: %w", err)
	}

	if post.Status != "approved" {
		return nil, fmt.Errorf("cannot reply to unapproved post")
	}

	reply := &domain.Reply{
		PostID:      postID,
		Content:     content,
		UserID:      userID,
		IsAnonymous: isAnonymous,
	}

	if err := u.replyRepo.Create(reply); err != nil {
		return nil, fmt.Errorf("failed to create reply: %w", err)
	}

	return reply, nil
}

func (u *postUsecase) GetReplies(postID int) ([]*domain.Reply, error) {
	replies, err := u.replyRepo.GetByPostID(postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get replies: %w", err)
	}

	return replies, nil
}

type adminUsecase struct {
	postRepo domain.PostRepository
	userRepo domain.UserRepository
}

func NewAdminUsecase(postRepo domain.PostRepository, userRepo domain.UserRepository) domain.AdminUsecase {
	return &adminUsecase{
		postRepo: postRepo,
		userRepo: userRepo,
	}
}

func (u *adminUsecase) ListAllPosts(status string) ([]*domain.Post, error) {
	posts, _, err := u.postRepo.List(status, 1, 1000)
	if err != nil {
		return nil, fmt.Errorf("failed to list posts: %w", err)
	}

	return posts, nil
}

func (u *adminUsecase) ApprovePost(id int) error {
	if err := u.postRepo.Approve(id); err != nil {
		return fmt.Errorf("failed to approve post: %w", err)
	}

	return nil
}

func (u *adminUsecase) RejectPost(id int) error {
	if err := u.postRepo.Reject(id); err != nil {
		return fmt.Errorf("failed to reject post: %w", err)
	}

	return nil
}

func (u *adminUsecase) DeletePost(id int) error {
	if err := u.postRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}

func (u *adminUsecase) ListUsers() ([]*domain.User, error) {
	users, err := u.userRepo.List()
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}

func (u *adminUsecase) DeactivateUser(id int) error {
	if err := u.userRepo.Deactivate(id); err != nil {
		return fmt.Errorf("failed to deactivate user: %w", err)
	}

	return nil
}