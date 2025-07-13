package handler

import (
	"encoding/json"
	"net/http"

	"posting-app/repository"
)

type UserHandler struct {
	userRepo *repository.UserRepository
}

func NewUserHandler(userRepo *repository.UserRepository) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

type UpdateProfileRequest struct {
	DisplayName string  `json:"display_name" validate:"max=100"`
	Bio         *string `json:"bio" validate:"omitempty,max=500"`
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Get full user data from database
	fullUser, err := h.userRepo.GetByID(user.ID)
	if err != nil {
		writeError(w, http.StatusNotFound, "User not found")
		return
	}

	// Remove sensitive data
	fullUser.PasswordHash = ""

	writeJSON(w, http.StatusOK, fullUser)
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := validate.Struct(req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get current user data
	fullUser, err := h.userRepo.GetByID(user.ID)
	if err != nil {
		writeError(w, http.StatusNotFound, "User not found")
		return
	}

	// Update fields
	if req.DisplayName != "" {
		fullUser.DisplayName = req.DisplayName
	}
	fullUser.Bio = req.Bio

	err = h.userRepo.Update(fullUser)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to update profile")
		return
	}

	// Remove sensitive data
	fullUser.PasswordHash = ""

	writeJSON(w, http.StatusOK, fullUser)
}

func (h *UserHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	err := h.userRepo.Deactivate(user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to deactivate account")
		return
	}

	writeJSON(w, http.StatusOK, Response{
		Message: "Account deactivated successfully",
	})
}
