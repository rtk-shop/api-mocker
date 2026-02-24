package order

// todo: make it as share error
type InvalidRequestError struct {
	Message string
	Reason  string
}

func (e *InvalidRequestError) Error() string {
	return e.Message
}

func NewInvalidRequestError(reason string) *InvalidRequestError {
	return &InvalidRequestError{
		Message: "invalid request",
		Reason:  reason,
	}
}
