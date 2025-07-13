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
	Title   string `json:"title" validate:"required,max=200"`
	Content string `json:"content" validate:"required,max=5000"`
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

	post, err := h.postUsecase.CreatePost(user.ID, title, content, thumbnailURL)
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

	post, err := h.postUsecase.UpdatePost(user.ID, postID, title, content, thumbnailURL)
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

	post, err := h.postUsecase.GetPost(postID)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, post)
}

func (h *PostHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	page := getQueryInt(r, "page", 1)
	limit := getQueryInt(r, "limit", 20)

	posts, total, err := h.postUsecase.GetApprovedPosts(page, limit)
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

func generateFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	return strconv.FormatInt(time.Now().UnixNano(), 10) + ext
}
