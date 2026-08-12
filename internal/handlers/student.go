package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"sre.com/internal/models"
)

type StudentHandler struct {
	DB *sql.DB
}

func NewStudentHandler(db *sql.DB) *StudentHandler {
	return &StudentHandler{
		DB: db,
	}
}

func (h *StudentHandler) CreateStudent(c *gin.Context) {

	log.Println("POST /api/v1/students")

	var student models.Student

	if err := c.ShouldBindJSON(&student); err != nil {
		log.Printf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	query := `
	INSERT INTO students(name, age, email)
	VALUES($1,$2,$3)
	RETURNING id
	`

	err := h.DB.QueryRow(
		query,
		student.Name,
		student.Age,
		student.Email,
	).Scan(&student.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, student)
}

func (h *StudentHandler) GetStudents(c *gin.Context) {

	log.Println("GET /api/v1/students")

	rows, err := h.DB.Query("SELECT id, name, age, email FROM students")

	if err != nil {
		log.Printf("Failed to create student: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Error closing rows: %v", err)
		}
	}()

	var students []models.Student

	for rows.Next() {

		var student models.Student

		err := rows.Scan(
			&student.ID,
			&student.Name,
			&student.Age,
			&student.Email,
		)

		if err != nil {

			log.Printf("Failed to fetch students: %v", err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		students = append(students, student)
	}

	log.Printf("Returned %d students", len(students))

	c.JSON(http.StatusOK, students)
}

func (h *StudentHandler) GetStudent(c *gin.Context) {

	id := c.Param("id")

	log.Printf("GET /api/v1/students/%s", id)

	var student models.Student

	query := `
	SELECT id, name, age, email
	FROM students
	WHERE id=$1
	`

	err := h.DB.QueryRow(query, id).Scan(
		&student.ID,
		&student.Name,
		&student.Age,
		&student.Email,
	)

	if err == sql.ErrNoRows {

		log.Printf("Student %s not found", id)

		c.JSON(http.StatusNotFound, gin.H{
			"error": "Student not found",
		})
		return
	}

	if err != nil {

		log.Printf("Database error: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Printf("Student %s fetched successfully", id)

	c.JSON(http.StatusOK, student)
}

func (h *StudentHandler) UpdateStudent(c *gin.Context) {

	id := c.Param("id")

	log.Printf("PUT /api/v1/students/%s", id)

	var student models.Student

	if err := c.ShouldBindJSON(&student); err != nil {

		log.Printf("Invalid request body: %v", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	query := `
	UPDATE students
	SET name=$1,
	    age=$2,
	    email=$3
	WHERE id=$4
	`

	result, err := h.DB.Exec(
		query,
		student.Name,
		student.Age,
		student.Email,
		id,
	)

	if err != nil {

		log.Printf("Failed to update student %s: %v", id, err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	rows, _ := result.RowsAffected()

	if rows == 0 {

		log.Printf("Student %s not found", id)

		c.JSON(http.StatusNotFound, gin.H{
			"error": "Student not found",
		})
		return
	}

	student.ID, _ = strconv.Atoi(id)

	log.Printf("Student %s updated successfully", id)

	c.JSON(http.StatusOK, student)
}

func (h *StudentHandler) DeleteStudent(c *gin.Context) {

	id := c.Param("id")

	log.Printf("DELETE /api/v1/students/%s", id)

	query := `
	DELETE FROM students
	WHERE id=$1
	`

	result, err := h.DB.Exec(query, id)

	if err != nil {

		log.Printf("Failed to delete student %s: %v", id, err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	rows, _ := result.RowsAffected()

	if rows == 0 {

		log.Printf("Student %s not found", id)

		c.JSON(http.StatusNotFound, gin.H{
			"error": "Student not found",
		})
		return
	}

	log.Printf("Student %s deleted successfully", id)

	c.JSON(http.StatusOK, gin.H{
		"message": "Student deleted successfully",
	})
}
