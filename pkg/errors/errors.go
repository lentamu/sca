package errors

type ErrNotFound struct {
	Msg string
}

func (e ErrNotFound) Error() string {
	return e.Msg
}

type ErrConflict struct {
	Msg string
}

func (e ErrConflict) Error() string {
	return e.Msg
}
