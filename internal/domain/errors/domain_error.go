package domerrs

type DomainError struct {
	message   string
	errorType ErrorType
}

func NewDomainError(message string, errorType ErrorType) *DomainError {
	return &DomainError{
		message:   message,
		errorType: errorType,
	}
}

func (e *DomainError) Error() string {
	return e.message
}

func (e *DomainError) Type() ErrorType {
	return e.errorType
}
