{{ define "enum" }}
{{- $e := .Data -}}
// {{ $e.GoName }} is the '{{ $e.SQLName }}' enum type from schema '{{ schema }}'.
type {{ $e.GoName }} uint16

// {{ $e.GoName }} values.
const (
{{ range $e.Values -}}
	// {{ $e.GoName }}{{ .GoName }} is the '{{ .SQLName }}' {{ $e.SQLName }}.
	{{ $e.GoName }}{{ .GoName }} {{ $e.GoName }} = {{ .ConstValue }}
{{ end -}}
)

// String satisfies the [fmt.Stringer] interface.
func ({{ short $e.GoName }} {{ $e.GoName }}) String() string {
	switch {{ short $e.GoName }} {
{{ range $e.Values -}}
	case {{ $e.GoName }}{{ .GoName }}:
		return "{{ .SQLName }}"
{{ end -}}
	}
	return fmt.Sprintf("{{ $e.GoName }}(%d)", {{ short $e.GoName }})
}

// MarshalText marshals [{{ $e.GoName }}] into text.
func ({{ short $e.GoName }} {{ $e.GoName }}) MarshalText() ([]byte, error) {
	return []byte({{ short $e.GoName }}.String()), nil
}

// UnmarshalText unmarshals [{{ $e.GoName }}] from text.
func ({{ short $e.GoName }} *{{ $e.GoName }}) UnmarshalText(buf []byte) error {
	switch str := string(buf); str {
{{ range $e.Values -}}
	case "{{ .SQLName }}":
		*{{ short $e.GoName }} = {{ $e.GoName }}{{ .GoName }}
{{ end -}}
	default:
		return ErrInvalid{{ $e.GoName }}(str)
	}
	return nil
}

// Value satisfies the [driver.Valuer] interface.
func ({{ short $e.GoName }} {{ $e.GoName }}) Value() (driver.Value, error) {
	return {{ short $e.GoName }}.String(), nil
}

// Scan satisfies the [sql.Scanner] interface.
func ({{ short $e.GoName }} *{{ $e.GoName }}) Scan(v any) error {
	switch x := v.(type) {
	case []byte:
		return {{ short $e.GoName }}.UnmarshalText(x)
	case string:
		return {{ short $e.GoName }}.UnmarshalText([]byte(x))
	}
	return ErrInvalid{{ $e.GoName }}(fmt.Sprintf("%T", v))
}

{{ $nullName := (printf "%s%s" "Null" $e.GoName) -}}
{{- $nullShort := (short $nullName) -}}
// {{ $nullName }} represents a null '{{ $e.SQLName }}' enum for schema '{{ schema }}'.
type {{ $nullName }} struct {
	{{ $e.GoName }} {{ $e.GoName }}
	// Valid is true if [{{ $e.GoName }}] is not null.
	Valid bool
}

// Value satisfies the [driver.Valuer] interface.
func ({{ $nullShort }} {{ $nullName }}) Value() (driver.Value, error) {
	if !{{ $nullShort }}.Valid {
		return nil, nil
	}
	return {{ $nullShort }}.{{ $e.GoName }}.Value()
}

// Scan satisfies the [sql.Scanner] interface.
func ({{ $nullShort }} *{{ $nullName }}) Scan(v any) error {
	if v == nil {
		{{ $nullShort }}.{{ $e.GoName }}, {{ $nullShort }}.Valid = 0, false
		return nil
	}
	err := {{ $nullShort }}.{{ $e.GoName }}.Scan(v)
	{{ $nullShort }}.Valid = err == nil
	return err
}

// ErrInvalid{{ $e.GoName }} is the invalid [{{ $e.GoName }}] error.
type ErrInvalid{{ $e.GoName }} string

// Error satisfies the error interface.
func (err ErrInvalid{{ $e.GoName }}) Error() string {
        return fmt.Sprintf("invalid {{ $e.GoName }}(%s)", string(err))
}
{{ end }}

