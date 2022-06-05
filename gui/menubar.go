package gui

import (
	"givematerial/givematlib"
	"log"

	"github.com/gotk3/gotk3/gtk"
)

func createMenuBar(config *givematlib.ApplicationConfig) (*gtk.MenuBar, error) {
	menu, _ := gtk.MenuBarNew()
	_, fileMenuItem, err := createFileMenu()
	if err != nil {
		return nil, err
	}
	menu.Append(fileMenuItem)

	_, actionsMenuItem, err := createActionMenu(config)
	if err != nil {
		return nil, err
	}
	menu.Append(actionsMenuItem)

	return menu, nil
}

func createFileMenu() (*gtk.Menu, *gtk.MenuItem, error) {
	menu, menuItem, err := createMenu("File")
	if err != nil {
		return nil, nil, err
	}

	quitItem, err := gtk.MenuItemNewWithLabel("Quit")
	if err != nil {
		return nil, nil, err
	}
	quitItem.Connect("activate", func() {
		gtk.MainQuit()
	})
	menu.Append(quitItem)

	return menu, menuItem, nil
}

func createActionMenu(config *givematlib.ApplicationConfig) (*gtk.Menu, *gtk.MenuItem, error) {
	menu, menuItem, err := createMenu("Actions")
	if err != nil {
		return nil, nil, err
	}

	updateAnkiItem, err := gtk.MenuItemNewWithLabel("Update Anki data")
	if err != nil {
		return nil, nil, err
	}
	updateAnkiItem.Connect("activate", func() {
		log.Print("Updating Anki data")

		statusMessages := make(chan string)
		go givematlib.UpdateAnki(
			config.AnkiFile,
			config.AnkiDecksForLanguage,
			statusMessages,
			true,
		)
		go displayStatusMessages(statusMessages, applicationStatusBar)
	})
	menu.Append(updateAnkiItem)

	updateWanikaniItem, err := gtk.MenuItemNewWithLabel("Update Wanikani data")
	if err != nil {
		return nil, nil, err
	}
	updateWanikaniItem.Connect("activate", func() {
		log.Print("Updating Wanikani data")

		statusMessages := make(chan string)
		go givematlib.UpdateWanikani(
			config.WanikaniApiToken,
			statusMessages,
			true,
		)
		go displayStatusMessages(statusMessages, applicationStatusBar)
	})
	menu.Append(updateWanikaniItem)

	return menu, menuItem, nil
}

func createMenu(label string) (*gtk.Menu, *gtk.MenuItem, error) {
	menu, err := gtk.MenuNew()
	if err != nil {
		return nil, nil, err
	}

	menuItem, err := gtk.MenuItemNewWithLabel(label)
	if err != nil {
		return nil, nil, err
	}

	menuItem.SetSubmenu(menu)
	return menu, menuItem, nil
}
