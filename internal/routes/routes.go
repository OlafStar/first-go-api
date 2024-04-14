package routes

import (
	"database/sql"
	"net/http"

	"example.com/jobboard/internal/api/user"
	"example.com/jobboard/internal/api/advertisment"
	"example.com/jobboard/internal/jwt"
	"example.com/jobboard/internal/middleware"
)

func NewRouter(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	setupUserAPI(mux, db)
	setupAdvertismentAPI(mux, db)

	return mux
}

func setupUserAPI(mux *http.ServeMux, db *sql.DB) {
	setupLoginRoute(mux, db)
	setupRegisterRoute(mux, db)
	setupIAMRoute(mux, db)
}

func setupAdvertismentAPI(mux *http.ServeMux, db *sql.DB) {
	setupCreateAdvertismentRoute(mux, db)
	setupGetAdvertismentRoute(mux, db)
}

func setupLoginRoute(mux *http.ServeMux, db *sql.DB) {
	mux.HandleFunc("/api/auth/login", middleware.Chain(func(w http.ResponseWriter, r *http.Request) {
		jwt.LoginHandler(db, w, r)
	}, middleware.Method("GET"), middleware.Logging()))
}

func setupRegisterRoute(mux *http.ServeMux, db *sql.DB) {
	mux.HandleFunc("/api/auth/register", middleware.Chain(func(w http.ResponseWriter, r *http.Request) {
		user.ReqisterUser(db, w, r)
	}, middleware.Method("POST"), middleware.Logging()))
}

func setupIAMRoute(mux *http.ServeMux, db *sql.DB) {
	mux.HandleFunc("/api/iam", middleware.Chain(jwt.ProtectedRequest(func(w http.ResponseWriter, r *http.Request) {
		user.Iam(db, w, r)
	}), middleware.Method("GET"), middleware.Logging()))
}

func setupCreateAdvertismentRoute(mux *http.ServeMux, db *sql.DB) {
	mux.HandleFunc("/api/advertisment/create", middleware.Chain(jwt.ProtectedRequest(func(w http.ResponseWriter, r *http.Request) {
		advertisment.CreateAdvertisment(db, w, r)
	}), middleware.Method("POST"), middleware.Logging()))
}

func setupGetAdvertismentRoute(mux *http.ServeMux, db *sql.DB) {
	mux.HandleFunc("/api/advertisment/get", middleware.Chain(jwt.ProtectedRequest(func(w http.ResponseWriter, r *http.Request) {
		advertisment.GetAdvertisment(db, w, r)
	}), middleware.Method("GET"), middleware.Logging()))
}