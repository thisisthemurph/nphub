package form

// FieldError represents an error for a single input field.
type FieldError struct {
	Field   string
	Message string
}

func NewFieldError(field, message string) FieldError {
	return FieldError{
		Field:   field,
		Message: message,
	}
}

// Error returns a string representation of the FieldError.
func (e FieldError) Error() string {
	return e.Message
}

// FormErrors represents a map of one or more FormError structs.
type FormErrors struct {
	Errors map[string]FieldError
}

func NewFormErrors() *FormErrors {
	return &FormErrors{
		Errors: make(map[string]FieldError),
	}
}

// Set an error for the specified input field.
func (fe *FormErrors) Set(field, message string) {
	fe.Errors[field] = NewFieldError(field, message)
}

// Get an error for the specified input field.
// If the error does not exist, the bool ok return value will be false.
func (fe *FormErrors) Get(field string) (FieldError, bool) {
	err, ok := fe.Errors[field]
	return err, ok
}

// All returns a map of all errors. The key is the input field name and the value is the FormError.
func (fe *FormErrors) All() map[string]FieldError {
	return fe.Errors
}

// HasError returns true if there are any errors present.
func (fe *FormErrors) HasError() bool {
	return len(fe.Errors) > 0
}
