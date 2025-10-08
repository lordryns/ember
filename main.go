package main

import (
	"ember/engine"
	"ember/filemenu"
	"ember/globals"
	"ember/helpers"
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/lusingander/colorpicker"
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
	var runButton = widget.NewButton("", func() {})
	runButton = widget.NewButton("Run Game", func() {
		var _, err = os.Stat(engine.PROJECT_PATH)
		if !os.IsNotExist(err) {
			go func() {
				go func() {
					if runButton.Text == "Run Game" {
						if err := engine.StartDevEngine(&engine.GAME_CONFIG, runButton); err != nil {
							fyne.CurrentApp().SendNotification(fyne.NewNotification("Error", fmt.Sprintf("Failed to run game! err: %v", err)))
							return
						}
					} else {
						engine.StopDevEngine(runButton)
					}

				}()

				fyne.Do(func() {
					projectPath.SetText(engine.PROJECT_PATH)
				})

			}()
		} else {
			fyne.CurrentApp().SendNotification(fyne.NewNotification("Error", "You must open a project before you can run anything!"))
		}

	})

	var saveButton = widget.NewButton("Save", func() {
		var path = projectPath.Text
		go func() {
			var _, err = os.Stat(engine.PROJECT_PATH)
			if !os.IsNotExist(err) {
				fyne.Do(func() {
					if err := helpers.WriteStructToFile(filepath.Join(engine.PROJECT_PATH, "ember.json"), &engine.GAME_CONFIG); err != nil {
						fyne.CurrentApp().SendNotification(fyne.NewNotification("Error", fmt.Sprintf("Unable to save to path! err: %v", err)))
						return
					} else {

						projectPath.SetText("Saved!")
					}
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
		setTabBasedOnId(currentTabID, mainContentBlock, window)
	})
	var fileButton = widget.NewButton("File menu", func() {
		filemenu.FileMenuWindow(window, projectPath, mainContentBlock, func() {
			setTabBasedOnId(currentTabID, mainContentBlock, window)
		})
	})
	return container.NewBorder(nil, nil, nil, container.NewHBox(runButton, saveButton, refreshButton, fileButton), projectPath)
}

func setTabBasedOnId(id widget.ListItemID, mainContent *fyne.Container, window fyne.Window) {
	switch id {
	case 1:
		mainContent.Objects = []fyne.CanvasObject{spritesContentTab()}
	case 2:
		mainContent.Objects = []fyne.CanvasObject{objectContentTab(mainContent, window)}
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
		setTabBasedOnId(id, mainContent, window)
	}

	list.Select(0)

	return list
}

func defaultContentTab(config *globals.GameConfig) *fyne.Container {
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

	var applyButton = widget.NewButton("", func() {})
	applyButton = widget.NewButton("Apply", func() {
		engine.GAME_CONFIG.Title = titleEntry.Text
		go func() {
			fyne.Do(func() {
				applyButton.SetText("Applied!")
			})

			time.Sleep(time.Second * 2)
			fyne.Do(func() {
				applyButton.SetText("Apply")
			})

		}()
	})
	var geometryContainer = container.New(layout.NewGridLayout(3), geometryLabel, XEntry, YEntry)
	return container.NewVBox(titleContainer, geometryContainer, applyButton)
}

func spritesContentTab() *fyne.Container {
	var textLabel = widget.NewLabel("Sprites tab")
	return container.NewHBox(textLabel)
}

func objectContentTab(mainContentBlock *fyne.Container, window fyne.Window) *fyne.Container {
	var objectList = widget.NewList(func() int { return len(engine.GAME_CONFIG.Objects) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(lii widget.ListItemID, co fyne.CanvasObject) {
			var r = co.(*widget.Label)
			r.SetText(engine.GAME_CONFIG.Objects[lii].ID)
		})

	var idEntry = widget.NewEntry()
	idEntry.SetPlaceHolder("Enter ID")

	var shapeOptions = []string{"Choose shape", "Rect", "Circle"}
	var shapeSelect = widget.NewSelect(shapeOptions, func(s string) {})
	shapeSelect.SetSelected("Choose shape")
	var idShapeRow = container.New(layout.NewGridLayout(2), idEntry, shapeSelect)

	var XPosEntry = widget.NewEntry()
	XPosEntry.SetPlaceHolder("X")

	var YPosEntry = widget.NewEntry()
	YPosEntry.SetPlaceHolder("Y")

	var positionRow = container.New(layout.NewGridLayout(3), widget.NewLabel("Position:"), XPosEntry, YPosEntry)

	var XSizeEntry = widget.NewEntry()
	XSizeEntry.SetPlaceHolder("X")

	var YSizeEntry = widget.NewEntry()
	YSizeEntry.SetPlaceHolder("Y")

	var sizeRow = container.New(layout.NewGridLayout(3), widget.NewLabel("Size:"), XSizeEntry, YSizeEntry)

	var colorEntry = widget.NewEntry()
	var colorPicker = colorpicker.New(200, colorpicker.StyleHue)
	colorPicker.SetOnChanged(func(c color.Color) {
		colorEntry.SetText(helpers.ColorToHex(c))
	})
	var colorButton = widget.NewButton("Set Color", func() {
		dialog.NewCustom("Pick Color", "Dismiss", container.NewCenter(colorPicker), window).Show()
	})

	// just realised i could do this instead of the other method, crazy right?
	var colorRow = container.NewGridWithColumns(3, widget.NewLabel("Set Color:"), colorEntry, colorButton)
	var isBodyCheck = widget.NewCheck("Is Body", func(b bool) {})
	var isAreaCheck = widget.NewCheck("Is Area", func(b bool) {})
	var isStaticCheck = widget.NewCheck("Is Static", func(b bool) {})

	var keyMapTopBar = container.NewBorder(nil, nil, nil, widget.NewButton("Add Key", func() {}))
	var keyMapMainContent = container.NewCenter(widget.NewLabel("Nothing to see here..."))
	var keyMapContainer = container.NewVBox(keyMapTopBar, keyMapMainContent)
	var keyPressButton = widget.NewButton("Key map", func() {
		var keyMapDialog = dialog.NewCustom("Key Map", "Close", keyMapContainer, window)
		keyMapDialog.Resize(fyne.NewSize(400, 400))
		keyMapDialog.Show()
	})
	var areaBodyKeymapRow = container.NewGridWithColumns(4, isBodyCheck, isAreaCheck, isStaticCheck, keyPressButton)

	var currentObjectID int = -1
	var deleteObjectButton = widget.NewButton("Delete", func() {

		if currentObjectID > -1 {
			dialog.NewConfirm("Delete Object?", fmt.Sprintf("Are you sure you want to delete the object '%v'?", engine.GAME_CONFIG.Objects[currentObjectID].ID), func(b bool) {
				if b {
					engine.GAME_CONFIG.Objects = append(engine.GAME_CONFIG.Objects[:currentObjectID], engine.GAME_CONFIG.Objects[currentObjectID+1:]...)
					setTabBasedOnId(currentTabID, mainContentBlock, window)
				}
			}, window).Show()
		}
	})
	var newObjectButton = widget.NewButton("Set Object", func() {
		var id = func() string {
			if len(idEntry.Text) > 0 {
				return idEntry.Text
			}
			return fmt.Sprintf("object_%v", len(engine.GAME_CONFIG.Objects))
		}()
		var shape = func() string {
			if shapeSelect.Selected == shapeOptions[0] {
				return shapeOptions[1]
			}
			return shapeSelect.Selected
		}()

		var pos = func() globals.Position {
			return globals.Position{X: helpers.CovertToInt(XPosEntry.Text), Y: helpers.CovertToInt(YPosEntry.Text)}
		}()

		var color = func() string {
			var t = colorEntry.Text
			if len(t) > 0 {
				return t
			}

			return "#ffffff"
		}()

		var size = func() globals.Size {
			return globals.Size{X: helpers.CovertToInt(XSizeEntry.Text), Y: helpers.CovertToInt(YSizeEntry.Text)}
		}()

		var objects = engine.GAME_CONFIG.Objects
		var object = globals.GameObject{ID: id, Shape: shape, Size: size, Pos: pos, Color: color, IsBody: isBodyCheck.Checked, HasArea: isAreaCheck.Checked, IsStatic: isStaticCheck.Checked}

		var canUpdate bool = true
		for i, obj := range objects {
			if obj.ID == object.ID {
				engine.GAME_CONFIG.Objects[i] = object

				canUpdate = false
				return
			}
		}

		if canUpdate {
			engine.GAME_CONFIG.Objects = append(engine.GAME_CONFIG.Objects, object)
		}
		setTabBasedOnId(currentTabID, mainContentBlock, window)

	})
	newObjectButton.Importance = widget.HighImportance

	objectList.OnSelected = func(id widget.ListItemID) {
		currentObjectID = id
		var c = engine.GAME_CONFIG.Objects[id]
		idEntry.SetText(c.ID)
		shapeSelect.SetSelected(c.Shape)
		XPosEntry.SetText(strconv.Itoa(c.Pos.X))
		YPosEntry.SetText(strconv.Itoa(c.Pos.Y))
		XSizeEntry.SetText(strconv.Itoa(c.Size.X))
		YSizeEntry.SetText(strconv.Itoa(c.Size.Y))
		colorEntry.SetText(c.Color)
		isBodyCheck.SetChecked(c.IsBody)
		isAreaCheck.SetChecked(c.HasArea)
		isStaticCheck.SetChecked(c.IsStatic)
	}

	var mainContent = container.NewVBox(container.NewCenter(container.NewHBox(newObjectButton, deleteObjectButton)), idShapeRow, positionRow, sizeRow, colorRow, areaBodyKeymapRow)
	var layoutSplit = container.NewHSplit(mainContent, objectList)
	layoutSplit.SetOffset(0.8)
	return container.New(layout.NewGridLayout(1), layoutSplit)
}

func functionContentTab() *fyne.Container {
	var textLabel = widget.NewLabel("Function tab")
	return container.NewHBox(textLabel)
}

func PhysicsContentTab() *fyne.Container {
	var gravityLabel = widget.NewLabel("Gravity: ")
	var gravityEntry = widget.NewEntry()
	gravityEntry.SetText(strconv.Itoa(engine.GAME_CONFIG.Gravity))

	var gravityContainer = container.NewBorder(nil, nil, gravityLabel, nil, gravityEntry)
	var applyButton = widget.NewButton("", func() {})
	applyButton = widget.NewButton("Apply", func() {

		go func() {
			engine.GAME_CONFIG.Gravity = helpers.CovertToInt(gravityEntry.Text)
			fyne.Do(func() {
				applyButton.SetText("Applied!")
			})

			time.Sleep(time.Second * 2)
			fyne.Do(func() {
				applyButton.SetText("Apply")
			})

		}()

	})
	return container.NewVBox(gravityContainer, applyButton)
}
