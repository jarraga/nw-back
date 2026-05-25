package customers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"nw-back/internal/postgres/db"

	"github.com/go-chi/chi/v5"
)

const (
	defaultLimit = 50
	maxLimit     = 100
)

type listParams struct {
	limit  int
	offset int
}

type debtParams struct {
	dueDay int
}

type debtListParams struct {
	dueDay          int
	sortBy          string
	sortDirection   string
	companyName     string
	companyTypes    []string
	includeReviewed bool
	limit           int
	offset          int
}

type monthlyDelinquencyParams struct {
	year   int
	dueDay int
}

func parseListParams(r *http.Request) (listParams, error) {
	limit, err := queryInt(r, "limit", defaultLimit)
	if err != nil {
		return listParams{}, err
	}

	offset, err := queryInt(r, "offset", 0)
	if err != nil {
		return listParams{}, err
	}

	if limit > maxLimit {
		limit = maxLimit
	}

	return listParams{
		limit:  limit,
		offset: offset,
	}, nil
}

func parseDebtParams(r *http.Request) (debtParams, error) {
	dueDay, err := queryInt(r, "dueDay", 10)
	if err != nil {
		return debtParams{}, err
	}

	if dueDay < 1 || dueDay > 31 {
		return debtParams{}, fmt.Errorf("dueDay must be between 1 and 31")
	}

	return debtParams{
		dueDay: dueDay,
	}, nil
}

func parseDebtListParams(r *http.Request) (debtListParams, error) {
	debt, err := parseDebtParams(r)
	if err != nil {
		return debtListParams{}, err
	}

	sortBy := r.URL.Query().Get("sortBy")
	if sortBy == "" {
		sortBy = "amount"
	}

	if sortBy != "amount" && sortBy != "months" {
		return debtListParams{}, fmt.Errorf("sortBy must be amount or months")
	}

	sortDirection := r.URL.Query().Get("sortDirection")
	if sortDirection == "" {
		sortDirection = "desc"
	}

	if sortDirection != "asc" && sortDirection != "desc" {
		return debtListParams{}, fmt.Errorf("sortDirection must be asc or desc")
	}

	list, err := parseListParams(r)
	if err != nil {
		return debtListParams{}, err
	}

	companyTypes, err := parseCompanyTypes(r)
	if err != nil {
		return debtListParams{}, err
	}

	companyName := strings.TrimSpace(r.URL.Query().Get("companyName"))
	includeReviewed, err := queryBool(r, "includeReviewed", false)
	if err != nil {
		return debtListParams{}, err
	}

	return debtListParams{
		dueDay:          debt.dueDay,
		sortBy:          sortBy,
		sortDirection:   sortDirection,
		companyName:     companyName,
		companyTypes:    companyTypes,
		includeReviewed: includeReviewed,
		limit:           list.limit,
		offset:          list.offset,
	}, nil
}

func parseMonthlyDelinquencyParams(r *http.Request) (monthlyDelinquencyParams, error) {
	year, err := requiredQueryInt(r, "year")
	if err != nil {
		return monthlyDelinquencyParams{}, err
	}

	if year < 1900 || year > 2100 {
		return monthlyDelinquencyParams{}, fmt.Errorf("year must be between 1900 and 2100")
	}

	debt, err := parseDebtParams(r)
	if err != nil {
		return monthlyDelinquencyParams{}, err
	}

	return monthlyDelinquencyParams{
		year:   year,
		dueDay: debt.dueDay,
	}, nil
}

func parseCompanyTypes(r *http.Request) ([]string, error) {
	values := r.URL.Query()["companyType"]
	companyTypes := []string{}

	for _, value := range values {
		for _, companyType := range strings.Split(value, ",") {
			companyType = strings.TrimSpace(companyType)
			if companyType == "" {
				continue
			}

			if companyType != "enterprise" && companyType != "pyme" && companyType != "startup" {
				return nil, fmt.Errorf("companyType must be enterprise, pyme or startup")
			}

			companyTypes = append(companyTypes, companyType)
		}
	}

	return companyTypes, nil
}

func parseCompanyType(value string) (db.CompanyType, error) {
	companyType := db.CompanyType(strings.TrimSpace(value))

	switch companyType {
	case db.CompanyTypeEnterprise,
		db.CompanyTypePyme,
		db.CompanyTypeStartup:
		return companyType, nil
	default:
		return "", fmt.Errorf("companyType must be enterprise, pyme or startup")
	}
}

func parseCustomerID(r *http.Request) (int64, error) {
	value := chi.URLParam(r, "customerID")
	if value == "" {
		return 0, fmt.Errorf("customerID is required")
	}

	customerID, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("customerID must be a number")
	}

	if customerID <= 0 {
		return 0, fmt.Errorf("customerID must be greater than 0")
	}

	return customerID, nil
}

func parseActionID(r *http.Request) (int64, error) {
	value := chi.URLParam(r, "actionID")
	if value == "" {
		return 0, fmt.Errorf("actionID is required")
	}

	actionID, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("actionID must be a number")
	}

	if actionID <= 0 {
		return 0, fmt.Errorf("actionID must be greater than 0")
	}

	return actionID, nil
}

func parseCustomerActionType(value string) (db.CustomerActionType, error) {
	actionType := db.CustomerActionType(strings.TrimSpace(value))

	switch actionType {
	case db.CustomerActionTypeCall,
		db.CustomerActionTypeEmail,
		db.CustomerActionTypePersonalVisit,
		db.CustomerActionTypeOther:
		return actionType, nil
	default:
		return "", fmt.Errorf("type must be call, email, personal_visit or other")
	}
}

func parsePaymentStatus(value string) (db.PaymentStatus, error) {
	if strings.TrimSpace(value) == "" {
		return db.PaymentStatusPaid, nil
	}

	status := db.PaymentStatus(strings.TrimSpace(value))

	switch status {
	case db.PaymentStatusPending,
		db.PaymentStatusPaid:
		return status, nil
	default:
		return "", fmt.Errorf("status must be pending or paid")
	}
}

func requiredQueryInt(r *http.Request, key string) (int, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return 0, fmt.Errorf("%s is required", key)
	}

	number, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be a number", key)
	}

	return number, nil
}

func queryInt(r *http.Request, key string, fallback int) (int, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return fallback, nil
	}

	number, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be a number", key)
	}

	if number < 0 {
		return 0, fmt.Errorf("%s must be greater than or equal to 0", key)
	}

	return number, nil
}

func queryBool(r *http.Request, key string, fallback bool) (bool, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return fallback, nil
	}

	boolean, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("%s must be true or false", key)
	}

	return boolean, nil
}
