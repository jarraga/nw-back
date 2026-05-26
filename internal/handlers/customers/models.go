package customers

import (
	"encoding/json"
	"time"

	"nw-back/internal/postgres/db"

	"github.com/jackc/pgx/v5/pgtype"
)

type response struct {
	ID               int64          `json:"id"`
	CompanyName      string         `json:"companyName"`
	CompanyType      db.CompanyType `json:"companyType"`
	Phone            string         `json:"phone"`
	Email            string         `json:"email"`
	MonthlyFee       int32          `json:"monthlyFee"`
	BillingStartedAt time.Time      `json:"billingStartedAt"`
	Comments         string         `json:"comments"`
	Review           reviewResponse `json:"review"`
	CreatedAt        time.Time      `json:"createdAt"`
}

type paginatedResponse[T any] struct {
	Items []T   `json:"items"`
	Total int32 `json:"total"`
}

type debtResponse struct {
	TotalDebt int64 `json:"totalDebt"`
}

type debtListResponse struct {
	ID               int64          `json:"id"`
	CompanyName      string         `json:"companyName"`
	CompanyType      db.CompanyType `json:"companyType"`
	MonthlyFee       int32          `json:"monthlyFee"`
	BillingStartedAt time.Time      `json:"billingStartedAt"`
	Comments         string         `json:"comments"`
	Review           reviewResponse `json:"review"`
	OverdueMonths    int32          `json:"overdueMonths"`
	OverdueAmount    int64          `json:"overdueAmount"`
}

type createCustomerRequest struct {
	CompanyName      string      `json:"companyName"`
	CompanyType      string      `json:"companyType"`
	Phone            string      `json:"phone"`
	Email            string      `json:"email"`
	MonthlyFee       json.Number `json:"monthlyFee"`
	BillingStartedAt string      `json:"billingStartedAt"`
	Comments         string      `json:"comments"`
}

type createActionRequest struct {
	Type          string  `json:"type"`
	Comments      string  `json:"comments"`
	InformantName *string `json:"informantName"`
}

type updateActionCommentsRequest struct {
	Comments string `json:"comments"`
}

type createPaymentRequest struct {
	Year   int32   `json:"year"`
	Month  int32   `json:"month"`
	Status string  `json:"status"`
	PaidAt *string `json:"paidAt"`
}

type updateCommentsRequest struct {
	Comments string `json:"comments"`
}

type updateCustomerRequest struct {
	Phone      string      `json:"phone"`
	Email      string      `json:"email"`
	MonthlyFee json.Number `json:"monthlyFee"`
}

type reviewCustomerRequest struct {
	Days       int32  `json:"days"`
	ReviewedBy string `json:"reviewedBy"`
}

type reviewResponse struct {
	ReviewedAt    *time.Time `json:"reviewedAt"`
	ReviewedUntil *time.Time `json:"reviewedUntil"`
	ReviewedBy    *string    `json:"reviewedBy"`
	IsReviewed    bool       `json:"isReviewed"`
}

type actionResponse struct {
	ID            int64                 `json:"id"`
	CustomerID    int64                 `json:"customerID"`
	Type          db.CustomerActionType `json:"type"`
	Comments      string                `json:"comments"`
	InformantName *string               `json:"informantName"`
	ActionDate    time.Time             `json:"actionDate"`
	CreatedAt     time.Time             `json:"createdAt"`
}

type paymentResponse struct {
	ID         int64            `json:"id"`
	CustomerID int64            `json:"customerID"`
	Year       int32            `json:"year"`
	Month      int32            `json:"month"`
	Status     db.PaymentStatus `json:"status"`
	PaidAt     *time.Time       `json:"paidAt"`
	CreatedAt  time.Time        `json:"createdAt"`
}

type detailResponse struct {
	Customer response          `json:"customer"`
	Actions  []actionResponse  `json:"actions"`
	Payments []paymentResponse `json:"payments"`
	Debt     customerDebt      `json:"debt"`
}

type customerDebt struct {
	DueDay        int32 `json:"dueDay"`
	OverdueMonths int32 `json:"overdueMonths"`
	OverdueAmount int64 `json:"overdueAmount"`
}

