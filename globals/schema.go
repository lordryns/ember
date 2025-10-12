package globals

type GameConfig struct {
	Title     string       `json:"title"`
	Color     string       `json:"color"`
	Gravity   int          `json:"gravity"`
	Objects   []GameObject `json:"objects"`
	Functions []GameFunc   `json:"functions"`
	Update    []GameFunc   `json:"update"`
}

type GameObject struct {
	ID           string   `json:"id"`
	Shape        string   `json:"shape"`
	Pos          Position `json:"pos"`
	Size         Size     `json:"size"`
	Mass         int      `json:"mass"`
	GravityScale int      `json:"gravity_scale"`
	Color        string   `json:"color"`
	IsBody       bool     `json:"is_body"`
	HasArea      bool     `json:"has_area"`
	IsStatic     bool     `json:"is_static"`
	KeyMap       []KeyMap `json:"key_map"`
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
	Key       string   `json:"key"`
	Func      GameFunc `json:"func"`
	PressType string   `json:"press_type"`
}

type GameFunc struct {
	ID   string `json:"id"`
	Args []Arg  `json:"args"`
	Src  string `json:"src"`
}

type Arg struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Value string `json:"values"`
}
type SupportedInputs struct {
	Keyboard []string
	Mouse    []string
	Touch    []string
	Gamepad  []string
}
