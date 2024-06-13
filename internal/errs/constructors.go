package errs

func ErrNotFound(msg string, err error) error {
	return base{Code: 404, Message: msg, Err: err}
}

func ErrBadRequest(msg string, err error) error {
	return base{Code: 400, Message: msg, Err: err}
}

func ErrConflict(msg string, err error) error {
	return base{Code: 409, Message: msg, Err: err}
}

func ErrUnauthorized(msg string, err error) error {
	return base{Code: 401, Message: msg, Err: err}
}

func ErrInternal(err error) error {
	return base{Code: 500, Message: "Internal server error", Err: err}
}

func ErrForbidden(msg string, err error) error {
	return base{Code: 403, Message: msg, Err: err}
}
