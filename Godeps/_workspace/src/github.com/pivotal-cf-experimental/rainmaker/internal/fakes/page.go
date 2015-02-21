package fakes

import (
	"encoding/json"
	"fmt"
)

type Pageable interface {
	Len() int
	Items() []interface{}
}

type Page struct {
	collection Pageable
	path       string
	number     int
	perPage    int
}

func NewPage(collection Pageable, path string, number, perPage int) Page {
	return Page{
		collection: collection,
		path:       path,
		number:     number,
		perPage:    perPage,
	}
}

func (page Page) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"total_results": page.collection.Len(),
		"total_pages":   page.PageCount(),
		"prev_url":      page.PrevURL(),
		"next_url":      page.NextURL(),
		"resources":     page.Resources(),
	})
}

func (page Page) Resources() []interface{} {
	var resources []interface{}
	start := (page.number - 1) * page.perPage
	end := (page.number * page.perPage)

	if end > page.collection.Len() {
		end = page.collection.Len()
	}

	for _, item := range page.collection.Items()[start:end] {
		resources = append(resources, item)
	}

	return resources
}

func (page Page) PageCount() int {
	length := page.collection.Len()
	pageCount := length / page.perPage

	if (length % page.perPage) != 0 {
		pageCount++
	}

	return pageCount
}

func (page Page) PrevURL() string {
	if page.number > 1 {
		return fmt.Sprintf("%s?page=%d&results-per-page=%d", page.path, page.number-1, page.perPage)
	}

	return ""
}

func (page Page) NextURL() string {
	if page.PageCount() > 1 && page.number < page.PageCount() {
		return fmt.Sprintf("%s?page=%d&results-per-page=%d", page.path, page.number+1, page.perPage)
	}

	return ""
}
