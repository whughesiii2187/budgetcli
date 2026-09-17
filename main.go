package main

import (
	"log"

	budgetcmd "github.com/whughesiii2187/wch-projects/Go/budgetcli/cmd"
	"github.com/whughesiii2187/wch-projects/Go/budgetcli/internal/database"
)

func main() {
	pool, err := database.DBConnect()
	if err != nil {
		log.Fatal("Could not connect to database:", err)
	}
	budgetcmd.DB = pool
	defer pool.Close()
	budgetcmd.Execute()
}
