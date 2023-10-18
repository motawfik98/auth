package database

type DuplicateEmailError struct{}

func (e *DuplicateEmailError) Error() string {
	return "Duplicate Email"
}

type NoUserFoundError struct{}

func (e *NoUserFoundError) Error() string {
	return "Wrong email"
}
