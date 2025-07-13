package usecase

import (
	"errors"
	"fmt"
	"log/slog"

	"posting-app/domain"
	"posting-app/repository"
)

type PostUsecase struct {
	postRepo *repository.PostRepository
	userRepo *repository.UserRepository
}

func NewPostUsecase(postRepo *repository.PostRepository, userRepo *repository.UserRepository) *PostUsecase {
	return &PostUsecase{
		postRepo: postRepo,
		userRepo: userRepo,
	}
}

func (u *PostUsecase) CreatePost(userID int, title, content string, thumbnailURL *string) (*domain.Post, error) {
	// Check if user has active subscription
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.SubscriptionStatus != domain.UserSubscriptionStatusActive {
		return nil, errors.New("active subscription required to create posts")
	}

	post := &domain.Post{
		Title:        title,
		Content:      content,
		ThumbnailURL: thumbnailURL,
		AuthorID:     userID,
		Status:       domain.PostStatusPending, // Requires admin approval
	}

	err = u.postRepo.Create(post)
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	// Get the post with author info
	createdPost, err := u.postRepo.GetByID(post.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get created post: %w", err)
	}

	slog.Info("Post created successfully", "post_id", post.ID, "user_id", userID)
	return createdPost, nil
}

func (u *PostUsecase) UpdatePost(userID, postID int, title, content string, thumbnailURL *string) (*domain.Post, error) {
	post, err := u.postRepo.GetByID(postID)
	if err != nil {
		return nil, errors.New("post not found")
	}

	if post.AuthorID != userID {
		return nil, errors.New("you can only edit your own posts")
	}

	// Can only edit if post is pending
	if post.Status != domain.PostStatusPending {
		return nil, errors.New("can only edit posts that are pending approval")
	}

	post.Title = title
	post.Content = content
	post.ThumbnailURL = thumbnailURL

	err = u.postRepo.Update(post)
	if err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	// Get updated post
	updatedPost, err := u.postRepo.GetByID(postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated post: %w", err)
	}

	slog.Info("Post updated successfully", "post_id", postID, "user_id", userID)
	return updatedPost, nil
}

func (u *PostUsecase) DeletePost(userID, postID int) error {
	post, err := u.postRepo.GetByID(postID)
	if err != nil {
		return errors.New("post not found")
	}

	if post.AuthorID != userID {
		return errors.New("you can only delete your own posts")
	}

	err = u.postRepo.Delete(postID)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	slog.Info("Post deleted successfully", "post_id", postID, "user_id", userID)
	return nil
}

func (u *PostUsecase) GetPost(postID int) (*domain.Post, error) {
	post, err := u.postRepo.GetByID(postID)
	if err != nil {
		return nil, errors.New("post not found")
	}

	// Only return approved posts or posts to their authors
	if post.Status != domain.PostStatusApproved {
		return nil, errors.New("post not available")
	}

	return post, nil
}

func (u *PostUsecase) GetApprovedPosts(page, limit int) ([]*domain.Post, int, error) {
	posts, total, err := u.postRepo.GetApproved(page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get posts: %w", err)
	}

	return posts, total, nil
}

func (u *PostUsecase) GetUserPosts(userID, page, limit int) ([]*domain.Post, int, error) {
	posts, total, err := u.postRepo.GetByUserID(userID, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user posts: %w", err)
	}

	return posts, total, nil
}

func (u *PostUsecase) CreateReply(userID, postID int, content string, isAnonymous bool) (*domain.Reply, error) {
	// Check if user has active subscription
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.SubscriptionStatus != domain.UserSubscriptionStatusActive {
		return nil, errors.New("active subscription required to create replies")
	}

	// Check if post exists and is approved
	post, err := u.postRepo.GetByID(postID)
	if err != nil {
		return nil, errors.New("post not found")
	}

	if post.Status != domain.PostStatusApproved {
		return nil, errors.New("can only reply to approved posts")
	}

	var authorID *int
	if !isAnonymous {
		authorID = &userID
	}

	reply := &domain.Reply{
		Content:     content,
		PostID:      postID,
		AuthorID:    authorID,
		IsAnonymous: isAnonymous,
	}

	err = u.postRepo.CreateReply(reply)
	if err != nil {
		return nil, fmt.Errorf("failed to create reply: %w", err)
	}

	// Set author info if not anonymous
	if !isAnonymous {
		reply.Author = user
	}

	slog.Info("Reply created successfully", "reply_id", reply.ID, "post_id", postID, "user_id", userID, "anonymous", isAnonymous)
	return reply, nil
}

// Admin functions
func (u *PostUsecase) GetPostsForAdmin(page, limit int, status *domain.PostStatus) ([]*domain.Post, int, error) {
	posts, total, err := u.postRepo.GetForAdmin(page, limit, status)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get posts for admin: %w", err)
	}

	return posts, total, nil
}

func (u *PostUsecase) ApprovePost(postID int) error {
	post, err := u.postRepo.GetByID(postID)
	if err != nil {
		return errors.New("post not found")
	}

	if post.Status == domain.PostStatusApproved {
		return errors.New("post is already approved")
	}

	err = u.postRepo.UpdateStatus(postID, domain.PostStatusApproved)
	if err != nil {
		return fmt.Errorf("failed to approve post: %w", err)
	}

	slog.Info("Post approved successfully", "post_id", postID)
	return nil
}

func (u *PostUsecase) RejectPost(postID int) error {
	post, err := u.postRepo.GetByID(postID)
	if err != nil {
		return errors.New("post not found")
	}

	if post.Status == domain.PostStatusRejected {
		return errors.New("post is already rejected")
	}

	err = u.postRepo.UpdateStatus(postID, domain.PostStatusRejected)
	if err != nil {
		return fmt.Errorf("failed to reject post: %w", err)
	}

	slog.Info("Post rejected successfully", "post_id", postID)
	return nil
}
