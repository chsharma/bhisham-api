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
	HandheldHandler  *handlers.HandheldHandler
}

func NewRouter(gbhishamHandler *handlers.BhishamHandler, userHandler *handlers.UserHandler, dashboardHandler *handlers.DashboardHandler, handheldHandler *handlers.HandheldHandler) *Router {
	return &Router{
		BhishamHandler:   gbhishamHandler,
		UserHandler:      userHandler,
		DashboardHandler: dashboardHandler,
		HandheldHandler:  handheldHandler,
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
	mux.Handle(url+"/user/get-users", middleware.JWTAuthentication(http.HandlerFunc(r.UserHandler.GetUsers)))
	mux.Handle(url+"/user/update-password", middleware.JWTAuthentication(http.HandlerFunc(r.UserHandler.UpdatePassword)))
	mux.Handle(url+"/user/update-user", middleware.JWTAuthentication(http.HandlerFunc(r.UserHandler.UpdateUser)))

	// Bhisham Routes
	mux.Handle(url+"/bhisham/create", middleware.JWTAuthentication(http.HandlerFunc(r.BhishamHandler.CreateBhisham)))
	mux.Handle(url+"/bhisham/create-data", middleware.JWTAuthentication(http.HandlerFunc(r.BhishamHandler.CreateBhishamData)))
	mux.Handle(url+"/bhisham/update-data", middleware.JWTAuthentication(http.HandlerFunc(r.BhishamHandler.UpdateBhishamData)))
	mux.Handle(url+"/bhisham/update-mapping-data", middleware.JWTAuthentication(http.HandlerFunc(r.BhishamHandler.UpdateBhishamMapping)))

	mux.Handle(url+"/bhisham/add-bhisham-data", middleware.JWTAuthentication(http.HandlerFunc(r.BhishamHandler.AddItemData)))
	mux.Handle(url+"/bhisham/delete-bhisham-data", middleware.JWTAuthentication(http.HandlerFunc(r.BhishamHandler.DeleteItemData)))

	mux.Handle(url+"/bhisham/add-mapping-data", middleware.JWTAuthentication(http.HandlerFunc(r.BhishamHandler.AddItemData)))
	mux.Handle(url+"/bhisham/delete-mapping-data", middleware.JWTAuthentication(http.HandlerFunc(r.BhishamHandler.DeleteItemData)))

	mux.Handle(url+"/bhisham/mark-update-data", middleware.JWTAuthentication(http.HandlerFunc(r.BhishamHandler.MarkUpdateBhishamData)))
	mux.Handle(url+"/bhisham/close-bhisham", middleware.JWTAuthentication(http.HandlerFunc(r.BhishamHandler.CloseBhisham)))

	// Dashboard Routes
	mux.Handle(url+"/dashboard/get-stats", middleware.JWTAuthentication(http.HandlerFunc(r.DashboardHandler.GetDashboardStats)))
	mux.Handle(url+"/dashboard/get-bhisham", middleware.JWTAuthentication(http.HandlerFunc(r.DashboardHandler.GetBhisham)))
	mux.Handle(url+"/dashboard/get-cubes", middleware.JWTAuthentication(http.HandlerFunc(r.DashboardHandler.GetChildCube)))
	mux.Handle(url+"/dashboard/get-kits", middleware.JWTAuthentication(http.HandlerFunc(r.DashboardHandler.GetChildKits)))
	mux.Handle(url+"/dashboard/get-items", middleware.JWTAuthentication(http.HandlerFunc(r.DashboardHandler.GetKitItems)))
	mux.Handle(url+"/dashboard/get-mapping-items", middleware.JWTAuthentication(http.HandlerFunc(r.DashboardHandler.GetMappingKitItems)))

	mux.Handle(url+"/dashboard/data-update-type", middleware.JWTAuthentication(http.HandlerFunc(r.DashboardHandler.GetUpdateType)))
	mux.Handle(url+"/dashboard/get-mapp-data", middleware.JWTAuthentication(http.HandlerFunc(r.DashboardHandler.GetAllMappingBhishamData)))

	// Handheld Routes
	mux.Handle(url+"/handheld/get-all-data", middleware.JWTAuthentication(http.HandlerFunc(r.HandheldHandler.GetAllBhishamData)))
	mux.Handle(url+"/handheld/get-bhishamid", middleware.JWTAuthentication(http.HandlerFunc(r.HandheldHandler.GetBhishamID)))

	return mux
}
