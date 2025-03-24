package main

import (
	"bhisham-api/config"
	"bhisham-api/internal/app/handlers"
	"bhisham-api/internal/app/repositories"
	"bhisham-api/internal/app/router"
	"bhisham-api/internal/app/services"
	"log"
	"net/http"
)

// CORS Middleware
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins (Change `*` to specific domain if needed)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// Initialize database
	database := config.ConnectDB()
	defer database.Close()

	// Initialize repositories
	bhishamRepo := &repositories.BhishamRepository{DB: database}
	userRepo := &repositories.UserRepository{DB: database}
	dashboardRepo := &repositories.DashboardRepository{DB: database}
	handheldRepo := &repositories.HandheldRepository{DB: database}

	// Initialize services
	bhishamService := &services.BhishamService{GameRepo: bhishamRepo}
	userService := &services.UserService{UserRepo: userRepo}
	dashboardService := &services.DashboardService{DashboardRepo: dashboardRepo}
	handheldService := &services.HandheldService{HandheldRepo: handheldRepo}
	// Initialize handlers
	bhishamHandler := &handlers.BhishamHandler{BhishamService: bhishamService}
	userHandler := &handlers.UserHandler{UserService: userService}
	dashboardHandler := &handlers.DashboardHandler{DashboardService: dashboardService}
	handheldHandler := &handlers.HandheldHandler{HandheldService: handheldService}

	appRouter := router.NewRouter(bhishamHandler, userHandler, dashboardHandler, handheldHandler)

	mux := appRouter.RegisterRoutes()

	handler := enableCORS(mux)

	// Start the HTTP server
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
