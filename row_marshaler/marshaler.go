package row_marshaler

import (
	"fmt"
	"reflect"
)

// Unmarshal parses a row literal string and stores the result in the value
// pointed to by v. If v is nil or not a pointer to a struct, Unmarshal
// returns an error.
//
// Unmarshal uses the "row" struct tag to determine field positions.
// Fields without a row tag are ignored.
//
// Example:
//
//	type Address struct {
//	    Street  string `row:"1"`
//	    City    string `row:"2"`
//	    ZipCode int    `row:"3"`
//	}
//
//	row := `("123 Main St","Springfield",12345)`
//	var addr Address
//	err := Unmarshal(row, &addr)
func Unmarshal(data string, v interface{}) error {
	// Check for nil
	if v == nil {
		return ErrNilPointer
	}

	// Get reflect value and check it's a pointer
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return ErrNotStruct
	}

	if rv.IsNil() {
		return ErrNilPointer
	}

	// Get the element (the struct)
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	// Parse the row literal
	tokens, err := ParseRowLiteral(data)
	if err != nil {
		return err
	}

	// Get struct field info
	fields, err := GetStructFields(rv.Type())
	if err != nil {
		return err
	}

	// Track which required fields have been set
	requiredFields := make(map[int]string)
	for pos, field := range fields {
		if field.Options.Required {
			requiredFields[pos] = field.Name
		}
	}

	// Assign values to fields
	for i, token := range tokens {
		pos := i + 1 // 1-indexed

		field, exists := fields[pos]
		if !exists {
			// No field mapped to this position, skip
			continue
		}

		// Remove from required tracking
		delete(requiredFields, pos)

		// Get the struct field value
		fieldValue := rv.Field(field.Index)
		if !fieldValue.CanSet() {
			continue
		}

		// Determine the value to convert
		valueStr := token.Value
		if valueStr == "" && field.Options.HasDefault {
			valueStr = field.Options.Default
		}

		// Convert the value
		converted, err := ConvertToType(valueStr, field.Type, field.Name, pos)
		if err != nil {
			return err
		}

		// Set the field value
		fieldValue.Set(converted)
	}

	// Apply defaults for missing fields and check required fields
	for pos, field := range fields {
		if pos > len(tokens) {
			// Field position is beyond the token count
			if field.Options.HasDefault {
				fieldValue := rv.Field(field.Index)
				if fieldValue.CanSet() {
					converted, err := ConvertToType(field.Options.Default, field.Type, field.Name, pos)
					if err != nil {
						return err
					}
					fieldValue.Set(converted)
				}
				// Remove from required tracking only if we handled it with a default
				delete(requiredFields, pos)
			}
			// Note: Required fields without defaults stay in requiredFields
		}
	}

	// Check for missing required fields
	if len(requiredFields) > 0 {
		// Get the first missing required field for error message
		for pos, name := range requiredFields {
			return NewValidationError(name, "", fmt.Sprintf("required field at position %d is missing", pos))
		}
	}

	return nil
}

