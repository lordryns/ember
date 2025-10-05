package filemenu

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func FileMenuWindow(window fyne.Window, projectPath *widget.Label, windowSidebar *widget.List) *dialog.CustomDialog {
	var customDialog *dialog.CustomDialog
	customDialog = dialog.NewCustom("File", "Close", container.NewBorder(nil, nil, nil, projectSelectContainer(window, projectPath, customDialog, windowSidebar)), window)
	customDialog.Show()

	return customDialog
}
