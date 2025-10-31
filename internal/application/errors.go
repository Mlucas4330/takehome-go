package application

type AppError struct {
	Status  int
	Code    string
	Message string
	Details any
}

func (e *AppError) Error() string { return e.Message }

func NewError(status int, code, msg string, details any) *AppError {
	return &AppError{Status: status, Code: code, Message: msg, Details: details}
}
