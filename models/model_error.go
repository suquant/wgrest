package models

type Error struct {

	// Error code
	Code string `json:"code"`

	// Error's short description
	Message string `json:"message"`

	// Error's detail description
	Detail string `json:"detail,omitempty"`
}
