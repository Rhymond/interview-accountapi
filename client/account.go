package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rhymond/interview-accountapi/models"
)

// AccountService holds API functionality for accounts API.
type AccountService struct {
	client *Client
}

// Create registers an existing bank account with Form3 or create a new one. The country attribute must be specified as a minimum. Depending on the country, other attributes such as bank_id and bic are mandatory.
func (s *AccountService) Create(ctx context.Context, account *models.Account) (*models.Account, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "v1/organisation/accounts", account)
	acc := &models.Account{}
	resp, err := s.client.Do(ctx, req, acc)
	if err != nil {
		return nil, resp, err
	}

	return acc, resp, nil
}

// List accounts with the ability to filter and page.
func (s *AccountService) List(ctx context.Context, pagination *Pagination) ([]models.Account, *Response, error) {
	path := fmt.Sprintf("v1/organisation/accounts")
	path, err := addOptions(path, pagination)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	var accounts []models.Account
	resp, err := s.client.Do(ctx, req, &accounts)
	if err != nil {
		return nil, resp, err
	}
	return accounts, resp, nil
}

// Fetch a single account using the account ID.
func (s *AccountService) Fetch(ctx context.Context, id string) (*models.Account, *Response, error) {
	path := fmt.Sprintf("v1/organisation/accounts/%s", id)
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	account := &models.Account{}
	resp, err := s.client.Do(ctx, req, account)
	if err != nil {
		return nil, resp, err
	}

	return account, resp, err
}

// Delete an account using the account ID.
func (s *AccountService) Delete(ctx context.Context, id string) (*Response, error) {
	path := fmt.Sprintf("v1/organisation/accounts/%s?version=0", id)
	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	var accounts []models.Account
	resp, err := s.client.Do(ctx, req, accounts)
	return resp, err
}
