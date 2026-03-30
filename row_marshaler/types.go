package row_marshaler

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
)

var (
	nullStringType  = reflect.TypeOf(sql.NullString{})
	nullBoolType    = reflect.TypeOf(sql.NullBool{})
	nullInt16Type   = reflect.TypeOf(sql.NullInt16{})
	nullInt32Type   = reflect.TypeOf(sql.NullInt32{})
	nullInt64Type   = reflect.TypeOf(sql.NullInt64{})
	nullByteType    = reflect.TypeOf(sql.NullByte{})
	nullFloat64Type = reflect.TypeOf(sql.NullFloat64{})
	nullTimeType    = reflect.TypeOf(sql.NullTime{})
)

// Stringer interface for custom string conversion.
type Stringer interface {
	String() string
}

// StringParser interface for custom parsing from string.
type StringParser interface {
	ParseString(s string) error
}

// ConvertToType converts a string value to the target reflect.Type.
// fieldName and fieldPos are used for error reporting.
func ConvertToType(value string, targetType reflect.Type, fieldName string, fieldPos int) (reflect.Value, error) {
	// Handle pointers
	if targetType.Kind() == reflect.Ptr {
		elemType := targetType.Elem()

		// Empty value for pointer means nil
		if value == "" {
			return reflect.Zero(targetType), nil
		}

		// Convert to the element type
		elemValue, err := ConvertToType(value, elemType, fieldName, fieldPos)
		if err != nil {
			return reflect.Value{}, err
		}

		// Create pointer and set value
		ptr := reflect.New(elemType)
		ptr.Elem().Set(elemValue)
		return ptr, nil
	}

	if converted, handled, err := convertSQLNull(value, targetType, fieldName, fieldPos); handled {
		return converted, err
	}

	// Handle basic types
	switch targetType.Kind() {
	case reflect.String:
		return reflect.ValueOf(value), nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return convertToInt(value, targetType, fieldName, fieldPos)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return convertToUint(value, targetType, fieldName, fieldPos)

	case reflect.Float32, reflect.Float64:
		return convertToFloat(value, targetType, fieldName, fieldPos)

	case reflect.Bool:
		return convertToBool(value, fieldName, fieldPos)

	case reflect.Struct:
		// Special handling for time.Time
		if targetType == reflect.TypeOf(time.Time{}) {
			return convertToTime(value, fieldName, fieldPos)
		}

		ptrType := reflect.PtrTo(targetType)

		// If the struct implements sql.Scanner, delegate parsing to it so
		// wrapper types (eg, generated Null composites) can participate.
		if ptrType.Implements(reflect.TypeOf((*sql.Scanner)(nil)).Elem()) {
			newVal := reflect.New(targetType)
			scanner := newVal.Interface().(sql.Scanner)
			if err := scanner.Scan([]byte(value)); err != nil {
				return reflect.Value{}, NewTypeMismatchError(fieldName, fieldPos, value, targetType.String(), err.Error())
			}
			return newVal.Elem(), nil
		}

		// Check if the type implements StringParser
		if ptrType.Implements(reflect.TypeOf((*StringParser)(nil)).Elem()) {
			newVal := reflect.New(targetType)
			parser := newVal.Interface().(StringParser)
			if err := parser.ParseString(value); err != nil {
				return reflect.Value{}, NewTypeMismatchError(fieldName, fieldPos, value, targetType.String(), err.Error())
			}
			return newVal.Elem(), nil
		}

		return reflect.Value{}, NewTypeMismatchError(fieldName, fieldPos, value, targetType.String(), "unsupported struct type")

	case reflect.Slice:
		return convertToSlice(value, targetType, fieldName, fieldPos)

	default:
		return reflect.Value{}, NewTypeMismatchError(fieldName, fieldPos, value, targetType.String(), "unsupported type")
	}
}

