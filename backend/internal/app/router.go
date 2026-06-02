package app

import (
	"time"

	"github.com/carkeeper/backend/config"
	"github.com/carkeeper/backend/database"
	"github.com/carkeeper/backend/internal/handler"
	authMiddleware "github.com/carkeeper/backend/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// NewRouter builds the HTTP router used by the API server and integration tests.
func NewRouter(handlers *handler.Handler, cfg *config.Config, db *database.DB) *chi.Mux {
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
			r.Post("/logout", handlers.Logout)
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
