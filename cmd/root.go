// Package cmd where it all gets started
package cmd

import (
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "budgetcli",
	Short: "A Go CLI tool for managing the budget",
}

var DB *pgxpool.Pool

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
