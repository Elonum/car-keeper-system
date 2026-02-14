package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/carkeeper/backend/config"
	"github.com/carkeeper/backend/database"
	"github.com/carkeeper/backend/internal/handler"
	authMiddleware "github.com/carkeeper/backend/internal/middleware"
	"github.com/carkeeper/backend/internal/repository"
	"github.com/carkeeper/backend/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.New(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Database connection established")

	repos := repository.New(db)
	services := service.New(repos, cfg)
	handlers := handler.New(services, cfg)
	router := setupRouter(handlers, cfg)

	server := &http.Server{
		Addr:         cfg.Server.Address(),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server starting on %s", cfg.Server.Address())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func setupRouter(handlers *handler.Handler, cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Route("/api", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", handlers.Register)
			r.Post("/login", handlers.Login)
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.AuthMiddleware(handlers.Services().Auth))
				r.Get("/me", handlers.GetMe)
			})
		})

		r.Route("/catalog", func(r chi.Router) {
			r.Get("/brands", handlers.GetBrands)
			r.Get("/models", handlers.GetModels)
			r.Get("/generations", handlers.GetGenerations)
			r.Get("/trims", handlers.GetTrims)
			r.Get("/trims/{id}", handlers.GetTrim)
			r.Get("/engine-types", handlers.GetEngineTypes)
			r.Get("/transmissions", handlers.GetTransmissions)
			r.Get("/drive-types", handlers.GetDriveTypes)
		})

		r.Route("/configurator", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.AuthMiddleware(handlers.Services().Auth))
				r.Get("/colors", handlers.GetColors)
				r.Get("/options", handlers.GetOptions)
				r.Post("/configurations", handlers.CreateConfiguration)
				r.Get("/configurations/{id}", handlers.GetConfiguration)
				r.Put("/configurations/{id}", handlers.UpdateConfiguration)
				r.Delete("/configurations/{id}", handlers.DeleteConfiguration)
			})
		})

		r.Route("/orders", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.AuthMiddleware(handlers.Services().Auth))
				r.Post("/", handlers.CreateOrder)
				r.Get("/", handlers.GetUserOrders)
				r.Get("/{id}", handlers.GetOrder)
				r.Patch("/{id}/status", handlers.UpdateOrderStatus)
			})
		})

		r.Route("/service", func(r chi.Router) {
			r.Get("/types", handlers.GetServiceTypes)
			r.Get("/branches", handlers.GetBranches)
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.AuthMiddleware(handlers.Services().Auth))
				r.Get("/user-cars", handlers.GetUserCars)
				r.Post("/appointments", handlers.CreateAppointment)
				r.Get("/appointments", handlers.GetUserAppointments)
				r.Get("/appointments/{id}", handlers.GetAppointment)
				r.Patch("/appointments/{id}/cancel", handlers.CancelAppointment)
			})
		})

		r.Route("/news", func(r chi.Router) {
			r.Get("/", handlers.GetNews)
			r.Get("/{id}", handlers.GetNewsByID)
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.AuthMiddleware(handlers.Services().Auth))
				r.Post("/", handlers.CreateNews)
				r.Put("/{id}", handlers.UpdateNews)
				r.Patch("/{id}/publish", handlers.PublishNews)
				r.Patch("/{id}/unpublish", handlers.UnpublishNews)
				r.Delete("/{id}", handlers.DeleteNews)
			})
		})

		r.Route("/profile", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.AuthMiddleware(handlers.Services().Auth))
				r.Get("/cars", handlers.GetUserCars)
				r.Post("/cars", handlers.CreateUserCar)
				r.Get("/cars/{id}", handlers.GetUserCar)
				r.Get("/configurations", handlers.GetUserConfigurations)
			})
		})
	})

	return r
}
