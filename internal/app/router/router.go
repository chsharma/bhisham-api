package router

import (
	"bhisham-api/internal/app/handlers"
	"bhisham-api/internal/middleware"
	"fmt"
	"net/http"
)

type Router struct {
	BhishamHandler   *handlers.BhishamHandler
	UserHandler      *handlers.UserHandler
	DashboardHandler *handlers.DashboardHandler
}

func NewRouter(gbhishamHandler *handlers.BhishamHandler, userHandler *handlers.UserHandler, dashboardHandler *handlers.DashboardHandler) *Router {
	return &Router{
		BhishamHandler:   gbhishamHandler,
		UserHandler:      userHandler,
		DashboardHandler: dashboardHandler,
	}
}

func (r *Router) RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	url := "/v1/api" // API base path

	// Corrected default handler
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome To Api")
	})

	// User Routes
	mux.Handle(url+"/user/login", http.HandlerFunc(r.UserHandler.AuthenticateUser))
	mux.Handle(url+"/user/create-user", middleware.JWTAuthentication(http.HandlerFunc(r.UserHandler.CreateUser)))

	// Bhisham Routes
	mux.Handle(url+"/bhisham/create", middleware.JWTAuthentication(http.HandlerFunc(r.BhishamHandler.CreateBhisham)))
	mux.Handle(url+"/bhisham/create-data", middleware.JWTAuthentication(http.HandlerFunc(r.BhishamHandler.CreateBhishamData)))
	mux.Handle(url+"/bhisham/update-data", middleware.JWTAuthentication(http.HandlerFunc(r.BhishamHandler.UpdateBhishamData)))

	// Dashboard Routes
	mux.Handle(url+"/dashboard/get-stats", middleware.JWTAuthentication(http.HandlerFunc(r.DashboardHandler.GetDashboardStats)))
	mux.Handle(url+"/dashboard/get-bhisham", middleware.JWTAuthentication(http.HandlerFunc(r.DashboardHandler.GetBhisham)))
	mux.Handle(url+"/dashboard/get-cubes", middleware.JWTAuthentication(http.HandlerFunc(r.DashboardHandler.GetChildCube)))
	mux.Handle(url+"/dashboard/get-kits", middleware.JWTAuthentication(http.HandlerFunc(r.DashboardHandler.GetChildKits)))

	return mux
}
