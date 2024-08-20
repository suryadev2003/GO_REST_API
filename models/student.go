package models

import (
	"database/sql"
	"time"
)

type Student struct {
	ID        int            `json:"id"`
	Name      string         `json:"name"`
	CreatedBy string         `json:"created_by"`
	CreatedOn time.Time      `json:"created_on"`
	UpdatedBy sql.NullString `json:"updated_by"`
	UpdatedOn sql.NullTime   `json:"updated_on"`
}

func (s *Student) Create(db *sql.DB) error {
	query := "INSERT INTO students (name, created_by, created_on) VALUES (?, ?, ?)"
	_, err := db.Exec(query, s.Name, s.CreatedBy, time.Now())
	return err
}

func (s *Student) Update(db *sql.DB) error {
	query := `UPDATE students SET name = ?, updated_by = ?, updated_on = ? WHERE id = ?`
	_, err := db.Exec(query, s.Name, s.UpdatedBy, s.UpdatedOn, s.ID)
	return err
}

func (s *Student) Delete(db *sql.DB, id int) error {
	query := "DELETE FROM students WHERE id = ?"
	_, err := db.Exec(query, id)
	return err
}

func GetStudents(db *sql.DB) ([]Student, error) {
	query := "SELECT id, name, created_by, created_on, updated_by, updated_on FROM students"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []Student
	for rows.Next() {
		var student Student
		if err := rows.Scan(&student.ID, &student.Name, &student.CreatedBy, &student.CreatedOn, &student.UpdatedBy, &student.UpdatedOn); err != nil {
			return nil, err
		}
		students = append(students, student)
	}
	return students, nil
}

func GetStudentByID(db *sql.DB, id int) (*Student, error) {
	query := "SELECT id, name, created_by, created_on, updated_by, updated_on FROM students WHERE id = ?"
	row := db.QueryRow(query, id)

	var student Student
	err := row.Scan(&student.ID, &student.Name, &student.CreatedBy, &student.CreatedOn, &student.UpdatedBy, &student.UpdatedOn)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &student, nil
}
