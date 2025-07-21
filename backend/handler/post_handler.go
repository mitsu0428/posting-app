package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"posting-app/usecase"
)

type PostHandler struct {
	postUsecase *usecase.PostUsecase
}

func NewPostHandler(postUsecase *usecase.PostUsecase) *PostHandler {
	return &PostHandler{
		postUsecase: postUsecase,
	}
}

type CreatePostRequest struct {
	Title       string `json:"title" validate:"required,max=200"`
	Content     string `json:"content" validate:"required,max=5000"`
	CategoryIDs []int  `json:"category_ids" validate:"max=5"`
	GroupID     *int   `json:"group_id"`
}

type CreateCategoryRequest struct {
	Name        string `json:"name" validate:"required,max=100"`
	Description string `json:"description" validate:"max=500"`
	Color       string `json:"color" validate:"required,hexcolor"`
}

type CreateGroupRequest struct {
	Name        string `json:"name" validate:"required,max=100"`
	Description string `json:"description" validate:"max=500"`
}

type AddGroupMemberRequest struct {
	UserID int `json:"user_id" validate:"required"`
}

type CreateReplyRequest struct {
	Content     string `json:"content" validate:"required,max=2000"`
	IsAnonymous bool   `json:"is_anonymous"`
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		writeError(w, http.StatusBadRequest, "Failed to parse form")
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	categoryIDsStr := r.FormValue("category_ids")
	groupIDStr := r.FormValue("group_id")

	if title == "" || content == "" {
		writeError(w, http.StatusBadRequest, "Title and content are required")
		return
	}

	// Parse category IDs
	var categoryIDs []int
	if categoryIDsStr != "" {
		categoryIDsParts := strings.Split(categoryIDsStr, ",")
		for _, idStr := range categoryIDsParts {
			id, err := strconv.Atoi(strings.TrimSpace(idStr))
			if err != nil {
				writeError(w, http.StatusBadRequest, "Invalid category ID")
				return
			}
			categoryIDs = append(categoryIDs, id)
		}
	}

	// Parse group ID
	var groupID *int
	if groupIDStr != "" {
		id, err := strconv.Atoi(groupIDStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "Invalid group ID")
			return
		}
		groupID = &id
	}

	if len(title) > 200 {
		writeError(w, http.StatusBadRequest, "Title must be less than 200 characters")
		return
	}

	if len(content) > 5000 {
		writeError(w, http.StatusBadRequest, "Content must be less than 5000 characters")
		return
	}

	var thumbnailURL *string

	// Handle file upload if present
	file, header, err := r.FormFile("thumbnail")
	if err == nil {
		defer file.Close()

		// Check file size (5MB limit)
		if header.Size > 5<<20 {
			writeError(w, http.StatusBadRequest, "File size must be less than 5MB")
			return
		}

		// Check file type
		contentType := header.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "image/jpeg") && !strings.HasPrefix(contentType, "image/png") {
			writeError(w, http.StatusBadRequest, "Only JPEG and PNG images are allowed")
			return
		}

		// Save file
		filename := generateFilename(header.Filename)
		uploadPath := filepath.Join("uploads", filename)

		// Create uploads directory if it doesn't exist
		os.MkdirAll("uploads", 0755)

		outFile, err := os.Create(uploadPath)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to save file")
			return
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, file)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to save file")
			return
		}

		url := "/uploads/" + filename
		thumbnailURL = &url
	}

	post, err := h.postUsecase.CreatePost(user.ID, title, content, thumbnailURL, categoryIDs, groupID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, post)
}

