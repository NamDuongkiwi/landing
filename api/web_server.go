package api

import (
	"github.com/gorilla/mux"
	auth_api "landing-page/api/auth-api"
	"landing-page/manager"
	user_manager "landing-page/manager/user-manager"
	"landing-page/pkg/middlewares"
)

type Server struct {
	UserManager user_manager.UserManager
	Middleware  middlewares.Middleware
}

func (s *Server) RunServer() {
	r := mux.NewRouter()
	r.Use(s.Middleware.CommonMiddleWare)
	authHandler := auth_api.AuthHandler{
		UserManager: s.UserManager,
	}
	r.HandleFunc("/login", authHandler.Login).Methods("POST")
	r.HandleFunc("/logout", authHandler.Logout).Methods("POST")
	//adminURL := r.PathPrefix("/admin").Subrouter()
}

func NewServer() (server Server) {
	db := manager.Connect()
	defer db.Close()
	server.UserManager = user_manager.NewUserManager(db)
	server.Middleware = middlewares.Middleware{
		UserManager: server.UserManager,
	}
	return server
}
