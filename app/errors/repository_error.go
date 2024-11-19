package errors

type UniqueViolationError struct {
	message string
}

func NewUniqueViolationError(message string) *UniqueViolationError {
	return &UniqueViolationError{message}
}

func (e *UniqueViolationError) Error() string {
	return e.message
}

type NotFoundError struct {
	message string
}

func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{message}
}

func (e *NotFoundError) Error() string {
	return e.message
}

type WrongLoginError struct {
	message string
}

func NewWrongLoginError(message string) *WrongLoginError {
	return &WrongLoginError{message}
}

func (e *WrongLoginError) Error() string {
	return e.message
}

type BadRequestError struct {
	message string
}

func NewBadRequestError(message string) *BadRequestError {
	return &BadRequestError{message}
}

func (e *BadRequestError) Error() string {
	return e.message
}
