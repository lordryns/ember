package helpers

import (
	"fyne.io/fyne/v2/widget"
)

var SidebarItems = []string{"Defaults", "Sprites", "Object", "Functions", "Physics"}

func SetSidebarContent(list *widget.List, items []string) {
	SidebarItems = items
	list.Refresh()
}
