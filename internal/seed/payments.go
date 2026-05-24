package seed

import (
	"context"
	"log"
	"time"

	"nw-back/internal/postgres/db"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5/pgtype"
)

func createCustomerPayments(ctx context.Context, queries *db.Queries, customers []db.Customer, config Config) error {
	dataTo := dataToMonth()
	paymentsCreated := 0

	for _, customer := range customers {
		dataFrom := customer.BillingStartedAt.Time
		payments, err := randomPayments(customer, dataFrom, dataTo, config)
		if err != nil {
			return err
		}

		for _, payment := range payments {
			_, err = queries.CreateCustomerPayment(ctx, payment)
			if err != nil {
				return err
			}

			paymentsCreated++
		}
	}

	log.Printf("%d customer payments created", paymentsCreated)
	return nil
}

func randomPayments(customer db.Customer, dataFrom time.Time, dataTo time.Time, config Config) ([]db.CreateCustomerPaymentParams, error) {
	payments := []db.CreateCustomerPaymentParams{}
	currentMonth := dataFrom
	now := time.Now()
	lastPaidAt := time.Time{}

	for !currentMonth.After(dataTo) {
		payment, paidAt, exists := randomPayment(customer, currentMonth, lastPaidAt, now, config)
		if !exists {
			currentMonth = currentMonth.AddDate(0, 1, 0)
			continue
		}

		payments = append(payments, payment)

		lastPaidAt = paidAt.Time

		currentMonth = currentMonth.AddDate(0, 1, 0)
	}

	return payments, nil
}

func randomPayment(customer db.Customer, month time.Time, lastPaidAt time.Time, now time.Time, config Config) (db.CreateCustomerPaymentParams, pgtype.Timestamptz, bool) {
	paidAt := randomPaidAt(customer.CompanyType, month, lastPaidAt, config)

	if paidAt.After(now) {
		return db.CreateCustomerPaymentParams{}, pgtype.Timestamptz{}, false
	}

	timestamptz := pgtype.Timestamptz{
		Time:  paidAt,
		Valid: true,
	}

	return db.CreateCustomerPaymentParams{
		CustomerID: customer.ID,
		Year:       int32(month.Year()),
		Month:      int32(month.Month()),
		Status:     db.PaymentStatusPaid,
		PaidAt:     timestamptz,
	}, timestamptz, true
}

func randomPaidAt(companyType db.CompanyType, month time.Time, lastPaidAt time.Time, config Config) time.Time {
	fromDays := config.GeneralPaymentDelayFromDays
	toDays := config.GeneralPaymentDelayToDays

	if companyType == db.CompanyTypeEnterprise {
		fromDays = config.EnterprisePaymentDelayFromDays
		toDays = config.EnterprisePaymentDelayToDays
	}

	delayDays := gofakeit.Number(fromDays, toDays)
	paidAt := month.AddDate(0, 0, delayDays)

	if !lastPaidAt.IsZero() && !paidAt.After(lastPaidAt) {
		daysAfterLastPayment := gofakeit.Number(1, 10)
		paidAt = lastPaidAt.AddDate(0, 0, daysAfterLastPayment)
	}

	return paidAt
}

func dataToMonth() time.Time {
	now := time.Now().UTC()
	return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
}
