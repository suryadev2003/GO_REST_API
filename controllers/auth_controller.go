package controllers

import (
	"Student_RESTAPI/models"
	"Student_RESTAPI/utils"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO users (UserID,Password) VALUES (?, ?)"
	_, err = db.Exec(query, user.UserID, user.Password)
	if err != nil {
		log.Printf("Error inserting user into database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	message := map[string]string{"message": "User registered successfully"}
	data, err := json.Marshal(message)
	if err != nil {
		http.Error(w, "Internal Server Error: Unable to marshal response", http.StatusInternalServerError)
		return
	}

	err = writeToFile(data)
	if err != nil {
		http.Error(w, "Internal Server Error: Unable to write to file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}

func Login(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var storedPassword string
	query := "SELECT Password FROM users WHERE UserID = ?"
	err = db.QueryRow(query, user.UserID).Scan(&storedPassword)
	if err != nil || storedPassword != user.Password {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateToken(user.UserID)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	tokenMessage := map[string]string{"token": token}
	data, err := json.Marshal(tokenMessage)
	if err != nil {
		http.Error(w, "Internal Server Error: Unable to marshal response", http.StatusInternalServerError)
		return
	}

	err = writeToFile(data)
	if err != nil {
		http.Error(w, "Internal Server Error: Unable to write to file", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tokenMessage)
}
