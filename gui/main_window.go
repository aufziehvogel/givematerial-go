package gui

import (
	"givematerial/givematlib"
	"log"

	"github.com/gotk3/gotk3/gtk"
)

var applicationStatusBar *gtk.Statusbar

func Init(config *givematlib.ApplicationConfig) {
	gtk.Init(nil)

	// Create a new toplevel window, set its title, and connect it to the
	// "destroy" signal to exit the GTK main loop when it is destroyed.
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle("Simple Example")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	treeView, listStore := createTextsTable()
	updateTextsTable(listStore)

	b, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	win.Add(b)
	b.SetVExpand(true)

	languages, err := loadLanguages()
	if err != nil {
		log.Panic("Could not load languages", err)
	}

	scrollableTreelist, _ := gtk.ScrolledWindowNew(nil, nil)
	scrollableTreelist.Add(treeView)
	scrollableTreelist.SetVExpand(true)

	bSelection, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	for language := range languages {
		button, _ := gtk.ButtonNewWithLabel(language.ShortCode())

		button.Connect("clicked", func(obj *gtk.Button) {
			label, _ := obj.GetLabel()
			selectedLanguage = label
			languageTableFilter.Refilter()
		})
		bSelection.Add(button)
	}
	menuBar, err := createMenuBar(config)
	if err != nil {
		log.Panic("Could not create menu bar", err)
	}

	statusBar, err := gtk.StatusbarNew()
	if err != nil {
		log.Panic("Could not create status bar", err)
	}
	applicationStatusBar = statusBar

	b.PackStart(menuBar, false, false, 0)
	b.PackStart(bSelection, false, false, 0)
	b.PackStart(scrollableTreelist, true, true, 0)
	b.PackStart(statusBar, false, false, 0)

	// Set the default window size.
	win.SetDefaultSize(800, 600)

	// Recursively show all widgets contained in this window.
	win.ShowAll()

	// Begin executing the GTK main loop.  This blocks until
	// gtk.MainQuit() is run.
	gtk.Main()
}

func loadLanguages() (map[givematlib.Language]struct{}, error) {
	texts, err := givematlib.ListTexts()
	if err != nil {
		return nil, err
	}

	languages := make(map[givematlib.Language]struct{})
	for _, textId := range texts {
		text, err := givematlib.LoadText(textId)
		if err != nil {
			return nil, err
		}

		language := givematlib.MakeLanguageFromShortCode(text.Language)
		languages[language] = struct{}{}
	}

	return languages, nil
}
