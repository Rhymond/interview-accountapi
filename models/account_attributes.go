package models

// AccountAttributes account attributes
type AccountAttributes struct {
	// Is the account business or personal?
	// Enum: [Personal Business]
	AccountClassification string `json:"account_classification,omitempty"`

	// Is the account opted out of account matching, e.g. CoP?
	AccountMatchingOptOut bool `json:"account_matching_opt_out,omitempty"`

	// Account number of the account. A unique number will automatically be generated if not provided.
	// Pattern: ^[A-Z0-9]{0,64}$
	AccountNumber string `json:"account_number,omitempty"`

	// Alternative account names. Used for Confirmation of Payee matching.
	// Max Items: 3
	AlternativeBankAccountNames []string `json:"alternative_bank_account_names,omitempty"`

	// Primary account name. Used for Confirmation of Payee matching. Required if confirmation_of_payee_enabled is true for the organisation.
	// Max Length: 140
	// Min Length: 1
	BankAccountName string `json:"bank_account_name,omitempty"`

	// Local country bank identifier. In the UK this is the sort code.
	// Pattern: ^[A-Z0-9]{0,16}$
	BankID string `json:"bank_id,omitempty"`

	// ISO 20022 code used to identify the type of bank ID being used
	// Pattern: ^[A-Z]{0,16}$
	BankIDCode string `json:"bank_id_code,omitempty"`

	// ISO 4217 code used to identify the base currency of the account
	// Pattern: ^[A-Z]{3}$
	BaseCurrency string `json:"base_currency,omitempty"`

	// SWIFT BIC in either 8 or 11 character format
	// Pattern: ^([A-Z]{6}[A-Z0-9]{2}|[A-Z]{6}[A-Z0-9]{5})$
	Bic string `json:"bic,omitempty"`

	// ISO 3166-1 code used to identify the domicile of the account
	// Required: true
	// Pattern: ^[A-Z]{2}$
	Country string `json:"country"`

	// A free-format reference that can be used to link this account to an external system
	// Pattern: ^[a-zA-Z0-9-$@., ]{0,256}$
	CustomerID string `json:"customer_id,omitempty"`

	// Customer first name.
	// Max Length: 40
	// Min Length: 1
	FirstName string `json:"first_name,omitempty"`

	// IBAN of the account. Will be calculated from other fields if not supplied.
	// Pattern: ^[A-Z]{2}[0-9]{2}[A-Z0-9]{0,64}$
	Iban string `json:"iban,omitempty"`

	// Is the account joint?
	JointAccount bool `json:"joint_account,omitempty"`

	// Secondary identification, e.g. building society roll number. Used for Confirmation of Payee.
	// Max Length: 140
	// Min Length: 1
	SecondaryIdentification string `json:"secondary_identification,omitempty"`

	// Customer title.
	// Max Length: 40
	// Min Length: 1
	Title string `json:"title,omitempty"`
}
