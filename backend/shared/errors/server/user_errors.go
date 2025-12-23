package server

type SessionErrType int
type UserErrType int

const (
	SessionExpired SessionErrType = iota
	SessionDiffIpAddr
	SessionDiffUserAgent
)

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

type SessionErr struct {
	Type    SessionErrType
	Message string
}

func (e *SessionErr) Error() string {
	return e.Message
}

func NewSessionErr(errType SessionErrType, message string) *SessionErr {
	return &SessionErr{Type: errType, Message: message}
}

type UserErr struct {
	Type    UserErrType
	Message string
}

func (e *UserErr) Error() string {
	return e.Message
}

func NewUserErr(errType UserErrType, message string) *UserErr {
	return &UserErr{Type: errType, Message: message}
}