type monthlyDelinquencyItemResponse struct {
	Month                 int32   `json:"month"`
	TotalCustomers        int32   `json:"totalCustomers"`
	OverdueCustomers      int32   `json:"overdueCustomers"`
	DelinquencyPercentage float64 `json:"delinquencyPercentage"`
}

type monthlyDelinquencyResponse struct {
	Year   int32                            `json:"year"`
	DueDay int32                            `json:"dueDay"`
	Items  []monthlyDelinquencyItemResponse `json:"items"`
}

type delinquencyRateResponse struct {
	DueDay                int32   `json:"dueDay"`
	TotalCustomers        int32   `json:"totalCustomers"`
	OverdueCustomers      int32   `json:"overdueCustomers"`
	DelinquencyPercentage float64 `json:"delinquencyPercentage"`
}

type customerMetricsResponse struct {
	DueDay         int32                        `json:"dueDay"`
	TotalCustomers int32                        `json:"totalCustomers"`
	CompanyTypes   []metricsCompanyTypeResponse `json:"companyTypes"`
	Debtors        metricsDebtorsResponse       `json:"debtors"`
}

type metricsCompanyTypeResponse struct {
	CompanyType string  `json:"companyType"`
	Customers   int32   `json:"customers"`
	Percentage  float64 `json:"percentage"`
}

type metricsDebtorsResponse struct {
	Customers     int32                        `json:"customers"`
	Percentage    float64                      `json:"percentage"`
	ByCompanyType []metricsCompanyTypeResponse `json:"byCompanyType"`
}

func newCustomerResponse(customer db.Customer) (response, error) {
	return response{
		ID:               customer.ID,
		CompanyName:      customer.CompanyName,
		CompanyType:      customer.CompanyType,
		Phone:            customer.Phone,
		Email:            customer.Email,
		MonthlyFee:       customer.MonthlyFee,
		BillingStartedAt: customer.BillingStartedAt.Time,
		Comments:         customer.Comments,
		Review:           newReviewResponse(customer.ReviewedAt, customer.ReviewedUntil, customer.ReviewedBy),
		CreatedAt:        customer.CreatedAt.Time,
	}, nil
}

