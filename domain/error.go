package domain

type CheckFailedError struct {
	Message string
}

func (e *CheckFailedError) Error() string {
	return e.Message
}