func convertSQLNull(value string, targetType reflect.Type, fieldName string, fieldPos int) (reflect.Value, bool, error) {
	if targetType == nullStringType {
		if value == "" || strings.EqualFold(value, "null") {
			return reflect.ValueOf(sql.NullString{Valid: false}), true, nil
		}
		return reflect.ValueOf(sql.NullString{String: value, Valid: true}), true, nil
	}

	if targetType == nullBoolType {
		if value == "" || strings.EqualFold(value, "null") {
			return reflect.ValueOf(sql.NullBool{Valid: false}), true, nil
		}
		parsed, err := convertToBool(value, fieldName, fieldPos)
		if err != nil {
			return reflect.Value{}, true, err
		}
		return reflect.ValueOf(sql.NullBool{Bool: parsed.Bool(), Valid: true}), true, nil
	}

	if targetType == nullInt16Type || targetType == nullInt32Type || targetType == nullInt64Type {
		if value == "" || strings.EqualFold(value, "null") {
			switch targetType {
			case nullInt16Type:
				return reflect.ValueOf(sql.NullInt16{Valid: false}), true, nil
			case nullInt32Type:
				return reflect.ValueOf(sql.NullInt32{Valid: false}), true, nil
			default:
				return reflect.ValueOf(sql.NullInt64{Valid: false}), true, nil
			}
		}

		intTarget := reflect.TypeOf(int64(0))
		switch targetType {
		case nullInt16Type:
			intTarget = reflect.TypeOf(int16(0))
		case nullInt32Type:
			intTarget = reflect.TypeOf(int32(0))
		}

		parsed, err := convertToInt(value, intTarget, fieldName, fieldPos)
		if err != nil {
			return reflect.Value{}, true, err
		}

		switch targetType {
		case nullInt16Type:
			return reflect.ValueOf(sql.NullInt16{Int16: int16(parsed.Int()), Valid: true}), true, nil
		case nullInt32Type:
			return reflect.ValueOf(sql.NullInt32{Int32: int32(parsed.Int()), Valid: true}), true, nil
		default:
			return reflect.ValueOf(sql.NullInt64{Int64: parsed.Int(), Valid: true}), true, nil
		}
	}

	if targetType == nullByteType {
		if value == "" || strings.EqualFold(value, "null") {
			return reflect.ValueOf(sql.NullByte{Valid: false}), true, nil
		}
		parsed, err := convertToUint(value, reflect.TypeOf(uint8(0)), fieldName, fieldPos)
		if err != nil {
			return reflect.Value{}, true, err
		}
		return reflect.ValueOf(sql.NullByte{Byte: uint8(parsed.Uint()), Valid: true}), true, nil
	}

	if targetType == nullFloat64Type {
		if value == "" || strings.EqualFold(value, "null") {
			return reflect.ValueOf(sql.NullFloat64{Valid: false}), true, nil
		}
		parsed, err := convertToFloat(value, reflect.TypeOf(float64(0)), fieldName, fieldPos)
		if err != nil {
			return reflect.Value{}, true, err
		}
		return reflect.ValueOf(sql.NullFloat64{Float64: parsed.Float(), Valid: true}), true, nil
	}

	if targetType == nullTimeType {
		if value == "" || strings.EqualFold(value, "null") {
			return reflect.ValueOf(sql.NullTime{Valid: false}), true, nil
		}
		parsed, err := convertToTime(value, fieldName, fieldPos)
		if err != nil {
			return reflect.Value{}, true, err
		}
		return reflect.ValueOf(sql.NullTime{Time: parsed.Interface().(time.Time), Valid: true}), true, nil
	}

	return reflect.Value{}, false, nil
}

func convertToInt(value string, targetType reflect.Type, fieldName string, fieldPos int) (reflect.Value, error) {
	if value == "" {
		return reflect.Zero(targetType), nil
	}

	value = strings.TrimSpace(value)

	// Determine bit size
	var bitSize int
	switch targetType.Kind() {
	case reflect.Int8:
		bitSize = 8
	case reflect.Int16:
		bitSize = 16
	case reflect.Int32:
		bitSize = 32
	case reflect.Int64:
		bitSize = 64
	default: // reflect.Int
		bitSize = strconv.IntSize
	}

	parsed, err := strconv.ParseInt(value, 10, bitSize)
	if err != nil {
		return reflect.Value{}, NewTypeMismatchError(fieldName, fieldPos, value, targetType.String(), err.Error())
	}

	result := reflect.New(targetType).Elem()
	result.SetInt(parsed)
	return result, nil
}

func convertToUint(value string, targetType reflect.Type, fieldName string, fieldPos int) (reflect.Value, error) {
	if value == "" {
		return reflect.Zero(targetType), nil
	}

	value = strings.TrimSpace(value)

	// Determine bit size
	var bitSize int
	switch targetType.Kind() {
	case reflect.Uint8:
		bitSize = 8
	case reflect.Uint16:
		bitSize = 16
	case reflect.Uint32:
		bitSize = 32
	case reflect.Uint64:
		bitSize = 64
	default: // reflect.Uint
		bitSize = strconv.IntSize
	}

	parsed, err := strconv.ParseUint(value, 10, bitSize)
	if err != nil {
		return reflect.Value{}, NewTypeMismatchError(fieldName, fieldPos, value, targetType.String(), err.Error())
	}

	result := reflect.New(targetType).Elem()
	result.SetUint(parsed)
	return result, nil
}

