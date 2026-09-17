package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/whughesiii2187/wch-projects/Go/budgetcli/internal/models"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a budget item",
	Run:   runAdd,
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func runAdd(cmd *cobra.Command, args []string) {
	for {
		var billType string
		var balanceInput string
		var payPeriodInput string
		var dueAmountInput string
		var dueRecurringDateInput string
		var annualInput string
		var bill models.Bill

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Enter Bill Name:").
					Prompt("?").
					Value(&bill.BillName),
				huh.NewSelect[string]().
					Title("Select the Bill Type:").
					Options(
						huh.NewOption("Recurring", "recurring"),
						huh.NewOption("Credit/Loan", "credit/loan"),
					).
					Value(&billType),
				huh.NewConfirm().
					Title("Is this bill an AutoPay?").
					Affirmative("Auto Pay").
					Negative("Manual Pay").
					Value(&bill.IsAutoPay),
				huh.NewSelect[string]().
					Title("Which Pay Period?").
					Options(
						huh.NewOption("A", "A"),
						huh.NewOption("B", "B"),
					).
					Value(&payPeriodInput),
				huh.NewSelect[string]().
					Title("Is this bill annual or monthly").
					Options(
						huh.NewOption("Annual", "annual"),
						huh.NewOption("Monthly", "monthly"),
					).
					Value(&annualInput),
				huh.NewInput().
					Title("Enter current due amount, if any").
					Prompt("$").
					Value(&dueAmountInput),
				huh.NewInput().
					Title("Enter the date due (day of the month)").
					Prompt(":").
					Value(&dueRecurringDateInput).
					Validate(func(s string) error {
						_, err := strconv.Atoi(s)
						if err != nil {
							return fmt.Errorf("please enter a valid date")
						}
						return nil
					}),
				huh.NewInput().
					Title("Additional Notes:").
					Prompt("?").
					Value(&bill.Notes),
			),
			huh.NewGroup(
				huh.NewInput().
					Title("Enter the current balance due:").
					Prompt("$").
					Value(&balanceInput).
					Validate(func(s string) error {
						_, err := strconv.ParseFloat(s, 64)
						if err != nil {
							return fmt.Errorf("please enter a valid amount")
						}
						return nil
					}),
			).WithHideFunc(func() bool { return billType != "credit/loan" }),
		)

		err := form.Run()
		if err != nil {
			fmt.Println(err)
			return
		}

		bill.PayPeriodPaid = models.PayPeriod(payPeriodInput)

		if balanceInput != "" {
			num, err := strconv.ParseFloat(balanceInput, 64)
			if err != nil {
				fmt.Println(err)
				return
			}
			bill.DueBalance = &num
		}
		bill.DueRecurringDate, _ = strconv.Atoi(dueRecurringDateInput)

		if dueAmountInput != "" {
			num, err := strconv.ParseFloat(dueAmountInput, 64)
			if err != nil {
				fmt.Println(err)
				return
			}
			bill.DueAmount = &num
		}

		if annualInput == "annual" {
			bill.Annual = true
		}

		var toSave bool
		var toContinue bool
		form = huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Do you wish to save this bill?").
					Affirmative("Save").
					Negative("Cancel").
					Value(&toSave),
			),
			huh.NewGroup(
				huh.NewConfirm().
					Title("Do you want to add another bill?").
					Affirmative("Yes").
					Negative("No").
					Value(&toContinue),
			).WithHideFunc(func() bool { return !toSave }),
		)

		err = form.Run()
		if err != nil {
			fmt.Println("Error displaying form")
		}

		if !toSave || !toContinue {
			return
		} else if toSave {
			err = insertAddData(&bill)
			if err != nil {
				fmt.Println("Error while inserting to database, please try again", err)
			}
		}
	}
}

func insertAddData(b *models.Bill) error {
	_, err := DB.Exec(
		context.Background(),
		"INSERT INTO bills (bill_name, due_date, pay_period, auto_pay, annual, notes, balance, amount_due) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)",
		b.BillName, b.DueRecurringDate, b.PayPeriodPaid, b.IsAutoPay, b.Annual, b.Notes, b.DueBalance, b.DueAmount,
	)
	if err != nil {
		return fmt.Errorf("Error inserting bill: %v", err)
	}
	fmt.Println("Bill added successfully")
	return nil
}
