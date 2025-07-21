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

func (u *PostUsecase) CreatePost(userID int, title, content string, thumbnailURL *string, categoryIDs []int, groupID *int) (*domain.Post, error) {
	// Check if user has active subscription
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.SubscriptionStatus != domain.UserSubscriptionStatusActive {
		return nil, errors.New("active subscription required to create posts")
	}

	// Validate category limit
	if len(categoryIDs) > 5 {
		return nil, errors.New("maximum 5 categories allowed")
	}

	// If groupID is provided, check if user is member of the group
	if groupID != nil {
		isMember, err := u.postRepo.IsGroupMember(*groupID, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to check group membership: %w", err)
		}
		if !isMember {
			return nil, errors.New("you are not a member of this group")
		}
	}

	post := &domain.Post{
		Title:        title,
		Content:      content,
		ThumbnailURL: thumbnailURL,
		AuthorID:     userID,
		Status:       domain.PostStatusPending, // Requires admin approval
		GroupID:      groupID,
	}

	err = u.postRepo.Create(post)
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	// Add categories if provided
	if len(categoryIDs) > 0 {
		err = u.postRepo.AddPostCategories(post.ID, categoryIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to add categories: %w", err)
		}
	}

	// Get the post with author info
	createdPost, err := u.postRepo.GetByID(post.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get created post: %w", err)
	}

	slog.Info("Post created successfully", "post_id", post.ID, "user_id", userID)
	return createdPost, nil
}

