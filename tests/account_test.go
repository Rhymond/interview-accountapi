// +build integration

package tests

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/google/uuid"
	"github.com/rhymond/interview-accountapi/client"
	"github.com/rhymond/interview-accountapi/models"
)

type apiFeature struct {
	client         *client.Client
	resp           *client.Response
	listedAccounts []models.Account
	fetchedAccount *models.Account
	createdAccount *models.Account
}

func (a *apiFeature) reset(interface{}) {
	a.resp = nil
	a.listedAccounts = nil
	a.fetchedAccount = nil
	a.createdAccount = nil

	accs, _, err := a.client.Account.List(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
	for _, acc := range accs {
		_, err := a.client.Account.Delete(context.TODO(), acc.ID)
		if err != nil {
			panic(err)
		}
	}
}

func (a *apiFeature) theResponseCodeShouldBe(code int) error {
	if code != a.resp.Response.StatusCode {
		return fmt.Errorf("expected response code to be: %d, but actual is: %d", code, a.resp.Response.StatusCode)
	}
	return nil
}

func (a *apiFeature) theResponseShouldMatchJSON(body *gherkin.DocString) (err error) {
	var expected, actual interface{}

	if err = json.Unmarshal([]byte(body.Content), &expected); err != nil {
		return
	}

	if err = json.Unmarshal(a.resp.Data, &actual); err != nil {
		return
	}

	if !reflect.DeepEqual(expected, actual) {
		return fmt.Errorf("expected JSON does not match actual, %v vs. %v", expected, actual)
	}
	return nil
}

func (a *apiFeature) iListAccounts() (err error) {
	accs, resp, err := a.client.Account.List(context.TODO(), nil)
	if err != nil {
		return
	}
	a.resp = resp
	a.listedAccounts = accs
	return
}

func (a *apiFeature) iListAccountsPerPage(perPage int) (err error) {
	accs, resp, err := a.client.Account.List(context.TODO(), &client.Pagination{
		PerPage: perPage,
	})
	if err != nil {
		return
	}
	a.resp = resp
	a.listedAccounts = accs
	return
}

func (a *apiFeature) iListAccountsPerPageInPage(perPage, page int) (err error) {
	accs, resp, err := a.client.Account.List(context.TODO(), &client.Pagination{
		PerPage: perPage,
		Page:    page - 1,
	})
	if err != nil {
		return
	}
	a.resp = resp
	a.listedAccounts = accs
	return
}

func (a *apiFeature) iDeleteAccount(id string) (err error) {
	resp, err := a.client.Account.Delete(context.TODO(), id)
	if err != nil {
		return
	}
	a.resp = resp
	return
}

func (a *apiFeature) iCreateAccount(id string) (err error) {
	acc, resp, err := a.client.Account.Create(context.TODO(), &models.Account{
		Attributes: models.AccountAttributes{
			Country: "GB",
		},
		ID:             id,
		OrganisationID: uuid.Must(uuid.NewUUID()).String(),
		Type:           "accounts",
	})
	if err != nil {
		return
	}
	a.createdAccount = acc
	a.resp = resp
	return
}

func (a *apiFeature) iCreateAccounts(count int) (err error) {
	for i := 1; i <= count; i++ {
		_, _, err = a.client.Account.Create(context.TODO(), &models.Account{
			Attributes: models.AccountAttributes{
				Country: "GB",
			},
			ID:             uuid.Must(uuid.NewUUID()).String(),
			OrganisationID: uuid.Must(uuid.NewUUID()).String(),
			Type:           "accounts",
		})
	}
	return
}

func (a *apiFeature) theCountOfAccounts(count int) (err error) {
	if len(a.listedAccounts) != count {
		return fmt.Errorf("expected count of account to be %d, but got %d", count, len(a.listedAccounts))
	}
	return
}

func (a *apiFeature) iFetchAccount(id string) (err error) {
	acc, resp, err := a.client.Account.Fetch(context.TODO(), id)
	if err != nil {
		if eerr, ok := err.(*client.ErrorResponse); ok {
			a.resp = &client.Response{
				Response: eerr.Response,
			}
			return nil
		}

		return err
	}
	a.fetchedAccount = acc
	a.resp = resp
	return
}

func (a *apiFeature) theFetchedAccountIdShouldBe(id string) (err error) {
	if a.fetchedAccount.ID != id {
		return fmt.Errorf("expected fetched account id to be %q, but got %q", id, a.fetchedAccount.ID)
	}
	return
}

func (a *apiFeature) theCreatedAccountIdShouldBe(id string) (err error) {
	if a.createdAccount.ID != id {
		return fmt.Errorf("expected fetched account id to be %q, but got %q", id, a.fetchedAccount.ID)
	}
	return
}

func (a *apiFeature) theCurrentPageIs(currentPage int) (err error) {
	actualCurrPage, err := a.resp.Links.CurrentPage()
	if err != nil {
		return
	}
	if actualCurrPage != currentPage {
		return fmt.Errorf("expected current page to be %d, but got %d", currentPage, actualCurrPage)
	}
	return
}

func (a *apiFeature) thePageIsNotLast() (err error) {
	if a.resp.Links.IsLastPage() {
		return errors.New("expected current page not to be last")
	}
	return
}

func (a *apiFeature) thePageIsLast() (err error) {
	if !a.resp.Links.IsLastPage() {
		return errors.New("expected current page to be last")
	}
	return
}

func FeatureContext(s *godog.Suite) {
	api := &apiFeature{}
	accountApiAddr := os.Getenv("ACCOUNT_API_ADDR")
	if accountApiAddr == "" {
		accountApiAddr = "http://localhost:8080"
	}
	u, err := url.Parse(accountApiAddr)
	if err != nil {
		panic(err)
	}

	api.client = client.NewClient(nil, u)

	s.BeforeScenario(api.reset)
	s.Step(`^I create (\d+) accounts$`, api.iCreateAccounts)
	s.Step(`^I list accounts$`, api.iListAccounts)
	s.Step(`^I list (\d+) accounts per page$`, api.iListAccountsPerPage)
	s.Step(`^I list (\d+) accounts per page in page (\d+)$`, api.iListAccountsPerPageInPage)
	s.Step(`^I create account with id "([^"]*)"$`, api.iCreateAccount)
	s.Step(`^I fetch account with id "([^"]*)"$`, api.iFetchAccount)
	s.Step(`^I delete account with id "([^"]*)"$`, api.iDeleteAccount)
	s.Step(`^the response code should be (\d+)$`, api.theResponseCodeShouldBe)
	s.Step(`^the count of accounts should be (\d+)$`, api.theCountOfAccounts)
	s.Step(`^the response should match json:$`, api.theResponseShouldMatchJSON)
	s.Step(`^the fetched account id should be "([^"]*)"$`, api.theFetchedAccountIdShouldBe)
	s.Step(`^the created account id should be "([^"]*)"$`, api.theCreatedAccountIdShouldBe)
	s.Step(`^the current page should be (\d+)$`, api.theCurrentPageIs)
	s.Step(`^the page should be the last$`, api.thePageIsLast)
	s.Step(`^the page should not be the last$`, api.thePageIsNotLast)
}

func TestMain(m *testing.M) {
	status := godog.RunWithOptions("godog", func(s *godog.Suite) {
		FeatureContext(s)
	}, godog.Options{
		Format:    "pretty",
		Paths:     []string{"features"},
		Randomize: time.Now().UTC().UnixNano(), // randomize scenario execution order
	})

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}
