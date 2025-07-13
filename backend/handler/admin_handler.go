package handler

import (
	"encoding/json"
	"net/http"

	"posting-app/domain"
	"posting-app/repository"
	"posting-app/usecase"
)

type AdminHandler struct {
	authUsecase *usecase.AuthUsecase
	postUsecase *usecase.PostUsecase
	userRepo    *repository.UserRepository
}

func NewAdminHandler(
	authUsecase *usecase.AuthUsecase,
	postUsecase *usecase.PostUsecase,
	userRepo *repository.UserRepository,
) *AdminHandler {
	return &AdminHandler{
		authUsecase: authUsecase,
		postUsecase: postUsecase,
		userRepo:    userRepo,
	}
}

func (h *AdminHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := validate.Struct(req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, token, err := h.authUsecase.AdminLogin(req.Email, req.Password)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Remove sensitive data
	user.PasswordHash = ""

	writeJSON(w, http.StatusOK, LoginResponse{
		User:        user,
		AccessToken: token,
	})
}

func (h *AdminHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	page := getQueryInt(r, "page", 1)
	limit := getQueryInt(r, "limit", 20)

	var status *domain.PostStatus
	statusParam := r.URL.Query().Get("status")
	if statusParam != "" {
		switch statusParam {
		case "pending":
			s := domain.PostStatusPending
			status = &s
		case "approved":
			s := domain.PostStatusApproved
			status = &s
		case "rejected":
			s := domain.PostStatusRejected
			status = &s
		}
	}

	posts, total, err := h.postUsecase.GetPostsForAdmin(page, limit, status)
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

func (h *AdminHandler) ApprovePost(w http.ResponseWriter, r *http.Request) {
	postID, err := getIntParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	err = h.postUsecase.ApprovePost(postID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, Response{
		Message: "Post approved successfully",
	})
}

func (h *AdminHandler) RejectPost(w http.ResponseWriter, r *http.Request) {
	postID, err := getIntParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	err = h.postUsecase.RejectPost(postID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, Response{
		Message: "Post rejected successfully",
	})
}

func (h *AdminHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	page := getQueryInt(r, "page", 1)
	limit := getQueryInt(r, "limit", 20)

	users, total, err := h.userRepo.GetAll(page, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Remove sensitive data
	for _, user := range users {
		user.PasswordHash = ""
	}

	response := PaginatedResponse{
		Data:  users,
		Total: total,
		Page:  page,
		Limit: limit,
	}

	writeJSON(w, http.StatusOK, response)
}

func (h *AdminHandler) BanUser(w http.ResponseWriter, r *http.Request) {
	userID, err := getIntParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	err = h.userRepo.Ban(userID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, Response{
		Message: "User banned successfully",
	})
}
