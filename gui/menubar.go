package gui

import "github.com/gotk3/gotk3/gtk"

func createMenuBar() *gtk.MenuBar {
	menu, _ := gtk.MenuBarNew()
	fileMenu, _ := gtk.MenuNew()
	fileMenuItem, _ := gtk.MenuItemNewWithLabel("File")
	fileMenuItem.SetSubmenu(fileMenu)

	quitItem, _ := gtk.MenuItemNewWithLabel("Quit")
	fileMenu.Append(quitItem)

	quitItem.Connect("activate", func() {
		gtk.MainQuit()
	})

	menu.Append(fileMenuItem)

	return menu
}
