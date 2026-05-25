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
	dueDay        int
	sortBy        string
	sortDirection string
	companyTypes  []string
	limit         int
	offset        int
}

type searchParams struct {
	companyName string
	limit       int
	offset      int
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

	return debtListParams{
		dueDay:        debt.dueDay,
		sortBy:        sortBy,
		sortDirection: sortDirection,
		companyTypes:  companyTypes,
		limit:         list.limit,
		offset:        list.offset,
	}, nil
}

func parseSearchParams(r *http.Request) (searchParams, error) {
	companyName := strings.TrimSpace(r.URL.Query().Get("companyName"))
	if companyName == "" {
		return searchParams{}, fmt.Errorf("companyName is required")
	}

	list, err := parseListParams(r)
	if err != nil {
		return searchParams{}, err
	}

	return searchParams{
		companyName: companyName,
		limit:       list.limit,
		offset:      list.offset,
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
