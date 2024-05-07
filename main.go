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

type Config struct {
	PORT string
	DB   *database.Queries
}

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

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Println("empty dbURL")

	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	dbQueries := database.New(db)

	apiConfig := Config{
		PORT: port,
		DB:   dbQueries,
	}

	server(&apiConfig)

}
