package cmd

import (
	"context"
	"fmt"
	"slices"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/whughesiii2187/wch-projects/Go/budgetcli/internal/models"
)

var updateCmd = &cobra.Command{
	Use:  "update",
	Long: "Make updates to bills that have been prepped or activate / deactivate bills",
	Run:  runUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) {
	var choice string
	formOptions := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Which update method to you choose?").
				Options(
					huh.NewOption("Update upcoming payment", "payment"),
					huh.NewOption("Modify a bill", "modify"),
				).Value(&choice),
		),
	)

	err := formOptions.Run()
	if err != nil {
		fmt.Printf("unable to display form: %v", err)
		return
	}

	if choice == "payment" {
		for {
			getUpdates, err := updatePaymentSelectForm()
			if err != nil {
				fmt.Printf("error displaying form: %v", err)
				return
			}
			if len(getUpdates) == 0 {
				continue
			}
			err = updatePaymentInsertForm(getUpdates)
			if err != nil {
				fmt.Printf("error updating: %v", err)
				return
			}
			break
		}
	}
	if choice == "modify" {
		for {
			modifyBill, err := updateBillsSelectForm()
			if err != nil {
				fmt.Printf("error displaying form: %v", err)
				return
			}
			if modifyBill.BillId < 1 {
				continue
			}
			err = updateBillInsertForm(modifyBill)
			if err != nil {
				fmt.Printf("error displaying form: %v", err)
			}
			break
		}
	}
}

func pullUpdates() ([]models.UpdatePayments, error) {
	rows, err := DB.Query(
		context.Background(),
		fmt.Sprintln("select a.payment_id,a.bill_id,b.bill_name,a.due_date,a.amount_due,a.paid from bill_history as a join bills as b on a.bill_id = b.bill_id WHERE a.paid = false"),
	)
	if err != nil {
		fmt.Println("Error querying bills:", err)
		return nil, err
	}
	defer rows.Close()

	var update []models.UpdatePayments

	for rows.Next() {
		var updates models.UpdatePayments

		err := rows.Scan(
			&updates.PaymentId,
			&updates.BillId,
			&updates.BillName,
			&updates.DueDate,
			&updates.AmountDue,
			&updates.Paid,
		)
		if err != nil {
			fmt.Println("Error reading rows:", err)
			return nil, err
		}
		update = append(update, updates)
	}
	return update, nil
}

func pullBills() ([]models.UpdateBill, error) {
	rows, err := DB.Query(
		context.Background(),
		fmt.Sprintln("select bill_id, bill_name, due_date, pay_period, balance, amount_due, auto_pay, annual, notes, active FROM bills"),
	)
	if err != nil {
		fmt.Printf("error querying bills: %v", err)
		return nil, err
	}
	defer rows.Close()

	var update []models.UpdateBill

	for rows.Next() {
		var updates models.UpdateBill

		err := rows.Scan(
			&updates.BillId,
			&updates.BillName,
			&updates.DueRecurringDate,
			&updates.PayPeriodPaid,
			&updates.DueBalance,
			&updates.DueAmount,
			&updates.IsAutoPay,
			&updates.Annual,
			&updates.Notes,
			&updates.IsActive,
		)
		if err != nil {
			fmt.Printf("error reading rows: %v", err)
			return nil, err
		}
		update = append(update, updates)
	}
	return update, nil
}

func updateBillsSelectForm() (models.UpdateBill, error) {
	var selection int
	var fields []huh.Option[int]

	billList, err := pullBills()
	if err != nil {
		fmt.Printf("error pulling bill data: %v", err)
		return models.UpdateBill{}, err
	}

	for _, bill := range billList {
		selectField := huh.NewOption(bill.BillName, bill.BillId)
		fields = append(fields, selectField)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("Which bill do you wish to modify").
				Options(fields...).Value(&selection),
		),
	)

	err = form.Run()
	if err != nil {
		fmt.Printf("error displaying form: %v", err)
	}

	var newSelection models.UpdateBill

	for _, choice := range billList {
		if choice.BillId == selection {
			newSelection = choice
		}
	}
	return newSelection, nil
}

