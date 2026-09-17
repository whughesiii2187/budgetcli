package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// func unsetEnv() {
// 	unsetTheseVars := []string{"DB_PORT", "DB_HOST", "DB_USER", "DB_PASSWORD", "DB_NAME"}
//
// 	for _, val := range unsetTheseVars {
// 		if os.Getenv(val) != "" {
// 			os.Unsetenv(val)
// 		}
// 	}
// }

func DBConnect() (*pgxpool.Pool, error) {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading info")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	pool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
