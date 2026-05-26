package xls

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xuri/excelize/v2"
)

type Customer struct {
	ID               int64
	CompanyName      string
	CompanyType      string
	Phone            string
	Email            string
	MonthlyFee       int32
	BillingStartedAt time.Time
	Comments         string
}

type Payment struct {
	CustomerID int64
	Year       int
	Month      int
	PaidAt     time.Time
}

func ExportCustomers(ctx context.Context, pool *pgxpool.Pool, outputPath string) (string, error) {
	file, err := CustomerWorkbook(ctx, pool)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if outputPath == "" {
		outputPath = defaultOutputPath()
	}

	err = os.MkdirAll(filepath.Dir(outputPath), 0755)
	if err != nil {
		return "", err
	}

	err = file.SaveAs(outputPath)
	if err != nil {
		return "", err
	}

	return outputPath, nil
}

func WriteCustomers(ctx context.Context, pool *pgxpool.Pool, writer io.Writer) error {
	file, err := CustomerWorkbook(ctx, pool)
	if err != nil {
		return err
	}
	defer file.Close()

	return file.Write(writer)
}

func CustomerWorkbook(ctx context.Context, pool *pgxpool.Pool) (*excelize.File, error) {
	customers, err := listCustomers(ctx, pool)
	if err != nil {
		return nil, err
	}

	payments, err := listPayments(ctx, pool)
	if err != nil {
		return nil, err
	}

	file := excelize.NewFile()

	err = buildWorkbook(file, customers, payments)
	if err != nil {
		file.Close()
		return nil, err
	}

	return file, nil
}

