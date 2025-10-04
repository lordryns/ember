package main

import (
	"ember/filemenu"

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

	var mainContent = container.NewVBox(widget.NewLabel("Hello from Ember!"))
	var rootContent = container.NewBorder(toolBar(window, &projectPath), nil, nil, nil, mainContent)

	window.SetContent(rootContent)
	window.ShowAndRun()
}

func toolBar(window fyne.Window, projectPath *string) *fyne.Container {
	var fileButton = widget.NewButton("File menu", func() {
		filemenu.FileMenuWindow(window, projectPath)
	})
	return container.NewBorder(nil, nil, nil, container.NewHBox(fileButton), widget.NewLabel("Ember v 0.1"))
}
