package client

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name               string
		givenHTTPClient    *http.Client
		expectedHTTPClient *http.Client
	}{
		{
			name:               "it should set httpclient to default client if given http client is nil",
			givenHTTPClient:    nil,
			expectedHTTPClient: http.DefaultClient,
		},
		{
			name: "it should set httpclient to a given http client",
			givenHTTPClient: &http.Client{
				Transport:     nil,
				CheckRedirect: nil,
				Jar:           nil,
				Timeout:       20,
			},
			expectedHTTPClient: &http.Client{
				Transport:     nil,
				CheckRedirect: nil,
				Jar:           nil,
				Timeout:       20,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualClient := NewClient(test.givenHTTPClient, nil)
			assert.Equal(t, actualClient.httpClient, test.expectedHTTPClient)
		})
	}
}

func TestAddOptions(t *testing.T) {
	tests := []struct {
		name            string
		givenPath       string
		givenOptions    *Pagination
		expectedPath    string
		expectedIsError bool
	}{
		{
			name:            "it should add pagination",
			givenPath:       "/action",
			givenOptions:    &Pagination{Page: 1},
			expectedPath:    "/action?page[number]=1",
			expectedIsError: false,
		},
		{
			name:            "it should add pagination existing parameters",
			givenPath:       "/action?scope=all",
			givenOptions:    &Pagination{Page: 1},
			expectedPath:    "/action?page[number]=1&scope=all",
			expectedIsError: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := addOptions(test.givenPath, test.givenOptions)
			if test.expectedIsError {
				assert.NotNil(t, err)
			}

			if !test.expectedIsError {
				assert.Nil(t, err)
			}

			gotURL, err := url.Parse(got)
			require.Nil(t, err)
			expectedURL, err := url.Parse(test.expectedPath)
			require.Nil(t, err)

			assert.Equal(t, expectedURL.Path, gotURL.Path)
			assert.Equal(t, expectedURL.Query(), gotURL.Query())
		})
	}
}

func TestDo(t *testing.T) {
	router, server, client := createTestServer()
	defer server.Close()

	type foo struct {
		A string
	}

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":{"A":"a"},"links":{}}`))
	}).Methods(http.MethodGet)

	ctx := context.Background()
	req, _ := client.NewRequest(ctx, http.MethodGet, "/", nil)
	body := foo{}
	_, err := client.Do(context.TODO(), req, &body)
	require.Nil(t, err)
	assert.Equal(t, foo{"a"}, body)
}

func TestNewRequest(t *testing.T) {
	u, _ := url.Parse("http://coolurl.com")
	c := NewClient(nil, u)

	type foo struct {
		A string
	}

	body := foo{A: "B"}
	req, err := c.NewRequest(context.TODO(), http.MethodGet, "/foo", &body)
	require.Nil(t, err)

	expectedURL := "http://coolurl.com/foo"
	assert.Equal(t, expectedURL, req.URL.String())

	reqBody, err := ioutil.ReadAll(req.Body)
	require.Nil(t, err)
	assert.Equal(t, "{\"data\":{\"A\":\"B\"}}\n", string(reqBody))
}
