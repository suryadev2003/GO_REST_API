package controllers

import (
	"Student_RESTAPI/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func writeToFile(data []byte) error {

	file, err := os.OpenFile("output.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	_, err = file.WriteString("\n")
	return err
}

func CreateStudent(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var student models.Student

	err := json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		http.Error(w, "Bad Request: Unable to parse request body", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "Unauthorized: Unable to retrieve user ID from context", http.StatusUnauthorized)
		return
	}

	student.CreatedOn = time.Now()
	student.CreatedBy = userID

	err = student.Create(db)
	if err != nil {
		http.Error(w, "Internal Server Error: Unable to create student", http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(student)
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
}

func UpdateStudent(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var student models.Student

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}
	student.ID = id

	existingStudent, err := models.GetStudentByID(db, id)
	if err != nil {
		http.Error(w, "Internal Server Error: Unable to check student existence", http.StatusInternalServerError)
		return
	}

	if existingStudent == nil {
		response := map[string]string{"message": "Student with this ID not present"}
		data, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Internal Server Error: Unable to marshal response", http.StatusInternalServerError)
			return
		}

		err = writeToFile(data)
		if err != nil {
			http.Error(w, "Internal Server Error: Unable to write to file", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		http.Error(w, "Bad Request: Unable to parse request body", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "Unauthorized: Unable to retrieve user ID from context", http.StatusUnauthorized)
		return
	}

	student.UpdatedOn = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	student.UpdatedBy = sql.NullString{
		String: userID,
		Valid:  true,
	}

	err = student.Update(db)
	if err != nil {
		http.Error(w, "Internal Server Error: Unable to update student", http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(student)
	if err != nil {
		http.Error(w, "Internal Server Error: Unable to marshal response", http.StatusInternalServerError)
		return
	}

	err = writeToFile(data)
	if err != nil {
		http.Error(w, "Internal Server Error: Unable to write to file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteStudent(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Bad Request: Invalid student ID", http.StatusBadRequest)
		return
	}

	student, err := models.GetStudentByID(db, id)
	if err != nil {
		http.Error(w, "Internal Server Error: Unable to check student existence", http.StatusInternalServerError)
		return
	}

	if student == nil {
		response := map[string]string{"message": "Student with this ID not present"}
		data, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Internal Server Error: Unable to marshal response", http.StatusInternalServerError)
			return
		}

		err = writeToFile(data)
		if err != nil {
			http.Error(w, "Internal Server Error: Unable to write to file", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = student.Delete(db, id)
	if err != nil {
		http.Error(w, "Internal Server Error: Unable to delete student", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Student deleted successfully"}
	data, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error: Unable to marshal response", http.StatusInternalServerError)
		return
	}

	err = writeToFile(data)
	if err != nil {
		http.Error(w, "Internal Server Error: Unable to write to file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetStudentByID(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Bad Request: Invalid student ID", http.StatusBadRequest)
		return
	}

	student, err := models.GetStudentByID(db, id)
	if err != nil {
		log.Printf("Error retrieving student: %v", err)
		http.Error(w, "Internal Server Error: Unable to retrieve student", http.StatusInternalServerError)
		return
	}

	if student == nil {

		errorMessage := map[string]string{"error": "Not Found: Student not found"}
		data, err := json.Marshal(errorMessage)
		if err != nil {
			http.Error(w, "Internal Server Error: Unable to marshal error response", http.StatusInternalServerError)
			return
		}

		err = writeToFile(data)
		if err != nil {
			http.Error(w, "Internal Server Error: Unable to write to file", http.StatusInternalServerError)
			return
		}

		http.Error(w, "Not Found: Student not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"id":         student.ID,
		"name":       student.Name,
		"created_by": student.CreatedBy,
		"created_on": student.CreatedOn,
		"updated_by": student.UpdatedBy.String,
		"updated_on": student.UpdatedOn.Time,
	}

	data, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error: Unable to marshal response", http.StatusInternalServerError)
		return
	}

	err = writeToFile(data)
	if err != nil {
		http.Error(w, "Internal Server Error: Unable to write to file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
