package model

import (
	"github.com/pkg/errors"
	"net/url"
	"strconv"
)

const (
	newerIdField = "newer_id"
	olderIdField = "older_id"
	countField   = "count"
)

var (
	ErrApplyingQueryParamsFailed = errors.New("bunq: failed to apply query parameters")
	ErrNoNextPage                = errors.Wrap(ErrApplyingQueryParamsFailed, "bunq: cannot get next page; there is none.")
	ErrNoPreviousPage            = errors.Wrap(ErrApplyingQueryParamsFailed, "bunq: cannot get previous page; there is none.")
	ErrNoFuturePage              = errors.Wrap(ErrApplyingQueryParamsFailed, "bunq: cannot get future page; there is none.")
	ErrNoPaginationId            = errors.Wrap(ErrApplyingQueryParamsFailed, "bunq: cannot get id for pagination; there is none.")
)

type QueryParam func(query url.Values) error

type Pagination struct {
	FutureURL string `json:"future_url"`
	NewerURL  string `json:"newer_url"`
	OlderURL  string `json:"older_url"`
	count     int
}

// HasNext returns true if there is a next page of results.
// If this is false and NextPage is called, it will query for future pages that may or may not exist.
func (p *Pagination) HasNext() bool {
	return p.NewerURL != ""
}

// HasFuture returns true if the current page has no next page, but there may be future pages.
func (p *Pagination) HasFuture() bool {
	return p.FutureURL != ""
}

// HasPrevious returns true if there is a previous page of results.
// If this is false, PreviousPage will panic with ErrNoNextPage.
func (p *Pagination) HasPrevious() bool {
	return p.OlderURL != ""
}

// SetCount overrides the count of items per page for the current pagination.
func (p *Pagination) SetCount(count int) *Pagination {
	p.count = count
	return p
}

// NextPage returns a QueryParam that can be used to get the next page of a paginated response.
// If Pagination.HasNext() returns false, this function will instead query future pages that may or may not exist.
func (p *Pagination) NextPage() QueryParam {
	if p == nil {
		return nil
	}
	if p.HasNext() {
		return pageParam(p, p.NewerURL, newerIdField, ErrNoNextPage)
	} else if p.HasFuture() {
		return pageParam(p, p.FutureURL, newerIdField, ErrNoFuturePage)
	}
	return func(query url.Values) error {
		return ErrNoNextPage
	}
}

// PreviousPage returns a QueryParam that can be used to get the previous page of a paginated response.
// If Pagination.HasPrevious() returns false, this function will panic with ErrNoPreviousPage.
func (p *Pagination) PreviousPage() QueryParam {
	if p == nil {
		return nil
	}
	//if !p.HasPrevious() {
	//	panic(ErrNoNextPage)
	//}
	return pageParam(p, p.OlderURL, olderIdField, ErrNoPreviousPage)
}

func pageParam(pagination *Pagination, uri, field string, errno error) QueryParam {
	// uri is nil if there is no next/future page
	if errno != nil && uri == "" {
		return func(query url.Values) error {
			return errno
		}
	}

	parsedURL, err := url.Parse(uri)
	if err != nil {
		panic(errors.Wrap(err, "failed to parse url"))
	}

	values := parsedURL.Query()

	// If there is no id, we cannot paginate
	if !values.Has(field) {
		return func(query url.Values) error {
			return ErrNoPaginationId
		}
	}

	id := values.Get(field)
	count := values.Get(countField)

	return func(query url.Values) error {
		query.Set(field, id)

		// Set count if it was previously set
		if count != "" {
			query.Set(countField, count)
		}

		// If count was overridden by current pagination, use that instead
		if pagination.count != 0 {
			query.Set(countField, strconv.Itoa(pagination.count))
		}

		return nil
	}
}