// Marshal returns the row literal encoding of v.
//
// Marshal uses the "row" struct tag to determine field positions and
// serializes fields in order of their positions.
//
// String values are enclosed in double quotes. Numeric values are not quoted.
// Boolean values are lowercase (true/false).
//
// If v is not a struct or pointer to a struct, Marshal returns an error.
//
// Example:
//
//	addr := Address{
//	    Street:  "123 Main St",
//	    City:    "Springfield",
//	    ZipCode: 12345,
//	}
//	result, err := Marshal(addr)
//	// result: `("123 Main St","Springfield",12345)`
func Marshal(v interface{}) (string, error) {
	// Check for nil
	if v == nil {
		return "", ErrNilPointer
	}

	// Get reflect value
	rv := reflect.ValueOf(v)

	// Dereference pointer if needed
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return "", ErrNilPointer
		}
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return "", ErrNotStruct
	}

	// Get struct field info
	fields, err := GetStructFields(rv.Type())
	if err != nil {
		return "", err
	}

	if len(fields) == 0 {
		return "()", nil
	}

	// Get sorted positions to output in order
	positions := GetSortedPositions(fields)

	// Build the output values
	var values []string
	lastIncludedPos := 0

	for _, pos := range positions {
		field := fields[pos]
		fieldValue := rv.Field(field.Index)

		// Check omitempty
		if field.Options.OmitEmpty && IsZeroValue(fieldValue) {
			continue
		}

		// Add placeholder empty values for any skipped positions
		for i := lastIncludedPos + 1; i < pos; i++ {
			// Check if there's a field at this position
			if skipField, exists := fields[i]; exists {
				if skipField.Options.OmitEmpty && IsZeroValue(rv.Field(skipField.Index)) {
					continue
				}
			}
			values = append(values, `""`)
		}

		// Format the value
		formatted, err := FormatValue(fieldValue)
		if err != nil {
			return "", NewMarshalError(field.Name, pos, fieldValue.Interface(), err.Error())
		}

		values = append(values, formatted)
		lastIncludedPos = pos
	}

	return BuildRowLiteral(values), nil
}

// MarshalWithOptions allows specifying custom marshaling behavior.
type MarshalOptions struct {
	// IncludeEmpty includes zero-value fields even if they have omitempty.
	IncludeEmpty bool
	// FillGaps includes empty values for gaps in position sequence.
	FillGaps bool
}

// MarshalWith returns the row literal encoding of v using the specified options.
func MarshalWith(v interface{}, opts MarshalOptions) (string, error) {
	if v == nil {
		return "", ErrNilPointer
	}

	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return "", ErrNilPointer
		}
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return "", ErrNotStruct
	}

	fields, err := GetStructFields(rv.Type())
	if err != nil {
		return "", err
	}

	if len(fields) == 0 {
		return "()", nil
	}

	maxPos := GetMaxPosition(fields)
	var values []string

	for pos := 1; pos <= maxPos; pos++ {
		field, exists := fields[pos]

		if !exists {
			if opts.FillGaps {
				values = append(values, `""`)
			}
			continue
		}

		fieldValue := rv.Field(field.Index)

		// Check omitempty (unless IncludeEmpty is set)
		if !opts.IncludeEmpty && field.Options.OmitEmpty && IsZeroValue(fieldValue) {
			if opts.FillGaps {
				values = append(values, `""`)
			}
			continue
		}

		formatted, err := FormatValue(fieldValue)
		if err != nil {
			return "", NewMarshalError(field.Name, pos, fieldValue.Interface(), err.Error())
		}

		values = append(values, formatted)
	}

	return BuildRowLiteral(values), nil
}

// UnmarshalStrict is like Unmarshal but returns an error if there are
// extra tokens in the row that don't map to any field.
func UnmarshalStrict(data string, v interface{}) error {
	if v == nil {
		return ErrNilPointer
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return ErrNotStruct
	}

	if rv.IsNil() {
		return ErrNilPointer
	}

	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	tokens, err := ParseRowLiteral(data)
	if err != nil {
		return err
	}

	fields, err := GetStructFields(rv.Type())
	if err != nil {
		return err
	}

	maxPos := GetMaxPosition(fields)
	if len(tokens) > maxPos {
		return NewValidationError("", "", "row has more values than struct has tagged fields")
	}

	// Use regular unmarshal for the rest
	return Unmarshal(data, v)
}

// Valid checks if a string is a valid row literal format.
func Valid(data string) bool {
	_, err := ParseRowLiteral(data)
	return err == nil
}

// TokenCount returns the number of tokens in a row literal string.
// Returns -1 if the string is not a valid row literal.
func TokenCount(data string) int {
	tokens, err := ParseRowLiteral(data)
	if err != nil {
		return -1
	}
	return len(tokens)
}
