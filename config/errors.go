package config

type ErrorWithStatus struct {
	Code    int    `json:"error-code"`
	Message string `json:"error-message"`
}

const (
	Base int = 1 + iota
	Binding
	JSONBody
)
