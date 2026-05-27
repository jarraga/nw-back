package xls

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"nw-back/internal/resetdb"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xuri/excelize/v2"
)

type ImportResult struct {
	CustomersCreated int `json:"customersCreated"`
	PaymentsCreated  int `json:"paymentsCreated"`
}

func ImportCustomers(ctx context.Context, pool *pgxpool.Pool, reader io.Reader) (ImportResult, error) {
	file, err := excelize.OpenReader(reader)
	if err != nil {
		return ImportResult{}, err
	}
	defer file.Close()

	result := ImportResult{}
	importedCustomers := map[int64]bool{}
	customers := []Customer{}
	payments := []Payment{}

	for _, sheet := range file.GetSheetList() {
		year, ok := sheetYear(sheet)
		if !ok {
			continue
		}

		rows, err := file.GetRows(sheet)
		if err != nil {
			return ImportResult{}, err
		}

		for rowIndex := 1; rowIndex < len(rows); rowIndex++ {
			rowNumber := rowIndex + 1
			customer, ok, err := readCustomer(file, sheet, rowNumber)
			if err != nil {
				return ImportResult{}, err
			}
			if !ok {
				continue
			}

			if !importedCustomers[customer.ID] {
				importedCustomers[customer.ID] = true
				customers = append(customers, customer)
			}

			rowPayments, err := readPayments(file, sheet, rowNumber, customer.ID, year)
			if err != nil {
				return ImportResult{}, err
			}

			payments = append(payments, rowPayments...)
		}
	}

	err = resetdb.Reset(ctx, pool)
	if err != nil {
		return ImportResult{}, err
	}

	err = copyImportCustomers(ctx, pool, customers)
	if err != nil {
		return ImportResult{}, err
	}

	err = copyImportPayments(ctx, pool, payments)
	if err != nil {
		return ImportResult{}, err
	}

	err = resetCustomerSequence(ctx, pool)
	if err != nil {
		return ImportResult{}, err
	}

	result.CustomersCreated = len(customers)
	result.PaymentsCreated = len(payments)

	return result, nil
}

func readCustomer(file *excelize.File, sheet string, row int) (Customer, bool, error) {
	id, ok, err := int64Cell(file, sheet, row, 1)
	if err != nil || !ok {
		return Customer{}, false, err
	}

	monthlyFee, ok, err := int32Cell(file, sheet, row, 6)
	if err != nil {
		return Customer{}, false, err
	}
	if !ok {
		return Customer{}, false, fmt.Errorf("%s!F%d monthly_fee is required", sheet, row)
	}

	billingStartedAt, ok, err := dateCell(file, sheet, row, 7)
	if err != nil {
		return Customer{}, false, err
	}
	if !ok {
		return Customer{}, false, fmt.Errorf("%s!G%d billing_started_at is required", sheet, row)
	}

	return Customer{
		ID:               id,
		CompanyName:      stringCell(file, sheet, row, 2),
		CompanyType:      stringCell(file, sheet, row, 3),
		Phone:            stringCell(file, sheet, row, 4),
		Email:            stringCell(file, sheet, row, 5),
		MonthlyFee:       monthlyFee,
		BillingStartedAt: billingStartedAt,
		Comments:         stringCell(file, sheet, row, 8),
	}, true, nil
}

func readPayments(file *excelize.File, sheet string, row int, customerID int64, year int) ([]Payment, error) {
	payments := []Payment{}

	for month := 1; month <= 12; month++ {
		paidAt, ok, err := dateCell(file, sheet, row, 8+month)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}

		payments = append(payments, Payment{
			CustomerID: customerID,
			Year:       year,
			Month:      month,
			PaidAt:     paidAt,
		})
	}

	return payments, nil
}

