package entity

// Error represents an API error response.
type Error struct {
	// Code is the error code
	Code string `json:"code"`

	// Message is the error's short description
	Message string `json:"message"`

	// Detail is the error's detailed description
	Detail string `json:"detail,omitempty"`
}

// Common error codes
const (
	ErrCodeDeviceNotFound     = "device_not_found"
	ErrCodeDeviceExists       = "device_exists"
	ErrCodePeerNotFound       = "peer_not_found"
	ErrCodeInvalidRequest     = "invalid_request"
	ErrCodeInternalError      = "internal_error"
	ErrCodeUnauthorized       = "unauthorized"
)
