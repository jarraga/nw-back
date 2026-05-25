package customers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"nw-back/internal/postgres/db"

	"github.com/jackc/pgx/v5/pgtype"
)

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var request createCustomerRequest

	decoder := json.NewDecoder(r.Body)
	decoder.UseNumber()

	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	params, err := newCreateCustomerParams(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	customer, err := h.queries.CreateCustomer(r.Context(), params)
	if err != nil {
		http.Error(w, "failed to create customer", http.StatusInternalServerError)
		return
	}

	response, err := newCustomerResponse(customer)
	if err != nil {
		http.Error(w, "failed to build customer response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func newCreateCustomerParams(request createCustomerRequest) (db.CreateCustomerParams, error) {
	companyName := strings.TrimSpace(request.CompanyName)
	if companyName == "" {
		return db.CreateCustomerParams{}, fmt.Errorf("companyName is required")
	}

	companyType, err := parseCompanyType(request.CompanyType)
	if err != nil {
		return db.CreateCustomerParams{}, err
	}

	phone := strings.TrimSpace(request.Phone)
	if phone == "" {
		return db.CreateCustomerParams{}, fmt.Errorf("phone is required")
	}

	email := strings.TrimSpace(request.Email)
	if email == "" {
		return db.CreateCustomerParams{}, fmt.Errorf("email is required")
	}

	monthlyFee, err := parseMonthlyFee(request.MonthlyFee)
	if err != nil {
		return db.CreateCustomerParams{}, err
	}

	billingStartedAt, err := parseBillingStartedAt(request.BillingStartedAt)
	if err != nil {
		return db.CreateCustomerParams{}, err
	}

	return db.CreateCustomerParams{
		CompanyName:      companyName,
		CompanyType:      companyType,
		Phone:            phone,
		Email:            email,
		MonthlyFee:       monthlyFee,
		BillingStartedAt: billingStartedAt,
		Comments:         strings.TrimSpace(request.Comments),
	}, nil
}

func parseMonthlyFee(value json.Number) (pgtype.Numeric, error) {
	if value.String() == "" {
		return pgtype.Numeric{}, fmt.Errorf("monthlyFee is required")
	}

	number, err := value.Float64()
	if err != nil {
		return pgtype.Numeric{}, fmt.Errorf("monthlyFee must be a number")
	}

	if number <= 0 {
		return pgtype.Numeric{}, fmt.Errorf("monthlyFee must be greater than 0")
	}

	monthlyFee := pgtype.Numeric{}

	err = monthlyFee.Scan(value.String())
	if err != nil {
		return pgtype.Numeric{}, fmt.Errorf("invalid monthlyFee")
	}

	return monthlyFee, nil
}

func parseBillingStartedAt(value string) (pgtype.Date, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return pgtype.Date{}, fmt.Errorf("billingStartedAt is required")
	}

	date, err := time.Parse("2006-01-02", value)
	if err != nil {
		return pgtype.Date{}, fmt.Errorf("billingStartedAt must use YYYY-MM-DD format")
	}

	return pgtype.Date{
		Time:  date,
		Valid: true,
	}, nil
}
