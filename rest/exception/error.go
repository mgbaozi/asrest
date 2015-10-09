package exception

type Error struct {
	Error     int    `json:"error"`
	Exception string `json:"exception"`
}

func NewError() *Error {
	return &Error{}
}

func NotFound() *Error {
	err := NewError()
	err.Error = 1
	err.Exception = "Not Found"
	return err
}
