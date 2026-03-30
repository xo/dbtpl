package row_marshaler

import (
	"reflect"
	"strconv"
	"strings"
)

// TagName is the struct tag key used for row marshaling configuration.
const TagName = "row"

// TagOptions represents the parsed options from a row struct tag.
type TagOptions struct {
	Position   int    // The 1-indexed position in the row
	OmitEmpty  bool   // If true, omit field if zero value during marshal
	Required   bool   // If true, fail unmarshal if field is missing
	Default    string // Default value if field is missing
	HasDefault bool   // Whether a default value was specified
}

// FieldInfo contains metadata about a struct field for marshaling/unmarshaling.
type FieldInfo struct {
	Name       string       // Field name
	Index      int          // Index in struct
	Type       reflect.Type // Field type
	Options    TagOptions   // Parsed tag options
	HasTag     bool         // Whether the field has a row tag
}

// ParseTag parses a row struct tag value and returns TagOptions.
// Tag format: "position[,option1,option2=value,...]"
// Examples: "1", "2,omitempty", "3,required", "4,default=0.0"
func ParseTag(tagValue string) (TagOptions, error) {
	opts := TagOptions{}

	if tagValue == "" || tagValue == "-" {
		return opts, nil
	}

	parts := splitTagParts(tagValue)
	if len(parts) == 0 {
		return opts, NewValidationError("", tagValue, "empty tag value")
	}

	// Parse position (first part)
	pos, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return opts, NewValidationError("", tagValue, "position must be a positive integer")
	}
	if pos < 1 {
		return opts, NewValidationError("", tagValue, "position must be >= 1 (1-indexed)")
	}
	opts.Position = pos

	// Parse options (remaining parts)
	for i := 1; i < len(parts); i++ {
		part := strings.TrimSpace(parts[i])
		if part == "" {
			continue
		}

		// Check for key=value format
		if idx := strings.Index(part, "="); idx != -1 {
			key := strings.TrimSpace(part[:idx])
			value := strings.TrimSpace(part[idx+1:])

			switch key {
			case "default":
				opts.Default = value
				opts.HasDefault = true
			default:
				return opts, NewValidationError("", tagValue, "unknown option: "+key)
			}
		} else {
			// Simple flag option
			switch part {
			case "omitempty":
				opts.OmitEmpty = true
			case "required":
				opts.Required = true
			default:
				return opts, NewValidationError("", tagValue, "unknown option: "+part)
			}
		}
	}

	// Validate conflicting options
	if opts.Required && opts.OmitEmpty {
		return opts, NewValidationError("", tagValue, "conflicting options: required and omitempty")
	}
	if opts.Required && opts.HasDefault {
		return opts, NewValidationError("", tagValue, "conflicting options: required and default")
	}

	return opts, nil
}

// splitTagParts splits a tag value by commas, respecting quoted values.
func splitTagParts(tag string) []string {
	var parts []string
	var current strings.Builder
	inQuote := false
	escaped := false

	for _, r := range tag {
		if escaped {
			current.WriteRune(r)
			escaped = false
			continue
		}

		switch r {
		case '\\':
			escaped = true
			current.WriteRune(r)
		case '"':
			inQuote = !inQuote
			current.WriteRune(r)
		case ',':
			if inQuote {
				current.WriteRune(r)
			} else {
				parts = append(parts, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}

// GetStructFields extracts field information from a struct type.
// It returns a map of position -> FieldInfo for fields with valid row tags.
func GetStructFields(t reflect.Type) (map[int]FieldInfo, error) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, ErrNotStruct
	}

	fields := make(map[int]FieldInfo)
	positions := make(map[int]string) // Track position to field name for duplicate detection

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Get the row tag
		tagValue, hasTag := field.Tag.Lookup(TagName)
		if !hasTag || tagValue == "-" {
			continue
		}

		// Parse the tag
		opts, err := ParseTag(tagValue)
		if err != nil {
			if ve, ok := err.(*ValidationError); ok {
				ve.FieldName = field.Name
			}
			return nil, err
		}

		// Check for duplicate positions
		if existingField, exists := positions[opts.Position]; exists {
			return nil, NewValidationError(field.Name, tagValue,
				"duplicate position "+strconv.Itoa(opts.Position)+" (also used by field "+existingField+")")
		}
		positions[opts.Position] = field.Name

		fields[opts.Position] = FieldInfo{
			Name:    field.Name,
			Index:   i,
			Type:    field.Type,
			Options: opts,
			HasTag:  true,
		}
	}

	return fields, nil
}

// GetMaxPosition returns the maximum position from a map of field infos.
func GetMaxPosition(fields map[int]FieldInfo) int {
	max := 0
	for pos := range fields {
		if pos > max {
			max = pos
		}
	}
	return max
}

// GetSortedPositions returns positions in sorted order.
func GetSortedPositions(fields map[int]FieldInfo) []int {
	if len(fields) == 0 {
		return nil
	}

	max := GetMaxPosition(fields)
	positions := make([]int, 0, len(fields))

	for i := 1; i <= max; i++ {
		if _, exists := fields[i]; exists {
			positions = append(positions, i)
		}
	}

	return positions
}
