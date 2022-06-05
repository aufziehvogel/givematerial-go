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

	b, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	win.Add(b)
	b.SetVExpand(true)

	// TODO: Re-generate table when learnables data has been updated
	learnablesStatus := givematlib.StatusCacheNew()

	texts, err := givematlib.ListTexts()
	if err != nil {
		panic(err)
	}

	languages := make(map[string]struct{})
	for _, textId := range texts {
		text, err := givematlib.LoadText(textId)
		if err != nil {
			log.Panic("Could not load text:", err)
		}

		languages[text.Language] = struct{}{}
		knownLearnables, _ := learnablesStatus.ReadLearnableStatus(text.Language)

		iter := listStore.Append()
		listStore.Set(iter, []int{0, 1, 2}, []interface{}{
			text.Title,
			text.Language,
			len(text.Unknown(knownLearnables)),
		})
	}

	scrollableTreelist, _ := gtk.ScrolledWindowNew(nil, nil)
	scrollableTreelist.Add(treeView)
	scrollableTreelist.SetVExpand(true)

	bSelection, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	for language, _ := range languages {
		button, _ := gtk.ButtonNewWithLabel(language)

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
