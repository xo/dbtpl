// Package cmd provides dbtpl command-line application logic.
package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	iofs "io/fs"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/kenshaw/snaker"
	"github.com/xo/dbtpl/loader"
	"github.com/xo/dbtpl/models"
	"github.com/xo/dbtpl/templates"
	xo "github.com/xo/dbtpl/types"
	"github.com/xo/dburl"
	"github.com/xo/dburl/passfile"
	"github.com/xo/ox"
	"github.com/yookoala/realpath"
)

// Args contains command-line arguments.
type Args struct {
	// Verbose enables verbose output.
	Verbose bool
	// LoaderParams are database loader parameters.
	LoaderParams LoaderParams
	// TemplateParams are template parameters.
	TemplateParams TemplateParams
	// QueryParams are query parameters.
	QueryParams QueryParams
	// SchemaParams are schema parameters.
	SchemaParams SchemaParams
	// OutParams are out parameters.
	OutParams OutParams
}

// NewArgs creates args for the provided template names.
func NewArgs(name string, names ...string) *Args {
	// default args
	return &Args{
		LoaderParams: LoaderParams{
			Flags: make(map[xo.ContextKey]*xo.Value),
		},
		TemplateParams: TemplateParams{
			Type:  xo.NewValue("string", name, "template type", names...),
			Flags: make(map[xo.ContextKey]*xo.Value),
		},
		SchemaParams: SchemaParams{
			FkMode:  xo.NewValue("string", "smart", "foreign key resolution mode", "smart", "parent", "field", "key"),
			Include: xo.NewValue("glob", "", "include types"),
			Exclude: xo.NewValue("glob", "", "exclude types"),
		},
	}
}

// LoaderParams are loader parameters.
type LoaderParams struct {
	// Schema is the name of the database schema.
	Schema string
	// Flags are additional loader flags.
	Flags map[xo.ContextKey]*xo.Value
}

// TemplateParams are template parameters.
type TemplateParams struct {
	// Type is the name of the template.
	Type *xo.Value
	// Src is the src directory of the template.
	Src string
	// Flags are additional template flags.
	Flags map[xo.ContextKey]*xo.Value
}

// QueryParams are query parameters.
type QueryParams struct {
	// Query is the query to introspect.
	Query string
	// Type is the type name.
	Type string
	// TypeComment is the type comment.
	TypeComment string
	// Func is the func name.
	Func string
	// FuncComment is the func comment.
	FuncComment string
	// Trim enables triming whitespace.
	Trim bool
	// Strip enables stripping the '::<type> AS <name>' in queries.
	Strip bool
	// One toggles the generated code to expect only one result.
	One bool
	// Flat toggles the generated code to return all scanned values directly.
	Flat bool
	// Exec toggles the generated code to do a db exec.
	Exec bool
	// Interpolate enables interpolation.
	Interpolate bool
	// Delimiter is the delimiter for parameterized values.
	Delimiter string
	// Fields are the fields to scan the result to.
	Fields string
	// AllowNulls enables results to have null types.
	AllowNulls bool
}

// SchemaParams are schema parameters.
type SchemaParams struct {
	// FkMode is the foreign resolution mode.
	FkMode *xo.Value
	// Include allows the user to specify which types should be included. Can
	// match multiple types via regex patterns.
	//
	// - When unspecified, all types are included.
	// - When specified, only types match will be included.
	// - When a type matches an exclude entry and an include entry,
	//   the exclude entry will take precedence.
	Include *xo.Value
	// Exclude allows the user to specify which types should be skipped. Can
	// match multiple types via regex patterns.
	//
	// When unspecified, all types are included in the schema.
	Exclude *xo.Value
	// UseIndexNames toggles using index names.
	//
	// This is not enabled by default, because index names are often generated
	// using database design software which often gives non-descriptive names
	// to indexes (for example, 'authors__b124214__u_idx' instead of the more
	// descriptive 'authors_title_idx').
	UseIndexNames bool
}

// OutParams are out parameters.
type OutParams struct {
	// Out is the out path.
	Out string
	// Single when true changes behavior so that output is to one file.
	Single string
	// Debug toggles direct writing of files to disk, skipping post processing.
	Debug bool
}

// Run runs the code generation.
func Run(ctx context.Context, name string) {
	dir := parseArg("--src", "-d", os.Args)
	template := parseArg("--template", "-t", os.Args)
	ts, err := NewTemplateSet(ctx, dir, template)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	// create args
	args := NewArgs(ts.Target(), ts.Targets()...)
	// create root command
	opts := RootCommand(ctx, name, ts, args)
	ox.RunContext(ctx, opts...)
}

