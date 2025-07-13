package handler

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"posting-app/infrastructure"
)

type Handlers struct {
	Auth         *AuthHandler
	Post         *PostHandler
	Admin        *AdminHandler
	Subscription *SubscriptionHandler
	User         *UserHandler
}

func NewRouter(handlers *Handlers, jwtService *infrastructure.JWTService) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://*.vercel.app"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Public routes
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", handlers.Auth.Register)
		r.Post("/login", handlers.Auth.Login)
		r.Post("/forgot-password", handlers.Auth.ForgotPassword)
		r.Post("/reset-password", handlers.Auth.ResetPassword)
	})

	// Admin auth
	r.Post("/admin/login", handlers.Admin.Login)

	// Subscription webhook (public)
	r.Post("/subscription/webhook", handlers.Subscription.HandleWebhook)

	// Serve uploaded files
	r.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads/"))))

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware(jwtService))

		// Auth
		r.Post("/auth/logout", handlers.Auth.Logout)

		// User routes
		r.Route("/user", func(r chi.Router) {
			r.Get("/profile", handlers.User.GetProfile)
			r.Put("/profile", handlers.User.UpdateProfile)
			r.Post("/change-password", handlers.Auth.ChangePassword)
			r.Post("/deactivate", handlers.User.Deactivate)
			r.Get("/posts", handlers.Post.GetUserPosts)
		})

		// Post routes
		r.Route("/posts", func(r chi.Router) {
			r.Get("/", handlers.Post.GetPosts)
			r.Post("/", handlers.Post.CreatePost)
			r.Get("/{id}", handlers.Post.GetPost)
			r.Put("/{id}", handlers.Post.UpdatePost)
			r.Delete("/{id}", handlers.Post.DeletePost)
			r.Post("/{id}/replies", handlers.Post.CreateReply)
		})

		// Subscription routes
		r.Route("/subscription", func(r chi.Router) {
			r.Get("/status", handlers.Subscription.GetStatus)
			r.Post("/create-checkout-session", handlers.Subscription.CreateCheckoutSession)
		})

		// Admin routes
		r.Group(func(r chi.Router) {
			r.Use(AdminMiddleware)

			r.Route("/admin", func(r chi.Router) {
				r.Get("/posts", handlers.Admin.GetPosts)
				r.Post("/posts/{id}/approve", handlers.Admin.ApprovePost)
				r.Post("/posts/{id}/reject", handlers.Admin.RejectPost)
				r.Get("/users", handlers.Admin.GetUsers)
				r.Post("/users/{id}/ban", handlers.Admin.BanUser)
			})
		})
	})

	return r
}
