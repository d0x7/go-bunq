package pagination

import (
	"github.com/d0x7/go-bunq/model"
	"net/url"
	"strconv"
)

// Count returns a query parameter, that will limit the amount of results returned to the given count.
func Count(count int) model.QueryParam {
	return func(query url.Values) error {
		query.Set("count", strconv.Itoa(count))
		return nil
	}
}

// NewerThan returns a query parameter, that will return only results newer than the given id.
func NewerThan(id int) model.QueryParam {
	return func(query url.Values) error {
		query.Set("newer_id", strconv.Itoa(id))
		return nil
	}
}

// OlderThan returns a query parameter, that will return only results older than the given id.
func OlderThan(id int) model.QueryParam {
	return func(query url.Values) error {
		query.Set("older_id", strconv.Itoa(id))
		return nil
	}
}
