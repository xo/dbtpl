//go:build dbtpl

package yaml

import (
	"context"
	"text/template"

	"github.com/goccy/go-yaml"
	xo "github.com/xo/dbtpl/types"
)

// Init registers the template.
func Init(ctx context.Context, f func(xo.TemplateType)) error {
	f(xo.TemplateType{
		Modes: []string{"query", "schema"},
		Flags: []xo.Flag{},
		Funcs: func(ctx context.Context, _ string) (template.FuncMap, error) {
			return template.FuncMap{
				// yaml marshals v as yaml.
				"yaml": func(v any) (string, error) {
					buf, err := yaml.MarshalWithOptions(v)
					if err != nil {
						return "", err
					}
					return string(buf), nil
				},
			}, nil
		},
		Process: func(ctx context.Context, _ string, set *xo.Set, emit func(xo.Template)) error {
			emit(xo.Template{
				Partial: "yaml",
				Dest:    "dbtpl.dbtpl.yaml",
				Data:    set,
			})
			return nil
		},
	})
	return nil
}
