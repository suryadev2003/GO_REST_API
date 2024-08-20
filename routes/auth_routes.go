package routes

import (
	"Student_RESTAPI/controllers"
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func AuthRoutes(r *mux.Router, db *sql.DB) {
	r.Handle("/register", negroni.New(
		negroni.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controllers.Register(w, r, db)
		})),
	)).Methods("POST")

	r.Handle("/login", negroni.New(
		negroni.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controllers.Login(w, r, db)
		})),
	)).Methods("POST")
}
