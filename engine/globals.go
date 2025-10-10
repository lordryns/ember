package engine

import "ember/globals"

var PROJECT_PATH string
var GAME_CONFIG globals.GameConfig = globals.GameConfig{Title: "Ember Game", Gravity: 0,
	Objects: []globals.GameObject{}, Functions: []globals.GameFunc{}}

var CUSTOM_TYPES = []string{"String", "Number", "Boolean", "Object"}

var ALL_INPUTS = globals.SupportedInputs{
	Keyboard: []string{
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m",
		"n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
		"up", "down", "left", "right",
		"space", "enter", "escape", "shift", "ctrl", "alt", "tab",
	},
	Mouse: []string{"left", "right", "middle", "wheel"},
	Touch: []string{"tap", "swipe", "pinch"},
	Gamepad: []string{
		"buttonSouth", "buttonEast", "buttonWest", "buttonNorth",
		"shoulderLeft", "shoulderRight",
		"triggerLeft", "triggerRight",
		"dpadUp", "dpadDown", "dpadLeft", "dpadRight",
		"start", "select",
	},
}
