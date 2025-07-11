package errors

type AppError interface {
	error
	Code() int
	Type() string
	Details() map[string]interface{}
}

type BaseError struct {
	ErrorCode    int                    `json:"code"`
	ErrorMessage string                 `json:"message"`
	ErrorType    string                 `json:"type"`
	ErrorDetails map[string]interface{} `json:"details,omitempty"`
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

func (e *BaseError) Details() map[string]interface{} {
	return e.ErrorDetails
}
