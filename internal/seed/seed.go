package seed

import (
	"context"
	"log"
	"time"

	"nw-back/internal/postgres"
	"nw-back/internal/postgres/db"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5/pgtype"
)

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

	customers, err = markReviewedDebtorCustomers(ctx, queries, customers, config)
	if err != nil {
		return err
	}

	err = createCustomerActions(ctx, queries, customers)
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

func markReviewedDebtorCustomers(ctx context.Context, queries *db.Queries, customers []db.Customer, config Config) ([]db.Customer, error) {
	if config.ReviewedCustomersPercentage <= 0 {
		return customers, nil
	}

	reviewedCustomers := 0

	for index, customer := range customers {
		debt, err := queries.GetCustomerDebtSummary(ctx, db.GetCustomerDebtSummaryParams{
			DueDay:     int32(config.DueDay),
			CustomerID: customer.ID,
		})
		if err != nil {
			return nil, err
		}

		if debt.OverdueAmount == 0 {
			continue
		}

		if gofakeit.Number(1, 100) > config.ReviewedCustomersPercentage {
			continue
		}

		reviewedCustomer, err := queries.MarkCustomerReviewed(ctx, db.MarkCustomerReviewedParams{
			ID:         customer.ID,
			Days:       int32(gofakeit.Number(5, 10)),
			ReviewedBy: randomReviewerName(),
		})
		if err != nil {
			return nil, err
		}

		customers[index] = reviewedCustomer
		reviewedCustomers++
	}

	log.Printf("%d debtor customers marked as reviewed", reviewedCustomers)
	return customers, nil
}

func randomReviewerName() string {
	names := []string{
		"Martina Alvarez",
		"Lucas Pereyra",
		"Valentina Sosa",
		"Nicolas Romero",
		"Camila Herrera",
		"Federico Molina",
		"Agustina Medina",
		"Joaquin Navarro",
		"Florencia Rivas",
		"Matias Cabrera",
	}

	return names[gofakeit.Number(0, len(names)-1)]
}

func randomCustomer(config Config) (db.CreateCustomerParams, error) {
	monthlyFee, err := randomMonthlyFee(config)
	if err != nil {
		return db.CreateCustomerParams{}, err
	}

	billingStartedAt := randomBillingStartedAt(config)

	return db.CreateCustomerParams{
		CompanyName:      gofakeit.Company(),
		CompanyType:      randomCompanyType(config),
		Phone:            gofakeit.Phone(),
		Email:            gofakeit.Email(),
		MonthlyFee:       monthlyFee,
		BillingStartedAt: billingStartedAt,
		Comments:         randomCustomerComment(),
	}, nil
}

func randomCustomerComment() string {
	comments := []string{
		"Cliente con buena predisposicion para resolver consultas administrativas.",
		"Suele responder mejor por la manana y prefiere mensajes breves.",
		"Empresa con crecimiento reciente y necesidades operativas cambiantes.",
		"Contacto principal atento, aunque requiere seguimiento para cerrar temas.",
		"Cliente sensible a cambios de precio y condiciones comerciales.",
		"Cuenta con procesos internos formales para aprobar pagos y novedades.",
		"Prefiere centralizar la comunicacion en una sola persona del equipo.",
		"Cliente historico con uso estable del servicio contratado.",
		"Necesita recordatorios periodicos para mantener documentacion al dia.",
		"Empresa con buena relacion comercial y baja friccion operativa.",
	}

	return comments[gofakeit.Number(0, len(comments)-1)]
}

func randomCompanyType(config Config) db.CompanyType {
	totalWeight := config.EnterpriseCompanyWeight + config.PymeCompanyWeight + config.StartupCompanyWeight
	randomWeight := gofakeit.Number(1, totalWeight)

	if randomWeight <= config.EnterpriseCompanyWeight {
		return db.CompanyTypeEnterprise
	}

	if randomWeight <= config.EnterpriseCompanyWeight+config.PymeCompanyWeight {
		return db.CompanyTypePyme
	}

	return db.CompanyTypeStartup
}

func randomMonthlyFee(config Config) (int32, error) {
	value := gofakeit.Number(config.MonthlyFeeFrom, config.MonthlyFeeTo)

	return int32(value), nil
}

func randomBillingStartedAt(config Config) pgtype.Date {
	firstMonth := time.Date(config.DataFromYear, time.January, 1, 0, 0, 0, 0, time.UTC)
	randomMonths := gofakeit.Number(0, config.CustomerStartVariationMonths)

	return pgtype.Date{
		Time:  firstMonth.AddDate(0, randomMonths, 0),
		Valid: true,
	}
}
