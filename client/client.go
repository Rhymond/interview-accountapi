package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"

	"github.com/google/go-querystring/query"
)

var contentType = "application/vnd.api+json"

// Client is API client.
type Client struct {
	BaseURL    *url.URL
	httpClient *http.Client

	Account *AccountService
}

// Pagination is a structure required to build query parameters for pagination.
type Pagination struct {
	Page    int `url:"page[number],omitempty"`
	PerPage int `url:"page[size],omitempty"`
}

// NewClient creates new API client instance.
func NewClient(httpClient *http.Client, baseURL *url.URL) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	c := &Client{
		BaseURL:    baseURL,
		httpClient: httpClient,
	}
	c.Account = &AccountService{client: c}
	return c
}

// addOptions adds query parameters to given path.
func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	origURL, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	origValues := origURL.Query()
	newValues, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	for k, v := range newValues {
		origValues[k] = v
	}
	origURL.RawQuery = origValues.Encode()
	return origURL.String(), nil
}

// NewRequest creates an API request. A relative URL can be provided in urlStr, which will be resolved to the
// BaseURL of the Client. If body parameter is specified, the value pointed to by
// body is JSON encoded and included in as the request body.
func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if body != nil {
		requestBody := struct {
			Data interface{} `json:"data"`
		}{
			Data: body,
		}

		err = json.NewEncoder(buf).Encode(requestBody)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Accept", contentType)
	return req, nil
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value
// pointed to by v, or returned as an error if an API error has occurred.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	req = req.WithContext(ctx)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if respErr := resp.Body.Close(); err == nil {
			err = respErr
		}
	}()

	r := NewResponse(resp)
	err = checkResponse(resp)
	if err != nil {
		return r, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil || len(data) == 0 {
		return nil, err
	}
	err = json.Unmarshal(data, r)
	if err != nil {
		return nil, err
	}

	if v != nil {
		err = json.Unmarshal(r.Data, v)
		if err != nil {
			return nil, err
		}
	}

	return r, err
}

// CheckResponse checks the API response for errors, and returns them if present. A response is considered an
// error if it has a status code outside the 200 range. API error responses are expected to have JSON response body that maps to ErrorResponse. Any other response body will be silently ignored.
func checkResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r, StatusCode: r.StatusCode}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, errorResponse)
		if err != nil {
			errorResponse.Message = string(data)
		}
	}

	return errorResponse
}
