package config

type ErrorWithStatus struct {
	Code    int    `json:"error-code"`
	Message string `json:"error-message"`
}

const (
	BaseError int = 1 + iota
	BindingError
	JSONBodyError
	StoreError
	GUIDNotFoundError
	TempFileError
	ReadError
	WriteError
	CloseError
	BlurError
)
