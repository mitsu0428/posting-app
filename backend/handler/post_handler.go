package handler

import (
	"encoding/json"
	"net/http"
	"posting-app/domain"
	"strconv"

	"github.com/gorilla/mux"
)

type CreatePostRequest struct {
	Title        string  `json:"title"`
	Content      string  `json:"content"`
	ThumbnailURL *string `json:"thumbnail_url,omitempty"`
}

type CreateReplyRequest struct {
	Content     string `json:"content"`
	IsAnonymous bool   `json:"is_anonymous"`
}

type PostListResponse struct {
	Posts []*domain.Post `json:"posts"`
	Total int            `json:"total"`
}

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(UserContextKey).(*domain.User)

	var req CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	post, err := h.postUsecase.CreatePost(user.ID, req.Title, req.Content, req.ThumbnailURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func (h *Handler) GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	post, err := h.postUsecase.GetPost(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func (h *Handler) ListPosts(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	limit := 20
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	posts, total, err := h.postUsecase.ListPosts(status, page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := PostListResponse{
		Posts: posts,
		Total: total,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetUserPosts(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(UserContextKey).(*domain.User)

	posts, err := h.postUsecase.GetUserPosts(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (h *Handler) CreateReply(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(UserContextKey).(*domain.User)
	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var req CreateReplyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	reply, err := h.postUsecase.CreateReply(postID, user.ID, req.Content, req.IsAnonymous)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(reply)
}

func (h *Handler) GetReplies(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	replies, err := h.postUsecase.GetReplies(postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(replies)
}