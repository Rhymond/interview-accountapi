package client

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/rhymond/interview-accountapi/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountService_CreateResponseError(t *testing.T) {
	tests := []struct {
		name            string
		givenResponse   string
		givenStatusCode int
		expectedError   string
	}{
		{
			name:            "it should return an error on malformed response",
			givenResponse:   `not-a-json`,
			givenStatusCode: http.StatusOK,
			expectedError:   "invalid character 'o' in literal null (expecting 'u')",
		},
		{
			name:            "it should return an error on empty response",
			givenResponse:   `{}`,
			givenStatusCode: http.StatusOK,
			expectedError:   "unexpected end of JSON input",
		},
		{
			name: "it should return custom api error on internal server error status",
			givenResponse: `{
				"error_message": "custom error message"
			}`,
			givenStatusCode: http.StatusInternalServerError,
			expectedError:   "code: 500, message: custom error message",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router, server, client := createTestServer()
			defer server.Close()

			router.HandleFunc("/v1/organisation/accounts", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(test.givenStatusCode)
				fmt.Fprintf(w, test.givenResponse)
			}).Methods(http.MethodPost)

			_, _, err := client.Account.Create(context.TODO(), &models.Account{})
			require.EqualError(t, err, test.expectedError)
		})
	}
}

func TestAccountService_CreateResponseSuccess(t *testing.T) {
	tests := []struct {
		name              string
		givenResponse     string
		givenStatusCode   int
		expectedAccountID string
	}{
		{
			name: "it should return given account id on valid response",
			givenResponse: `{
				"data": {
					"attributes": {},
					"created_on": "2019-10-02T13:34:32.324Z",
					"id": "b8952241-a065-462e-a7d2-6a9c94010f0f",
					"modified_on": "2019-10-02T13:34:32.324Z",
					"organisation_id": "efab8098-d2e7-47f0-9db3-1c318920f71d",
					"type": "accounts",
					"version": 0
				},
				"links": {
					"self": "/v1/organisation/accounts/b8952241-a065-462e-a7d2-6a9c94010f0f"
				}
			}`,
			expectedAccountID: "b8952241-a065-462e-a7d2-6a9c94010f0f",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router, server, client := createTestServer()
			defer server.Close()

			router.HandleFunc("/v1/organisation/accounts", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, test.givenResponse)
			}).Methods(http.MethodPost)

			acc, _, err := client.Account.Create(context.TODO(), &models.Account{})
			if assert.Nil(t, err) {
				assert.Equal(t, test.expectedAccountID, acc.ID)
			}
		})
	}
}