{{ define "composite" }}
{{- $c := .Data -}}
{{- if $c.Comment -}}
// {{ $c.Comment | eval $c.GoName }}
{{- else -}}
// {{ $c.GoName }} is the '{{ $c.SQLName }}' composite type from schema '{{ schema }}'.
{{- end }}
type {{ $c.GoName }} struct {
{{- range $i, $f := $c.Fields -}}
        {{- $comment := $f.SQLName -}}
        {{- if $f.Comment }}
                {{- $comment = $f.Comment -}}
        {{- end }}
        {{ $f.GoName }} {{ type $f.Type }} `json:"{{ $f.SQLName }}" row:"{{ inc $i }}"` // {{ $comment }}
{{- end }}
}

{{ $nullName := (printf "Null%s" $c.GoName) -}}
// {{ $nullName }} represents a null '{{ $c.SQLName }}' composite for schema '{{ schema }}'.
type {{ $nullName }} struct {
        {{ $c.GoName }} {{ $c.GoName }}
        // Valid is true if [{{ $c.GoName }}] is not null.
        Valid bool
}

// Value satisfies the [driver.Valuer] interface.
func ({{ short $nullName }} {{ $nullName }}) Value() (driver.Value, error) {
        if !{{ short $nullName }}.Valid {
                return nil, nil
        }
        return row_marshaler.Marshal({{ short $nullName }}.{{ $c.GoName }})
}

// Scan satisfies the [sql.Scanner] interface.
func ({{ short $nullName }} *{{ $nullName }}) Scan(v any) error {
        if v == nil {
                {{ short $nullName }}.{{ $c.GoName }}, {{ short $nullName }}.Valid = {{ $c.GoName }}{}, false
                return nil
        }
        if err := scan{{ $c.GoName }}Literal(v, &{{ short $nullName }}.{{ $c.GoName }}); err != nil {
                return err
        }
        {{ short $nullName }}.Valid = true
        return nil
}

// {{ $c.GoName }}Array represents a PostgreSQL array of '{{ $c.SQLName }}'.
type {{ $c.GoName }}Array []{{ $c.GoName }}

// Value satisfies the [driver.Valuer] interface.
func ({{ short $c.GoName }}a {{ $c.GoName }}Array) Value() (driver.Value, error) {
        if {{ short $c.GoName }}a == nil {
                return nil, nil
        }
        parts := make([]string, len({{ short $c.GoName }}a))
        for i, elem := range {{ short $c.GoName }}a {
                literal, err := row_marshaler.Marshal(elem)
                if err != nil {
                        return nil, err
                }
                // Escape for PostgreSQL array literal: \ -> \\, " -> \"
                escaped := strings.ReplaceAll(literal, "\\", "\\\\")
                escaped = strings.ReplaceAll(escaped, "\"", "\\\"")
                parts[i] = "\"" + escaped + "\""
        }
        return "{" + strings.Join(parts, ",") + "}", nil
}

// Scan satisfies the [sql.Scanner] interface.
func ({{ short $c.GoName }}a *{{ $c.GoName }}Array) Scan(v any) error {
        if v == nil {
                *{{ short $c.GoName }}a = nil
                return nil
        }
        var data string
        switch x := v.(type) {
        case []byte:
                data = string(x)
        case string:
                data = x
        default:
                return fmt.Errorf("cannot scan %T into {{ $c.GoName }}Array", v)
        }
        elements, err := parse{{ $c.GoName }}Array(data)
        if err != nil {
                return err
        }
        *{{ short $c.GoName }}a = elements
        return nil
}

// MarshalJSON marshals [{{ $c.GoName }}] to JSON.
func ({{ short $c.GoName }} {{ $c.GoName }}) MarshalJSON() ([]byte, error) {
        type alias {{ $c.GoName }}
        return json.Marshal(alias({{ short $c.GoName }}))
}

// UnmarshalJSON unmarshals [{{ $c.GoName }}] from JSON.
func ({{ short $c.GoName }} *{{ $c.GoName }}) UnmarshalJSON(data []byte) error {
        type alias {{ $c.GoName }}
        return json.Unmarshal(data, (*alias)({{ short $c.GoName }}))
}

