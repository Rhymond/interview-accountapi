package models

type Account struct {
	// attributes
	// Required: true
	Attributes AccountAttributes `json:"attributes"`

	// Unique resource ID
	// Required: true
	// Format: uuid
	ID string `json:"id"`

	// Unique ID of the organisation this resource is created by
	// Required: true
	// Format: uuid
	OrganisationID string `json:"organisation_id"`

	// Name of the resource type
	// Pattern: ^[A-Za-z]*$
	Type string `json:"type,omitempty"`

	// Version number
	// Minimum: 0
	Version int64 `json:"version,omitempty"`
}