func TestAccountService_CreateRequest(t *testing.T) {
	tests := []struct {
		name         string
		givenAccount *models.Account
		expectedBody string
	}{
		{
			name: "it should include given account to the request body",
			givenAccount: &models.Account{
				Attributes: models.AccountAttributes{
					AccountNumber: "1234",
					Country:       "GB",
				},
				ID:             "account-id",
				OrganisationID: "organisation-id",
				Type:           "account-type",
			},
			expectedBody: `{"data":{"attributes":{"country":"GB","account_number":"1234"},"id":"account-id","organisation_id":"organisation-id","type":"account-type","version":0}}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router, server, client := createTestServer()
			defer server.Close()
			var isCalled bool
			router.Path("/v1/organisation/accounts").
				Methods(http.MethodPost).
				HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					isCalled = true
					data, err := ioutil.ReadAll(r.Body)
					if assert.Nil(t, err) {
						data = data[:len(data)-1] // remove /n
						assert.Equal(t, test.expectedBody, string(data))
					}
				})

			_, _, _ = client.Account.Create(context.TODO(), test.givenAccount)
			assert.True(t, isCalled)
		})
	}
}

func TestAccountService_FetchResponseError(t *testing.T) {
	tests := []struct {
		name            string
		givenResponse   string
		givenStatusCode int
		expectedError   string
	}{
		{
			name:            "it should return an error on malformed response",
			givenResponse:   `not-a-json`,
			givenStatusCode: http.StatusOK,
			expectedError:   "invalid character 'o' in literal null (expecting 'u')",
		},
		{
			name:            "it should return an error on empty response",
			givenResponse:   `{}`,
			givenStatusCode: http.StatusOK,
			expectedError:   "unexpected end of JSON input",
		},
		{
			name: "it should return custom api error on internal server error status",
			givenResponse: `{
				"error_message": "custom error message"
			}`,
			givenStatusCode: http.StatusInternalServerError,
			expectedError:   "code: 500, message: custom error message",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router, server, client := createTestServer()
			defer server.Close()

			router.
				Path("/v1/organisation/accounts/{id}").
				Methods(http.MethodGet).
				HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(test.givenStatusCode)
					fmt.Fprintf(w, test.givenResponse)
				})

			ctx := context.Background()
			_, _, err := client.Account.Fetch(ctx, "account-id")
			require.EqualError(t, err, test.expectedError)
		})
	}
}

func TestAccountService_FetchResponseSuccess(t *testing.T) {
	tests := []struct {
		name              string
		givenResponse     string
		expectedAccountID string
	}{
		{
			name: "it should return account id on valid response",
			givenResponse: `{
				"data": {
					"attributes": {},
					"created_on": "2019-10-02T13:34:32.324Z",
					"id": "b8952241-a065-462e-a7d2-6a9c94010f0f",
					"modified_on": "2019-10-02T13:34:32.324Z",
					"organisation_id": "efab8098-d2e7-47f0-9db3-1c318920f71d",
					"type": "accounts",
					"version": 0
				},
				"links": {
					"self": "/v1/organisation/accounts/b8952241-a065-462e-a7d2-6a9c94010f0f"
				}
			}`,
			expectedAccountID: "b8952241-a065-462e-a7d2-6a9c94010f0f",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router, server, client := createTestServer()
			defer server.Close()

			router.
				Path("/v1/organisation/accounts/{id}").
				Methods(http.MethodGet).
				HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					fmt.Fprintf(w, test.givenResponse)
				})

			ctx := context.Background()
			acc, _, err := client.Account.Fetch(ctx, "account-id")
			if assert.Nil(t, err) {
				assert.Equal(t, test.expectedAccountID, acc.ID)
			}
		})
	}
}

func TestAccountService_FetchRequest(t *testing.T) {
	tests := []struct {
		name           string
		givenAccountID string
		expectedBody   string
	}{
		{
			name:           "it should include given account to the request body",
			givenAccountID: "account-id",
			expectedBody:   ``,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router, server, client := createTestServer()
			defer server.Close()
			var isCalled bool

			router.Path("/v1/organisation/accounts/{id}").
				Methods(http.MethodGet).
				HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					isCalled = true
					data, err := ioutil.ReadAll(r.Body)
					if assert.Nil(t, err) {
						assert.Equal(t, test.expectedBody, string(data))
						vars := mux.Vars(r)
						assert.Equal(t, test.givenAccountID, vars["id"])
					}
				})

			_, _, _ = client.Account.Fetch(context.TODO(), test.givenAccountID)
			assert.True(t, isCalled)
		})
	}
}

func TestAccountService_ListResponseError(t *testing.T) {
	tests := []struct {
		name            string
		givenResponse   string
		givenStatusCode int
		expectedError   string
	}{
		{
			name:            "it should return error on malformed response",
			givenResponse:   `not-a-json`,
			givenStatusCode: http.StatusOK,
			expectedError:   "invalid character 'o' in literal null (expecting 'u')",
		},
		{
			name:            "it should return error on empty response",
			givenResponse:   `{}`,
			givenStatusCode: http.StatusOK,
			expectedError:   "unexpected end of JSON input",
		},
		{
			name: "it should return custom api error on internal server error status",
			givenResponse: `{
				"error_message": "custom error message"
			}`,
			givenStatusCode: http.StatusInternalServerError,
			expectedError:   "code: 500, message: custom error message",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router, server, client := createTestServer()
			defer server.Close()

			router.HandleFunc("/v1/organisation/accounts", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(test.givenStatusCode)
				fmt.Fprintf(w, test.givenResponse)
			}).Methods(http.MethodGet)

			_, _, err := client.Account.List(context.TODO(), nil)
			assert.EqualError(t, err, test.expectedError)
		})
	}
}

func TestAccountService_ListResponseSuccess(t *testing.T) {
	tests := []struct {
		name                 string
		givenResponse        string
		expectedAccountCount int
		expectedAccountIDs   []string
		expectedError        string
	}{
		{
			name: "it should return empty list of accounts on valid response with null data",
			givenResponse: `{
				"data": null
			}`,
			expectedAccountCount: 0,
			expectedAccountIDs:   []string{},
		},
		{
			name: "it should return list of accounts on given valid response",
			givenResponse: `{
				"data": [
					{
						"attributes": {
							"account_classification": "Personal",
							"account_matching_opt_out": false,
							"account_number": "10000004",
							"alternative_bank_account_names": [],
							"bank_id": "400302",
							"bank_id_code": "GBDSC",
							"base_currency": "GBP",
							"bic": "NWBKGB42",
							"country": "GB",
							"customer_id": "234",
							"first_name": "Mary-Jane Doe",
							"iban": "GB28NWBK40030212764204",
							"joint_account": false,
							"title": "Sir"
						},
						"created_on": "2019-10-03T18:09:42.257Z",
						"id": "bdf9e1a8-481e-483f-b54c-7103cffceb21",
						"modified_on": "2019-10-03T18:09:42.257Z",
						"organisation_id": "7c3d20ff-ed78-45c4-aae0-0184cf6d3060",
						"type": "accounts",
						"version": 0
					},
					{
						"attributes": {
							"account_classification": "Personal",
							"account_matching_opt_out": false,
							"account_number": "10000004",
							"alternative_bank_account_names": [],
							"bank_id": "400302",
							"bank_id_code": "GBDSC",
							"base_currency": "GBP",
							"bic": "NWBKGB42",
							"country": "GB",
							"customer_id": "234",
							"first_name": "Mary-Jane Doe",
							"iban": "GB28NWBK40030212764204",
							"joint_account": false,
							"title": "Sir"
						},
						"created_on": "2019-10-03T18:09:43.501Z",
						"id": "63c0a226-5b6c-4ef9-a0bb-436dd39d45bb",
						"modified_on": "2019-10-03T18:09:43.501Z",
						"organisation_id": "8d0f39dd-6c17-4720-915e-7a015dda9f47",
						"type": "accounts",
						"version": 0
					}
				],
				"links": {
					"self": "/v1/organisation/accounts/b8952241-a065-462e-a7d2-6a9c94010f0f"
				}
			}`,
			expectedAccountCount: 2,
			expectedAccountIDs:   []string{"bdf9e1a8-481e-483f-b54c-7103cffceb21", "63c0a226-5b6c-4ef9-a0bb-436dd39d45bb"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router, server, client := createTestServer()
			defer server.Close()

			router.HandleFunc("/v1/organisation/accounts", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, test.givenResponse)
			}).Methods(http.MethodGet)

			accs, _, err := client.Account.List(context.TODO(), nil)
			if assert.Nil(t, err) {
				assert.Equal(t, test.expectedAccountCount, len(accs))

				accsids := make([]string, len(accs))
				for i, acc := range accs {
					accsids[i] = acc.ID
				}

				assert.Equal(t, test.expectedAccountIDs, accsids)
			}
		})
	}
}

func TestAccountService_DeleteResponse(t *testing.T) {
	tests := []struct {
		name            string
		givenResponse   string
		givenStatusCode int
		expectedError   string
	}{
		{
			name:            "it should return an error on malformed response",
			givenResponse:   `not-a-json`,
			givenStatusCode: http.StatusOK,
			expectedError:   "invalid character 'o' in literal null (expecting 'u')",
		},
		{
			name:            "it should return an error on empty response",
			givenResponse:   `{}`,
			givenStatusCode: http.StatusOK,
			expectedError:   "unexpected end of JSON input",
		},
		{
			name: "it should return custom api error on internal server error status",
			givenResponse: `{
				"error_message": "custom error message"
			}`,
			givenStatusCode: http.StatusInternalServerError,
			expectedError:   "code: 500, message: custom error message",
		},
		{
			name:            "it should return no error with no content",
			givenStatusCode: http.StatusNoContent,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router, server, client := createTestServer()
			defer server.Close()

			router.HandleFunc("/v1/organisation/accounts/{id}", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(test.givenStatusCode)
				fmt.Fprintf(w, test.givenResponse)
			}).Methods(http.MethodDelete)

			ctx := context.Background()
			_, err := client.Account.Delete(ctx, "account-id")
			if test.expectedError != "" {
				assert.EqualError(t, err, test.expectedError)
			}
			if test.expectedError == "" {
				assert.Nil(t, err)
			}
		})
	}
}

func TestAccountService_DeleteRequest(t *testing.T) {
	tests := []struct {
		name           string
		givenAccountID string
		expectedBody   string
	}{
		{
			name:           "it should include given account id into the request context",
			givenAccountID: "account-id",
			expectedBody:   ``,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router, server, client := createTestServer()
			defer server.Close()
			var isCalled bool

			router.Path("/v1/organisation/accounts/{id}").
				Methods(http.MethodDelete).
				HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					isCalled = true
					data, err := ioutil.ReadAll(r.Body)
					if assert.Nil(t, err) {
						assert.Equal(t, test.expectedBody, string(data))
						vars := mux.Vars(r)
						assert.Equal(t, test.givenAccountID, vars["id"])
					}
				})

			_, _ = client.Account.Delete(context.TODO(), test.givenAccountID)
			assert.True(t, isCalled)
		})
	}
}