func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	postID, err := getIntParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	// Parse multipart form
	err = r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		writeError(w, http.StatusBadRequest, "Failed to parse form")
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")

	if title == "" || content == "" {
		writeError(w, http.StatusBadRequest, "Title and content are required")
		return
	}

	if len(title) > 200 {
		writeError(w, http.StatusBadRequest, "Title must be less than 200 characters")
		return
	}

	if len(content) > 5000 {
		writeError(w, http.StatusBadRequest, "Content must be less than 5000 characters")
		return
	}

	var thumbnailURL *string

	// Handle file upload if present
	file, header, err := r.FormFile("thumbnail")
	if err == nil {
		defer file.Close()

		// Check file size (5MB limit)
		if header.Size > 5<<20 {
			writeError(w, http.StatusBadRequest, "File size must be less than 5MB")
			return
		}

		// Check file type
		contentType := header.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "image/jpeg") && !strings.HasPrefix(contentType, "image/png") {
			writeError(w, http.StatusBadRequest, "Only JPEG and PNG images are allowed")
			return
		}

		// Save file
		filename := generateFilename(header.Filename)
		uploadPath := filepath.Join("uploads", filename)

		// Create uploads directory if it doesn't exist
		os.MkdirAll("uploads", 0755)

		outFile, err := os.Create(uploadPath)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to save file")
			return
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, file)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to save file")
			return
		}

		url := "/uploads/" + filename
		thumbnailURL = &url
	}

	// Parse category IDs
	categoryIDsStr := r.FormValue("category_ids")
	var categoryIDs []int
	if categoryIDsStr != "" {
		categoryIDsParts := strings.Split(categoryIDsStr, ",")
		for _, idStr := range categoryIDsParts {
			id, err := strconv.Atoi(strings.TrimSpace(idStr))
			if err != nil {
				writeError(w, http.StatusBadRequest, "Invalid category ID")
				return
			}
			categoryIDs = append(categoryIDs, id)
		}
	}

	// Parse group ID
	groupIDStr := r.FormValue("group_id")
	var groupID *int
	if groupIDStr != "" {
		id, err := strconv.Atoi(groupIDStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "Invalid group ID")
			return
		}
		groupID = &id
	}

	post, err := h.postUsecase.UpdatePost(user.ID, postID, title, content, thumbnailURL, categoryIDs, groupID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, post)
}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	postID, err := getIntParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	err = h.postUsecase.DeletePost(user.ID, postID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	postID, err := getIntParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	var userID *int
	user := GetUserFromContext(r.Context())
	if user != nil {
		userID = &user.ID
	}

	post, err := h.postUsecase.GetPost(postID, userID)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, post)
}

func (h *PostHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	page := getQueryInt(r, "page", 1)
	limit := getQueryInt(r, "limit", 20)

	var userID *int
	user := GetUserFromContext(r.Context())
	if user != nil {
		userID = &user.ID
	}

	posts, total, err := h.postUsecase.GetApprovedPosts(page, limit, userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := PaginatedResponse{
		Data:  posts,
		Total: total,
		Page:  page,
		Limit: limit,
	}

	writeJSON(w, http.StatusOK, response)
}

func (h *PostHandler) GetUserPosts(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	page := getQueryInt(r, "page", 1)
	limit := getQueryInt(r, "limit", 20)

	posts, total, err := h.postUsecase.GetUserPosts(user.ID, page, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := PaginatedResponse{
		Data:  posts,
		Total: total,
		Page:  page,
		Limit: limit,
	}

	writeJSON(w, http.StatusOK, response)
}

func (h *PostHandler) CreateReply(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	postID, err := getIntParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	var req CreateReplyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := validate.Struct(req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	reply, err := h.postUsecase.CreateReply(user.ID, postID, req.Content, req.IsAnonymous)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, reply)
}

// Category handlers
func (h *PostHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.postUsecase.GetAllCategories()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, categories)
}

func (h *PostHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Only admin can create categories
	if user.Role != "admin" {
		writeError(w, http.StatusForbidden, "Admin access required")
		return
	}

	var req CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := validate.Struct(req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	category, err := h.postUsecase.CreateCategory(req.Name, req.Description, req.Color)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, category)
}

// Like handlers
func (h *PostHandler) ToggleLike(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	postID, err := getIntParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	err = h.postUsecase.ToggleLike(postID, user.ID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Group handlers
func (h *PostHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := validate.Struct(req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	group, err := h.postUsecase.CreateGroup(user.ID, req.Name, req.Description)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, group)
}

func (h *PostHandler) GetUserGroups(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	groups, err := h.postUsecase.GetUserGroups(user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, groups)
}

func (h *PostHandler) AddGroupMember(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	groupID, err := getIntParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid group ID")
		return
	}

	var req AddGroupMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := validate.Struct(req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.postUsecase.AddGroupMember(groupID, user.ID, req.UserID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *PostHandler) GetGroupPosts(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	groupID, err := getIntParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid group ID")
		return
	}

	page := getQueryInt(r, "page", 1)
	limit := getQueryInt(r, "limit", 20)

	posts, total, err := h.postUsecase.GetGroupPosts(groupID, user.ID, page, limit)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := PaginatedResponse{
		Data:  posts,
		Total: total,
		Page:  page,
		Limit: limit,
	}

	writeJSON(w, http.StatusOK, response)
}

func (h *PostHandler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	groupID, err := getIntParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid group ID")
		return
	}

	var req CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	group, err := h.postUsecase.UpdateGroup(groupID, user.ID, req.Name, req.Description)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, group)
}

func (h *PostHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}
	
	groupID, err := getIntParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid group ID")
		return
	}
	
	err = h.postUsecase.DeleteGroup(groupID, user.ID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	writeJSON(w, http.StatusOK, map[string]string{"message": "Group deleted successfully"})
}

