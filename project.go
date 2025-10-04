package main

import (
	"errors"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func projectSelectContainer(window fyne.Window, projectPath *string) *fyne.Container {
	var dirLabel = widget.NewLabel("No project is open...yet")
	var openButton = widget.NewButton("open", func() {
		var folderDialog = dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, window)
				return
			}

			if uri != nil {
				if err := isValidProject(uri.Path()); err != nil {
					dialog.ShowError(err, window)
					return
				}
				*projectPath = uri.Path()
				dirLabel.SetText(*projectPath)
			}
		}, window)

		folderDialog.Show()
	})
	var createButton = widget.NewButton("create", func() {

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

					if err := createProject(path, name); err != nil {
						dialog.ShowError(err, window)
						return
					}

					*projectPath = filepath.Join(path, name)
					dirLabel.SetText(*projectPath)
					dialog.ShowInformation("Project info", "Project created successfully!", window)
				}
			}, window)

		createFolderForm.Resize(fyne.NewSize(400, 200))
		createFolderForm.Show()
	})

	return container.NewBorder(nil, nil, nil, container.NewHBox(openButton, createButton), dirLabel)
}