func convertToFloat(value string, targetType reflect.Type, fieldName string, fieldPos int) (reflect.Value, error) {
	if value == "" {
		return reflect.Zero(targetType), nil
	}

	value = strings.TrimSpace(value)

	var bitSize int
	switch targetType.Kind() {
	case reflect.Float32:
		bitSize = 32
	default: // reflect.Float64
		bitSize = 64
	}

	parsed, err := strconv.ParseFloat(value, bitSize)
	if err != nil {
		return reflect.Value{}, NewTypeMismatchError(fieldName, fieldPos, value, targetType.String(), err.Error())
	}

	result := reflect.New(targetType).Elem()
	result.SetFloat(parsed)
	return result, nil
}

func convertToBool(value string, fieldName string, fieldPos int) (reflect.Value, error) {
	if value == "" {
		return reflect.ValueOf(false), nil
	}

	value = strings.TrimSpace(strings.ToLower(value))

	switch value {
	case "true", "1", "yes", "on", "t", "y":
		return reflect.ValueOf(true), nil
	case "false", "0", "no", "off", "f", "n":
		return reflect.ValueOf(false), nil
	default:
		return reflect.Value{}, NewTypeMismatchError(fieldName, fieldPos, value, "bool", "must be true/false")
	}
}

func convertToTime(value string, fieldName string, fieldPos int) (reflect.Value, error) {
	if value == "" {
		return reflect.ValueOf(time.Time{}), nil
	}

	value = strings.TrimSpace(value)

	// Try RFC3339 first
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		// Try RFC3339Nano
		t, err = time.Parse(time.RFC3339Nano, value)
		if err != nil {
			// Try other common formats
			formats := []string{
				"2006-01-02",
				"2006-01-02 15:04:05",
				"2006-01-02T15:04:05",
				"01/02/2006",
				"01/02/2006 15:04:05",
			}
			for _, format := range formats {
				t, err = time.Parse(format, value)
				if err == nil {
					return reflect.ValueOf(t), nil
				}
			}
			return reflect.Value{}, NewTypeMismatchError(fieldName, fieldPos, value, "time.Time", "invalid time format")
		}
	}
	return reflect.ValueOf(t), nil
}

func convertToSlice(value string, targetType reflect.Type, fieldName string, fieldPos int) (reflect.Value, error) {
	if value == "" {
		return reflect.MakeSlice(targetType, 0, 0), nil
	}

	value = strings.TrimSpace(value)

	// Accept either JSON-like [a,b] or PostgreSQL-style {a,b} array notation
	if len(value) < 2 || (value[0] != '[' && value[0] != '{') || (value[len(value)-1] != ']' && value[len(value)-1] != '}') {
		return reflect.Value{}, NewTypeMismatchError(fieldName, fieldPos, value, targetType.String(), "slice must be in [val1,val2,...] or {val1,val2,...} format")
	}

	// Extract content
	content := value[1 : len(value)-1]
	if content == "" {
		return reflect.MakeSlice(targetType, 0, 0), nil
	}

	// Parse elements
	elements := splitSliceElements(content)
	elemType := targetType.Elem()
	slice := reflect.MakeSlice(targetType, 0, len(elements))

	for i, elem := range elements {
		elem = strings.TrimSpace(elem)
		elemValue, err := ConvertToType(elem, elemType, fmt.Sprintf("%s[%d]", fieldName, i), fieldPos)
		if err != nil {
			return reflect.Value{}, err
		}
		slice = reflect.Append(slice, elemValue)
	}

	return slice, nil
}

// splitSliceElements splits comma-separated slice elements, respecting quotes.
func splitSliceElements(content string) []string {
	var elements []string
	var current strings.Builder
	inQuote := false
	escaped := false

	for i := 0; i < len(content); i++ {
		ch := content[i]

		if escaped {
			current.WriteByte(ch)
			escaped = false
			continue
		}

		switch ch {
		case '\\':
			escaped = true
			current.WriteByte(ch)
		case '"':
			inQuote = !inQuote
			current.WriteByte(ch)
		case ',':
			if inQuote {
				current.WriteByte(ch)
			} else {
				elements = append(elements, current.String())
				current.Reset()
			}
		default:
			current.WriteByte(ch)
		}
	}

	if current.Len() > 0 {
		elements = append(elements, current.String())
	}

	return elements
}

