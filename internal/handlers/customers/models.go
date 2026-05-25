package customers

import (
	"time"

	"nw-back/internal/postgres/db"
)

type response struct {
	ID               int64          `json:"id"`
	CompanyName      string         `json:"companyName"`
	CompanyType      db.CompanyType `json:"companyType"`
	Phone            string         `json:"phone"`
	Email            string         `json:"email"`
	MonthlyFee       float64        `json:"monthlyFee"`
	BillingStartedAt time.Time      `json:"billingStartedAt"`
	CreatedAt        time.Time      `json:"createdAt"`
}

type paginatedResponse[T any] struct {
	Items []T   `json:"items"`
	Total int32 `json:"total"`
}

type debtResponse struct {
	TotalDebt float64 `json:"totalDebt"`
}

type debtListResponse struct {
	ID               int64          `json:"id"`
	CompanyName      string         `json:"companyName"`
	CompanyType      db.CompanyType `json:"companyType"`
	Phone            string         `json:"phone"`
	Email            string         `json:"email"`
	MonthlyFee       float64        `json:"monthlyFee"`
	BillingStartedAt time.Time      `json:"billingStartedAt"`
	OverdueMonths    int32          `json:"overdueMonths"`
	OverdueAmount    float64        `json:"overdueAmount"`
}

func newListResponse(customers []db.Customer) ([]response, error) {
	items := make([]response, 0, len(customers))

	for _, customer := range customers {
		monthlyFee, err := customer.MonthlyFee.Float64Value()
		if err != nil {
			return nil, err
		}

		items = append(items, response{
			ID:               customer.ID,
			CompanyName:      customer.CompanyName,
			CompanyType:      customer.CompanyType,
			Phone:            customer.Phone,
			Email:            customer.Email,
			MonthlyFee:       monthlyFee.Float64,
			BillingStartedAt: customer.BillingStartedAt.Time,
			CreatedAt:        customer.CreatedAt.Time,
		})
	}

	return items, nil
}

func newPaginatedListResponse(customers []db.Customer, total int32) (paginatedResponse[response], error) {
	items, err := newListResponse(customers)
	if err != nil {
		return paginatedResponse[response]{}, err
	}

	return paginatedResponse[response]{
		Items: items,
		Total: total,
	}, nil
}

func newDebtResponse(totalDebt float64) debtResponse {
	return debtResponse{
		TotalDebt: totalDebt,
	}
}

func newDebtListResponse(customers []db.ListCustomersDebtRow) ([]debtListResponse, error) {
	items := make([]debtListResponse, 0, len(customers))

	for _, customer := range customers {
		monthlyFee, err := customer.MonthlyFee.Float64Value()
		if err != nil {
			return nil, err
		}

		overdueAmount, err := customer.OverdueAmount.Float64Value()
		if err != nil {
			return nil, err
		}

		items = append(items, debtListResponse{
			ID:               customer.ID,
			CompanyName:      customer.CompanyName,
			CompanyType:      customer.CompanyType,
			Phone:            customer.Phone,
			Email:            customer.Email,
			MonthlyFee:       monthlyFee.Float64,
			BillingStartedAt: customer.BillingStartedAt.Time,
			OverdueMonths:    customer.OverdueMonths,
			OverdueAmount:    overdueAmount.Float64,
		})
	}

	return items, nil
}

func newPaginatedDebtListResponse(customers []db.ListCustomersDebtRow, total int32) (paginatedResponse[debtListResponse], error) {
	items, err := newDebtListResponse(customers)
	if err != nil {
		return paginatedResponse[debtListResponse]{}, err
	}

	return paginatedResponse[debtListResponse]{
		Items: items,
		Total: total,
	}, nil
}
