package cmd

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/jackc/pgx/v5"

	"github.com/spf13/cobra"
	"github.com/whughesiii2187/wch-projects/Go/budgetcli/internal/models"
)

var prepCmd = &cobra.Command{
	Use:   "prep",
	Short: "Insert amounts due for a budget item",
	Run:   runPrep,
}

func init() {
	rootCmd.AddCommand(prepCmd)
}

func runPrep(cmd *cobra.Command, args []string) {
	var fields []huh.Field
	var prepMonth string
	now := time.Now()
	year := now.Year()
	month := now.Month()

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Which month do you want to run prep for?").
				Options(
					huh.NewOption("Current Month", "current"),
					huh.NewOption("Next Month", "next"),
				).
				Value(&prepMonth),
		),
	)

	err := form.Run()
	if err != nil {
		fmt.Println("Error attemtping to run form:", err)
		return
	}

	if prepMonth == "next" {
		month = month + 1
	}

	upNext, err := retrieveData()
	if err != nil {
		fmt.Println("Unable to retrieve bill data:", err)
		return
	}

	itemsToInsert := make([]models.InsertPayments, 0, len(upNext))

	for i, bill := range upNext {
		var existingId int
		newDate := calculateDueDate(bill.DueDate, year, month, bill.PayPeriod)
		err := DB.QueryRow(context.Background(), "select bill_id from bill_history where bill_id = $1 and due_date = $2;", upNext[i].BillId, newDate).Scan(&existingId)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				itemsToInsert = append(itemsToInsert, upNext[i])
			} else {
				fmt.Printf("error happened checking for dupe: %v", err)
				continue
			}
		} else {
			fmt.Printf("A record already exists for %v, not adding.\n", upNext[i].BillName)
			continue
		}

		if bill.AmountDue == nil {
			inputField := huh.NewInput().
				Title("Enter the amount due for " + bill.BillName + " if known").
				Prompt("$").
				Value(&itemsToInsert[len(itemsToInsert)-1].AmountInput).
				Validate(func(s string) error {
					_, err := strconv.ParseFloat(s, 64)
					if err != nil {
						return fmt.Errorf("please enter a valid amount")
					}
					return nil
				})
			fields = append(fields, inputField)
		}
	}

	if len(fields) == 0 {
		fmt.Println("Nothing to prepare, all bills accounted for")
		return
	}

	form = huh.NewForm(
		huh.NewGroup(fields...),
	)

	err = form.Run()
	if err != nil {
		fmt.Println(err)
		return
	}

	for i, bill := range itemsToInsert {
		zero := 0.0
		if bill.AmountInput != "" {
			num, err := strconv.ParseFloat(bill.AmountInput, 64)
			if err != nil {
				fmt.Println(err)
				return
			}
			itemsToInsert[i].AmountDue = &num
		} else if bill.AmountDue == nil && bill.AmountInput == "" {
			itemsToInsert[i].AmountDue = &zero
		}

		newDate := calculateDueDate(itemsToInsert[i].DueDate, year, month, bill.PayPeriod)
		err := insertData(itemsToInsert[i].BillId, newDate, *itemsToInsert[i].AmountDue, itemsToInsert[i].PayPeriod, itemsToInsert[i].Autopay)
		if err != nil {
			fmt.Printf("Trouble inserting data %v", err)
		}
		fmt.Printf("New month data for %v added successfully!\n", itemsToInsert[i].BillName)
	}
}

func calculateDueDate(dueDate, year int, month time.Month, payPeriod string) time.Time {
	if dueDate <= 10 && strings.ToLower(payPeriod) == "b" {
		month++
	}
	if month > 12 {
		month -= 12
		year++
	}
	return time.Date(year, month, dueDate, 0, 0, 0, 0, time.UTC)
}

func insertData(bill_id int, due_date time.Time, amount_due float64, pay_period string, auto_pay bool) error {
	_, err := DB.Exec(context.Background(),
		"INSERT INTO bill_history (bill_id,amount_due,due_date,pay_period,auto_pay) VALUES ($1,$2,$3,$4,$5)",
		bill_id, amount_due, due_date, pay_period, auto_pay,
	)
	if err != nil {
		return fmt.Errorf("failed to insert %w", err)
	}
	return nil
}

func retrieveData() ([]models.InsertPayments, error) {
	rows, err := DB.Query(context.Background(),
		"select bill_id,bill_name,due_date,amount_due,balance,pay_period,auto_pay from bills where active = true;",
	)
	if err != nil {
		fmt.Println("Error querying bills:", err)
		return nil, err
	}
	defer rows.Close()

	var payme []models.InsertPayments

	for rows.Next() {
		var bill models.InsertPayments
		err := rows.Scan(
			&bill.BillId,
			&bill.BillName,
			&bill.DueDate,
			&bill.AmountDue,
			&bill.Balance,
			&bill.PayPeriod,
			&bill.Autopay,
		)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return nil, err
		}
		payme = append(payme, bill)
	}
	return payme, nil
}