func listCustomers(ctx context.Context, pool *pgxpool.Pool) ([]Customer, error) {
	rows, err := pool.Query(ctx, `
		SELECT
			id,
			company_name,
			company_type::text,
			phone,
			email,
			monthly_fee,
			billing_started_at,
			comments
		FROM customers
		WHERE deactivated = FALSE
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	customers := []Customer{}

	for rows.Next() {
		var customer Customer
		err = rows.Scan(
			&customer.ID,
			&customer.CompanyName,
			&customer.CompanyType,
			&customer.Phone,
			&customer.Email,
			&customer.MonthlyFee,
			&customer.BillingStartedAt,
			&customer.Comments,
		)
		if err != nil {
			return nil, err
		}

		customers = append(customers, customer)
	}

	return customers, rows.Err()
}

func listPayments(ctx context.Context, pool *pgxpool.Pool) ([]Payment, error) {
	rows, err := pool.Query(ctx, `
		SELECT
			customer_id,
			year,
			month,
			paid_at
		FROM customer_payments
		WHERE status = 'paid'
		  AND paid_at IS NOT NULL
		ORDER BY customer_id, year, month
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	payments := []Payment{}

	for rows.Next() {
		var payment Payment
		err = rows.Scan(
			&payment.CustomerID,
			&payment.Year,
			&payment.Month,
			&payment.PaidAt,
		)
		if err != nil {
			return nil, err
		}

		payments = append(payments, payment)
	}

	return payments, rows.Err()
}

func buildWorkbook(file *excelize.File, customers []Customer, payments []Payment) error {
	years := workbookYears(customers, payments)
	paymentDates := paymentDatesByCustomer(payments)

	firstSheet := true
	for _, year := range years {
		sheet := fmt.Sprintf("%d", year)
		if firstSheet {
			err := file.SetSheetName("Sheet1", sheet)
			if err != nil {
				return err
			}
			firstSheet = false
		} else {
			_, err := file.NewSheet(sheet)
			if err != nil {
				return err
			}
		}

		err := writeYearSheet(file, sheet, year, customers, paymentDates)
		if err != nil {
			return err
		}
	}

	return nil
}

func writeYearSheet(file *excelize.File, sheet string, year int, customers []Customer, paymentDates map[int64]map[int]map[int]time.Time) error {
	headers := []string{
		"customer_id",
		"company_name",
		"company_type",
		"phone",
		"email",
		"monthly_fee",
		"billing_started_at",
		"comments",
		"january",
		"february",
		"march",
		"april",
		"may",
		"june",
		"july",
		"august",
		"september",
		"october",
		"november",
		"december",
	}

	for column, header := range headers {
		cell, err := excelize.CoordinatesToCellName(column+1, 1)
		if err != nil {
			return err
		}

		err = file.SetCellValue(sheet, cell, header)
		if err != nil {
			return err
		}
	}

	dateStyle, err := file.NewStyle(&excelize.Style{
		NumFmt: 14,
	})
	if err != nil {
		return err
	}

	headerStyle, err := file.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Color: "FFFFFF",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"4472C4"},
			Pattern: 1,
		},
	})
	if err != nil {
		return err
	}

	for rowIndex, customer := range customers {
		row := rowIndex + 2

		values := []any{
			customer.ID,
			customer.CompanyName,
			customer.CompanyType,
			customer.Phone,
			customer.Email,
			customer.MonthlyFee,
			customer.BillingStartedAt,
			customer.Comments,
		}

		for column, value := range values {
			cell, err := excelize.CoordinatesToCellName(column+1, row)
			if err != nil {
				return err
			}

			err = file.SetCellValue(sheet, cell, value)
			if err != nil {
				return err
			}
		}

		billingCell, err := excelize.CoordinatesToCellName(7, row)
		if err != nil {
			return err
		}

		err = file.SetCellStyle(sheet, billingCell, billingCell, dateStyle)
		if err != nil {
			return err
		}

		for month := 1; month <= 12; month++ {
			cell, err := excelize.CoordinatesToCellName(8+month, row)
			if err != nil {
				return err
			}

			paidAt, ok := paymentDates[customer.ID][year][month]
			if !ok {
				err = file.SetCellValue(sheet, cell, "")
				if err != nil {
					return err
				}

				continue
			}

			err = file.SetCellValue(sheet, cell, paidAt)
			if err != nil {
				return err
			}

			err = file.SetCellStyle(sheet, cell, cell, dateStyle)
			if err != nil {
				return err
			}
		}
	}

	lastRow := len(customers) + 1
	lastColumn, err := excelize.ColumnNumberToName(len(headers))
	if err != nil {
		return err
	}

	filterRange := fmt.Sprintf("A1:%s%d", lastColumn, lastRow)
	err = file.SetCellStyle(sheet, "A1", fmt.Sprintf("%s1", lastColumn), headerStyle)
	if err != nil {
		return err
	}

	if lastRow > 1 {
		err = file.SetCellStyle(sheet, "G2", fmt.Sprintf("G%d", lastRow), dateStyle)
		if err != nil {
			return err
		}

		err = file.SetCellStyle(sheet, "I2", fmt.Sprintf("%s%d", lastColumn, lastRow), dateStyle)
		if err != nil {
			return err
		}
	}

	err = file.AutoFilter(sheet, filterRange, nil)
	if err != nil {
		return err
	}

	err = file.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	})
	if err != nil {
		return err
	}

	for column := 1; column <= len(headers); column++ {
		name, err := excelize.ColumnNumberToName(column)
		if err != nil {
			return err
		}

		width := 14.0
		if column == 2 || column == 8 {
			width = 32
		}

		err = file.SetColWidth(sheet, name, name, width)
		if err != nil {
			return err
		}
	}

	return nil
}

func workbookYears(customers []Customer, payments []Payment) []int {
	yearSet := map[int]bool{}

	for _, customer := range customers {
		yearSet[customer.BillingStartedAt.Year()] = true
	}

	for _, payment := range payments {
		yearSet[payment.Year] = true
	}

	if len(yearSet) == 0 {
		yearSet[time.Now().Year()] = true
	}

	minYear := time.Now().Year()
	maxYear := minYear
	first := true

	for year := range yearSet {
		if first || year < minYear {
			minYear = year
		}
		if first || year > maxYear {
			maxYear = year
		}
		first = false
	}

	years := make([]int, 0, maxYear-minYear+1)
	for year := minYear; year <= maxYear; year++ {
		years = append(years, year)
	}

	return years
}

func paymentDatesByCustomer(payments []Payment) map[int64]map[int]map[int]time.Time {
	paymentDates := map[int64]map[int]map[int]time.Time{}

	for _, payment := range payments {
		if paymentDates[payment.CustomerID] == nil {
			paymentDates[payment.CustomerID] = map[int]map[int]time.Time{}
		}

		if paymentDates[payment.CustomerID][payment.Year] == nil {
			paymentDates[payment.CustomerID][payment.Year] = map[int]time.Time{}
		}

		paymentDates[payment.CustomerID][payment.Year][payment.Month] = payment.PaidAt
	}

	return paymentDates
}

func defaultOutputPath() string {
	return filepath.Join("exports", fmt.Sprintf("customers-%s.xlsx", time.Now().Format("2006-01-02-150405")))
}
