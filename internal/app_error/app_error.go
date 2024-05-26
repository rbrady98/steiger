package apperror

type AppError struct {
	Msg  string `json:"error"`
	Code int    `json:"-"`
}

func New(message string, code int) *AppError {
	return &AppError{
		Msg:  message,
		Code: code,
	}
}

func (e *AppError) Error() string {
	return e.Msg
}
