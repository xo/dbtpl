package createdbtpl

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"text/template"
	"unicode"

	"github.com/xo/xo/templates"
	xo "github.com/xo/xo/types"
)

func init() {
	formatterPath, _ := exec.LookPath("sql-formatter")
	var formatterOptions []string
	if formatterPath != "" {
		formatterOptions = []string{"-u", "-l={{ . }}", "-i=2", "--lines-between-queries=2"}
	}
	var funcs *Funcs
	templates.Register("createdb", &templates.TemplateSet{
		Files:   Files,
		For:     []string{"schema"},
		FileExt: ".xo.sql",
		Flags: []xo.Flag{
			{
				ContextKey:  FmtKey,
				Desc:        fmt.Sprintf("fmt command (default: %s)", formatterPath),
				Default:     formatterPath,
				PlaceHolder: "<path>",
				Value:       "",
			},
			{
				ContextKey:  FmtOptsKey,
				Desc:        fmt.Sprintf("fmt options (default: %s)", strings.Join(formatterOptions, ", ")),
				Default:     strings.Join(formatterOptions, ","),
				PlaceHolder: "<opts>",
				Value:       []string{},
			},
			{
				ContextKey: ConstraintKey,
				Desc:       "enable constraint name in output (postgres, mysql, sqlite3)",
				Default:    "false",
				Value:      false,
			},
			{
				ContextKey:  EscKey,
				Desc:        "escape mode (none, types, all; default: none)",
				PlaceHolder: "none",
				Default:     "none",
				Value:       "",
				Enums:       []string{"none", "types", "all"},
			},
			{
				ContextKey:  EngineKey,
				Desc:        "mysql table engine (default: InnoDB)",
				Default:     "InnoDB",
				PlaceHolder: `""`,
				Value:       "",
			},
			{
				ContextKey:  TrimCommentKey,
				Desc:        "trim leading comment from views and procs (--no-createdb-trim-comment)",
				Default:     "true",
				PlaceHolder: ``,
				Value:       false,
			},
		},
		Funcs: func(ctx context.Context) (template.FuncMap, error) {
			return funcs.FuncMap(), nil
		},
		FileName: func(ctx context.Context, tpl *templates.Template) string {
			return tpl.Name
		},
		Process: func(ctx context.Context, _ bool, set *templates.TemplateSet, v *xo.XO) error {
			if len(v.Schemas) == 0 {
				return errors.New("createdb template must be passed at least one schema")
			}
			for _, schema := range v.Schemas {
				schema.Tables = sortTables(schema.Tables)
				funcs = NewFuncs(ctx, schema.Enums)
				if err := set.Emit(ctx, &templates.Template{
					Name:     "xo",
					Template: "xo",
					Data:     schema,
				}); err != nil {
					return err
				}
			}
			return nil
		},
		Post: func(ctx context.Context, buf []byte) ([]byte, error) {
			formatterPath, lang := Fmt(ctx), Lang(ctx)
			if formatterPath == "" {
				return cleanEnd(cleanRE.ReplaceAll(buf, []byte("$1\n\n--"))), nil
			}
			// build options
			opts := FmtOpts(ctx)
			for i, o := range opts {
				tpl, err := template.New(fmt.Sprintf("option %d", i)).Parse(o)
				if err != nil {
					return nil, err
				}
				b := new(bytes.Buffer)
				if err := tpl.Execute(b, lang); err != nil {
					return nil, err
				}
				opts[i] = b.String()
			}
			// execute
			stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
			cmd := exec.Command(formatterPath, opts...)
			cmd.Stdin, cmd.Stdout, cmd.Stderr = bytes.NewReader(buf), stdout, stderr
			if err := cmd.Run(); err != nil {
				return nil, fmt.Errorf("unable to execute %s: %v: %s", formatterPath, err, stderr.String())
			}
			return cleanEnd(stdout.Bytes()), nil
		},
		Order: []string{"xo"},
	})
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
	tableMap := make(map[string]xo.Table)
	for _, table := range tables {
		tableMap[table.Name] = table
	}
	seen := make(map[string]bool)
	var sorted []xo.Table
	for _, table := range tables {
		sorted = sortAppendTable(tableMap, seen, sorted, table)
	}
	return sorted
}

// sortAppendTable appends and returns the list of foreign key dependencies for
// the table if not already in seen.
func sortAppendTable(tableMap map[string]xo.Table, seen map[string]bool, sorted []xo.Table, table xo.Table) []xo.Table {
	if seen[table.Name] {
		return sorted
	}
	for _, fk := range table.ForeignKeys {
		if t := tableMap[fk.RefTable]; table.Name != t.Name && !seen[t.Name] {
			sorted = sortAppendTable(tableMap, seen, sorted, t)
		}
	}
	seen[table.Name] = true
	return append(sorted, table)
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
	driver, _, _ := xo.DriverSchemaNthParam(ctx)
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

// Files are the embedded SQL templates.
//
//go:embed *.tpl
var Files embed.FS