// Value satisfies the [driver.Valuer] interface.
func ({{ short $c.GoName }} {{ $c.GoName }}) Value() (driver.Value, error) {
        return row_marshaler.Marshal({{ short $c.GoName }})
}

// Scan satisfies the [sql.Scanner] interface.
func ({{ short $c.GoName }} *{{ $c.GoName }}) Scan(v any) error {
        if v == nil {
                *{{ short $c.GoName }} = {{ $c.GoName }}{}
                return nil
        }
        return scan{{ $c.GoName }}Literal(v, {{ short $c.GoName }})
}

func scan{{ $c.GoName }}Literal(v any, dest *{{ $c.GoName }}) error {
        var record string
        switch x := v.(type) {
        case []byte:
                record = string(x)
        case string:
                record = x
        default:
                return fmt.Errorf("cannot scan %T into {{ $c.GoName }}", v)
        }

        if strings.HasPrefix(record, "{") {
                return json.Unmarshal([]byte(record), dest)
        }

        return row_marshaler.Unmarshal(record, dest)
}

// parse{{ $c.GoName }}Array converts a PostgreSQL array literal into its composite elements.
func parse{{ $c.GoName }}Array(arrayLiteral string) ([]{{ $c.GoName }}, error) {
        arrayLiteral = strings.TrimSpace(arrayLiteral)
        if arrayLiteral == "" || arrayLiteral == "{}" {
                return []{{ $c.GoName }}{}, nil
        }
        if len(arrayLiteral) < 2 || arrayLiteral[0] != '{' || arrayLiteral[len(arrayLiteral)-1] != '}' {
                return nil, fmt.Errorf("invalid array format: %s", arrayLiteral)
        }

        inner := arrayLiteral[1 : len(arrayLiteral)-1]
        if inner == "" {
                return []{{ $c.GoName }}{}, nil
        }

        rawElements := split{{ $c.GoName }}ArrayElements(inner)
        result := make([]{{ $c.GoName }}, 0, len(rawElements))
        for i, elem := range rawElements {
                elem = strings.TrimSpace(elem)
                if elem == "NULL" {
                        result = append(result, {{ $c.GoName }}{})
                        continue
                }
                if strings.HasPrefix(elem, "\"") && strings.HasSuffix(elem, "\"") {
                        elem = elem[1 : len(elem)-1]
                        // Unescape PostgreSQL array literal: \" -> ", \\ -> \
                        elem = strings.ReplaceAll(elem, "\\\"", "\"")
                        elem = strings.ReplaceAll(elem, "\\\\", "\\")
                }

                var item {{ $c.GoName }}
                if err := scan{{ $c.GoName }}Literal(elem, &item); err != nil {
                        return nil, fmt.Errorf("element %d: %w", i, err)
                }
                result = append(result, item)
        }

        return result, nil
}

// split{{ $c.GoName }}ArrayElements splits a PostgreSQL array inner literal into individual composite elements.
func split{{ $c.GoName }}ArrayElements(s string) []string {
        var result []string
        var current strings.Builder
        depth := 0
        inQuote := false
        escaped := false

        for i := 0; i < len(s); i++ {
                ch := s[i]
                if escaped {
                        current.WriteByte(ch)
                        escaped = false
                        continue
                }
                switch ch {
                case '\\':
                        current.WriteByte(ch)
                        escaped = true
                        continue
                case '"':
                        current.WriteByte(ch)
                        inQuote = !inQuote
                        continue
                case '(':
                        if !inQuote {
                                depth++
                        }
                case ')':
                        if !inQuote && depth > 0 {
                                depth--
                        }
                case ',':
                        if !inQuote && depth == 0 {
                                result = append(result, current.String())
                                current.Reset()
                                continue
                        }
                }
                current.WriteByte(ch)
        }

        if current.Len() > 0 {
                result = append(result, current.String())
        }

        return result
}
{{ end }}