func updateBillInsertForm(chosenBill models.UpdateBill) error {
	var amount, balance, annual, dueDate, payPeriod, active string
	if chosenBill.DueAmount != nil {
		amount = strconv.FormatFloat(*chosenBill.DueAmount, 'f', 2, 64)
	} else {
		amount = "0.00"
	}
	if chosenBill.DueBalance != nil {
		balance = strconv.FormatFloat(*chosenBill.DueBalance, 'f', 2, 64)
	} else {
		balance = "0.00"
	}
	if chosenBill.Annual != false {
		annual = "annual"
	} else {
		annual = "monthly"
	}
	if chosenBill.IsActive != true {
		active = "notactive"
	} else {
		active = "active"
	}
	dueDate = strconv.Itoa(chosenBill.DueRecurringDate)
	payPeriod = string(chosenBill.PayPeriodPaid)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Bill Name:").
				Prompt("?").
				Value(&chosenBill.BillName),
			huh.NewConfirm().
				Title("AutoPay").
				Affirmative("Auto Pay").
				Negative("Manual Pay").
				Value(&chosenBill.IsAutoPay),
			huh.NewSelect[string]().
				Title("Pay Period?").
				Options(
					huh.NewOption("A", "A"),
					huh.NewOption("B", "B"),
				).
				Value(&payPeriod),
			huh.NewSelect[string]().
				Title("Annual or Monthly").
				Options(
					huh.NewOption("Annual", "annual"),
					huh.NewOption("Monthly", "monthly"),
				).
				Value(&annual),
			huh.NewInput().
				Title("Enter current due amount, if any").
				Prompt("$").
				Value(&amount).
				Validate(func(s string) error {
					_, err := strconv.ParseFloat(s, 64)
					if err != nil {
						return fmt.Errorf("please enter a valid amount")
					}
					return nil
				},
				),
			huh.NewInput().
				Title("Enter the current balance due:").
				Prompt("$").
				Value(&balance).
				Validate(func(s string) error {
					_, err := strconv.ParseFloat(s, 64)
					if err != nil {
						return fmt.Errorf("please enter a valid amount")
					}
					return nil
				},
				),
			huh.NewInput().
				Title("Enter the date due (day of the month)").
				Prompt(":").
				Value(&dueDate).
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
				Value(&chosenBill.Notes),
			huh.NewSelect[string]().
				Title("Deactivate or leave Active").
				Options(
					huh.NewOption("Activate", "active"),
					huh.NewOption("Deactivate", "notactive"),
				).
				Value(&active),
		),
	)

	err := form.Run()
	if err != nil {
		fmt.Printf("error running form: %v", err)
	}

	newAmount, _ := strconv.ParseFloat(amount, 64)
	if newAmount != 0 {
		chosenBill.DueAmount = &newAmount
	}
	newBalance, _ := strconv.ParseFloat(balance, 64)
	if newBalance != 0 {
		chosenBill.DueBalance = &newBalance
	}
	chosenBill.DueRecurringDate, _ = strconv.Atoi(dueDate)
	chosenBill.PayPeriodPaid = models.PayPeriod(payPeriod)
	if annual == "annual" {
		chosenBill.Annual = true
	}
	if active == "notactive" {
		chosenBill.IsActive = false
	}

	_, err = DB.Exec(context.Background(),
		"UPDATE bills SET bill_name = $1, due_date = $2, pay_period = $3, balance =$4, amount_due = $5, auto_pay = $6, annual = $7, notes = $8, active = $9  WHERE bill_id = $10",
		chosenBill.BillName, chosenBill.DueRecurringDate, chosenBill.PayPeriodPaid, chosenBill.DueBalance, chosenBill.DueAmount, chosenBill.IsAutoPay, chosenBill.Annual, chosenBill.Notes, chosenBill.IsActive, chosenBill.BillId,
	)
	if err != nil {
		return fmt.Errorf("insert update failed: %v", err)
	}
	return nil
}

func updatePaymentSelectForm() ([]models.UpdatePayments, error) {
	var selections []int
	var fields []huh.Option[int]
	updateList, err := pullUpdates()
	if err != nil {
		fmt.Printf("error pulling update data: %v", err)
		return nil, err
	}

	for _, pymt := range updateList {
		selectField := huh.NewOption(pymt.BillName+" ("+pymt.DueDate.Format("2006/01/02")+")", pymt.PaymentId)
		fields = append(fields, selectField)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[int]().
				Title("Which payment(s) do you want to update?\nPress SPACE to select one or more bills").
				Options(fields...).Value(&selections),
		),
	)

	if err := form.Run(); err != nil {
		fmt.Printf("failed to run the form: %v", err)
		return nil, err
	}

	var newUpdatePayment []models.UpdatePayments

	for _, newPymt := range updateList {
		if slices.Contains(selections, newPymt.PaymentId) {
			newUpdatePayment = append(newUpdatePayment, newPymt)
		}
	}

	return newUpdatePayment, nil
}

func updatePaymentInsertForm(chosenBills []models.UpdatePayments) error {
	var fields []huh.Field

	for i, bill := range chosenBills {
		if *chosenBills[i].AmountDue == 0.00 {
			amountField := huh.NewInput().
				Title("Enter amount due for " + bill.BillName).
				Prompt("$").
				Value(&chosenBills[i].AmountInput)
			fields = append(fields, amountField)
		}

		formattedAmount := fmt.Sprintf("%.2f", *chosenBills[i].AmountDue)
		paidField := huh.NewConfirm().
			Title("Mark " + bill.BillName + " ($" + formattedAmount + ") as paid?").
			Affirmative("Paid").
			Negative("Not Paid").
			Value(&chosenBills[i].Paid)
		fields = append(fields, paidField)
	}

	form := huh.NewForm(
		huh.NewGroup(fields...),
	)
	if err := form.Run(); err != nil {
		fmt.Printf("failed to run the form: %v", err)
		return err
	}

	for i, bill := range chosenBills {
		if bill.AmountInput != "" {
			num, err := strconv.ParseFloat(bill.AmountInput, 64)
			if err != nil {
				fmt.Printf("invalid amount for %s: %v", bill.BillName, err)
				continue
			}
			chosenBills[i].AmountDue = &num
		}

		if bill.AmountInput == "" && bill.Paid {
			chosenBills[i].Paid = false
			fmt.Println("Amount due is still $0, there for unable to be marked as paid")
		}
		// UPDATE SQL goes here, using chosenBills[i].PaymentId, chosenBills[i].AmountDue, chosenBills[i].Paid
		_, err := DB.Exec(context.Background(),
			"UPDATE bill_history SET amount_due = $1, paid = $2 WHERE payment_id = $3",
			chosenBills[i].AmountDue, chosenBills[i].Paid, chosenBills[i].PaymentId,
		)
		if err != nil {
			return fmt.Errorf("insert update failed: %v", err)
		}
	}

	return nil
}