// NewTemplateSet creates a new templates set.
func NewTemplateSet(ctx context.Context, dir, template string) (*templates.Set, error) {
	// build template ts
	ts := templates.NewDefaultTemplateSet(ctx)
	switch {
	case dir == "" && template == "":
		// show all default templates
		if err := ts.LoadDefaults(ctx); err != nil {
			return nil, err
		}
	case template != "":
		// only load the selected default template
		if err := ts.LoadDefault(ctx, template); err != nil {
			return nil, err
		}
		ts.Use(template)
	default:
		// load specified template
		s := snaker.SnakeToCamel(filepath.Base(dir))
		s = strings.ReplaceAll(strings.ToLower(s), "_", "-")
		// add template
		var err error
		if s, err = ts.Add(ctx, s, os.DirFS(dir), false); err != nil {
			return nil, err
		}
		// use
		ts.Use(s)
	}
	return ts, nil
}

// RootCommand creates the root command.
func RootCommand(ctx context.Context, name string, ts *templates.Set, args *Args, cmdargs ...string) []ox.Option {
	// root
	opts := []ox.Option{
		ox.Usage(name, "the templated code generator for databases."),
		ox.Defaults(),
	}
	// add sub commands
	for _, f := range []func(*templates.Set, *Args) ([]ox.Option, error){
		QueryCommand,
		SchemaCommand,
		DumpCommand,
	} {
		subopts, err := f(ts, args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		opts = append(opts, ox.Sub(subopts...))
	}
	return opts
}

// QueryCommand builds the query command options.
func QueryCommand(ts *templates.Set, args *Args) ([]ox.Option, error) {
	// query flags
	fs := ox.Flags()
	fs = databaseFlags(fs, args)
	fs = outFlags(fs, args)
	fs = fs.String("query", "custom database query (uses stdin if not provided)", ox.Bind(&args.QueryParams.Query), ox.Short("Q")).
		String("type", "type name", ox.Bind(&args.QueryParams.Type), ox.Short("T")).
		String("type-comment", "type comment", ox.Bind(&args.QueryParams.TypeComment)).
		String("func", "func name", ox.Bind(&args.QueryParams.Func), ox.Short("F")).
		String("func-comment", "func comment", ox.Bind(&args.QueryParams.FuncComment)).
		Bool("trim", "enable trimming whitespace", ox.Bind(&args.QueryParams.Trim), ox.Short("M")).
		Bool("strip", "enable stripping type casts", ox.Bind(&args.QueryParams.Strip), ox.Short("B")).
		Bool("one", "enable returning single (only one) result", ox.Bind(&args.QueryParams.One), ox.Short("1")).
		Bool("flat", "enable returning unstructured (flat) values", ox.Bind(&args.QueryParams.Flat), ox.Short("l")).
		Bool("exec", "enable exec (disables query introspection)", ox.Bind(&args.QueryParams.Exec), ox.Short("X")).
		Bool("interpolate", "enable interpolation of embedded params", ox.Bind(&args.QueryParams.Interpolate), ox.Short("I")).
		String("delimiter", "delimiter used for embedded params", ox.Bind(&args.QueryParams.Delimiter), ox.Short("L"), ox.Default("%%")).
		String("fields", "override field names for results", ox.Bind(&args.QueryParams.Fields), ox.Short("Z")).
		Bool("allow-nulls", "allow result fields with NULL values", ox.Bind(&args.QueryParams.AllowNulls), ox.Short("U"))
	var err error
	if fs, err = templateFlags(fs, ts, true, args); err != nil {
		return nil, err
	}
	return []ox.Option{
		ox.Usage("query", "generate code for a database query from a template"),
		ox.Banner("Generate code for a database query from a template."),
		ox.ValidArgs(1, 1),
		fs,
		ox.Exec(Exec("query", ts, args)),
	}, nil
}

// SchemaCommand builds the schema command options.
func SchemaCommand(ts *templates.Set, args *Args) ([]ox.Option, error) {
	// schema flags
	fs := ox.Flags()
	fs = databaseFlags(fs, args)
	fs = outFlags(fs, args)
	fs.Var("fk-mode", args.SchemaParams.FkMode.Desc(), ox.Bind(args.SchemaParams.FkMode), ox.Short("k")).
		Var("include", args.SchemaParams.Include.Desc(), ox.Bind(args.SchemaParams.Include), ox.Short("i")).
		Var("exclude", args.SchemaParams.Exclude.Desc(), ox.Bind(args.SchemaParams.Exclude), ox.Short("e")).
		Bool("use-index-names", "use index names as defined in schema for generated code", ox.Bind(&args.SchemaParams.UseIndexNames), ox.Short("j"))
	var err error
	if fs, err = templateFlags(fs, ts, true, args); err != nil {
		return nil, err
	}
	if fs, err = loaderFlags(fs, args); err != nil {
		return nil, err
	}
	return []ox.Option{
		ox.Usage("schema", "generate code for a database schema from a template"),
		ox.Banner("Generate code for a database schema from a template."),
		ox.ValidArgs(1, 1),
		fs,
		ox.Exec(Exec("schema", ts, args)),
	}, nil
}

// DumpCommand builds the dump command options.
func DumpCommand(ts *templates.Set, args *Args) ([]ox.Option, error) {
	// dump flags
	fs := ox.Flags()
	fs, err := templateFlags(fs, ts, false, args)
	if err != nil {
		return nil, err
	}
	return []ox.Option{
		ox.Usage("dump", "dump template to path"),
		ox.Banner("Dump template to path."),
		ox.ValidArgs(1, 1),
		fs,
		ox.Exec(func(ctx context.Context, v []string) error {
			// set template
			ts.Use(args.TemplateParams.Type.AsString())
			// get template src
			src, err := ts.Src()
			if err != nil {
				return err
			}
			// ensure out dir exists
			if err := checkDir(v[0]); err != nil {
				return err
			}
			// dump
			return iofs.WalkDir(src, ".", func(n string, d iofs.DirEntry, err error) error {
				switch {
				case err != nil:
					return err
				case d.IsDir():
					return os.MkdirAll(filepath.Join(v[0], n), 0o755)
				}
				buf, err := iofs.ReadFile(src, n)
				if err != nil {
					return err
				}
				return os.WriteFile(filepath.Join(v[0], n), buf, 0o644)
			})
		}),
	}, nil
}

// databaseFlags adds database flags to the flag set.
func databaseFlags(fs *ox.FlagSet, args *Args) *ox.FlagSet {
	return fs.
		String("schema", "database schema name", ox.Bind(&args.LoaderParams.Schema), ox.Short("s"))
	// cmd.SetUsageTemplate(cmd.UsageTemplate() + "\nArgs:\n  <database url>  database url (e.g., postgres://user:pass@localhost:port/dbname, mysql://... )\n\n")
}

// outFlags adds out flags to the flag set.
func outFlags(fs *ox.FlagSet, args *Args) *ox.FlagSet {
	return fs.
		String("out", "out path", ox.Bind(&args.OutParams.Out), ox.Short("o"), ox.Default("models")).
		Bool("debug", "debug generated code (writes generated code to disk without post processing)", ox.Bind(&args.OutParams.Debug), ox.Short("D")).
		String("single", "output all contents to the specified file", ox.Bind(&args.OutParams.Single), ox.Short("S"))
}

// loaderFlags adds database loader flags to the flag set.
func loaderFlags(fs *ox.FlagSet, args *Args) (*ox.FlagSet, error) {
	var err error
	for _, flag := range loader.Flags() {
		if fs, err = flag.Add(fs, args.LoaderParams.Flags); err != nil {
			return nil, err
		}
	}
	return fs, nil
}

// templateFlags adds template flags to the flag set.
func templateFlags(fs *ox.FlagSet, ts *templates.Set, extra bool, args *Args) (*ox.FlagSet, error) {
	fs = fs.Var("template", args.TemplateParams.Type.Desc(), ox.Bind(args.TemplateParams.Type), ox.Short("t"))
	if extra {
		fs = fs.String("src", "template source directory", ox.Bind(&args.TemplateParams.Src), ox.Short("d"))
		var err error
		for _, name := range ts.Targets() {
			for _, flag := range ts.Flags(name) {
				if fs, err = flag.Add(fs, args.TemplateParams.Flags); err != nil {
					return nil, err
				}
			}
		}
	}
	return fs, nil
}

func parseArg(short, full string, args []string) (s string) {
	defer func() {
		s = strings.TrimSpace(s)
	}()
	for i := range args {
		switch s := strings.TrimSpace(args[i]); {
		case s == short, s == full:
			if i < len(args)-1 {
				return args[i+1]
			}
		case strings.HasPrefix(s, short):
			return strings.TrimPrefix(s, short)
		case strings.HasPrefix(s, full):
			return strings.TrimPrefix(s, full)
		}
	}
	return ""
}

// Exec handles the execution for query and schema.
func Exec(mode string, ts *templates.Set, args *Args) func(context.Context, []string) error {
	return func(ctx context.Context, cmdargs []string) error {
		// setup args
		if err := checkArgs(ctx, mode, ts, args); err != nil {
			return err
		}
		// set template
		ts.Use(args.TemplateParams.Type.AsString())
		// build context
		ctx = BuildContext(ctx, args)
		// enable verbose output for sql queries
		if args.Verbose {
			models.SetLogger(func(str string, v ...any) {
				s, z := "SQL: %s\n", []any{str}
				if len(v) != 0 {
					s, z = s+"PARAMS: %v\n", append(z, v)
				}
				fmt.Printf(s+"\n", z...)
			})
		}
		// open database
		var err error
		if ctx, err = open(ctx, cmdargs[0], args.LoaderParams.Schema); err != nil {
			return err
		}
		// load
		set, err := load(ctx, mode, ts, args)
		if err != nil {
			return err
		}
		return Generate(ctx, mode, ts, set, args)
	}
}

// Generate generates the dbtpl files with the provided templates, data, and
// arguments.
func Generate(ctx context.Context, mode string, ts *templates.Set, set *xo.Set, args *Args) error {
	// create set context
	ctx = ts.NewContext(ctx, mode)
	if err := displayErrors(ts); err != nil {
		return err
	}
	// preprocess
	ts.Pre(ctx, args.OutParams.Out, mode, set)
	if err := displayErrors(ts); err != nil {
		return err
	}
	// process
	ts.Process(ctx, args.OutParams.Out, mode, set)
	if err := displayErrors(ts); err != nil {
		return err
	}
	// post
	if !args.OutParams.Debug {
		ts.Post(ctx, mode)
		if err := displayErrors(ts); err != nil {
			return err
		}
	}
	// dump
	ts.Dump(args.OutParams.Out)
	if err := displayErrors(ts); err != nil {
		return err
	}
	return nil
}

// checkArgs sets up and checks args.
func checkArgs(ctx context.Context, mode string, ts *templates.Set, args *Args) error {
	// check template is available for the mode
	if err := ts.For(mode); err != nil {
		return err
	}
	// check --src and --template are exclusive
	/*
		if cmd.Flags().Lookup("src").Changed && cmd.Flags().Lookup("template").Changed {
			return errors.New("--src and --template cannot be used together")
		}
	*/
	// read query string from stdin if not provided via --query
	if mode == "query" && args.QueryParams.Query == "" {
		buf, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		args.QueryParams.Query = string(bytes.TrimRight(buf, "\r\n"))
	}
	// check out path
	if args.OutParams.Out != "" {
		var err error
		if args.OutParams.Out, err = realpath.Realpath(args.OutParams.Out); err != nil {
			return err
		}
		if err := checkDir(args.OutParams.Out); err != nil {
			return err
		}
	}
	return nil
}

// BuildContext builds a context for the mode and template.
func BuildContext(ctx context.Context, args *Args) context.Context {
	// add loader flags
	for k, v := range args.LoaderParams.Flags {
		ctx = context.WithValue(ctx, k, v.Interface())
	}
	// add template flags
	for k, v := range args.TemplateParams.Flags {
		ctx = context.WithValue(ctx, k, v.Interface())
	}
	// add out
	ctx = context.WithValue(ctx, xo.OutKey, args.OutParams.Out)
	ctx = context.WithValue(ctx, xo.SingleKey, args.OutParams.Single)
	return ctx
}

// open opens a connection to the database, returning a context for use in
// template generation.
func open(ctx context.Context, urlstr, schema string) (context.Context, error) {
	v, err := user.Current()
	if err != nil {
		return nil, err
	}
	// parse dsn
	u, err := dburl.Parse(urlstr)
	if err != nil {
		return nil, err
	}
	// open database
	db, err := passfile.OpenURL(u, v.HomeDir, "xopass")
	if err != nil {
		return nil, err
	}
	// add driver to context
	ctx = context.WithValue(ctx, xo.DriverKey, u.Driver)
	// add db to context
	ctx = context.WithValue(ctx, xo.DbKey, db)
	// determine schema
	if schema == "" {
		if schema, err = loader.Schema(ctx); err != nil {
			return nil, err
		}
	}
	// add schema to context
	ctx = context.WithValue(ctx, xo.SchemaKey, schema)
	return ctx, nil
}

// load loads a set of queries or schemas.
func load(ctx context.Context, mode string, _ *templates.Set, args *Args) (*xo.Set, error) {
	f := LoadSchema
	if mode == "query" {
		f = LoadQuery
	}
	set := new(xo.Set)
	if err := f(ctx, set, args); err != nil {
		return nil, err
	}
	return set, nil
}

// displayErrors displays collected errors from the set.
func displayErrors(ts *templates.Set) error {
	if errors := ts.Errors(); len(errors) != 0 {
		for _, err := range errors {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		}
		return fmt.Errorf("%d errors encountered", len(errors))
	}
	return nil
}

// checkDir checks that dir exists.
func checkDir(dir string) error {
	if !isDir(dir) {
		return fmt.Errorf("%s must exist and must be a directory", dir)
	}
	return nil
}

// isDir determines if dir is a directory.
func isDir(dir string) bool {
	if fi, err := os.Stat(dir); err == nil {
		return fi.IsDir()
	}
	return false
}