func (u *PostUsecase) UpdatePost(userID, postID int, title, content string, thumbnailURL *string, categoryIDs []int, groupID *int) (*domain.Post, error) {
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

	// Validate category limit
	if len(categoryIDs) > 5 {
		return nil, errors.New("maximum 5 categories allowed")
	}

	// If groupID is provided, check if user is member of the group
	if groupID != nil {
		isMember, err := u.postRepo.IsGroupMember(*groupID, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to check group membership: %w", err)
		}
		if !isMember {
			return nil, errors.New("you are not a member of this group")
		}
	}

	post.Title = title
	post.Content = content
	post.ThumbnailURL = thumbnailURL
	post.GroupID = groupID

	err = u.postRepo.Update(post)
	if err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	// Update categories
	err = u.postRepo.AddPostCategories(post.ID, categoryIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to update categories: %w", err)
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

	// Check if user is admin or post author
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if user.Role != domain.UserRoleAdmin && post.AuthorID != userID {
		return errors.New("you can only delete your own posts")
	}

	err = u.postRepo.Delete(postID)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	slog.Info("Post deleted successfully", "post_id", postID, "user_id", userID, "is_admin", user.Role == domain.UserRoleAdmin)
	return nil
}

func (u *PostUsecase) GetPost(postID int, userID *int) (*domain.Post, error) {
	post, err := u.postRepo.GetByID(postID)
	if err != nil {
		return nil, errors.New("post not found")
	}

	// Only return approved posts
	if post.Status != domain.PostStatusApproved {
		return nil, errors.New("post not available")
	}

	// If it's a group post, check if user is member of the group
	if post.GroupID != nil {
		if userID == nil {
			return nil, errors.New("authentication required")
		}
		isMember, err := u.postRepo.IsGroupMember(*post.GroupID, *userID)
		if err != nil {
			return nil, fmt.Errorf("failed to check group membership: %w", err)
		}
		if !isMember {
			return nil, errors.New("you are not a member of this group")
		}
	}

	// Set like status if user is provided
	if userID != nil {
		isLiked, err := u.postRepo.IsLikedByUser(postID, *userID)
		if err != nil {
			return nil, fmt.Errorf("failed to check like status: %w", err)
		}
		post.IsLiked = isLiked
	}

	return post, nil
}

func (u *PostUsecase) GetApprovedPosts(page, limit int, userID *int) ([]*domain.Post, int, error) {
	posts, total, err := u.postRepo.GetApproved(page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get posts: %w", err)
	}

	// Set like status if user is provided
	if userID != nil {
		for _, post := range posts {
			isLiked, err := u.postRepo.IsLikedByUser(post.ID, *userID)
			if err != nil {
				return nil, 0, fmt.Errorf("failed to check like status: %w", err)
			}
			post.IsLiked = isLiked
		}
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

// Category functions
func (u *PostUsecase) GetAllCategories() ([]domain.Category, error) {
	categories, err := u.postRepo.GetAllCategories()
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	return categories, nil
}

func (u *PostUsecase) CreateCategory(name, description, color string) (*domain.Category, error) {
	category := &domain.Category{
		Name:        name,
		Description: description,
		Color:       color,
	}
	
	err := u.postRepo.CreateCategory(category)
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}
	
	slog.Info("Category created successfully", "category_id", category.ID, "name", name)
	return category, nil
}

// Like functions
func (u *PostUsecase) ToggleLike(postID, userID int) error {
	// Check if post exists and is approved
	post, err := u.postRepo.GetByID(postID)
	if err != nil {
		return errors.New("post not found")
	}
	
	if post.Status != domain.PostStatusApproved {
		return errors.New("can only like approved posts")
	}
	
	// If it's a group post, check if user is member of the group
	if post.GroupID != nil {
		isMember, err := u.postRepo.IsGroupMember(*post.GroupID, userID)
		if err != nil {
			return fmt.Errorf("failed to check group membership: %w", err)
		}
		if !isMember {
			return errors.New("you are not a member of this group")
		}
	}
	
	err = u.postRepo.ToggleLike(postID, userID)
	if err != nil {
		return fmt.Errorf("failed to toggle like: %w", err)
	}
	
	slog.Info("Like toggled successfully", "post_id", postID, "user_id", userID)
	return nil
}

// Group functions
func (u *PostUsecase) CreateGroup(userID int, name, description string) (*domain.Group, error) {
	// Check if user has reached the limit of 3 groups
	groupCount, err := u.postRepo.GetUserGroupCount(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user group count: %w", err)
	}
	
	if groupCount >= 3 {
		return nil, errors.New("maximum 3 groups allowed per user")
	}
	
	group := &domain.Group{
		Name:        name,
		Description: description,
		OwnerID:     userID,
	}
	
	err = u.postRepo.CreateGroup(group)
	if err != nil {
		return nil, fmt.Errorf("failed to create group: %w", err)
	}
	
	slog.Info("Group created successfully", "group_id", group.ID, "user_id", userID, "name", name)
	return group, nil
}

func (u *PostUsecase) GetUserGroups(userID int) ([]domain.Group, error) {
	groups, err := u.postRepo.GetUserGroups(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user groups: %w", err)
	}
	return groups, nil
}

func (u *PostUsecase) AddGroupMember(groupID, ownerID, userID int) error {
	// Check if requester is the group owner
	groups, err := u.postRepo.GetUserGroups(ownerID)
	if err != nil {
		return fmt.Errorf("failed to get user groups: %w", err)
	}
	
	isOwner := false
	for _, group := range groups {
		if group.ID == groupID && group.OwnerID == ownerID {
			isOwner = true
			break
		}
	}
	
	if !isOwner {
		return errors.New("only group owner can add members")
	}
	
	err = u.postRepo.AddGroupMember(groupID, userID)
	if err != nil {
		return fmt.Errorf("failed to add group member: %w", err)
	}
	
	slog.Info("Group member added successfully", "group_id", groupID, "user_id", userID, "owner_id", ownerID)
	return nil
}

func (u *PostUsecase) GetGroupPosts(groupID, userID, page, limit int) ([]*domain.Post, int, error) {
	posts, total, err := u.postRepo.GetGroupPosts(groupID, userID, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get group posts: %w", err)
	}
	
	// Set like status
	for _, post := range posts {
		isLiked, err := u.postRepo.IsLikedByUser(post.ID, userID)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to check like status: %w", err)
		}
		post.IsLiked = isLiked
	}
	
	return posts, total, nil
}

func (u *PostUsecase) UpdateGroup(groupID, ownerID int, name, description string) (*domain.Group, error) {
	// Check if requester is the group owner
	groups, err := u.postRepo.GetUserGroups(ownerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user groups: %w", err)
	}
	
	var targetGroup *domain.Group
	for _, group := range groups {
		if group.ID == groupID && group.OwnerID == ownerID {
			targetGroup = &group
			break
		}
	}
	
	if targetGroup == nil {
		return nil, errors.New("only group owner can update group")
	}
	
	// Update group details
	targetGroup.Name = name
	targetGroup.Description = description
	
	err = u.postRepo.UpdateGroup(targetGroup)
	if err != nil {
		return nil, fmt.Errorf("failed to update group: %w", err)
	}
	
	slog.Info("Group updated successfully", "group_id", groupID, "owner_id", ownerID, "name", name)
	return targetGroup, nil
}

func (u *PostUsecase) DeleteGroup(groupID, ownerID int) error {
	// Check if requester is the group owner
	groups, err := u.postRepo.GetUserGroups(ownerID)
	if err != nil {
		return fmt.Errorf("failed to get user groups: %w", err)
	}
	
	var targetGroup *domain.Group
	for _, group := range groups {
		if group.ID == groupID && group.OwnerID == ownerID {
			targetGroup = &group
			break
		}
	}
	
	if targetGroup == nil {
		return errors.New("only group owner can delete group")
	}
	
	// Delete group (this will cascade delete members and posts)
	err = u.postRepo.DeleteGroup(groupID)
	if err != nil {
		return fmt.Errorf("failed to delete group: %w", err)
	}
	
	slog.Info("Group deleted successfully", "group_id", groupID, "owner_id", ownerID)
	return nil
}

func (u *PostUsecase) AddGroupMemberByDisplayName(groupID, ownerID int, displayName string) error {
	// Check if requester is the group owner
	groups, err := u.postRepo.GetUserGroups(ownerID)
	if err != nil {
		return fmt.Errorf("failed to get user groups: %w", err)
	}
	
	isOwner := false
	for _, group := range groups {
		if group.ID == groupID && group.OwnerID == ownerID {
			isOwner = true
			break
		}
	}
	
	if !isOwner {
		return errors.New("only group owner can add members")
	}
	
	// Find user by display name
	user, err := u.userRepo.GetByDisplayName(displayName)
	if err != nil {
		return fmt.Errorf("user with display name '%s' not found", displayName)
	}
	
	// Check if user is already a member
	isMember, err := u.postRepo.IsGroupMember(groupID, user.ID)
	if err != nil {
		return fmt.Errorf("failed to check group membership: %w", err)
	}
	
	if isMember {
		return fmt.Errorf("user '%s' is already a member of this group", displayName)
	}
	
	err = u.postRepo.AddGroupMember(groupID, user.ID)
	if err != nil {
		return fmt.Errorf("failed to add group member: %w", err)
	}
	
	slog.Info("Group member added successfully", "group_id", groupID, "user_id", user.ID, "display_name", displayName, "owner_id", ownerID)
	return nil
}

func (u *PostUsecase) SearchUsersByDisplayName(query string) ([]domain.User, error) {
	if len(query) < 2 {
		return nil, errors.New("search query must be at least 2 characters")
	}
	
	users, err := u.userRepo.SearchByDisplayName(query)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}
	
	return users, nil
}

