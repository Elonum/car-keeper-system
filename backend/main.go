package main

import (
	"context"
	"log"
	"log/slog"
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
	"github.com/carkeeper/backend/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	initLogging(cfg)

	db, err := database.New(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	slog.Info("database connection established")

	fileStore, err := storage.NewLocal(cfg.Storage.RootPath)
	if err != nil {
		log.Fatalf("Document storage: %v", err)
	}

	repos := repository.New(db)
	services := service.New(repos, cfg, fileStore)
	handlers := handler.New(services, cfg)
	router := setupRouter(handlers, cfg, db)

	server := &http.Server{
		Addr:         cfg.Server.Address(),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("server starting", "addr", cfg.Server.Address())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	slog.Info("server exited")
}

func initLogging(cfg *config.Config) {
	var h slog.Handler
	if cfg.Env == "production" {
		h = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	} else {
		h = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	}
	slog.SetDefault(slog.New(h))
}

func setupRouter(handlers *handler.Handler, cfg *config.Config, db *database.DB) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(authMiddleware.LimitRequestBody(cfg.Server.MaxJSONBodyBytes))
	r.Use(authMiddleware.RequestLogger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.CORSOrigins(),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", handler.Health(db))

	r.Route("/api", func(r chi.Router) {
		r.Get("/order-statuses", handlers.GetOrderStatuses)

		r.Route("/admin", func(r chi.Router) {
			r.Use(authMiddleware.AuthMiddleware(handlers.Services().Auth))
			r.Get("/roles", handlers.AdminListRoleDefinitions)
			r.Get("/orders", handlers.AdminListAllOrders)
			r.Get("/appointments", handlers.AdminListAllAppointments)
			r.Patch("/branches/{id}", handlers.AdminUpdateBranch)
			r.Route("/catalog", func(r chi.Router) {
				r.Route("/brands", func(r chi.Router) {
					r.Post("/", handlers.AdminCreateBrand)
					r.Patch("/{id}", handlers.AdminUpdateBrand)
					r.Delete("/{id}", handlers.AdminDeleteBrand)
				})
				r.Route("/models", func(r chi.Router) {
					r.Get("/", handlers.AdminListModels)
					r.Post("/", handlers.AdminCreateModel)
					r.Patch("/{id}", handlers.AdminUpdateModel)
					r.Delete("/{id}", handlers.AdminDeleteModel)
					r.Post("/{id}/image", handlers.AdminUploadModelImage)
				})
			})
			r.Route("/service", func(r chi.Router) {
				r.Route("/types", func(r chi.Router) {
					r.Post("/", handlers.AdminCreateServiceType)
					r.Patch("/{id}", handlers.AdminUpdateServiceType)
					r.Delete("/{id}", handlers.AdminDeleteServiceType)
				})
			})
			r.Route("/order-statuses", func(r chi.Router) {
				r.Get("/", handlers.AdminListOrderStatuses)
				r.Post("/", handlers.AdminCreateOrderStatus)
				r.Patch("/{id}", handlers.AdminUpdateOrderStatus)
				r.Delete("/{id}", handlers.AdminDeleteOrderStatus)
			})
		})

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
			r.Get("/models/{id}/image", handlers.GetModelImage)
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
				r.Get("/branches/{branchID}/availability", handlers.GetBranchAvailability)
				r.Get("/user-cars", handlers.GetUserCars)
				r.Post("/appointments", handlers.CreateAppointment)
				r.Get("/appointments", handlers.GetUserAppointments)
				r.Patch("/appointments/{id}/reschedule", handlers.RescheduleAppointment)
				r.Get("/appointments/{id}", handlers.GetAppointment)
				r.Patch("/appointments/{id}/cancel", handlers.CancelAppointment)
			})
		})

		r.Route("/news", func(r chi.Router) {
			r.Use(authMiddleware.OptionalAuthMiddleware(handlers.Services().Auth))
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
				r.Patch("/me", handlers.UpdateProfile)
				r.Post("/me/password", handlers.ChangePassword)
				r.Get("/cars", handlers.GetUserCars)
				r.Post("/cars", handlers.CreateUserCar)
				r.Delete("/cars/{id}", handlers.DeleteUserCar)
				r.Get("/cars/{id}", handlers.GetUserCar)
				r.Get("/configurations", handlers.GetUserConfigurations)
			})
		})

		r.Route("/documents", func(r chi.Router) {
			r.Use(authMiddleware.AuthMiddleware(handlers.Services().Auth))
			r.Post("/", handlers.CreateDocument)
			r.Get("/", handlers.ListDocuments)
			r.Get("/{documentID}/file", handlers.DownloadDocument)
			r.Get("/{documentID}", handlers.GetDocument)
			r.Delete("/{documentID}", handlers.DeleteDocument)
		})
	})

	return r
}
