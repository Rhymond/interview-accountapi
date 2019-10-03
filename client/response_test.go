package client

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLinks_CurrentPage(t *testing.T) {
	tests := []struct {
		name                string
		givenLinks          *Links
		expectedCurrentPage int
		expectedError       bool
	}{
		{
			name:                "it should set current page to 1 when links are not set",
			givenLinks:          nil,
			expectedCurrentPage: 1,
		},
		{
			name: "it should set current page to 1 when prev link is not set",
			givenLinks: &Links{
				Self:  "/subscriptions?page[number]=1&page[size]=2",
				Next:  "/subscriptions?page[number]=2&page[size]=2",
				First: "/subscriptions?page[number]=1&page[size]=2",
				Last:  "/subscriptions?page[number]=13&page[size]=2",
			},
			expectedCurrentPage: 1,
		},
		{
			name: "it should set current page to 3 when prev page is set",
			givenLinks: &Links{
				Self:  "/subscriptions?page[number]=3&page[size]=2",
				Prev:  "/subscriptions?page[number]=2&page[size]=2",
				Next:  "/subscriptions?page[number]=4&page[size]=2",
				First: "/subscriptions?page[number]=1&page[size]=2",
				Last:  "/subscriptions?page[number]=13&page[size]=2",
			},
			expectedCurrentPage: 3,
		},
		{
			name: "it should throw an error when prev url is malformed",
			givenLinks: &Links{
				Self:  "/subscriptions?page[number]=3&page[size]=2",
				Prev:  "some-malformed-url",
				Next:  "/subscriptions?page[number]=4&page[size]=2",
				First: "/subscriptions?page[number]=1&page[size]=2",
				Last:  "/subscriptions?page[number]=13&page[size]=2",
			},
			expectedError:       true,
			expectedCurrentPage: 0,
		},
		{
			name: "it should throw an error when page number is not set in prev link",
			givenLinks: &Links{
				Self:  "/subscriptions?page[number]=3&page[size]=2",
				Prev:  "/subscriptions?page[size]=2",
				Next:  "/subscriptions?page[number]=4&page[size]=2",
				First: "/subscriptions?page[number]=1&page[size]=2",
				Last:  "/subscriptions?page[number]=13&page[size]=2",
			},
			expectedError:       true,
			expectedCurrentPage: 0,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			curr, err := test.givenLinks.CurrentPage()
			if test.expectedError {
				require.NotNil(t, err)
			}
			if !test.expectedError {
				require.Nil(t, err)
			}
			assert.Equal(t, test.expectedCurrentPage, curr)
		})
	}
}

func TestLinks_IsLastPage(t *testing.T) {
	tests := []struct {
		name           string
		givenLinks     *Links
		expectedIsLast bool
	}{
		{
			name:           "it should set last page to true when no links are given",
			givenLinks:     nil,
			expectedIsLast: true,
		},
		{
			name: "it should set last page to false when last page is given",
			givenLinks: &Links{
				Self:  "/subscriptions?page[number]=1&page[size]=2",
				Next:  "/subscriptions?page[number]=2&page[size]=2",
				First: "/subscriptions?page[number]=1&page[size]=2",
				Last:  "/subscriptions?page[number]=13&page[size]=2",
			},
			expectedIsLast: false,
		},
		{
			name: "it should set last page to true when last page is not set",
			givenLinks: &Links{
				Self:  "/subscriptions?page[number]=1&page[size]=2",
				Next:  "/subscriptions?page[number]=2&page[size]=2",
				First: "/subscriptions?page[number]=1&page[size]=2",
			},
			expectedIsLast: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expectedIsLast, test.givenLinks.IsLastPage())
		})
	}
}

func TestErrorResponse_Error(t *testing.T) {
	tests := []struct {
		name            string
		givenError      *ErrorResponse
		expectedMessage string
	}{
		{
			name: "it should return response error message with status and message included",
			givenError: &ErrorResponse{
				Response:   nil,
				StatusCode: http.StatusInternalServerError,
				Message:    "error message",
			},
			expectedMessage: "code: 500, message: error message",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.EqualError(t, test.givenError, test.expectedMessage)
			actualError := test.givenError.Error()
			if test.expectedMessage != actualError {
				t.Errorf("expected error to be `%s`, but got `%s`", test.expectedMessage, actualError)
			}
		})
	}
}
