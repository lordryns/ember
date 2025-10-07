package main

import (
	"ember/engine"
	"ember/filemenu"
	"ember/helpers"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var currentTabID int

func main() {
	var projectPath = widget.NewLabel("No project open...yet")
	var root = app.NewWithID("com.lordryns.ember")
	var window = root.NewWindow("Ember v 0.1")
	window.Resize(fyne.NewSize(800, 500))

	var mainContent = container.NewStack()
	var windowSidebar = sideBar(window, mainContent)

	var alignContent = container.NewBorder(toolBar(window, projectPath, mainContent), nil, nil, nil, mainContent)
	var rootContent = container.NewHSplit(windowSidebar, alignContent)

	rootContent.SetOffset(0.2)
	window.SetContent(rootContent)
	window.ShowAndRun()
}
func toolBar(window fyne.Window, projectPath *widget.Label, mainContentBlock *fyne.Container) *fyne.Container {
	var runButton = widget.NewButton("Run Game", func() {})

	var saveButton = widget.NewButton("Save", func() {
		var path = projectPath.Text
		go func() {
			var _, err = os.Stat(engine.PROJECT_PATH)
			if !os.IsNotExist(err) {
				fyne.Do(func() {
					projectPath.SetText("Saved!")
				})

				time.Sleep(time.Second * 2)

				fyne.Do(func() {
					projectPath.SetText(path)
				})

			} else {
				fyne.CurrentApp().SendNotification(fyne.NewNotification("Error", "You must open a project before you can save anything!"))
			}
		}()

	})
	saveButton.Importance = widget.HighImportance
	var refreshButton = widget.NewButton("Refresh", func() {
		setTabBasedOnId(currentTabID, mainContentBlock)
	})
	var fileButton = widget.NewButton("File menu", func() {
		filemenu.FileMenuWindow(window, projectPath, mainContentBlock, func() {
			setTabBasedOnId(currentTabID, mainContentBlock)
		})
	})
	return container.NewBorder(nil, nil, nil, container.NewHBox(runButton, saveButton, refreshButton, fileButton), projectPath)
}

func setTabBasedOnId(id widget.ListItemID, mainContent *fyne.Container) {
	switch id {
	case 1:
		mainContent.Objects = []fyne.CanvasObject{spritesContentTab()}
	case 2:
		mainContent.Objects = []fyne.CanvasObject{objectContentTab()}
	case 3:
		mainContent.Objects = []fyne.CanvasObject{functionContentTab()}
	case 4:
		mainContent.Objects = []fyne.CanvasObject{PhysicsContentTab()}
	default:
		mainContent.Objects = []fyne.CanvasObject{defaultContentTab(&engine.GAME_CONFIG)}
	}

	currentTabID = id
	mainContent.Refresh()
}

func sideBar(window fyne.Window, mainContent *fyne.Container) *widget.List {
	var list = widget.NewList(func() int { return len(helpers.SidebarItems) }, func() fyne.CanvasObject { return widget.NewLabel("") },
		func(lii widget.ListItemID, co fyne.CanvasObject) {
			co.(*widget.Label).SetText(helpers.SidebarItems[lii])
		})

	list.OnSelected = func(id widget.ListItemID) {
		setTabBasedOnId(id, mainContent)
	}

	list.Select(0)

	return list
}

func defaultContentTab(config *engine.GameConfig) *fyne.Container {
	var titleLabel = widget.NewLabel("Game title: ")
	var titleEntry = widget.NewEntry()
	titleEntry.SetText(config.Title)
	var titleContainer = container.NewBorder(nil, nil, titleLabel, nil, titleEntry)

	var geometryLabel = widget.NewLabel("Geometry: ")
	var XEntry = widget.NewEntry()
	XEntry.SetPlaceHolder("X")
	XEntry.SetText("width()")
	var YEntry = widget.NewEntry()
	YEntry.SetText("height()")
	YEntry.SetPlaceHolder("Y")

	var geometryContainer = container.New(layout.NewGridLayout(3), geometryLabel, XEntry, YEntry)
	return container.NewVBox(titleContainer, geometryContainer)
}

func spritesContentTab() *fyne.Container {
	var textLabel = widget.NewLabel("Sprites tab")
	return container.NewHBox(textLabel)
}

func objectContentTab() *fyne.Container {
	var newObjectButton = widget.NewButton("New Object", func() {})
	var objectList = widget.NewList(func() int { return len(engine.GAME_CONFIG.Objects) }, func() fyne.CanvasObject { return widget.NewLabel("") },
		func(lii widget.ListItemID, co fyne.CanvasObject) {
			co.(*widget.Label).SetText(engine.GAME_CONFIG.Objects[lii].ID)
		})

	var mainContent = container.NewVBox(container.NewCenter(newObjectButton))
	var layoutSplit = container.NewHSplit(mainContent, objectList)
	layoutSplit.SetOffset(0.8)
	return container.New(layout.NewGridLayout(1), layoutSplit)
}

func functionContentTab() *fyne.Container {
	var textLabel = widget.NewLabel("Function tab")
	return container.NewHBox(textLabel)
}

func PhysicsContentTab() *fyne.Container {
	var textLabel = widget.NewLabel("Physics tab")
	return container.NewHBox(textLabel)
}