{{ define "foreignkey" }}
{{- $k := .Data -}}
// {{ func_name_context $k }} returns the {{ $k.RefTable }} associated with the [{{ $k.Table.GoName }}]'s ({{ names "" $k.Fields }}).
//
// Generated from foreign key '{{ $k.SQLName }}'.
{{ recv_context $k.Table $k }} {
	return {{ foreign_key_context $k }}
}
{{- if context_both }}

// {{ func_name $k }} returns the {{ $k.RefTable }} associated with the {{ $k.Table }}'s ({{ names "" $k.Fields }}).
//
// Generated from foreign key '{{ $k.SQLName }}'.
{{ recv $k.Table $k }} {
	return {{ foreign_key $k }}
}
{{- end }}
{{ end }}

{{ define "index" }}
{{- $i := .Data -}}
// {{ func_name_context $i }} retrieves a row from '{{ schema $i.Table.SQLName }}' as a [{{ $i.Table.GoName }}].
//
// Generated from index '{{ $i.SQLName }}'.
{{ func_context $i }} {
	// query
	{{ sqlstr "index" $i }}
	// run
	logf(sqlstr, {{ params $i.Fields false }})
{{- if $i.IsUnique }}
	{{ short $i.Table }} := {{ $i.Table.GoName }}{
	{{- if $i.Table.PrimaryKeys }}
		_exists: true,
	{{ end -}}
	}
	if err := {{ db "QueryRow"  $i }}.Scan({{ names (print "&" (short $i.Table) ".") $i.Table }}); err != nil {
		return nil, logerror(err)
	}
	return &{{ short $i.Table }}, nil
{{- else }}
	rows, err := {{ db "Query" $i }}
	if err != nil {
		return nil, logerror(err)
	}
	defer rows.Close()
	// process
	var res []*{{ $i.Table.GoName }}
	for rows.Next() {
		{{ short $i.Table }} := {{ $i.Table.GoName }}{
		{{- if $i.Table.PrimaryKeys }}
			_exists: true,
		{{ end -}}
		}
		// scan
		if err := rows.Scan({{ names_ignore (print "&" (short $i.Table) ".")  $i.Table }}); err != nil {
			return nil, logerror(err)
		}
		res = append(res, &{{ short $i.Table }})
	}
	if err := rows.Err(); err != nil {
		return nil, logerror(err)
	}
	return res, nil
{{- end }}
}

{{ if context_both -}}
// {{ func_name $i }} retrieves a row from '{{ schema $i.Table.SQLName }}' as a [{{ $i.Table.GoName }}].
//
// Generated from index '{{ $i.SQLName }}'.
{{ func $i }} {
	return {{ func_name_context $i }}({{ names "" "context.Background()" "db" $i }})
}
{{- end }}

{{end}}

