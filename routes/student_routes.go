package routes

import (
	"Student_RESTAPI/controllers"
	"Student_RESTAPI/middleware"
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func StudentRoutes(r *mux.Router, db *sql.DB) {
	r.Handle("/students/create", negroni.New(
		negroni.HandlerFunc(middleware.AuthMiddleware),
		negroni.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controllers.CreateStudent(w, r, db)
		})),
	)).Methods("POST")

	r.Handle("/students/{id}", negroni.New(
		negroni.HandlerFunc(middleware.AuthMiddleware),
		negroni.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controllers.GetStudentByID(w, r, db)
		})),
	)).Methods("GET")

	r.Handle("/students/{id}", negroni.New(
		negroni.HandlerFunc(middleware.AuthMiddleware),
		negroni.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controllers.UpdateStudent(w, r, db)
		})),
	)).Methods("PUT")

	r.Handle("/students/{id}", negroni.New(
		negroni.HandlerFunc(middleware.AuthMiddleware),
		negroni.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controllers.DeleteStudent(w, r, db)
		})),
	)).Methods("DELETE")
}