func newListResponse(customers []db.Customer) ([]response, error) {
	items := make([]response, 0, len(customers))

	for _, customer := range customers {
		item, err := newCustomerResponse(customer)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
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

func newDebtResponse(totalDebt int64) debtResponse {
	return debtResponse{
		TotalDebt: totalDebt,
	}
}

func newDebtListResponse(customers []db.ListCustomersDebtRow) ([]debtListResponse, error) {
	items := make([]debtListResponse, 0, len(customers))

	for _, customer := range customers {
		items = append(items, debtListResponse{
			ID:               customer.ID,
			CompanyName:      customer.CompanyName,
			CompanyType:      customer.CompanyType,
			MonthlyFee:       customer.MonthlyFee,
			BillingStartedAt: customer.BillingStartedAt.Time,
			Comments:         customer.Comments,
			Review:           newReviewResponse(customer.ReviewedAt, customer.ReviewedUntil, customer.ReviewedBy),
			OverdueMonths:    customer.OverdueMonths,
			OverdueAmount:    customer.OverdueAmount,
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

func newReviewResponse(reviewedAt pgtype.Timestamptz, reviewedUntil pgtype.Timestamptz, reviewedBy pgtype.Text) reviewResponse {
	var reviewedAtValue *time.Time
	if reviewedAt.Valid {
		reviewedAtValue = &reviewedAt.Time
	}

	var reviewedUntilValue *time.Time
	if reviewedUntil.Valid {
		reviewedUntilValue = &reviewedUntil.Time
	}

	var reviewedByValue *string
	if reviewedBy.Valid {
		reviewedByValue = &reviewedBy.String
	}

	return reviewResponse{
		ReviewedAt:    reviewedAtValue,
		ReviewedUntil: reviewedUntilValue,
		ReviewedBy:    reviewedByValue,
		IsReviewed:    reviewedUntil.Valid && reviewedUntil.Time.After(time.Now()),
	}
}

func newActionResponse(action db.CustomerAction) actionResponse {
	var informantName *string

	if action.InformantName.Valid {
		informantName = &action.InformantName.String
	}

	return actionResponse{
		ID:            action.ID,
		CustomerID:    action.CustomerID,
		Type:          action.Type,
		Comments:      action.Comments,
		InformantName: informantName,
		ActionDate:    action.ActionDate.Time,
		CreatedAt:     action.CreatedAt.Time,
	}
}

func newActionResponses(actions []db.CustomerAction) []actionResponse {
	items := make([]actionResponse, 0, len(actions))

	for _, action := range actions {
		items = append(items, newActionResponse(action))
	}

	return items
}

func newPaymentResponse(payment db.CustomerPayment) paymentResponse {
	var paidAt *time.Time

	if payment.PaidAt.Valid {
		paidAt = &payment.PaidAt.Time
	}

	return paymentResponse{
		ID:         payment.ID,
		CustomerID: payment.CustomerID,
		Year:       payment.Year,
		Month:      payment.Month,
		Status:     payment.Status,
		PaidAt:     paidAt,
		CreatedAt:  payment.CreatedAt.Time,
	}
}

func newPaymentResponses(payments []db.CustomerPayment) []paymentResponse {
	items := make([]paymentResponse, 0, len(payments))

	for _, payment := range payments {
		items = append(items, newPaymentResponse(payment))
	}

	return items
}

func newDetailResponse(customer db.Customer, actions []db.CustomerAction, payments []db.CustomerPayment, debt db.GetCustomerDebtSummaryRow, dueDay int32) (detailResponse, error) {
	customerResponse, err := newCustomerResponse(customer)
	if err != nil {
		return detailResponse{}, err
	}

	return detailResponse{
		Customer: customerResponse,
		Actions:  newActionResponses(actions),
		Payments: newPaymentResponses(payments),
		Debt: customerDebt{
			DueDay:        dueDay,
			OverdueMonths: debt.OverdueMonths,
			OverdueAmount: debt.OverdueAmount,
		},
	}, nil
}

func newMonthlyDelinquencyResponse(year int32, dueDay int32, rows []db.GetMonthlyDelinquencyRateRow) monthlyDelinquencyResponse {
	items := make([]monthlyDelinquencyItemResponse, 0, len(rows))

	for _, row := range rows {
		items = append(items, monthlyDelinquencyItemResponse{
			Month:                 row.Month,
			TotalCustomers:        row.TotalCustomers,
			OverdueCustomers:      row.OverdueCustomers,
			DelinquencyPercentage: row.DelinquencyPercentage,
		})
	}

	return monthlyDelinquencyResponse{
		Year:   year,
		DueDay: dueDay,
		Items:  items,
	}
}

func newDelinquencyRateResponse(dueDay int32, row db.GetCustomerDelinquencyRateRow) delinquencyRateResponse {
	return delinquencyRateResponse{
		DueDay:                dueDay,
		TotalCustomers:        row.TotalCustomers,
		OverdueCustomers:      row.OverdueCustomers,
		DelinquencyPercentage: row.DelinquencyPercentage,
	}
}

func newCustomerMetricsResponse(dueDay int32, rows []db.GetCustomerMetricsRow) customerMetricsResponse {
	response := customerMetricsResponse{
		DueDay:       dueDay,
		CompanyTypes: []metricsCompanyTypeResponse{},
		Debtors: metricsDebtorsResponse{
			ByCompanyType: []metricsCompanyTypeResponse{},
		},
	}

	if len(rows) == 0 {
		return response
	}

	response.TotalCustomers = rows[0].TotalCustomers
	response.Debtors.Customers = rows[0].DebtorCustomers
	response.Debtors.Percentage = rows[0].DebtorPercentage

	for _, row := range rows {
		response.CompanyTypes = append(response.CompanyTypes, metricsCompanyTypeResponse{
			CompanyType: row.CompanyType,
			Customers:   row.TypeCustomers,
			Percentage:  row.TypePercentage,
		})

		response.Debtors.ByCompanyType = append(response.Debtors.ByCompanyType, metricsCompanyTypeResponse{
			CompanyType: row.CompanyType,
			Customers:   row.TypeDebtorCustomers,
			Percentage:  row.TypeDebtorPercentage,
		})
	}

	return response
}
