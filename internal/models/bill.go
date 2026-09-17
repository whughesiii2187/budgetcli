package models

import "time"

type PayPeriod string

const (
	PayPeriodA PayPeriod = "A"
	PayPeriodB PayPeriod = "B"
)

type Bill struct {
	BillName         string
	DueRecurringDate int
	DueAmount        *float64
	DueBalance       *float64
	PayPeriodPaid    PayPeriod
	IsAutoPay        bool
	Notes            string
	Annual           bool
}

type Payments struct {
	BillName      string
	DueDate       time.Time
	DueAmount     *float64
	DueBalance    *float64
	PayPeriodPaid PayPeriod
	IsAutoPay     bool
	Paid          bool
}

type InsertPayments struct {
	BillId      int
	BillName    string
	DueDate     int
	AmountDue   *float64
	Balance     *float64
	PayPeriod   string
	Autopay     bool
	AmountInput string // placeholder for form amount entered
}

type UpdatePayments struct {
	PaymentId   int
	BillId      int
	BillName    string
	DueDate     time.Time
	AmountDue   *float64
	Paid        bool
	AmountInput string // placeholder for form amount entered
}

type UpdateBill struct {
	BillId           int
	BillName         string
	DueRecurringDate int
	DueAmount        *float64
	DueBalance       *float64
	PayPeriodPaid    PayPeriod
	IsAutoPay        bool
	Notes            string
	Annual           bool
	AmountInput      string
	BalanceInput     string
	IsActive         bool
}
