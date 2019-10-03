package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type Links struct {
	First string `json:"first"`
	Last  string `json:"last"`
	Prev  string `json:"prev"`
	Next  string `json:"next"`
	Self  string `json:"self"`
}

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

func (l *Links) IsLastPage() bool {
	if l == nil {
		return true
	}
	return l.Last == ""
}

type Response struct {
	// Original response
	Response *http.Response
	Data     json.RawMessage `json:"data"`
	Links    Links           `json:"links"`
}

// NewResponse creates a new Response for the provided http.Response
func NewResponse(r *http.Response) *Response {
	response := Response{Response: r}

	return &response
}

type ErrorResponse struct {
	// HTTP response that caused this error
	Response   *http.Response
	StatusCode int
	Code       string `json:"error_code"`
	Message    string `json:"error_message"`
}

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
