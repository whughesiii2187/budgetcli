package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"
	"github.com/whughesiii2187/wch-projects/Go/budgetcli/internal/models"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all budget items",
	Run:   runList,
}

func init() {
	listCmd.Flags().Bool("month", false, "Show all bills for the current month")
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) {
	month, _ := cmd.Flags().GetBool("month")

	listOutput(month)
}

func listOutput(month bool) {
	timeSpan := "CURRENT_DATE AND CURRENT_DATE+14"
	if month {
		timeSpan = "DATE_TRUNC('month', CURRENT_DATE) AND (DATE_TRUNC('month', CURRENT_DATE) + INTERVAL '2 month - 1 day')"
	}
	rows, err := DB.Query(
		context.Background(),
		fmt.Sprintf("select b.bill_name,a.due_date,a.amount_due, b.balance, a.pay_period, a.auto_pay, a.paid from bill_history as a join bills as b on a.bill_id = b.bill_id WHERE a.due_date BETWEEN %v;", timeSpan),
	)
	if err != nil {
		fmt.Println("Error querying bills:", err)
		return
	}
	defer rows.Close()

	var payme []models.Payments

	for rows.Next() {
		var bill models.Payments
		err := rows.Scan(
			&bill.BillName,
			&bill.DueDate,
			&bill.DueAmount,
			&bill.DueBalance,
			&bill.PayPeriodPaid,
			&bill.IsAutoPay,
			&bill.Paid,
		)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}
		payme = append(payme, bill)
	}

	outputColorTable(payme)
}

func outputColorTable(c []models.Payments) {
	colorCfg := renderer.ColorizedConfig{
		Header: renderer.Tint{
			FG: renderer.Colors{color.FgGreen, color.Bold}, // Green bold headers
			BG: renderer.Colors{color.BgHiWhite},
		},
		Column: renderer.Tint{
			FG: renderer.Colors{color.FgCyan}, // Default cyan for rows
			Columns: []renderer.Tint{
				{FG: renderer.Colors{color.FgMagenta}}, // Magenta for column 0
				{},                                     // Inherit default (cyan)
				{FG: renderer.Colors{color.FgHiRed}},   // High-intensity red for column 2
			},
		},
		Footer: renderer.Tint{
			FG: renderer.Colors{color.FgYellow, color.Bold}, // Yellow bold footer
			Columns: []renderer.Tint{
				{},                                      // Inherit default
				{FG: renderer.Colors{color.FgHiYellow}}, // High-intensity yellow for column 1
				{},                                      // Inherit default
			},
		},
		Border:    renderer.Tint{FG: renderer.Colors{color.FgWhite}}, // White borders
		Separator: renderer.Tint{FG: renderer.Colors{color.FgWhite}}, // White separators
	}

	table := tablewriter.NewTable(
		os.Stdout,
		tablewriter.WithRenderer(renderer.NewColorized(colorCfg)),
		tablewriter.WithConfig(tablewriter.Config{
			Row: tw.CellConfig{
				Formatting:   tw.CellFormatting{AutoWrap: tw.WrapNormal}, // Wrap long content
				Alignment:    tw.CellAlignment{Global: tw.AlignLeft},     // Left-align rows
				ColMaxWidths: tw.CellWidth{Global: 25},
			},
			Footer: tw.CellConfig{
				Alignment: tw.CellAlignment{Global: tw.AlignRight},
			},
		}),
	)

	table.Header([]string{"Bill", "Due Date", "Amount Due", "Balance", "Pay Period", "Auto Pay", "Paid"})
	var display []string
	for _, i := range c {
		dueDate := i.DueDate.Format("2006/01/02")
		amountDue := "N/A"
		balanceDue := "N/A"
		autoPay := "N"
		paid := "N"
		payPeriod := string(i.PayPeriodPaid)
		billName := i.BillName
		if i.DueAmount != nil {
			amountDue = fmt.Sprintf("%.2f", *i.DueAmount)
		}
		if i.DueBalance != nil {
			balanceDue = fmt.Sprintf("%.2f", *i.DueBalance)
		}
		if i.IsAutoPay {
			autoPay = "Y"
		}
		if i.Paid {
			paid = "Y"
		}
		display = []string{billName, dueDate, amountDue, balanceDue, payPeriod, autoPay, paid}
		table.Append(display)
	}
	table.Render()
}
