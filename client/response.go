package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// Links holds information about response pagination.
type Links struct {
	First string `json:"first"`
	Last  string `json:"last"`
	Prev  string `json:"prev"`
	Next  string `json:"next"`
	Self  string `json:"self"`
}

// CurrentPage returns current page from given response.
// if prev link URL is invalid it throws an error.
func (l *Links) CurrentPage() (int, error) {
	switch {
	case l == nil:
		return 1, nil
	case l.Prev == "" && l.Next != "":
		return 1, nil
	case l.Prev != "":
		prevPage, err := pageForURL(l.Prev)
		if err != nil {
			return 0, err
		}

		return prevPage + 1, nil
	}

	return 0, nil
}

// IsLastPage check if page is last.
func (l *Links) IsLastPage() bool {
	if l == nil {
		return true
	}
	return l.Last == ""
}

// Response is API HTTP response.
type Response struct {
	Response *http.Response
	Data     json.RawMessage `json:"data"`
	Links    Links           `json:"links"`
}

// NewResponse creates a new Response for the provided http.Response
func NewResponse(r *http.Response) *Response {
	response := Response{Response: r}

	return &response
}

// ErrorResponse is a custom error structure for API errors.
// It hold HTTP response that caused error and has all given details about an error.
type ErrorResponse struct {
	Response   *http.Response
	StatusCode int
	Code       string `json:"error_code"`
	Message    string `json:"error_message"`
}

// Error is required to be implemented to meet error interface
func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.StatusCode, e.Message)
}

func pageForURL(urlText string) (int, error) {
	u, err := url.ParseRequestURI(urlText)
	if err != nil {
		return 0, err
	}

	pageStr := u.Query().Get("page[number]")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return 0, err
	}

	return page, nil
}
