package types

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/gobwas/glob"
	"github.com/xo/ox"
)

// Flag is a option flag.
type Flag struct {
	ContextKey ContextKey
	Type       string
	Desc       string
	Default    any
	Short      string
	Enums      []string
	Aliases    []string
	Hidden     bool
}

// FlagSet is a set of option flags.
type FlagSet struct {
	Type string
	Name string
	Flag Flag
}

// Add adds the flag to the cmd.
func (flag FlagSet) Add(fs *ox.FlagSet, values map[ContextKey]*Value) (*ox.FlagSet, error) {
	typ := ox.StringT
	switch flag.Flag.Type {
	case "bool":
		typ = ox.BoolT
	case "int":
		typ = ox.IntT
	case "string":
	case "[]string":
		typ = ox.SliceT
	case "glob":
		typ = ox.GlobT
	default:
		return nil, fmt.Errorf("unknown flag type %s", flag.Flag.Type)
	}
	// create value
	if _, ok := values[flag.Flag.ContextKey]; !ok {
		values[flag.Flag.ContextKey] = NewValue(flag.Flag.Type, flag.Flag.Default, flag.Flag.Desc, flag.Flag.Enums...)
	}
	opts := []ox.Option{
		ox.Bind(values[flag.Flag.ContextKey]),
		ox.Hidden(flag.Flag.Hidden),
		typ,
	}
	if flag.Flag.Short != "" {
		opts = append(opts, ox.Short(flag.Flag.Short))
	}
	if flag.Flag.Default != nil {
		opts = append(opts, ox.Default(flag.Flag.Default))
	}
	if len(flag.Flag.Aliases) != 0 {
		opts = append(opts, ox.Aliases(flag.Flag.Aliases...))
	}
	return fs.Var(flag.Type+"-"+flag.Name, values[flag.Flag.ContextKey].Desc(), opts...), nil
}

// Value wraps a flag value.
type Value struct {
	typ   string
	def   any
	desc  string
	enums []string
	set   bool
	v     any
}

// NewValue creates a new flag value.
func NewValue(typ string, def any, desc string, enums ...string) *Value {
	var z any
	switch typ {
	case "bool":
		var b bool
		z = b
	case "int":
		var i int
		z = i
	case "string":
		var s string
		z = s
	case "[]string":
		var s []string
		z = s
	case "glob":
		var v []glob.Glob
		z = v
	}
	v := &Value{
		typ:   typ,
		def:   def,
		desc:  desc,
		enums: enums,
		v:     z,
	}
	/*
		if v.def != nil {
			if err := v.Set(v.def); err != nil {
				panic(err)
			}
			v.set = false
		}
	*/
	return v
}

// Desc returns the usage description for the flag value.
func (v *Value) Desc() string {
	if v.enums != nil {
		return v.desc + " <" + strings.Join(v.enums, "|") + ">"
	}
	return v.desc
}

// Set satisfies the [ox.Value] interface.
func (v *Value) Set(s string) error {
	v.set = true
	if v.enums != nil {
		if !slices.Contains(v.enums, s) {
			return fmt.Errorf("invalid value %q", s)
		}
	}
	switch v.typ {
	case "bool":
		b, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		v.v = b
	case "int":
		i, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		v.v = i
	case "string":
		v.v = s
	case "[]string":
		v.v = append(v.v.([]string), strings.Split(s, ",")...)
	case "glob":
		g, err := glob.Compile(s)
		if err != nil {
			return err
		}
		v.v = append(v.v.([]glob.Glob), g)
	}
	return nil
}

// Interface returns the value.
func (v *Value) Interface() any {
	if v.v == nil {
		panic("v should not be nil!")
	}
	return v.v
}

// AsBool returns the value as a bool.
func (v *Value) AsBool() bool {
	b, _ := v.v.(bool)
	return b
}

// AsInt returns the value as a int.
func (v *Value) AsInt() int {
	i, _ := v.v.(int)
	return i
}

// AsString returns the value as a string.
func (v *Value) AsString() string {
	s, _ := v.v.(string)
	return s
}

// AsStringSlice returns the value as a string slice.
func (v *Value) AsStringSlice() []string {
	z, _ := v.v.([]string)
	return z
}

// AsGlob returns the value as a glob slice.
func (v *Value) AsGlob() []glob.Glob {
	z, _ := v.v.([]glob.Glob)
	return z
}
