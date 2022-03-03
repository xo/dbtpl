//go:build xotpl

package createdb

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"text/template"
	"unicode"

	xo "github.com/xo/xo/types"
)

// Init registers the template.
func Init(ctx context.Context, f func(xo.TemplateType)) error {
	fmtCmd, _ := exec.LookPath("sql-formatter")
	var fmtOpts []string
	if fmtCmd != "" {
		fmtOpts = []string{"-u", "-l={{ . }}", "-i=2", "--lines-between-queries=2"}
	}
	f(xo.TemplateType{
		Modes: []string{"schema"},
		Flags: []xo.Flag{
			{
				ContextKey: FmtKey,
				Type:       "string",
				Desc:       "fmt command",
				Default:    fmtCmd,
			},
			{
				ContextKey: FmtOptsKey,
				Type:       "[]string",
				Desc:       "fmt options",
				Default:    strings.Join(fmtOpts, ","),
			},
			{
				ContextKey: ConstraintKey,
				Type:       "bool",
				Desc:       "enable constraint name in output",
				Default:    "false",
			},
			{
				ContextKey: EscKey,
				Type:       "string",
				Desc:       "escape mode",
				Default:    "none",
				Enums:      []string{"none", "types", "all"},
			},
			{
				ContextKey: EngineKey,
				Type:       "string",
				Desc:       "mysql table engine",
				Default:    "InnoDB",
			},
			{
				ContextKey: TrimCommentKey,
				Type:       "bool",
				Desc:       "trim leading comment from views and procs",
				Default:    "true",
			},
		},
		Funcs: NewFuncs,
		Process: func(ctx context.Context, _ string, set *xo.Set, emit func(xo.Template)) error {
			if len(set.Schemas) == 0 {
				return errors.New("createdb template must be passed at least one schema")
			}
			for _, schema := range set.Schemas {
				schema.Tables = sortTables(schema.Tables)
				emit(xo.Template{
					Partial:  "createdb",
					Dest:     "xo.xo.sql",
					SortName: schema.Name,
					Data:     schema,
				})
			}
			return nil
		},
		Post: func(ctx context.Context, mode string, files map[string][]byte, emit func(string, []byte)) error {
			// build options
			fmtPath, lang, opts := Fmt(ctx), Lang(ctx), FmtOpts(ctx)
			for i, o := range opts {
				tpl, err := template.New(fmt.Sprintf("option %d", i)).Parse(o)
				if err != nil {
					return err
				}
				b := new(bytes.Buffer)
				if err := tpl.Execute(b, lang); err != nil {
					return err
				}
				opts[i] = b.String()
			}
			// post-process files
			for file, content := range files {
				// skip
				if fmtPath == "" {
					emit(file, cleanEnd(cleanRE.ReplaceAll(content, []byte("$1\n\n--"))))
					continue
				}
				// execute
				stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
				cmd := exec.Command(fmtPath, opts...)
				cmd.Stdin, cmd.Stdout, cmd.Stderr = bytes.NewReader(content), stdout, stderr
				if err := cmd.Run(); err != nil {
					return fmt.Errorf("unable to execute %s: %v: %s", fmtPath, err, stderr.String())
				}
				emit(file, cleanEnd(stdout.Bytes()))
			}
			return nil
		},
	})
	return nil
}

// cleanRE matches empty lines.
var cleanRE = regexp.MustCompile(`([\.;])\n{2,}--`)

// cleanEnd trims the end of any spaces, ensuring it ends with exactly one
// newline.
func cleanEnd(buf []byte) []byte {
	return append(bytes.TrimRightFunc(buf, unicode.IsSpace), '\n')
}

// sortTables sorts tables.
func sortTables(tables []xo.Table) []xo.Table {
	m := make(map[string]xo.Table)
	for _, table := range tables {
		m[table.Name] = table
	}
	seen := make(map[string]bool)
	var sorted []xo.Table
	for _, table := range tables {
		sorted = sortAppendTable(m, seen, sorted, table)
	}
	return sorted
}

// sortAppendTable appends and returns the list of foreign key dependencies for
// the table if not already in seen.
func sortAppendTable(m map[string]xo.Table, seen map[string]bool, sorted []xo.Table, table xo.Table) []xo.Table {
	if seen[table.Name] {
		return sorted
	}
	for _, fk := range table.ForeignKeys {
		if t := m[fk.RefTable]; table.Name != t.Name && !seen[t.Name] {
			sorted = sortAppendTable(m, seen, sorted, t)
		}
	}
	seen[table.Name] = true
	return append(sorted, table)
}

