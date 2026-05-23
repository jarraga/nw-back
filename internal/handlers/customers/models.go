package customers

import (
	"time"

	"nw-back/internal/postgres/db"
)

type response struct {
	ID          int64          `json:"id"`
	CompanyName string         `json:"companyName"`
	CompanyType db.CompanyType `json:"companyType"`
	Phone       string         `json:"phone"`
	Email       string         `json:"email"`
	MonthlyFee  float64        `json:"monthlyFee"`
	CreatedAt   time.Time      `json:"createdAt"`
}

func newListResponse(customers []db.Customer) ([]response, error) {
	items := make([]response, 0, len(customers))

	for _, customer := range customers {
		monthlyFee, err := customer.MonthlyFee.Float64Value()
		if err != nil {
			return nil, err
		}

		items = append(items, response{
			ID:          customer.ID,
			CompanyName: customer.CompanyName,
			CompanyType: customer.CompanyType,
			Phone:       customer.Phone,
			Email:       customer.Email,
			MonthlyFee:  monthlyFee.Float64,
			CreatedAt:   customer.CreatedAt.Time,
		})
	}

	return items, nil
}
