-- BudgetCLI Database Schema

CREATE TABLE bills (
    bill_id      SERIAL PRIMARY KEY,
    bill_name    VARCHAR(255) NOT NULL,
    due_date     INTEGER NOT NULL,
    pay_period   CHAR(1) NOT NULL,
    balance      NUMERIC(10,2),
    amount_due   NUMERIC(10,2),
    auto_pay     BOOLEAN NOT NULL DEFAULT FALSE,
    annual       BOOLEAN NOT NULL,
    notes        VARCHAR(255),
    active       BOOLEAN NOT NULL DEFAULT TRUE,
   last_updated  TIMESTAMP DEFAULT NOW()
);

CREATE TABLE bill_history (
    payment_id    SERIAL PRIMARY KEY,
    bill_id       INTEGER NOT NULL REFERENCES bills(bill_id),
    amount_due    NUMERIC(10,2) NOT NULL,
    due_date      DATE NOT NULL,
    pay_period    CHAR(1) NOT NULL,
    auto_pay      BOOLEAN NOT NULL DEFAULT FALSE,
    paid          BOOLEAN NOT NULL DEFAULT FALSE,
    prepared_date TIMESTAMP DEFAULT NOW()
);