// Funcs is a set of template funcs.
type Funcs struct {
	driver      string
	constraint  bool
	escCols     bool
	escTypes    bool
	engine      string
	trimComment bool
}

// NewFuncs creates custom template funcs for the context.
func NewFuncs(ctx context.Context, _ string) (template.FuncMap, error) {
	driver, _, _ := xo.DriverDbSchema(ctx)
	funcs := &Funcs{
		driver:      driver,
		constraint:  Constraint(ctx),
		escCols:     Esc(ctx, "columns"),
		escTypes:    Esc(ctx, "types"),
		engine:      Engine(ctx),
		trimComment: TrimComment(ctx),
	}
	return template.FuncMap{
		"coldef":          funcs.coldef,
		"viewdef":         funcs.viewdef,
		"procdef":         funcs.procdef,
		"driver":          funcs.driverfn,
		"constraint":      funcs.constraintfn,
		"esc":             funcs.escType,
		"fields":          funcs.fields,
		"engine":          funcs.enginefn,
		"literal":         funcs.literal,
		"isEndConstraint": funcs.isEndConstraint,
		"comma":           comma,
	}, nil
}

// coldef generates a column definition.
func (f *Funcs) coldef(table xo.Table, field xo.Field) string {
	// normalize type
	typ := f.normalize(field.Type)
	// add sequence definition
	if field.IsSequence {
		typ = f.resolveSequence(typ, field)
	}
	// column def
	def := []string{f.escCol(field.Name), typ}
	// add default value
	if field.Default != "" && !field.IsSequence {
		def = append(def, "DEFAULT", f.alterDefault(field.Default))
	}
	if !field.Type.Nullable && !field.IsSequence {
		def = append(def, "NOT NULL")
	}
	// add constraints
	if fk := f.colFKey(table, field); fk != "" {
		def = append(def, fk)
	}
	return strings.Join(def, " ")
}

// alterDefault parses and alters default column values based on the driver.
func (f *Funcs) alterDefault(s string) string {
	switch f.driver {
	case "postgres":
		if m := postgresDefaultCastRE.FindStringSubmatch(s); m != nil {
			return m[1]
		}
	case "mysql":
		if v := strings.ToUpper(s); v == "CURRENT_TIMESTAMP()" {
			return "CURRENT_TIMESTAMP"
		}
	case "sqlite3":
		if !sqliteDefaultNeedsParenRE.MatchString(s) {
			return "(" + s + ")"
		}
	}
	return s
}

// postgresDefaultCastRE is the regexp to strip the datatype cast from the
// postgres default value.
var postgresDefaultCastRE = regexp.MustCompile(`(.*)::[a-zA-Z_ ]*(\[\])?$`)

// sqliteDefaultNeedsParen is the regexp to test whether the given value is
// correctly surrounded with parenthesis
//
// If it starts and ends with a parenthesis or a single or double quote, it
// does not need to be quoted with parenthesis.
var sqliteDefaultNeedsParenRE = regexp.MustCompile(`^([\('"].*[\)'"]|\d+)$`)

// resolveSequence resolves a sequence name.
func (f *Funcs) resolveSequence(typ string, field xo.Field) string {
	switch f.driver {
	case "postgres":
		switch typ {
		case "SMALLINT":
			return "SMALLSERIAL"
		case "INTEGER":
			return "SERIAL"
		case "BIGINT":
			return "BIGSERIAL"
		}
	case "mysql":
		return typ + " AUTO_INCREMENT"
	case "sqlite3":
		ext := " PRIMARY KEY AUTOINCREMENT"
		if !field.Type.Nullable {
			ext = " NOT NULL" + ext
		}
		return typ + ext
	case "sqlserver":
		return typ + " IDENTITY(1, 1)"
	case "oracle":
		return typ + " GENERATED ALWAYS AS IDENTITY"
	}
	return ""
}

// colFKey
func (f *Funcs) colFKey(table xo.Table, field xo.Field) string {
	for _, fk := range table.ForeignKeys {
		if len(fk.Fields) == 1 && fk.Fields[0] == field {
			tblName, fieldName := f.escType(fk.RefTable), fk.RefFields[0].Name
			return fmt.Sprintf("%sREFERENCES %s (%s)", f.constraintfn(fk.Name), tblName, fieldName)
		}
	}
	return ""
}

