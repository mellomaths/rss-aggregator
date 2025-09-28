package models

import (
	"errors"
	"net/http"
	"strconv"
)

type PaginatedParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (p *PaginatedParams) Decode(r *http.Request) error {
	limitString := r.URL.Query().Get("limit")
	offsetString := r.URL.Query().Get("offset")
	limit, err := strconv.Atoi(limitString)
	if err != nil {
		return err
	}
	offset, err := strconv.Atoi(offsetString)
	if err != nil {
		return err
	}
	p.Limit = int32(limit)
	p.Offset = int32(offset)
	return nil
}

func (p *PaginatedParams) Validate() error {
	if p.Limit <= 0 {
		return errors.New("limit must be greater than 0")
	}
	if p.Offset < 0 {
		return errors.New("offset must be greater than 0")
	}
	return nil
}

type Paginated[T any] struct {
	Data   []T `json:"data"`
	Total  int `json:"total"`
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}
