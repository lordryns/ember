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
	"slices"
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
	var fileButton = widget.NewButton("Project menu", func() {
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
		mainContent.Objects = []fyne.CanvasObject{functionContentTab(window)}
	case 4:
		mainContent.Objects = []fyne.CanvasObject{PhysicsContentTab()}
	case 5:
		mainContent.Objects = []fyne.CanvasObject{UpdatesContentTab(window)}

	default:
		mainContent.Objects = []fyne.CanvasObject{defaultContentTab(window, &engine.GAME_CONFIG)}
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

func defaultContentTab(window fyne.Window, config *globals.GameConfig) *fyne.Container {
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

	var backgroundEntry = widget.NewEntry()
	backgroundEntry.SetPlaceHolder("Set background color")
	backgroundEntry.SetText(engine.GAME_CONFIG.Color)
	var colorPicker = colorpicker.New(200, colorpicker.StyleHue)
	colorPicker.SetOnChanged(func(c color.Color) {
		backgroundEntry.SetText(helpers.ColorToHex(c))
	})

	var backgroundButton = widget.NewButton("Set Color", func() {
		dialog.NewCustom("Pick Color", "Dismiss", container.NewCenter(colorPicker), window).Show()
	})

	var backgroundContainer = container.NewGridWithColumns(3, widget.NewLabel("Background color:"), backgroundEntry, backgroundButton)
	//container.NewBorder(nil, nil, widget.NewLabel("Background color:"), backgroundButton, backgroundEntry)
	var applyButton = widget.NewButton("", func() {})
	applyButton = widget.NewButton("Apply", func() {
		engine.GAME_CONFIG.Title = titleEntry.Text
		engine.GAME_CONFIG.Color = backgroundEntry.Text
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

	var bottomStack = container.NewStack(container.NewBorder(nil, container.NewVBox(widget.NewSeparator(), widget.NewLabel(engine.PROJECT_PATH)), nil, nil, nil))
	return container.NewStack(container.NewVBox(titleContainer, geometryContainer, backgroundContainer, applyButton), bottomStack)
}

func spritesContentTab() *fyne.Container {
	var textLabel = widget.NewLabel("Sprites tab")
	return container.NewHBox(textLabel)
}

func keyMapListContainer(obi int) *widget.List {
	var items = engine.GAME_CONFIG.Objects[obi].KeyMap
	var list = widget.NewList(
		func() int { return len(items) },
		func() fyne.CanvasObject {
			return widget.NewButton("", nil)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			b := o.(*widget.Button)
			b.SetText(items[i].Key)
			b.OnTapped = func() {
				fmt.Println("Clicked:", items[i])
			}
		},
	)

	return list
}

type ArgStruct struct {
	Label *widget.Label
	Entry *widget.Entry
}

func objectContentTab(mainContentBlock *fyne.Container, window fyne.Window) *fyne.Container {
	var currentObjectID int = -1
	var tempKeyMap []globals.KeyMap

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
	shapeSelect.SetSelectedIndex(0)
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

	var keyMapMainContent = container.NewStack(widget.NewLabel("Nothing to see here..."))

	// update keymap display
	// this has been an absolute pain in the ass
	var updateKeyMapContent = func() {
		if len(tempKeyMap) > 0 {
			var list *widget.List
			list = widget.NewList(
				func() int { return len(tempKeyMap) },
				func() fyne.CanvasObject {
					return widget.NewButton("", nil)
				},
				func(i widget.ListItemID, o fyne.CanvasObject) {
					b := o.(*widget.Button)
					b.SetText(fmt.Sprintf("%s -> %s -> %s", tempKeyMap[i].Key, tempKeyMap[i].PressType, tempKeyMap[i].Func.ID))
					b.OnTapped = func() {
						tempKeyMap = append(tempKeyMap[:i], tempKeyMap[i+1:]...)
						list.Refresh()
					}
				},
			)
			keyMapMainContent.Objects = []fyne.CanvasObject{list}
		} else {
			keyMapMainContent.Objects = []fyne.CanvasObject{container.NewCenter(widget.NewLabel("Nothing to see here..."))}
		}
		keyMapMainContent.Refresh()
	}
	var parseFunctionIDs = func() []string {
		var s []string
		var fn = engine.GAME_CONFIG.Functions

		for _, f := range fn {
			s = append(s, f.ID)
		}

		return s

	}

	var inputSelect = widget.NewSelect(engine.ALL_INPUTS.Keyboard, func(s string) {})
	var funcSelect = widget.NewSelect(parseFunctionIDs(), func(s string) {})
	var pressTypeSelect = widget.NewSelect(engine.INPUT_PRESS_TYPE, func(s string) {})
	pressTypeSelect.SetSelectedIndex(0)
	var keyMapTopBar = container.NewBorder(nil, nil, widget.NewLabel("Tap to remove"), widget.NewButton("Add Key", func() {
		dialog.NewCustomConfirm("Add Key", "Add", "Dismiss",
			container.NewVBox(inputSelect, funcSelect, pressTypeSelect), func(b bool) {
				if b {
					var ii = inputSelect.SelectedIndex()
					var fi = funcSelect.SelectedIndex()
					var pt = pressTypeSelect.SelectedIndex()

					if ii > -1 && fi > -1 {
						var cur globals.KeyMap
						cur.Key = engine.ALL_INPUTS.Keyboard[ii]
						cur.PressType = engine.INPUT_PRESS_TYPE[pt]

						var argEntries []ArgStruct

						var curF = engine.GAME_CONFIG.Functions[fi]

						for _, args := range curF.Args {
							argEntries = append(argEntries, ArgStruct{Label: widget.NewLabel(fmt.Sprintf("%v:", args.ID)), Entry: widget.NewEntry()})
						}
						var funcContainer = container.NewVBox()
						for _, astr := range argEntries {
							funcContainer.Objects = append(funcContainer.Objects, container.NewBorder(nil, nil, astr.Label, nil, astr.Entry))
						}
						dialog.ShowCustomConfirm("Set Args", "Set", "Cancel", funcContainer, func(b bool) {

							for i := range argEntries {
								curF.Args[i].Value = argEntries[i].Entry.Text
							}

							cur.Func = curF
							tempKeyMap = append(tempKeyMap, cur)
							updateKeyMapContent()
						}, window)
					}
				}
			}, window).Show()
	}))
	var keyMapContainer = container.NewBorder(keyMapTopBar, nil, nil, nil, keyMapMainContent)
	var keyPressButton = widget.NewButton("Key map", func() {
		var keyMapDialog = dialog.NewCustom("Key Map", "Close", keyMapContainer, window)
		keyMapDialog.Resize(fyne.NewSize(400, 400))
		keyMapDialog.Show()
	})

	var massEntry = widget.NewEntry()
	massEntry.SetPlaceHolder("Set Object Mass")

	var gravityScaleEntry = widget.NewEntry()
	gravityScaleEntry.SetPlaceHolder("Set Gravity scale")

	var massScaleContainer = container.NewGridWithColumns(2, massEntry, gravityScaleEntry)

	var areaBodyKeymapRow = container.NewGridWithColumns(4, isBodyCheck, isAreaCheck, isStaticCheck, keyPressButton)

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
	deleteObjectButton.Importance = widget.DangerImportance
	var newObjectButton = widget.NewButton("Set Object", func() {
		var id = func() string {
			if len(idEntry.Text) > 0 {
				return helpers.RemoveWhiteSpaceAndIllegals(idEntry.Text)
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

		var mass = func() int {
			return helpers.CovertToInt(massEntry.Text)
		}()

		var gravityScale = func() int {
			return helpers.CovertToInt(gravityScaleEntry.Text)
		}()

		var objects = engine.GAME_CONFIG.Objects
		var object = globals.GameObject{ID: id, Shape: shape, Size: size,
			Pos: pos, Color: color, IsBody: isBodyCheck.Checked, HasArea: isAreaCheck.Checked, IsStatic: isStaticCheck.Checked, Mass: mass, GravityScale: gravityScale,
			KeyMap: tempKeyMap}

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
		tempKeyMap = engine.GAME_CONFIG.Objects[id].KeyMap
		updateKeyMapContent()
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
		massEntry.SetText(strconv.Itoa(c.Mass))
		gravityScaleEntry.SetText(strconv.Itoa(c.GravityScale))
		tempKeyMap = c.KeyMap
	}

	var mainContent = container.NewVBox(container.NewCenter(container.NewHBox(newObjectButton, deleteObjectButton)), idShapeRow, positionRow, sizeRow, colorRow, areaBodyKeymapRow, massScaleContainer)
	var layoutSplit = container.NewHSplit(mainContent, objectList)
	layoutSplit.SetOffset(0.8)
	return container.New(layout.NewGridLayout(1), layoutSplit)
}

func AddFunctionDialog(window fyne.Window, tType *string, functionsList *widget.List, func_name string, func_args []globals.Arg, func_src string) {
	var tempArgListID int = -1
	var tempArgList = func_args
	var tempCode = func_src

	var addFuncIDEntry = widget.NewEntry()
	addFuncIDEntry.SetPlaceHolder("ID...")
	addFuncIDEntry.SetText(func_name)

	var addFuncIdContainer = container.NewBorder(nil, nil, widget.NewLabel("Enter ID: "), nil, addFuncIDEntry)

	var argFuncLabel = widget.NewLabel("Arg:")
	var addFuncArgsEntry = widget.NewEntry()
	addFuncArgsEntry.SetPlaceHolder("Args...")

	var addFuncArgButton = widget.NewButton("Add Arg", func() {
		var t = addFuncArgsEntry.Text
		if len(t) > 0 {
			var temp = helpers.RemoveWhiteSpaceAndIllegals(t)
			var getArgID = func() []string {
				var l []string
				for _, te := range tempArgList {
					l = append(l, te.ID)
				}

				return l
			}
			if slices.Contains(getArgID(), temp) {
				dialog.ShowError(fmt.Errorf("Argument '%v' already exists in this function!", temp), window)
				return
			}

			var tEntry = widget.NewSelect(engine.CUSTOM_TYPES, func(s string) {
				*tType = s
			})
			tEntry.SetSelectedIndex(0)
			var tc = dialog.NewCustomConfirm("Set Type", "Set Type", "Cancel", tEntry, func(b bool) {
				if b {
					tempArgList = append(tempArgList, globals.Arg{ID: temp, Type: *tType})
					argFuncLabel.SetText(fmt.Sprintf("Arg(%v):", len(tempArgList)))
					addFuncArgsEntry.SetText("")

				}

			}, window)
			tc.Show()
		}
	})
	var showFuncArgButton = widget.NewButton("Show Args", func() {
		var list = widget.NewList(func() int { return len(tempArgList) }, func() fyne.CanvasObject { return widget.NewLabel("") }, func(lii widget.ListItemID, co fyne.CanvasObject) {
			// thanks gpt bro for this snippet
			idx := len(tempArgList) - 1 - lii
			cur := tempArgList[idx]
			co.(*widget.Label).SetText(fmt.Sprintf("%v:%v", cur.ID, cur.Type))
		})

		list.OnSelected = func(id widget.ListItemID) {
			tempArgListID = id
		}

		var deleteButton = widget.NewButton("Delete", func() {
			if tempArgListID > -1 && len(tempArgList) > 0 {
				tempArgList = append(tempArgList[:tempArgListID], tempArgList[tempArgListID+1:]...)
				tempArgListID = -1
				list.Refresh()
			}
		})
		deleteButton.Importance = widget.DangerImportance
		var d = dialog.NewCustom("Arguments", "Close", container.NewBorder(container.NewVBox(deleteButton), nil, nil, nil, list), window)
		d.Resize(fyne.NewSize(200, 300))
		d.Show()
	})
	var addFuncArgsContainer = container.NewBorder(nil, nil, argFuncLabel, container.NewHBox(addFuncArgButton, showFuncArgButton), addFuncArgsEntry)

	var addSrcButton = widget.NewButton("Source Code", func() {
		var editor = widget.NewMultiLineEntry()
		var template string
		for _, arg := range tempArgList {
			template += fmt.Sprintf("let  _%v = %v; // type: %v\n", arg.ID, arg.ID, arg.Type)
		}
		if len(tempCode) < 1 {
			editor.SetText(template)
		} else {
			editor.SetText(tempCode)
		}
		var d = dialog.NewCustomConfirm("Function Code", "Add", "Dismiss", editor, func(b bool) {
			if b {
				tempCode = editor.Text
			}
		}, window)
		d.Resize(fyne.NewSize(400, 400))
		d.Show()
	})

	var d = dialog.NewCustomConfirm("New Function", "Add Function", "Cancel", container.NewVBox(addFuncIdContainer, addFuncArgsContainer, addSrcButton), func(b bool) {
		if b {
			var parseFunc = func() string {
				var t = addFuncIDEntry.Text
				if len(t) < 1 {
					return fmt.Sprintf("func_%v", len(engine.GAME_CONFIG.Functions))
				}
				return helpers.RemoveWhiteSpaceAndIllegals(t)
			}
			var newFunc = globals.GameFunc{ID: parseFunc(), Args: tempArgList, Src: tempCode}
			var checkIfIDAlreadyExists = func() (bool, int) {
				for i, fn := range engine.GAME_CONFIG.Functions {
					if fn.ID == newFunc.ID {
						return true, i
					}

				}

				return false, -1
			}

			var exists, index = checkIfIDAlreadyExists()
			if !exists {
				engine.GAME_CONFIG.Functions = append(engine.GAME_CONFIG.Functions, newFunc)
			} else {
				engine.GAME_CONFIG.Functions[index] = newFunc
			}
			functionsList.Refresh()
		}
		tempArgList = []globals.Arg{}
		argFuncLabel.SetText("Arg:")
	}, window)
	d.Resize(fyne.NewSize(400, 200))
	d.Show()

}

func functionContentTab(window fyne.Window) *fyne.Container {
	var tType string
	var functionsList *widget.List
	var currentFuncID = -1

	functionsList = widget.NewList(func() int { return len(engine.GAME_CONFIG.Functions) }, func() fyne.CanvasObject { return helpers.NewClickableLabel("", func() {}) }, func(lii widget.ListItemID, co fyne.CanvasObject) {
		var e = engine.GAME_CONFIG.Functions[lii]
		var b = co.(*helpers.ClickableLabel)

		// Style for selection
		if lii == currentFuncID {
			b.TextStyle = fyne.TextStyle{Bold: true}

			b.SetText(fmt.Sprintf("+ %v:		%v", e.ID, e.Args))
		} else {
			b.TextStyle = fyne.TextStyle{}

			b.SetText(fmt.Sprintf("%v:		%v", e.ID, e.Args))
		}
		b.Refresh()

		b.OnTapped = func() {
			currentFuncID = lii
			functionsList.Refresh()
			AddFunctionDialog(window, &tType, functionsList, e.ID, e.Args, e.Src)
		}
	})
	var deleteFuncButton = widget.NewButton("Delete Function", func() {
		if currentFuncID > -1 {
			engine.GAME_CONFIG.Functions = append(engine.GAME_CONFIG.Functions[:currentFuncID], engine.GAME_CONFIG.Functions[currentFuncID+1:]...)
			functionsList.Refresh()
		}
	})
	var addFuncButton = widget.NewButton("Open Menu", func() {

		dialog.ShowCustom("Function Menu", "Close", container.NewHBox(widget.NewButton("Add Function", func() {
			AddFunctionDialog(window, &tType, functionsList, "", []globals.Arg{}, "")
		}), deleteFuncButton), window)

	})

	deleteFuncButton.Importance = widget.DangerImportance
	addFuncButton.Importance = widget.HighImportance

	var mainContainer = container.NewBorder(container.NewCenter(addFuncButton), nil, nil, nil, functionsList)
	return mainContainer
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

// claude assisted this fix
func UpdatesContentTab(window fyne.Window) *fyne.Container {
	var updateList *widget.List
	var extractStringFromFuncStruct = func() []string {
		var t []string
		for _, fn := range engine.GAME_CONFIG.Functions {
			t = append(t, fn.ID)
		}
		return t
	}
	var updateEntry = widget.NewSelect(extractStringFromFuncStruct(), func(s string) {})
	var addUpdateButton = widget.NewButton("Add to Update", func() {
		dialog.ShowCustomConfirm("Add to Update", "Add", "Cancel", updateEntry, func(b bool) {
			if b {
				var ui = updateEntry.SelectedIndex()
				if ui > -1 {
					var cur = engine.GAME_CONFIG.Functions[ui]
					var curCopy = globals.GameFunc{
						ID:   cur.ID,
						Src:  cur.Src,
						Args: make([]globals.Arg, len(cur.Args)),
					}
					// Deep copy the args
					for i, arg := range cur.Args {
						curCopy.Args[i] = globals.Arg{
							ID:    arg.ID,
							Type:  arg.Type,
							Value: arg.Value,
						}
					}

					if len(curCopy.Args) > 0 {
						var argEntries []ArgStruct
						for _, args := range curCopy.Args {
							argEntries = append(argEntries, ArgStruct{
								Label: widget.NewLabel(fmt.Sprintf("%v:", args.ID)),
								Entry: widget.NewEntry(),
							})
						}
						var funcContainer = container.NewVBox()
						for _, astr := range argEntries {
							funcContainer.Objects = append(funcContainer.Objects,
								container.NewBorder(nil, nil, astr.Label, nil, astr.Entry))
						}
						dialog.ShowCustomConfirm("Args", "Confirm", "Cancel", funcContainer, func(b bool) {
							if b {
								for i, entry := range argEntries {
									curCopy.Args[i].Value = entry.Entry.Text
								}
								engine.GAME_CONFIG.Update = append(engine.GAME_CONFIG.Update, curCopy)
								fmt.Println(len(engine.GAME_CONFIG.Update))
								updateList.Refresh()
							}
						}, window)
					} else {
						// If no args, add directly
						engine.GAME_CONFIG.Update = append(engine.GAME_CONFIG.Update, curCopy)
						updateList.Refresh()
					}
				}
			}
		}, window)
	})
	addUpdateButton.Importance = widget.HighImportance

	updateList = widget.NewList(
		func() int {
			return len(engine.GAME_CONFIG.Update)
		},
		func() fyne.CanvasObject {
			return widget.NewButton("", func() {})
		},
		func(lii widget.ListItemID, co fyne.CanvasObject) {
			var update = engine.GAME_CONFIG.Update[lii]
			co.(*widget.Button).SetText(fmt.Sprintf("%v -> %v", update.ID, update.Args))
			index := lii
			co.(*widget.Button).OnTapped = func() {
				if index < len(engine.GAME_CONFIG.Update) {
					engine.GAME_CONFIG.Update = append(
						engine.GAME_CONFIG.Update[:index],
						engine.GAME_CONFIG.Update[index+1:]...)
					updateList.Refresh()
				}
			}
		})

	return container.NewBorder(container.NewCenter(addUpdateButton), nil, nil, nil, updateList)
}
