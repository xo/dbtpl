{{ define "dot" -}}
{{ $s := .Data -}}
// Generated by dbtpl for the {{ $s.Name }} schema.
digraph {{ $s.Name|normalize }} {
	{{ if defaults -}}
	// Defaults
	{{ range defaults -}}
	{{ . }}
	{{ end }}
	{{ end -}}

	// Nodes (tables)
	{{- range $s.Tables }}
	{{ schema .Name }} [ label=<
		<table border="0" cellborder="1" cellspacing="0" cellpadding="4">
		<tr>{{ header (schema .Name) }}</tr>
		{{ range .Columns -}}
		<tr>{{ row . }}</tr>
		{{ end -}}
		</table>> ]
	{{ end }}

	{{- range $s.Tables }}
	{{- $t := .  -}}
	{{- range $t.ForeignKeys -}}
	{{- $fkey := . -}}
	{{- range $i, $field := .Fields }}
	{{ edge $t $fkey $i }} [
		headlabel={{ quotes $fkey.Name }}]
	{{- end }}
	{{- end -}}
	{{- end }}
}
{{ end }}
