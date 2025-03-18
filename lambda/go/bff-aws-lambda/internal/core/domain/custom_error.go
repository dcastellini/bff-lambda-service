package domain

import "net/http"

type CustomError struct {
	MessageError   string `json:"message"`
	HTTPStatusCode int    `json:"-"`
	MessageCode    int    `json:"code"`
}

var (
	ErrGeneric = CustomError{
		MessageCode:    1,
		MessageError:   "Something went wrong",
		HTTPStatusCode: http.StatusBadRequest,
	}
	ErrNotImplemented = CustomError{
		MessageCode:    2,
		MessageError:   "Country not implemented",
		HTTPStatusCode: http.StatusInternalServerError,
	}
	ErrInvalidRequest = CustomError{
		MessageCode:    3,
		MessageError:   "There are one or more missing or invalid required fields.",
		HTTPStatusCode: http.StatusBadRequest,
	}
)

func (e CustomError) Error() string {
	return e.MessageError
}

func BuildCustomError(err error) CustomError {
	return CustomError{
		MessageError:   err.Error(),
		MessageCode:    1,
		HTTPStatusCode: http.StatusBadRequest,
	}
}
