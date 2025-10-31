package apperror

import "net/http"

type Error struct {
	Status  int
	Code    string
	Message string
	Details any
}

func (e *Error) Error() string { return e.Message }

func New(status int, code, message string, details any) *Error {
	return &Error{Status: status, Code: code, Message: message, Details: details}
}

func GetStatus(err error) int {
	if e, ok := err.(*Error); ok {
		return e.Status
	}
	return http.StatusInternalServerError
}

func GetCode(err error) string {
	if e, ok := err.(*Error); ok {
		return e.Code
	}
	return ""
}

func GetMessage(err error) string {
	if e, ok := err.(*Error); ok {
		return e.Message
	}
	return ""
}

func GetDetails(err error) any {
	if e, ok := err.(*Error); ok {
		return e.Details
	}
	return nil
}

func ErrInvalidJSON() *Error {
	return New(http.StatusBadRequest, "invalid_json", "JSON malformado ou inválido", nil)
}

func ErrInvalidQueryParams(details any) *Error {
	return New(http.StatusBadRequest, "invalid_query_params", "Parâmetros de query inválidos", details)
}
