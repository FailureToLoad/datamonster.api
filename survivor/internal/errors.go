package repo

import "fmt"

type DuplicateNameError struct {
	msg string
}

func (e DuplicateNameError) Error() string {
	return e.msg
}

func NewDuplicateNameError(survivorName string) DuplicateNameError {
	return DuplicateNameError{
		msg: fmt.Sprintf("Survivor with name %s already exists", survivorName),
	}
}