// FormatValue formats a value for inclusion in a row literal.
// Strings are quoted, numbers are not.
func FormatValue(v reflect.Value) (string, error) {
	// Handle nil pointers
	if !v.IsValid() {
		return `""`, nil
	}

	// Dereference pointers
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return `""`, nil
		}
		v = v.Elem()
	}

	if formatted, handled, err := formatSQLNull(v); handled {
		return formatted, err
	}

	// Handle driver.Valuer after nulls to preserve sql.Null* semantics while
	// allowing custom wrapper types (eg, generated Null composites).
	if v.CanInterface() {
		if valuer, ok := v.Interface().(driver.Valuer); ok {
			literal, err := valuer.Value()
			if err != nil {
				return "", err
			}
			if literal == nil {
				return `""`, nil
			}

			switch lit := literal.(type) {
			case []byte:
				return EscapeString(string(lit)), nil
			case string:
				return EscapeString(lit), nil
			default:
				return fmt.Sprint(lit), nil
			}
		}
	}

	switch v.Kind() {
	case reflect.String:
		return EscapeString(v.String()), nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10), nil

	case reflect.Float32:
		return strconv.FormatFloat(v.Float(), 'f', -1, 32), nil

	case reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64), nil

	case reflect.Bool:
		if v.Bool() {
			return "true", nil
		}
		return "false", nil

	case reflect.Struct:
		// Handle time.Time
		if t, ok := v.Interface().(time.Time); ok {
			if t.IsZero() {
				return `""`, nil
			}
			return EscapeString(t.Format(time.RFC3339)), nil
		}

		// Check for Stringer interface
		if stringer, ok := v.Interface().(fmt.Stringer); ok {
			return EscapeString(stringer.String()), nil
		}

		// Attempt to marshal nested structs using row tags.
		if v.CanInterface() {
			if literal, err := Marshal(v.Interface()); err == nil {
				return literal, nil
			}
		}

		return "", fmt.Errorf("unsupported struct type: %s", v.Type().String())

	case reflect.Slice:
		return formatSlice(v)

	default:
		return "", fmt.Errorf("unsupported type: %s", v.Kind().String())
	}
}

func formatSQLNull(v reflect.Value) (string, bool, error) {
	switch nv := v.Interface().(type) {
	case sql.NullString:
		if !nv.Valid {
			return "NULL", true, nil
		}
		return EscapeString(nv.String), true, nil
	case sql.NullBool:
		if !nv.Valid {
			return "NULL", true, nil
		}
		if nv.Bool {
			return "true", true, nil
		}
		return "false", true, nil
	case sql.NullInt16:
		if !nv.Valid {
			return "NULL", true, nil
		}
		return strconv.FormatInt(int64(nv.Int16), 10), true, nil
	case sql.NullInt32:
		if !nv.Valid {
			return "NULL", true, nil
		}
		return strconv.FormatInt(int64(nv.Int32), 10), true, nil
	case sql.NullInt64:
		if !nv.Valid {
			return "NULL", true, nil
		}
		return strconv.FormatInt(nv.Int64, 10), true, nil
	case sql.NullByte:
		if !nv.Valid {
			return "NULL", true, nil
		}
		return strconv.FormatUint(uint64(nv.Byte), 10), true, nil
	case sql.NullFloat64:
		if !nv.Valid {
			return "NULL", true, nil
		}
		return strconv.FormatFloat(nv.Float64, 'f', -1, 64), true, nil
	case sql.NullTime:
		if !nv.Valid {
			return "NULL", true, nil
		}
		if nv.Time.IsZero() {
			return `""`, true, nil
		}
		return EscapeString(nv.Time.Format(time.RFC3339)), true, nil
	}

	return "", false, nil
}

func formatSlice(v reflect.Value) (string, error) {
	if v.IsNil() || v.Len() == 0 {
		return "{}", nil
	}

	if v.CanInterface() {
		if literal, err := pq.Array(v.Interface()).Value(); err == nil {
			switch lit := literal.(type) {
			case string:
				return EscapeString(lit), nil
			case []byte:
				return EscapeString(string(lit)), nil
			}
		}
	}

	var elements []string
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		formatted, err := FormatValue(elem)
		if err != nil {
			return "", err
		}
		if elem.Kind() == reflect.String {
			formatted = strings.Trim(formatted, `"`)
		}
		elements = append(elements, formatted)
	}

	return "{" + strings.Join(elements, ",") + "}", nil
}

// IsZeroValue checks if a reflect.Value is a zero value.
func IsZeroValue(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}

	switch v.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return v.IsNil()
	default:
		return v.IsZero()
	}
}
