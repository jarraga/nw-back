package customers

import (
	"fmt"
	"net/http"
	"strconv"
)

const (
	defaultLimit = 50
	maxLimit     = 100
)

type listParams struct {
	limit  int
	offset int
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