func (u *PostUsecase) GetGroupMembers(groupID, requestUserID int) ([]domain.User, error) {
	// Check if requester is a member of the group or owner
	groups, err := u.postRepo.GetUserGroups(requestUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user groups: %w", err)
	}
	
	var isOwnerOrMember bool
	for _, group := range groups {
		if group.ID == groupID {
			isOwnerOrMember = true
			break
		}
	}
	
	// Also check if user is a member
	if !isOwnerOrMember {
		isMember, err := u.postRepo.IsGroupMember(groupID, requestUserID)
		if err != nil {
			return nil, fmt.Errorf("failed to check group membership: %w", err)
		}
		isOwnerOrMember = isMember
	}
	
	if !isOwnerOrMember {
		return nil, errors.New("only group members can view member list")
	}
	
	members, err := u.postRepo.GetGroupMembers(groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group members: %w", err)
	}
	
	return members, nil
}

func (u *PostUsecase) RemoveGroupMember(groupID, ownerID, memberID int) error {
	// Check if requester is the group owner
	groups, err := u.postRepo.GetUserGroups(ownerID)
	if err != nil {
		return fmt.Errorf("failed to get user groups: %w", err)
	}
	
	var isOwner bool
	for _, group := range groups {
		if group.ID == groupID && group.OwnerID == ownerID {
			isOwner = true
			break
		}
	}
	
	if !isOwner {
		return errors.New("only group owner can remove members")
	}
	
	// Prevent owner from removing themselves (use LeaveGroup instead)
	if ownerID == memberID {
		return errors.New("group owner cannot remove themselves, use delete group instead")
	}
	
	// Check if member exists in group
	isMember, err := u.postRepo.IsGroupMember(groupID, memberID)
	if err != nil {
		return fmt.Errorf("failed to check group membership: %w", err)
	}
	
	if !isMember {
		return errors.New("user is not a member of this group")
	}
	
	err = u.postRepo.RemoveGroupMember(groupID, memberID)
	if err != nil {
		return fmt.Errorf("failed to remove group member: %w", err)
	}
	
	slog.Info("Group member removed successfully", "group_id", groupID, "member_id", memberID, "owner_id", ownerID)
	return nil
}

func (u *PostUsecase) LeaveGroup(groupID, userID int) error {
	// Check if user is a member of the group
	isMember, err := u.postRepo.IsGroupMember(groupID, userID)
	if err != nil {
		return fmt.Errorf("failed to check group membership: %w", err)
	}
	
	if !isMember {
		return errors.New("user is not a member of this group")
	}
	
	// Check if user is the owner
	groups, err := u.postRepo.GetUserGroups(userID)
	if err != nil {
		return fmt.Errorf("failed to get user groups: %w", err)
	}
	
	for _, group := range groups {
		if group.ID == groupID && group.OwnerID == userID {
			return errors.New("group owner cannot leave group, delete group instead")
		}
	}
	
	err = u.postRepo.RemoveGroupMember(groupID, userID)
	if err != nil {
		return fmt.Errorf("failed to leave group: %w", err)
	}
	
	slog.Info("User left group successfully", "group_id", groupID, "user_id", userID)
	return nil
}