{{ define "procs" }}
{{- $ps := .Data -}}
{{- range $p := $ps -}}
// {{ func_name_context $p }} calls the stored {{ $p.Type }} '{{ $p.Signature }}' on db.
{{ func_context $p }} {
{{- if and (driver "mysql") (eq $p.Type "procedure") (not $p.Void) }}
	// At the moment, the Go MySQL driver does not support stored procedures
	// with out parameters
	return {{ zero $p.Returns }}, fmt.Errorf("unsupported")
{{- else }}
	// call {{ schema $p.SQLName }}
	{{ sqlstr "proc" $p }}
	// run
{{- if not $p.Void }}
{{- range $p.Returns }}
	var {{ check_name .GoName }} {{ type .Type }}
{{- end }}
	logf(sqlstr, {{ params $p.Params false }})
{{- if and (driver "sqlserver" "oracle") (eq $p.Type "procedure")}}
	if _, err := {{ db_named "Exec" $p }}; err != nil {
{{- else }}
	if err := {{ db "QueryRow" $p }}.Scan({{ names "&" $p.Returns }}); err != nil {
{{- end }}
		return {{ zero $p.Returns }}, logerror(err)
	}
	return {{ range $p.Returns }}{{ check_name .GoName }}, {{ end }}nil
{{- else }}
	logf(sqlstr)
{{- if driver "sqlserver" "oracle" }}
	if _, err := {{ db_named "Exec" $p }}; err != nil {
{{- else }}
	if _, err := {{ db "Exec" $p }}; err != nil {
{{- end }}
		return logerror(err)
	}
	return nil
{{- end }}
{{- end }}
}

{{ if context_both -}}
// {{ func_name $p }} calls the {{ $p.Type }} '{{ $p.Signature }}' on db.
{{ func $p }} {
	return {{ func_name_context $p }}({{ names_all "" "context.Background()" "db" $p.Params }})
}
{{- end -}}
{{- end }}
{{ end }}

{{ define "typedef" }}
{{- $t := .Data -}}
{{- if $t.Comment -}}
// {{ $t.Comment | eval $t.GoName }}
{{- else -}}
// {{ $t.GoName }} represents a row from '{{ schema $t.SQLName }}'.
{{- end }}
type {{ $t.GoName }} struct {
{{ range $t.Fields -}}
	{{ field . }}
{{ end }}
{{- if $t.PrimaryKeys -}}
	// xo fields
	_exists, _deleted bool
{{ end -}}
}

{{ if $t.PrimaryKeys -}}
// Exists returns true when the [{{ $t.GoName }}] exists in the database.
func ({{ short $t }} *{{ $t.GoName }}) Exists() bool {
	return {{ short $t }}._exists
}

// Deleted returns true when the [{{ $t.GoName }}] has been marked for deletion
// from the database.
func ({{ short $t }} *{{ $t.GoName }}) Deleted() bool {
	return {{ short $t }}._deleted
}

// {{ func_name_context "Insert" }} inserts the [{{ $t.GoName }}] to the database.
{{ recv_context $t "Insert" }} {
	switch {
	case {{ short $t }}._exists: // already exists
		return logerror(&ErrInsertFailed{ErrAlreadyExists})
	case {{ short $t }}._deleted: // deleted
		return logerror(&ErrInsertFailed{ErrMarkedForDeletion})
	}
{{ if $t.Manual -}}
	// insert (manual)
	{{ sqlstr "insert_manual" $t }}
	// run
	{{ logf $t }}
	if _, err := {{ db_prefix "Exec" false $t }}; err != nil {
		return logerror(err)
	}
{{- else -}}
	// insert (primary key generated and returned by database)
	{{ sqlstr "insert" $t }}
	// run
	{{ logf $t $t.PrimaryKeys }}
{{ if (driver "postgres") -}}
	if err := {{ db_prefix "QueryRow" true $t }}.Scan(&{{ short $t }}.{{ (index $t.PrimaryKeys 0).GoName }}); err != nil {
		return logerror(err)
	}
{{- else if (driver "sqlserver") -}}
	rows, err := {{ db_prefix "Query" true $t }}
	if err != nil {
		return logerror(err)
	}
	defer rows.Close()
	// retrieve id
	var id int64
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return logerror(err)
		}
	}
	if err := rows.Err(); err != nil {
		return logerror(err)
	}
{{- else if (driver "oracle") -}}
	var id int64
	if _, err := {{ db_prefix "Exec" true $t (named "pk" "&id" true) }}; err != nil {
		return logerror(err)
	}
{{- else -}}
	res, err := {{ db_prefix "Exec" true $t }}
	if err != nil {
		return logerror(err)
	}
	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return logerror(err)
	}
{{- end -}}
{{ if not (driver "postgres") -}}
	// set primary key
	{{ short $t }}.{{ (index $t.PrimaryKeys 0).GoName }} = {{ (index $t.PrimaryKeys 0).Type }}(id)
{{- end }}
{{- end }}
	// set exists
	{{ short $t }}._exists = true
	return nil
}

{{ if context_both -}}
// Insert inserts the [{{ $t.GoName }}] to the database.
{{ recv $t "Insert" }} {
	return {{ short $t }}.InsertContext(context.Background(), db)
}
{{- end }}


