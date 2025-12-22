package users

type TokenErr struct {
	Code    int
	Message string
}

func (e *TokenErr) Error() string {
	return e.Message
}

func NewTokenErr(code int, message string) *TokenErr {
	return &TokenErr{Code: code, Message: message}
}
