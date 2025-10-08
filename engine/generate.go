package engine

import (
	"bytes"
	"ember/globals"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

func GenerateSourceFromConfig(config *globals.GameConfig) (string, error) {
	var mainString = `
<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>Good game</title>
</head>
<body>
<script src="https://unpkg.com/kaplay@3001.0.19/dist/kaplay.js"></script>
<script src="game.js"></script>
</body>
<script>
	kaplay()
	setGravity({{.Gravity}})

	{{range .Objects}}
		var {{.ID}} = add([rect({{.Size.X}}, {{.Size.Y}}), pos({{.Pos.X}}, {{.Pos.Y}}), color("{{.Color}}"), {{if .IsBody}}body({isStatic: {{.IsStatic}} }), {{end}} {{if .HasArea}}area(){{end}}]);
		{{.ID}}.gravityScale = {{.Weight}};
		{{end}}
</script>
</html>
`
	var tmpl, err2 = template.New("base").Parse(mainString)

	if err2 != nil {
		return "", err2
	}

	var buf bytes.Buffer
	tmpl.Execute(&buf, config)

	fmt.Println(buf.String())
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
