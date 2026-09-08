// Package row_marshaler provides functions for marshaling and unmarshaling
// Go structs to and from a row-literal format using custom struct tags.
package row_marshaler

import (
	"fmt"
)

// ParseError represents an error that occurred during parsing of a row literal.
type ParseError struct {
	Position int    // Position in the input string where the error occurred
	Message  string // Description of the parsing error
	Input    string // The original input that caused the error
}

func (e *ParseError) Error() string {
	if e.Position >= 0 {
		return fmt.Sprintf("parse error at position %d: %s", e.Position, e.Message)
	}
	return fmt.Sprintf("parse error: %s", e.Message)
}

// NewParseError creates a new ParseError with the given details.
func NewParseError(position int, message string, input string) *ParseError {
	return &ParseError{
		Position: position,
		Message:  message,
		Input:    input,
	}
}

// TypeMismatchError represents an error when a value cannot be converted to the target type.
type TypeMismatchError struct {
	FieldName    string // Name of the struct field
	FieldPos     int    // Position in the row (from tag)
	Value        string // The actual value that failed conversion
	ExpectedType string // The expected Go type
	Reason       string // Additional details about why conversion failed
}

func (e *TypeMismatchError) Error() string {
	return fmt.Sprintf("unmarshal error: field %q (position %d): cannot parse %q as %s: %s",
		e.FieldName, e.FieldPos, e.Value, e.ExpectedType, e.Reason)
}

// NewTypeMismatchError creates a new TypeMismatchError with the given details.
func NewTypeMismatchError(fieldName string, fieldPos int, value string, expectedType string, reason string) *TypeMismatchError {
	return &TypeMismatchError{
		FieldName:    fieldName,
		FieldPos:     fieldPos,
		Value:        value,
		ExpectedType: expectedType,
		Reason:       reason,
	}
}

// ValidationError represents an error in struct tag validation or field configuration.
type ValidationError struct {
	FieldName string // Name of the struct field
	TagValue  string // The value of the tag that caused the error
	Message   string // Description of the validation error
}

func (e *ValidationError) Error() string {
	if e.FieldName != "" {
		return fmt.Sprintf("validation error: field %q with tag %q: %s", e.FieldName, e.TagValue, e.Message)
	}
	return fmt.Sprintf("validation error: %s", e.Message)
}

// NewValidationError creates a new ValidationError with the given details.
func NewValidationError(fieldName string, tagValue string, message string) *ValidationError {
	return &ValidationError{
		FieldName: fieldName,
		TagValue:  tagValue,
		Message:   message,
	}
}

// MarshalError represents an error that occurred during marshaling.
type MarshalError struct {
	FieldName string // Name of the struct field
	FieldPos  int    // Position in the row (from tag)
	Value     any    // The value that failed to marshal
	Reason    string // Description of the error
}

func (e *MarshalError) Error() string {
	return fmt.Sprintf("marshal error: field %q (position %d): cannot marshal value: %s",
		e.FieldName, e.FieldPos, e.Reason)
}

// NewMarshalError creates a new MarshalError with the given details.
func NewMarshalError(fieldName string, fieldPos int, value any, reason string) *MarshalError {
	return &MarshalError{
		FieldName: fieldName,
		FieldPos:  fieldPos,
		Value:     value,
		Reason:    reason,
	}
}

// ErrNilPointer is returned when a nil pointer is passed to Marshal or Unmarshal.
var ErrNilPointer = fmt.Errorf("cannot marshal/unmarshal nil pointer")

// ErrNotStruct is returned when a non-struct type is passed.
var ErrNotStruct = fmt.Errorf("value must be a struct or pointer to struct")

// ErrEmptyInput is returned when an empty string is passed to Unmarshal.
var ErrEmptyInput = fmt.Errorf("input string is empty")
