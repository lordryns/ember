package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"text/template"
)

type GameConfig struct {
	Title     string       `json:"title"`
	Gravity   int          `json:"gravity"`
	Objects   []GameObject `json:"objects"`
	Functions []GameFunc   `json:"functions"`
}

type GameObject struct {
	ID      string   `json:"id"`
	Shape   string   `json:"shape"`
	Pos     Position `json:"pos"`
	Size    Size     `json:"size"`
	Color   string   `json:"color"`
	IsBody  bool     `json:"is_body"`
	HasArea bool     `json:"has_area"`
	KeyMap  []KeyMap `json:"key_map"`
}

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}
type Size struct {
	X int `json:"x"`
	Y int `json:"y"`
}
type KeyMap struct {
	Key  string `json:"key"`
	Func string `json:"func"`
}

type GameFunc struct {
	ID          string `json:"id"`
	UseExternal bool   `json:"use_external"`
	Src         string `json:"src"`
}

func LoadConfig(path string) (GameConfig, error) {
	var config GameConfig

	var fileBytes, err = os.ReadFile(path)
	if err != nil {
		return config, err
	}

	if err := json.Unmarshal(fileBytes, &config); err != nil {
		return config, err
	}

	return config, err
}

// var x = engine.GAME_CONFIG.Objects[0].KeyMap[0].PressType

func fmain() {
	generateSource()
}

func generateSource() {
	var config, err = LoadConfig("schema.json")
	if err != nil {
		panic(err)
	}

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

	{{range .Objects}}
	var {{.ID}} = add([rect({{.Size.X}}, {{.Size.Y}}), pos({{.Pos.X}}, {{.Pos.Y}}), color("{{.Color}}"), {{if .IsBody}}body(), {{end}} {{if .HasArea}}area(){{end}}]);
	{{end}}
</script>
</html>
`
	var tmpl, err2 = template.New("base").Parse(mainString)

	if err2 != nil {
		panic(err2)
	}

	var buf bytes.Buffer
	tmpl.Execute(&buf, config)

	fmt.Println(buf.String())
}
