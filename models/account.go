package models

// Account represents a bank account that is registered with Form3.
// It is used to validate and allocate inbound payments.
type Account struct {
	Attributes     AccountAttributes `json:"attributes"`
	ID             string            `json:"id"`
	OrganisationID string            `json:"organisation_id"`
	Type           string            `json:"type"`
	Version        int               `json:"version"`
}

// AccountAttributes represents account attributes.
type AccountAttributes struct {
	Country                     string   `json:"country"`
	BaseCurrency                string   `json:"base_currency,omitempty"`
	AccountNumber               string   `json:"account_number,omitempty"`
	BankID                      string   `json:"bank_id,omitempty"`
	BankIDCode                  string   `json:"bank_id_code,omitempty"`
	Bic                         string   `json:"bic,omitempty"`
	Iban                        string   `json:"iban,omitempty"`
	Title                       string   `json:"title,omitempty"`
	FirstName                   string   `json:"first_name,omitempty"`
	BankAccountName             string   `json:"bank_account_name,omitempty"`
	AlternativeBankAccountNames []string `json:"alternative_bank_account_names,omitempty"`
	AccountClassification       string   `json:"account_classification,omitempty"`
	JointAccount                bool     `json:"joint_account,omitempty"`
	AccountMatchingOptOut       bool     `json:"account_matching_opt_out,omitempty"`
	SecondaryIdentification     string   `json:"secondary_identification,omitempty"`
}
