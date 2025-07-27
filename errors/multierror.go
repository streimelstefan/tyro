package errors

import (
	"strings"
)

// MultiError accumulates multiple errors and implements the error interface.
type MultiError struct {
	errs []error
}

// New creates a new MultiError.
func New() *MultiError {
	return &MultiError{errs: make([]error, 0)}
}

// Add appends an error to the MultiError. Ignores nil errors.
func (m *MultiError) Add(err error) {
	if err != nil {
		m.errs = append(m.errs, err)
	}
}

// Error implements the error interface.
func (m *MultiError) Error() string {
	if len(m.errs) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("Multiple errors occurred:")
	for _, err := range m.errs {
		sb.WriteString("\n- ")
		sb.WriteString(err.Error())
	}
	return sb.String()
}

// Errors returns the slice of errors.
func (m *MultiError) Errors() []error {
	return m.errs
}

// HasErrors returns true if there are any errors accumulated.
func (m *MultiError) HasErrors() bool {
	return len(m.errs) > 0
}
