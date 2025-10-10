package engine

import (
	"bytes"
	"ember/globals"
	"os"
	"path/filepath"
	"text/template"
)

func GenerateSourceFromConfig(config *globals.GameConfig) (string, error) {

	// claude assisted in the generation of this template
	var mainString = `
<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>{{.Title}}</title>
</head>
<body>
<script src="https://unpkg.com/kaplay@3001.0.19/dist/kaplay.js"></script>
<script>
const game = kaplay({
	{{- if .Color }}
	background: "{{.Color}}",
	{{- end }}
});
setGravity({{.Gravity}});
{{range .Objects}}
// Create object {{.ID}}
var {{.ID}} = add([
	rect({{.Size.X}}, {{.Size.Y}}),
	pos({{.Pos.X}}, {{.Pos.Y}}),
	color("{{.Color}}"),
	{{- if .IsBody }}body({isStatic: {{.IsStatic}}}),{{- end }}
	{{- if .HasArea }}area(),{{- end }}
]);
{{.ID}}.gravityScale = {{.Weight}};
{{$objID := .ID}}

{{range .KeyMap}}
// Function for key {{.Key}} on object {{$objID}}

{{if eq .PressType "on_key_down"}}
onKeyDown("{{.Key}}", function() {
	{{.Func.ID}}({{range $i, $arg := .Func.Args}}{{if $i}}, {{end}}{{if eq $arg.Value "this"}}{{$objID}}{{else}}{{$arg.Value}}{{end}}{{end}});
});
{{else if eq .PressType "on_key_up"}}
onKeyRelease("{{.Key}}", function() {
	{{.Func.ID}}({{range $i, $arg := .Func.Args}}{{if $i}}, {{end}}{{if eq $arg.Value "this"}}{{$objID}}{{else}}{{$arg.Value}}{{end}}{{end}});
});
{{else if eq .PressType "on_key_press"}}
onKeyPress("{{.Key}}", function() {
	{{.Func.ID}}({{range $i, $arg := .Func.Args}}{{if $i}}, {{end}}{{if eq $arg.Value "this"}}{{$objID}}{{else}}{{$arg.Value}}{{end}}{{end}});
});
{{else if eq .PressType "is_key_pressed"}}
onUpdate(function() {
	if (isKeyDown("{{.Key}}")) {
		{{.Func.ID}}({{range $i, $arg := .Func.Args}}{{if $i}}, {{end}}{{if eq $arg.Value "this"}}{{$objID}}{{else}}{{$arg.Value}}{{end}}{{end}});
	}
});
{{end}}
{{end}}
{{end}}

{{range .Functions}}
function {{.ID}}({{range $i, $arg := .Args}}{{if $i}}, {{end}}{{$arg.ID}}{{end}}) {
	{{.Src}}
}
{{end}}
{{if .Update}}
onUpdate(() => {
	{{range .Update}}
	{{.ID}}({{range $i, $arg := .Args}}{{if $i}}, {{end}}{{$arg.Value}}{{end}});
	{{end}}
});
{{end}}
</script>
</body>
</html>
`
	var tmpl, err2 = template.New("base").Parse(mainString)

	if err2 != nil {
		return "", err2
	}

	var buf bytes.Buffer
	tmpl.Execute(&buf, config)

	return buf.String(), nil
}

func GenerateFile(path string, config *globals.GameConfig) error {
	var fullpath = filepath.Join(path, "index.html")
	var file, err = os.Create(fullpath)

	if err != nil {
		return err
	}

	var gen, err2 = GenerateSourceFromConfig(config)

	if err2 != nil {
		return err2
	}

	file.WriteString(gen)
	defer file.Close()
	return nil
}
