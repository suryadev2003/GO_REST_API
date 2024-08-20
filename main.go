package main

import (
	"Student_RESTAPI/routes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func main() {
	InitDB()
	defer db.Close()

	r := mux.NewRouter()

	n := negroni.New()

	n.Use(negroni.NewLogger())
	n.Use(negroni.NewRecovery())

	n.UseHandler(r)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the API!"))
	}).Methods("GET")

	routes.AuthRoutes(r, db)
	routes.StudentRoutes(r, db)

	log.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", n); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
