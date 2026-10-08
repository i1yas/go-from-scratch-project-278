package links

import (
	"fmt"
	"strings"
)

// FieldError contains information about field error
type FieldError struct {
	Field string
	Err   error
}

// Error forms field error message
func (e *FieldError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Err.Error())
}

// Unwrap unwraps field error
func (e *FieldError) Unwrap() error {
	return e.Err
}

// ValidationError contains field errors
type ValidationError struct {
	Fields []FieldError
}

// Error forms validation error message
func (v *ValidationError) Error() string {
	parts := make([]string, len(v.Fields))

	for i, field := range v.Fields {
		parts[i] = field.Error()
	}

	return fmt.Sprintf("validation failed: %s", strings.Join(parts, "; "))
}

// Unwrap unwraps field errors
func (v *ValidationError) Unwrap() []error {
	errs := make([]error, len(v.Fields))

	for i, field := range v.Fields {
		errs[i] = &field
	}

	return errs
}

func (v *ValidationError) add(field string, err error) {
	if err != nil {
		v.Fields = append(v.Fields, FieldError{
			Field: field,
			Err:   err,
		})
	}
}

func (v *ValidationError) err() error {
	if len(v.Fields) == 0 {
		return nil
	}

	return v
}
