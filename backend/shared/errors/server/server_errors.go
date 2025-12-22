package server

type ServerErr struct {
	Code    int
	Message string
}

func (e *ServerErr) Error() string {
	return e.Message
}

func NewServerErr(code int, message string) *ServerErr {
	return &ServerErr{Code: code, Message: message}
}