{{ if eq (len $t.Fields) (len $t.PrimaryKeys) -}}
// ------ NOTE: Update statements omitted due to lack of fields other than primary key ------
{{- else -}}
// {{ func_name_context "Update" }} updates a [{{ $t.GoName }}] in the database.
{{ recv_context $t "Update" }} {
	switch {
	case !{{ short $t }}._exists: // doesn't exist
		return logerror(&ErrUpdateFailed{ErrDoesNotExist})
	case {{ short $t }}._deleted: // deleted
		return logerror(&ErrUpdateFailed{ErrMarkedForDeletion})
	}
	// update with {{ if driver "postgres" }}composite {{ end }}primary key
	{{ sqlstr "update" $t }}
	// run
	{{ logf_update $t }}
	if _, err := {{ db_update "Exec" $t }}; err != nil {
		return logerror(err)
	}
	return nil
}

{{ if context_both -}}
// Update updates a [{{ $t.GoName }}] in the database.
{{ recv $t "Update" }} {
	return {{ short $t }}.UpdateContext(context.Background(), db)
}
{{- end }}

// {{ func_name_context "Save" }} saves the [{{ $t.GoName }}] to the database.
{{ recv_context $t "Save" }} {
	if {{ short $t }}.Exists() {
		return {{ short $t }}.{{ func_name_context "Update" }}({{ if context }}ctx, {{ end }}db)
	}
	return {{ short $t }}.{{ func_name_context "Insert" }}({{ if context }}ctx, {{ end }}db)
}

{{ if context_both -}}
// Save saves the [{{ $t.GoName }}] to the database.
{{ recv $t "Save" }} {
	if {{ short $t }}._exists {
		return {{ short $t }}.UpdateContext(context.Background(), db)
	}
	return {{ short $t }}.InsertContext(context.Background(), db)
}
{{- end }}

// {{ func_name_context "Upsert" }} performs an upsert for [{{ $t.GoName }}].
{{ recv_context $t "Upsert" }} {
	switch {
	case {{ short $t }}._deleted: // deleted
		return logerror(&ErrUpsertFailed{ErrMarkedForDeletion})
	}
	// upsert
	{{ sqlstr "upsert" $t }}
	// run
	{{ logf $t }}
	if _, err := {{ db_prefix "Exec" false $t }}; err != nil {
		return logerror(err)
	}
	// set exists
	{{ short $t }}._exists = true
	return nil
}

{{ if context_both -}}
// Upsert performs an upsert for [{{ $t.GoName }}].
{{ recv $t "Upsert" }} {
	return {{ short $t }}.UpsertContext(context.Background(), db)
}
{{- end -}}
{{- end }}

// {{ func_name_context "Delete" }} deletes the [{{ $t.GoName }}] from the database.
{{ recv_context $t "Delete" }} {
	switch {
	case !{{ short $t }}._exists: // doesn't exist
		return nil
	case {{ short $t }}._deleted: // deleted
		return nil
	}
{{ if eq (len $t.PrimaryKeys) 1 -}}
	// delete with single primary key
	{{ sqlstr "delete" $t }}
	// run
	{{ logf_pkeys $t }}
	if _, err := {{ db "Exec" (print (short $t) "." (index $t.PrimaryKeys 0).GoName) }}; err != nil {
		return logerror(err)
	}
{{- else -}}
	// delete with composite primary key
	{{ sqlstr "delete" $t }}
	// run
	{{ logf_pkeys $t }}
	if _, err := {{ db "Exec" (names (print (short $t) ".") $t.PrimaryKeys) }}; err != nil {
		return logerror(err)
	}
{{- end }}
	// set deleted
	{{ short $t }}._deleted = true
	return nil
}

{{ if context_both -}}
// Delete deletes the [{{ $t.GoName }}] from the database.
{{ recv $t "Delete" }} {
	return {{ short $t }}.DeleteContext(context.Background(), db)
}
{{- end -}}
{{- end }}
{{ end }}
