package main

import (
	"backend/database"
	"backend/routers"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file:, assuming variables are set in the environment.... %v", err)

	}
	dbUser := os.Getenv("DB_USER")

	dbPass := os.Getenv("DB_PASS")

	jwtkey := []byte(os.Getenv("JWT_KEY"))

	dbHost := os.Getenv("DB_HOST")
	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatalf("Invalid DB_PORT: %v", err)
	}
	dbName := os.Getenv("DB_NAME")
	frontendURLs := os.Getenv("FRONTEND_URL")

	allowedOrigins := []string{"http://localhost:5173"}

	if frontendURLs != "" {
		urls := strings.Split(frontendURLs, ",")
		for _, url := range urls {
			allowedOrigins = append(allowedOrigins, strings.TrimSpace(url))
		}
	}
	db := database.InitDB(dbUser, dbPass, dbHost, dbPort, dbName)
	router := routers.SetupRouter(db, jwtkey)
	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		Debug:            false,
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	handler := c.Handler(router)
	http.ListenAndServe(":"+port, handler)

	fmt.Println("Database connected!", db)
}
