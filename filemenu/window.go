package filemenu

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
)

func FileMenuWindow(window fyne.Window, projectPath *string) *dialog.CustomDialog {
	var customDialog = dialog.NewCustom("File", "Close", container.NewBorder(nil, nil, nil, projectSelectContainer(window, projectPath)), window)
	customDialog.Show()

	return customDialog
}
