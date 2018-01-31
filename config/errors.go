package config

// ErrorWithStatus represent the JSON error message delivered to the user in case of failure
type ErrorWithStatus struct {
	Code    int    `json:"error-code"`
	Message string `json:"error-message"`
}

// Error codes
const (
	JSONBodyError int = 1 + iota
	StoreError
	GUIDNotFoundError
	TempFileError
	ReadError
	WriteError
	CloseError
	BlurError
)
