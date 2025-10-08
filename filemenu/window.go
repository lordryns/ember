package filemenu

import (
	"ember/engine"
	"ember/helpers"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func FileMenuWindow(window fyne.Window, projectPath *widget.Label, mainContentBlock *fyne.Container, refreshWindow func()) *dialog.CustomDialog {
	var customDialog *dialog.CustomDialog
	var fileRow = container.NewGridWithColumns(2, projectSelectContainer(window, projectPath, customDialog, &engine.GAME_CONFIG, mainContentBlock, refreshWindow))

	var portEntry = widget.NewEntry()
	portEntry.SetPlaceHolder("Set Port")
	portEntry.SetText(strconv.Itoa(engine.PORT))
	var portButton = widget.NewButton("Set Port", func() {
		engine.PORT = helpers.ValidatePort(portEntry.Text)
	})

	var portRow = container.NewBorder(nil, nil, widget.NewLabel("PORT:"), portButton, portEntry)
	var mainContainer = container.NewVBox(fileRow, portRow)
	customDialog = dialog.NewCustom("Project", "Close", mainContainer, window)

	customDialog.Show()

	return customDialog
}
