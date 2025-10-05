package main

import (
	"ember/engine"
	"ember/filemenu"
	"ember/helpers"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	var configData engine.GameConfig
	var projectPath = widget.NewLabel("No project open...yet")
	var root = app.NewWithID("com.lordryns.ember")
	var window = root.NewWindow("Ember v 0.1")
	window.Resize(fyne.NewSize(800, 500))

	var mainContent = container.NewCenter()
	var windowSidebar = sideBar(window, mainContent, &configData)

	var alignContent = container.NewBorder(toolBar(window, projectPath, &configData, mainContent), nil, nil, nil, mainContent)
	var rootContent = container.NewHSplit(windowSidebar, alignContent)

	rootContent.SetOffset(0.2)
	window.SetContent(rootContent)
	window.ShowAndRun()
}
func toolBar(window fyne.Window, projectPath *widget.Label, config *engine.GameConfig, mainContentBlock *fyne.Container) *fyne.Container {
	var refreshButton = widget.NewButton("Refresh", func() {
		projectPath.Refresh()
		mainContentBlock.Refresh()
		for _, obj := range mainContentBlock.Objects {
			obj.Refresh()
		}
	})
	var fileButton = widget.NewButton("File menu", func() {
		filemenu.FileMenuWindow(window, projectPath, config, mainContentBlock)
	})
	return container.NewBorder(nil, nil, nil, container.NewHBox(refreshButton, fileButton), projectPath)
}

func sideBar(window fyne.Window, mainContent *fyne.Container, config *engine.GameConfig) *widget.List {
	var list = widget.NewList(func() int { return len(helpers.SidebarItems) }, func() fyne.CanvasObject { return widget.NewLabel("") },
		func(lii widget.ListItemID, co fyne.CanvasObject) {
			co.(*widget.Label).SetText(helpers.SidebarItems[lii])
		})

	list.OnSelected = func(id widget.ListItemID) {
		switch id {
		case 0:
			mainContent.Objects = []fyne.CanvasObject{defaultContentTab(config)}
		case 1:
			mainContent.Objects = []fyne.CanvasObject{spritesContentTab()}
		case 2:
			mainContent.Objects = []fyne.CanvasObject{objectContentTab()}
		case 3:
			mainContent.Objects = []fyne.CanvasObject{functionContentTab()}
		case 4:
			mainContent.Objects = []fyne.CanvasObject{PhysicsContentTab()}
		default:
			mainContent.Objects = []fyne.CanvasObject{}
		}

		mainContent.Refresh()
	}

	list.Select(0)

	return list
}

func defaultContentTab(config *engine.GameConfig) *fyne.Container {
	var textLabel = widget.NewLabel(fmt.Sprintf("Project title: %v", config.Title))
	return container.NewHBox(textLabel)
}

func spritesContentTab() *fyne.Container {
	var textLabel = widget.NewLabel("Sprites tab")
	return container.NewHBox(textLabel)
}

func objectContentTab() *fyne.Container {
	var textLabel = widget.NewLabel("Object tab")
	return container.NewHBox(textLabel)
}

func functionContentTab() *fyne.Container {
	var textLabel = widget.NewLabel("Function tab")
	return container.NewHBox(textLabel)
}

func PhysicsContentTab() *fyne.Container {
	var textLabel = widget.NewLabel("Physics tab")
	return container.NewHBox(textLabel)
}
