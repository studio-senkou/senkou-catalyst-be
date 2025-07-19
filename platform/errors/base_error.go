package errors

type AppError struct {
	Code    int `json:"code"`
	Details any `json:"details,omitempty"`
}

func NewAppError(code int, details any) *AppError {
	return &AppError{
		Code:    code,
		Details: details,
	}
}

type BaseError struct {
	ErrorCode    int    `json:"code"`
	ErrorMessage string `json:"message"`
	ErrorType    string `json:"type"`
	ErrorDetails any    `json:"details,omitempty"`
}

func (e *BaseError) Error() string {
	return e.ErrorMessage
}

func (e *BaseError) Code() int {
	return e.ErrorCode
}

func (e *BaseError) Type() string {
	return e.ErrorType
}

func (e *BaseError) Details() any {
	return e.ErrorDetails
}
