package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	var projectPath string
	var root = app.NewWithID("com.lordryns.ember")
	var window = root.NewWindow("Ember v 0.1")
	window.Resize(fyne.NewSize(800, 500))

	var mainContent = container.NewVBox(widget.NewLabel("Ember v 0.1"))
	var rootContent = container.NewBorder(projectSelectContainer(window, &projectPath), nil, nil, nil, mainContent)

	window.SetContent(rootContent)
	window.ShowAndRun()
}
