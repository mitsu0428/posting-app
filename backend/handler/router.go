package handler

import (
	"github.com/gorilla/mux"
)

func (h *Handler) SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/auth/register", h.CORSMiddleware(h.Register)).Methods("POST", "OPTIONS")
	api.HandleFunc("/auth/login", h.CORSMiddleware(h.Login)).Methods("POST", "OPTIONS")
	api.HandleFunc("/auth/logout", h.CORSMiddleware(h.AuthMiddleware(h.Logout))).Methods("POST", "OPTIONS")
	api.HandleFunc("/auth/forgot-password", h.CORSMiddleware(h.ForgotPassword)).Methods("POST", "OPTIONS")
	api.HandleFunc("/auth/reset-password", h.CORSMiddleware(h.ResetPassword)).Methods("POST", "OPTIONS")

	api.HandleFunc("/admin/login", h.CORSMiddleware(h.AdminLogin)).Methods("POST", "OPTIONS")

	api.HandleFunc("/posts", h.CORSMiddleware(h.AuthMiddleware(h.ListPosts))).Methods("GET", "OPTIONS")
	api.HandleFunc("/posts", h.CORSMiddleware(h.AuthMiddleware(h.CreatePost))).Methods("POST", "OPTIONS")
	api.HandleFunc("/posts/{id}", h.CORSMiddleware(h.AuthMiddleware(h.GetPost))).Methods("GET", "OPTIONS")
	api.HandleFunc("/posts/{id}/replies", h.CORSMiddleware(h.AuthMiddleware(h.GetReplies))).Methods("GET", "OPTIONS")
	api.HandleFunc("/posts/{id}/replies", h.CORSMiddleware(h.AuthMiddleware(h.CreateReply))).Methods("POST", "OPTIONS")

	api.HandleFunc("/me/posts", h.CORSMiddleware(h.AuthMiddleware(h.GetUserPosts))).Methods("GET", "OPTIONS")

	api.HandleFunc("/admin/posts", h.CORSMiddleware(h.AdminMiddleware(h.AdminListPosts))).Methods("GET", "OPTIONS")
	api.HandleFunc("/admin/posts/{id}/approve", h.CORSMiddleware(h.AdminMiddleware(h.ApprovePost))).Methods("POST", "OPTIONS")
	api.HandleFunc("/admin/posts/{id}/reject", h.CORSMiddleware(h.AdminMiddleware(h.RejectPost))).Methods("POST", "OPTIONS")
	api.HandleFunc("/admin/posts/{id}", h.CORSMiddleware(h.AdminMiddleware(h.DeletePost))).Methods("DELETE", "OPTIONS")
	api.HandleFunc("/admin/users", h.CORSMiddleware(h.AdminMiddleware(h.ListUsers))).Methods("GET", "OPTIONS")
	api.HandleFunc("/admin/users/{id}/deactivate", h.CORSMiddleware(h.AdminMiddleware(h.DeactivateUser))).Methods("POST", "OPTIONS")

	api.HandleFunc("/subscription/create-checkout-session", h.CORSMiddleware(h.AuthMiddleware(h.CreateCheckoutSession))).Methods("POST", "OPTIONS")
	api.HandleFunc("/subscription/webhook", h.CORSMiddleware(h.StripeWebhook)).Methods("POST", "OPTIONS")

	return r
}