package handlers

import (
	"bytes"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupRouter(handler *StudentHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.Default()

	router.POST("/students", handler.CreateStudent)
	router.GET("/students", handler.GetStudents)
	router.GET("/students/:id", handler.GetStudent)
	router.PUT("/students/:id", handler.UpdateStudent)
	router.DELETE("/students/:id", handler.DeleteStudent)

	return router
}

func TestCreateStudent(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}
	defer db.Close()

	handler := NewStudentHandler(db)
	router := setupRouter(handler)

	mock.ExpectQuery("INSERT INTO students").
		WithArgs(
			"Chinmay",
			24,
			"chinmay@gmail.com",
		).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(1),
		)

	body := []byte(`{
		"name":"Chinmay",
		"age":24,
		"email":"chinmay@gmail.com"
	}`)

	req, _ := http.NewRequest(
		http.MethodPost,
		"/students",
		bytes.NewBuffer(body),
	)

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected %d got %d",
			http.StatusCreated,
			w.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestGetStudent(t *testing.T) {

	db, mock, _ := sqlmock.New()
	defer db.Close()

	handler := NewStudentHandler(db)
	router := setupRouter(handler)

	rows := sqlmock.NewRows(
		[]string{"id", "name", "age", "email"},
	)

	rows.AddRow(
		1,
		"Chinmay",
		24,
		"chinmay@gmail.com",
	)

	mock.ExpectQuery("SELECT id, name, age, email").
		WithArgs("1").
		WillReturnRows(rows)

	req, _ := http.NewRequest(
		http.MethodGet,
		"/students/1",
		nil,
	)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected %d got %d",
			http.StatusOK,
			w.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestGetStudents(t *testing.T) {

	db, mock, _ := sqlmock.New()
	defer db.Close()

	handler := NewStudentHandler(db)
	router := setupRouter(handler)

	rows := sqlmock.NewRows(
		[]string{"id", "name", "age", "email"},
	)

	rows.AddRow(
		1,
		"Chinmay",
		24,
		"chinmay@gmail.com",
	)

	mock.ExpectQuery("SELECT id, name, age, email FROM students").
		WillReturnRows(rows)

	req, _ := http.NewRequest(
		http.MethodGet,
		"/students",
		nil,
	)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected %d got %d",
			http.StatusOK,
			w.Code)
	}
}

func TestUpdateStudent(t *testing.T) {

	db, mock, _ := sqlmock.New()
	defer db.Close()

	handler := NewStudentHandler(db)
	router := setupRouter(handler)

	mock.ExpectExec("UPDATE students").
		WithArgs(
			"Rahul",
			25,
			"rahul@gmail.com",
			"1",
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	body := []byte(`{
		"name":"Rahul",
		"age":25,
		"email":"rahul@gmail.com"
	}`)

	req, _ := http.NewRequest(
		http.MethodPut,
		"/students/1",
		bytes.NewBuffer(body),
	)

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected %d got %d",
			http.StatusOK,
			w.Code)
	}
}

func TestDeleteStudent(t *testing.T) {

	db, mock, _ := sqlmock.New()
	defer db.Close()

	handler := NewStudentHandler(db)
	router := setupRouter(handler)

	mock.ExpectExec("DELETE FROM students").
		WithArgs("1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	req, _ := http.NewRequest(
		http.MethodDelete,
		"/students/1",
		nil,
	)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected %d got %d",
			http.StatusOK,
			w.Code)
	}
}
