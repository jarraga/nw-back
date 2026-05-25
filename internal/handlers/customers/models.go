package customers

import (
	"encoding/json"
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
	Comments         string         `json:"comments"`
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
	Comments         string         `json:"comments"`
	OverdueMonths    int32          `json:"overdueMonths"`
	OverdueAmount    float64        `json:"overdueAmount"`
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
	Type     string `json:"type"`
	Comments string `json:"comments"`
}

type createPaymentRequest struct {
	Year   int32   `json:"year"`
	Month  int32   `json:"month"`
	Status string  `json:"status"`
	PaidAt *string `json:"paidAt"`
}

type actionResponse struct {
	ID         int64                 `json:"id"`
	CustomerID int64                 `json:"customerID"`
	Type       db.CustomerActionType `json:"type"`
	Comments   string                `json:"comments"`
	ActionDate time.Time             `json:"actionDate"`
	CreatedAt  time.Time             `json:"createdAt"`
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
}

func newCustomerResponse(customer db.Customer) (response, error) {
	monthlyFee, err := customer.MonthlyFee.Float64Value()
	if err != nil {
		return response{}, err
	}

	return response{
		ID:               customer.ID,
		CompanyName:      customer.CompanyName,
		CompanyType:      customer.CompanyType,
		Phone:            customer.Phone,
		Email:            customer.Email,
		MonthlyFee:       monthlyFee.Float64,
		BillingStartedAt: customer.BillingStartedAt.Time,
		Comments:         customer.Comments,
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
			Comments:         customer.Comments,
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

func newActionResponse(action db.CustomerAction) actionResponse {
	return actionResponse{
		ID:         action.ID,
		CustomerID: action.CustomerID,
		Type:       action.Type,
		Comments:   action.Comments,
		ActionDate: action.ActionDate.Time,
		CreatedAt:  action.CreatedAt.Time,
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

func newDetailResponse(customer db.Customer, actions []db.CustomerAction, payments []db.CustomerPayment) (detailResponse, error) {
	customerResponse, err := newCustomerResponse(customer)
	if err != nil {
		return detailResponse{}, err
	}

	return detailResponse{
		Customer: customerResponse,
		Actions:  newActionResponses(actions),
		Payments: newPaymentResponses(payments),
	}, nil
}
