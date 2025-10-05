package main

import (
	"ember/filemenu"
	"ember/helpers"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	var projectPath = widget.NewLabel("No project open...yet")
	var root = app.NewWithID("com.lordryns.ember")
	var window = root.NewWindow("Ember v 0.1")
	window.Resize(fyne.NewSize(800, 500))

	var windowSidebar = sideBar(window)
	var mainContent = container.NewVBox(widget.NewLabel("Hello from Ember!"))

	var alignContent = container.NewBorder(toolBar(window, projectPath, windowSidebar), nil, nil, nil, mainContent)
	var rootContent = container.NewHSplit(windowSidebar, alignContent)

	rootContent.SetOffset(0.2)
	window.SetContent(rootContent)
	window.ShowAndRun()
}
func toolBar(window fyne.Window, projectPath *widget.Label, windowSidebar *widget.List) *fyne.Container {
	var fileButton = widget.NewButton("File menu", func() {
		filemenu.FileMenuWindow(window, projectPath, windowSidebar)
	})
	return container.NewBorder(nil, nil, nil, container.NewHBox(fileButton), projectPath)
}

func sideBar(window fyne.Window) *widget.List {
	var list = widget.NewList(func() int { return len(helpers.SidebarItems) }, func() fyne.CanvasObject { return widget.NewLabel("") },
		func(lii widget.ListItemID, co fyne.CanvasObject) {
			co.(*widget.Label).SetText(helpers.SidebarItems[lii])
		})

	list.OnSelected = func(id widget.ListItemID) {
		dialog.ShowInformation("info", helpers.SidebarItems[id], window)
	}

	return list
}
