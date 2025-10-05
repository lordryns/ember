package helpers

import "fyne.io/fyne/v2/widget"

var SidebarItems = []string{}

func SetSidebarContent(list *widget.List, items []string) {
	SidebarItems = items
	list.Refresh()
}