// viewdef generates a view definition.
func (f *Funcs) viewdef(view xo.Table) string {
	def := view.Definition
	switch f.driver {
	case "postgres", "mysql", "oracle":
		def = fmt.Sprintf("CREATE VIEW %s AS\n%s", f.escType(view.Name), view.Definition)
	}
	if f.trimComment {
		if strings.HasPrefix(def, "--") {
			def = def[strings.Index(def, "\n")+1:]
		}
	}
	return strings.TrimSuffix(def, ";")
}

// procdef generates a proc definition.
func (f *Funcs) procdef(proc xo.Proc) string {
	def := f.cleanProcDef(proc.Definition)
	// prepend signature definition
	if f.driver == "postgres" || f.driver == "mysql" {
		def = f.procSignature(proc) + "\n" + def
	}
	return def
}

// celanProcDef cleans a proc definition.
func (f *Funcs) cleanProcDef(def string) string {
	switch f.driver {
	// nothing needs to be done for postgres
	// only add the query language suffix
	case "postgres":
		return def + "\n$$ LANGUAGE plpgsql"
	// the trailing semicolon shouldn't be escaped for sqlserver
	case "sqlserver":
		def = strings.TrimSuffix(def, ";")
	// oracle only just needs the CREATE prefix
	case "oracle":
		def = "CREATE " + def
	}
	if f.trimComment {
		if strings.HasPrefix(def, "--") {
			def = def[strings.Index(def, "\n")+1:]
		}
	}
	return strings.ReplaceAll(def, ";", "\\;")
}

// procSignature generates a proc signature.
func (f *Funcs) procSignature(proc xo.Proc) string {
	// create function signature
	typ := "PROCEDURE"
	if proc.Type == "function" {
		typ = "FUNCTION"
	}
	var params []string
	var end string
	// add params
	for _, field := range proc.Params {
		params = append(params, fmt.Sprintf("%s %s", f.escCol(field.Name), f.normalize(field.Type)))
	}
	// add return values
	if len(proc.Returns) == 1 && proc.Returns[0].Name == "r0" {
		end += " RETURNS " + f.normalize(proc.Returns[0].Type)
	} else {
		for _, field := range proc.Returns {
			params = append(params, fmt.Sprintf("OUT %s %s", f.escCol(field.Name), f.normalize(field.Type)))
		}
	}
	signature := fmt.Sprintf("CREATE %s %s(%s)%s", typ, f.escType(proc.Name), strings.Join(params, ", "), end)
	if f.driver == "postgres" {
		signature += " AS $$"
	}
	return signature
}

// driverfn determines if a driver is allowed.
func (f *Funcs) driverfn(allowed ...string) bool {
	for _, d := range allowed {
		if f.driver == d {
			return true
		}
	}
	return false
}

// contstraintfn
func (f *Funcs) constraintfn(name string) string {
	if f.constraint || f.driver == "sqlserver" || f.driver == "oracle" {
		return fmt.Sprintf("CONSTRAINT %s ", f.escType(name))
	}
	return ""
}

// fields
func (f *Funcs) fields(v interface{}) string {
	switch x := v.(type) {
	case []xo.Field:
		var fs []string
		for _, field := range x {
			fs = append(fs, f.escCol(field.Name))
		}
		return strings.Join(fs, ", ")
	}
	return fmt.Sprintf("[[ UNKNOWN TYPE %T ]]", v)
}

// esc escapes s.
func (f *Funcs) esc(s string, esc bool) string {
	if !esc {
		return s
	}
	var start, end string
	switch f.driver {
	case "postgres", "sqlite3", "oracle":
		start, end = `"`, `"`
	case "mysql":
		start, end = "`", "`"
	case "sqlserver":
		start, end = "[", "]"
	}
	return start + s + end
}

// escType
func (f *Funcs) escType(value string) string {
	return f.esc(value, f.escTypes)
}

// escCol
func (f *Funcs) escCol(value string) string {
	return f.esc(value, f.escCols)
}

// enginefn returns the engine for the database (mysql).
func (f *Funcs) enginefn() string {
	if f.driver != "mysql" || f.engine == "" {
		return ""
	}
	return fmt.Sprintf(" ENGINE=%s", f.engine)
}

