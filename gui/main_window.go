package gui

import (
	"givematerial/givematlib"
	"log"

	"github.com/gotk3/gotk3/gtk"
)

// Initialize GTK without parsing any command line arguments.

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

	texts, err := givematlib.ListTexts()
	if err != nil {
		panic(err)
	}
	for _, textId := range texts {
		text, err := givematlib.LoadText(textId)
		if err != nil {
			log.Panic("Could not load text:", err)
		}
		iter := listStore.Append()
		listStore.Set(iter, []int{0, 1, 2}, []interface{}{
			text.Title,
			text.Language,
			text.Unknown([]string{}),
		})
	}

	scrollableTreelist, _ := gtk.ScrolledWindowNew(nil, nil)
	scrollableTreelist.Add(treeView)
	scrollableTreelist.SetVExpand(true)

	bSelection, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	for _, language := range []string{"en", "ja", "es"} {
		button, _ := gtk.ButtonNewWithLabel(language)

		button.Connect("clicked", func(obj *gtk.Button) {
			label, _ := obj.GetLabel()
			selectedLanguage = label
			languageTableFilter.Refilter()
		})
		bSelection.Add(button)
	}
	b.PackStart(createMenuBar(), false, false, 0)
	b.PackStart(bSelection, false, false, 0)
	b.PackStart(scrollableTreelist, true, true, 0)

	// Set the default window size.
	win.SetDefaultSize(800, 600)

	// Recursively show all widgets contained in this window.
	win.ShowAll()

	// Begin executing the GTK main loop.  This blocks until
	// gtk.MainQuit() is run.
	gtk.Main()
}
