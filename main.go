package main

import (
	"database/sql"
	"fmt"
	"log"
	"muhammaddev/internal/database"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Yha Allah, Please make this a successful Project.")

	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("there is no port provided kindly provide a port.")
		return
	}
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Println("empty API_KEY")
		return

	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Println("empty dbURL")
		return
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	dbQueries := database.New(db)

	apiConfig := Config{
		PORT:                          port,
		DB:                            dbQueries,
		API_KEY:                       apiKey,
		AccessTokenExpirationMinutes:  15,
		RefreshTokenExpirationMinutes: 24 * 60,
	}
	// refreshToken, err := auth.GenerateRefreshToken()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(ref reshToken)
	// fmt.Println(apiConfig)
	server(&apiConfig)

}
