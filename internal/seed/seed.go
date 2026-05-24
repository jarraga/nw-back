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

	err = resetDatabase(ctx)
	if err != nil {
		return err
	}

	log.Println("database reset OK")

	queries := db.New(postgres.DB)

	customers, err := createCustomers(ctx, queries, config)
	if err != nil {
		return err
	}

	err = createCustomerPayments(ctx, queries, customers, config)
	if err != nil {
		return err
	}

	return nil
}

func createCustomers(ctx context.Context, queries *db.Queries, config Config) ([]db.Customer, error) {
	gofakeit.Seed(0)
	customers := make([]db.Customer, 0, config.ActiveCustomers)

	for range config.ActiveCustomers {
		customer, err := randomCustomer(config)
		if err != nil {
			return nil, err
		}

		createdCustomer, err := queries.CreateCustomer(ctx, customer)
		if err != nil {
			return nil, err
		}

		customers = append(customers, createdCustomer)
	}

	log.Printf("%d customers created", config.ActiveCustomers)
	return customers, nil
}

func randomCustomer(config Config) (db.CreateCustomerParams, error) {
	monthlyFee, err := randomMonthlyFee(config)
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

func randomMonthlyFee(config Config) (pgtype.Numeric, error) {
	monthlyFee := pgtype.Numeric{}
	value := gofakeit.Number(config.MonthlyFeeFrom, config.MonthlyFeeTo)
	formattedValue := strconv.Itoa(value)

	err := monthlyFee.Scan(formattedValue)
	if err != nil {
		return pgtype.Numeric{}, fmt.Errorf("invalid monthly fee: %w", err)
	}

	return monthlyFee, nil
}
