package helpers

import (
	"fyne.io/fyne/v2/widget"
)

var SidebarItems = []string{"Defaults", "Sprites", "Objects", "Functions", "Physics"}

func SetSidebarContent(list *widget.List, items []string) {
	SidebarItems = items
	list.Refresh()
}
