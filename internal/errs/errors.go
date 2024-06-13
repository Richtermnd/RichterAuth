package errs

const (
	BadRequestType   = "BadRequest"
	ConflictType     = "Conflict"
	NotFoundType     = "NotFound"
	UnknownErrorType = "UnknownError"
	InternalType     = "Internal"
	UnauthorizedType = "Unauthorized"
)

// var errsMap = map[string]interface{}{}

type base struct {
	Code    int    `json:"-"`
	Type    string `json:"type"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (b base) Error() string {
	return b.Message
}

func (b base) Unwrap() error {
	return b.Err
}

func (b base) HttpCode() int {
	return b.Code
}
