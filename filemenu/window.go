package filemenu

import (
	"ember/engine"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func FileMenuWindow(window fyne.Window, projectPath *widget.Label, mainContentBlock *fyne.Container, refreshWindow func()) *dialog.CustomDialog {
	var customDialog *dialog.CustomDialog
	customDialog = dialog.NewCustom("File", "Close", container.NewBorder(nil, nil, nil, projectSelectContainer(window, projectPath, customDialog, &engine.GAME_CONFIG, mainContentBlock, refreshWindow)), window)
	customDialog.Show()

	return customDialog
}