func copyImportCustomers(ctx context.Context, pool *pgxpool.Pool, customers []Customer) error {
	if len(customers) == 0 {
		return nil
	}

	rows := make([][]any, 0, len(customers))

	for _, customer := range customers {
		rows = append(rows, []any{
			customer.ID,
			customer.CompanyName,
			customer.CompanyType,
			customer.Phone,
			customer.Email,
			customer.MonthlyFee,
			customer.BillingStartedAt,
			customer.Comments,
		})
	}

	_, err := pool.CopyFrom(
		ctx,
		pgx.Identifier{"customers"},
		[]string{
			"id",
			"company_name",
			"company_type",
			"phone",
			"email",
			"monthly_fee",
			"billing_started_at",
			"comments",
		},
		pgx.CopyFromRows(rows),
	)
	return err
}

func copyImportPayments(ctx context.Context, pool *pgxpool.Pool, payments []Payment) error {
	if len(payments) == 0 {
		return nil
	}

	rows := make([][]any, 0, len(payments))

	for _, payment := range payments {
		rows = append(rows, []any{
			payment.CustomerID,
			payment.Year,
			payment.Month,
			"paid",
			payment.PaidAt,
		})
	}

	_, err := pool.CopyFrom(
		ctx,
		pgx.Identifier{"customer_payments"},
		[]string{
			"customer_id",
			"year",
			"month",
			"status",
			"paid_at",
		},
		pgx.CopyFromRows(rows),
	)
	return err
}

func resetCustomerSequence(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
		SELECT setval(
			pg_get_serial_sequence('customers', 'id'),
			COALESCE((SELECT MAX(id) FROM customers), 1),
			(SELECT COUNT(*) > 0 FROM customers)
		)
	`)
	return err
}

func sheetYear(sheet string) (int, bool) {
	year, err := strconv.Atoi(strings.TrimSpace(sheet))
	if err != nil {
		return 0, false
	}

	return year, true
}

func stringCell(file *excelize.File, sheet string, row int, column int) string {
	cell, err := excelize.CoordinatesToCellName(column, row)
	if err != nil {
		return ""
	}

	value, err := file.GetCellValue(sheet, cell)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(value)
}

func int64Cell(file *excelize.File, sheet string, row int, column int) (int64, bool, error) {
	value := stringCell(file, sheet, row, column)
	if value == "" {
		return 0, false, nil
	}

	number, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, false, fmt.Errorf("%s!%s must be an integer", sheet, cellName(column, row))
	}

	return number, true, nil
}

func int32Cell(file *excelize.File, sheet string, row int, column int) (int32, bool, error) {
	number, ok, err := int64Cell(file, sheet, row, column)
	if err != nil || !ok {
		return 0, ok, err
	}

	return int32(number), true, nil
}

func dateCell(file *excelize.File, sheet string, row int, column int) (time.Time, bool, error) {
	cell := cellName(column, row)

	raw, err := file.GetCellValue(sheet, cell, excelize.Options{RawCellValue: true})
	if err != nil {
		return time.Time{}, false, err
	}

	raw = strings.TrimSpace(raw)
	if raw == "" {
		return time.Time{}, false, nil
	}

	if serial, err := strconv.ParseFloat(raw, 64); err == nil {
		date, err := excelize.ExcelDateToTime(serial, false)
		if err != nil {
			return time.Time{}, false, err
		}

		return date, true, nil
	}

	formatted, err := file.GetCellValue(sheet, cell)
	if err != nil {
		return time.Time{}, false, err
	}

	date, err := parseDate(strings.TrimSpace(formatted))
	if err != nil {
		return time.Time{}, false, fmt.Errorf("%s!%s must be a date", sheet, cell)
	}

	return date, true, nil
}

func parseDate(value string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02",
		"1/2/06",
		"1/2/2006",
		"2/1/06",
		"2/1/2006",
	}

	var lastErr error
	for _, format := range formats {
		date, err := time.Parse(format, value)
		if err == nil {
			return date, nil
		}
		lastErr = err
	}

	return time.Time{}, lastErr
}

func cellName(column int, row int) string {
	cell, err := excelize.CoordinatesToCellName(column, row)
	if err != nil {
		return ""
	}

	return cell
}
