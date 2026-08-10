package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"sre.com/internal/db"
	"sre.com/internal/handlers"
)

func main() {

	log.Println("Loading environment variables...")

	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using system environment variables")
	}

	log.Println("Environment variables loaded")

	log.Println("Connecting to PostgreSQL...")

	// Connect to PostgreSQL
	database, err := db.Connect()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to PostgreSQL")

	defer database.Close()

	router := gin.Default()

	router.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	studentHandler := handlers.NewStudentHandler(database)

	api := router.Group("/api/v1")
	{
		api.POST("/students", studentHandler.CreateStudent)
		api.GET("/students", studentHandler.GetStudents)
		api.GET("/students/:id", studentHandler.GetStudent)
		api.PUT("/students/:id", studentHandler.UpdateStudent)
		api.DELETE("/students/:id", studentHandler.DeleteStudent)
	}

	port := os.Getenv("PORT")

	print("Server started on port: " + port + "\n")

	log.Printf("Server started on :%s", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