func (h *PostHandler) AddGroupMemberByDisplayName(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	groupID, err := getIntParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid group ID")
		return
	}

	var req struct {
		DisplayName string `json:"display_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.DisplayName == "" {
		writeError(w, http.StatusBadRequest, "Display name is required")
		return
	}

	err = h.postUsecase.AddGroupMemberByDisplayName(groupID, user.ID, req.DisplayName)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Member added successfully"})
}

func (h *PostHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		writeError(w, http.StatusBadRequest, "Search query is required")
		return
	}

	users, err := h.postUsecase.SearchUsersByDisplayName(query)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Remove sensitive fields
	var safeUsers []map[string]interface{}
	for _, u := range users {
		safeUser := map[string]interface{}{
			"id":           u.ID,
			"display_name": u.DisplayName,
			"bio":          u.Bio,
		}
		safeUsers = append(safeUsers, safeUser)
	}

	writeJSON(w, http.StatusOK, safeUsers)
}

func (h *PostHandler) GetGroupMembers(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}
	
	groupID, err := getIntParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid group ID")
		return
	}
	
	members, err := h.postUsecase.GetGroupMembers(groupID, user.ID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	// Remove sensitive fields
	var safeMembers []map[string]interface{}
	for _, member := range members {
		safeMember := map[string]interface{}{
			"id":           member.ID,
			"display_name": member.DisplayName,
			"bio":          member.Bio,
		}
		safeMembers = append(safeMembers, safeMember)
	}
	
	writeJSON(w, http.StatusOK, safeMembers)
}

func (h *PostHandler) RemoveGroupMember(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}
	
	groupID, err := getIntParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid group ID")
		return
	}
	
	memberID, err := getIntParam(r, "memberId")
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid member ID")
		return
	}
	
	err = h.postUsecase.RemoveGroupMember(groupID, user.ID, memberID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	writeJSON(w, http.StatusOK, map[string]string{"message": "Member removed successfully"})
}

func (h *PostHandler) LeaveGroup(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}
	
	groupID, err := getIntParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid group ID")
		return
	}
	
	err = h.postUsecase.LeaveGroup(groupID, user.ID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	writeJSON(w, http.StatusOK, map[string]string{"message": "Left group successfully"})
}

func generateFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	return strconv.FormatInt(time.Now().UnixNano(), 10) + ext
}
