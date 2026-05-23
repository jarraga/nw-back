package seed

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"nw-back/internal/postgres"
	"nw-back/internal/postgres/db"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5/pgtype"
)

const customersAmount = 100

var companyTypes = []db.CompanyType{
	db.CompanyTypeEnterprise,
	db.CompanyTypePyme,
	db.CompanyTypeStartup,
}

func Run(ctx context.Context) error {
	err := postgres.Connect(ctx)
	if err != nil {
		return err
	}
	defer postgres.Close()

	log.Println("postgres connection OK")

	queries := db.New(postgres.DB)

	err = createCustomers(ctx, queries)
	if err != nil {
		return err
	}

	return nil
}

func createCustomers(ctx context.Context, queries *db.Queries) error {
	gofakeit.Seed(0)

	for range customersAmount {
		customer, err := randomCustomer()
		if err != nil {
			return err
		}

		_, err = queries.CreateCustomer(ctx, customer)
		if err != nil {
			return err
		}
	}

	log.Printf("%d customers created", customersAmount)
	return nil
}

func randomCustomer() (db.CreateCustomerParams, error) {
	monthlyFee, err := randomMonthlyFee()
	if err != nil {
		return db.CreateCustomerParams{}, err
	}

	return db.CreateCustomerParams{
		CompanyName: gofakeit.Company(),
		CompanyType: randomCompanyType(),
		Phone:       gofakeit.Phone(),
		Email:       gofakeit.Email(),
		MonthlyFee:  monthlyFee,
	}, nil
}

func randomCompanyType() db.CompanyType {
	index := gofakeit.Number(0, len(companyTypes)-1)
	return companyTypes[index]
}

func randomMonthlyFee() (pgtype.Numeric, error) {
	monthlyFee := pgtype.Numeric{}
	value := gofakeit.Number(200, 15000)
	formattedValue := strconv.Itoa(value)

	err := monthlyFee.Scan(formattedValue)
	if err != nil {
		return pgtype.Numeric{}, fmt.Errorf("invalid monthly fee: %w", err)
	}

	return monthlyFee, nil
}
