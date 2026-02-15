package validation

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	Field   string
	Message string
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	errors := make([]string, len(v))
	for i, err := range v {
		errors[i] = fmt.Sprintf("%s: %s", err.Field, err.Message)
	}
	return strings.Join(errors, ", ")
}

func (v ValidationErrors) Add(field, message string) ValidationErrors {
	return append(v, ValidationError{Field: field, Message: message})
}

func (v ValidationErrors) HasErrors() bool {
	return len(v) > 0
}