// normalize normalizes a datatype.
func (f *Funcs) normalize(datatype xo.Type) string {
	typ := f.convert(datatype)
	if datatype.Scale > 0 && !omitPrecision[f.driver][typ] {
		typ += fmt.Sprintf("(%d, %d)", datatype.Prec, datatype.Scale)
	} else if datatype.Prec > 0 && !omitPrecision[f.driver][typ] {
		typ += fmt.Sprintf("(%d)", datatype.Prec)
	}
	if datatype.Unsigned {
		typ += " UNSIGNED"
	}
	if datatype.IsArray {
		typ += "[]"
	}
	return typ
}

// convert converts
func (f *Funcs) convert(datatype xo.Type) string {
	// mysql enums
	if f.driver == "mysql" && datatype.Enum != nil {
		var enums []string
		for _, v := range datatype.Enum.Values {
			enums = append(enums, fmt.Sprintf("'%s'", v.Name))
		}
		return fmt.Sprintf("ENUM(%s)", strings.Join(enums, ", "))
	}
	// check aliases
	typ := datatype.Type
	if alias, ok := typeAliases[f.driver][typ]; ok {
		typ = alias
	}
	return strings.ToUpper(typ)
}

// literal properly escapes string literals within single quotes
// (Used for enum values in postgres)
func (f *Funcs) literal(literal string) string {
	return fmt.Sprint("'", strings.ReplaceAll(literal, "'", "''"), "'")
}

func (f *Funcs) isEndConstraint(idx xo.Index) bool {
	if f.driver == "sqlite3" && idx.Fields[0].IsSequence {
		return false
	}
	return idx.IsPrimary || idx.IsUnique
}

var typeAliases = map[string]map[string]string{
	"postgres": {
		"character varying":           "varchar",
		"character":                   "char",
		"time without time zone":      "time",
		"timestamp without time zone": "timestamp",
		"time with time zone":         "timetz",
		"timestamp with time zone":    "timestamptz",
	},
}

var omitPrecision = map[string]map[string]bool{
	"sqlserver": {
		"TINYINT":        true,
		"SMALLINT":       true,
		"INT":            true,
		"BIGINT":         true,
		"REAL":           true,
		"SMALLMONEY":     true,
		"MONEY":          true,
		"BIT":            true,
		"DATE":           true,
		"TIME":           true,
		"DATETIME":       true,
		"DATETIME2":      true,
		"SMALLDATETIME":  true,
		"DATETIMEOFFSET": true,
	},
	"oracle": {
		"TIMESTAMP":                      true,
		"TIMESTAMP WITH TIME ZONE":       true,
		"TIMESTAMP WITH LOCAL TIME ZONE": true,
	},
}

func comma(i int, v interface{}) string {
	var l int
	switch x := v.(type) {
	case []xo.Field:
		l = len(x)
	}
	if i+1 < l {
		return ","
	}
	return ""
}

// Context keys.
const (
	FmtKey         xo.ContextKey = "fmt"
	FmtOptsKey     xo.ContextKey = "fmt-opts"
	ConstraintKey  xo.ContextKey = "constraint"
	EscKey         xo.ContextKey = "escape"
	EngineKey      xo.ContextKey = "engine"
	TrimCommentKey xo.ContextKey = "trim-comment"
)

// Fmt returns fmt from the context.
func Fmt(ctx context.Context) string {
	s, _ := ctx.Value(FmtKey).(string)
	return s
}

// FmtOpts returns fmt-opts from the context.
func FmtOpts(ctx context.Context) []string {
	v, _ := ctx.Value(FmtOptsKey).([]string)
	return v
}

// Constraint returns constraint from the context.
func Constraint(ctx context.Context) bool {
	b, _ := ctx.Value(ConstraintKey).(bool)
	return b
}

// Esc returns esc from the context.
func Esc(ctx context.Context, esc string) bool {
	v, _ := ctx.Value(EscKey).(string)
	return v == "all" || v == esc
}

// Engine returns engine from the context.
func Engine(ctx context.Context) string {
	s, _ := ctx.Value(EngineKey).(string)
	return s
}

// TrimComment returns trim-comments from the context.
func TrimComment(ctx context.Context) bool {
	b, _ := ctx.Value(TrimCommentKey).(bool)
	return b
}

// Lang returns the sql-formatter language to use from the context based on the
// context driver.
func Lang(ctx context.Context) string {
	driver, _, _ := xo.DriverDbSchema(ctx)
	switch driver {
	case "postgres", "sqlite3":
		return "postgresql"
	case "mysql":
		return "mysql"
	case "sqlserver":
		return "tsql"
	case "oracle":
		return "plsql"
	}
	return "sql"
}
