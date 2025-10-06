package filemenu

import (
	"ember/engine"
	"ember/helpers"
	"errors"
	"fmt"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func projectSelectContainer(window fyne.Window, projectPath *widget.Label, directWindow *dialog.CustomDialog, config *engine.GameConfig, mainContentBlock *fyne.Container, refreshWindow func()) *fyne.Container {
	var openButton = widget.NewButton("Open Existing", func() {
		var folderDialog = dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, window)
				return
			}

			if uri != nil {
				if err := helpers.IsValidProject(uri.Path()); err != nil {
					dialog.ShowError(err, window)
					return
				}
				var conf, err = engine.LoadConfig(uri.Path())
				if err != nil {
					dialog.ShowError(fmt.Errorf("Unable to load game configuration! err: %s", err), window)
				}

				*config = conf
				projectPath.SetText(uri.Path())
				mainContentBlock.Objects[0].Refresh()
				refreshWindow()
			}
		}, window)

		folderDialog.Show()
	})
	var createButton = widget.NewButton("New Project", func() {

		var newFolderEntry = widget.NewEntry()
		newFolderEntry.SetPlaceHolder("Set project path")
		newFolderEntry.Disable()
		var newFolderButton = widget.NewButton("Browse", func() {
			var folderDialog = dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
				if err != nil {
					dialog.ShowError(err, window)
				}

				if uri != nil {
					newFolderEntry.SetText(uri.Path())
				}

			}, window)

			folderDialog.Show()
		})

		var newFolderNameEntry = widget.NewEntry()
		newFolderNameEntry.SetPlaceHolder("Enter project title")

		var createFolderForm = dialog.NewForm("New Project",
			"Create", "Dismiss",
			[]*widget.FormItem{widget.NewFormItem("Title", newFolderNameEntry),
				widget.NewFormItem("Path", container.NewBorder(nil, nil, nil, newFolderButton, newFolderEntry)),
			}, func(ok bool) {
				if ok {
					var name = newFolderNameEntry.Text
					var path = newFolderEntry.Text

					if len(name) < 1 {
						dialog.ShowError(errors.New("Project name cannot be left empty!"), window)
						return
					}

					if len(path) < 1 {
						dialog.ShowError(errors.New("Select a valid path!"), window)
						return
					}

					if err := helpers.CreateProject(path, name, config); err != nil {
						dialog.ShowError(err, window)
						return
					}

					projectPath.SetText(filepath.Join(path, name))
					dialog.ShowInformation("Project info", "Project created successfully!", window)
					mainContentBlock.Refresh()
					refreshWindow()
				}
			}, window)

		createFolderForm.Resize(fyne.NewSize(400, 200))
		createFolderForm.Show()
	})

	return container.NewBorder(nil, nil, nil, container.NewHBox(createButton, openButton), layout.NewSpacer())
}
