package router

import (
	"bhisham-api/internal/app/handlers"
	"bhisham-api/internal/app/middleware"
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
	url := "/v1/api" // Corrected variable declaration

	// User Routes
	mux.Handle(url+"/user/login", http.HandlerFunc(r.UserHandler.AuthenticateUser))
	mux.Handle(url+"/user/create-user", middleware.JWTAuthentication(http.HandlerFunc(r.UserHandler.CreateUser)))

	// Bhisham Routes
	mux.Handle(url+"/bhisham/create", middleware.JWTAuthentication(http.HandlerFunc(r.BhishamHandler.CreateBhisham)))

	// Dashboard Routes
	mux.Handle(url+"/dashboard/get-stats", middleware.JWTAuthentication(http.HandlerFunc(r.DashboardHandler.GetDashboardStats)))
	mux.Handle(url+"/dashboard/get-bhisham", middleware.JWTAuthentication(http.HandlerFunc(r.DashboardHandler.GetBhisham)))
	mux.Handle(url+"/dashboard/get-cubes", middleware.JWTAuthentication(http.HandlerFunc(r.DashboardHandler.GetChildCube)))
	mux.Handle(url+"/dashboard/get-kits", middleware.JWTAuthentication(http.HandlerFunc(r.DashboardHandler.GetChildKits)))

	return mux
}
