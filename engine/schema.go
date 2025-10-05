package engine

type GameConfig struct {
	Title     string     `json:"title"`
	Sprites   []Sprite   `json:"sprites"`
	Objects   []Object   `json:"objects"`
	Functions []Function `json:"functions"`
}

type Sprite struct {
	ID          string       `json:"id"`
	Type        string       `json:"type"`
	Path        string       `json:"path"`
	Spritesheet *Spritesheet `json:"spritesheet"`
}

type Spritesheet struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Frames int `json:"frames"`
}

type Object struct {
	ID      string            `json:"id"`
	Sprite  string            `json:"sprite"`
	X       int               `json:"x"`
	Y       int               `json:"y"`
	Dynamic bool              `json:"dynamic"`
	UseKey  bool              `json:"useKey"`
	KeyMap  map[string]string `json:"key_map"`
}

type Function struct {
	ID          string `json:"id"`
	UseExternal bool   `json:"use_external"`
	Src         string `json:"src"`
}
