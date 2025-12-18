package main

import (
	"backend/database"
	"backend/routers"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	dbUser := os.Getenv("DB_USER")

	dbPass := os.Getenv("DB_PASS")

	jwtkey := []byte(os.Getenv("JWT_KEY"))

	db := database.InitDB(dbUser, dbPass)
	router := routers.SetupRouter(db, jwtkey)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		Debug:            false,
	})
	handler := c.Handler(router)
	http.ListenAndServe(":8080", handler)

	fmt.Println("Database connected!", db)
}
