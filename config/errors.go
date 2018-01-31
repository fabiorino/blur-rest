package config

// ErrorWithStatus represent the JSON error message delivered to the user in case of failure
type ErrorWithStatus struct {
	Code    int    `json:"error-code"`
	Message string `json:"error-message"`
}

// Error codes
const (
	BindingError int = 1 + iota
	JSONBodyError
	StoreError
	GUIDNotFoundError
	TempFileError
	ReadError
	WriteError
	CloseError
	BlurError
)
