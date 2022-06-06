package gui

import (
	"givematerial/givematlib"
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var selectedLanguage string
var languageTableFilter *gtk.TreeModelFilter
var textTableModel *gtk.ListStore

const (
	COLUMN_TITLE = iota
	COLUMN_LANGUAGE
	COLUMN_UNKNOWN_VOCABS
)

func createTextsTable() (*gtk.TreeView, *gtk.ListStore) {
	treeView, _ := gtk.TreeViewNew()
	listStore, _ := gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_INT)
	textTableModel = listStore

	filter, _ := listStore.FilterNew(nil)
	filter.SetVisibleFunc(filterByLanguage)
	sort, _ := gtk.TreeModelSortNew(filter)
	languageTableFilter = filter

	treeView.SetModel(sort)
	treeView.AppendColumn(createColumn("Title", COLUMN_TITLE))
	treeView.AppendColumn(createColumn("Language", COLUMN_LANGUAGE))
	treeView.AppendColumn(createColumn("Unknown Vocabulary", COLUMN_UNKNOWN_VOCABS))

	// Just for testing
	combinedColumn, _ := gtk.TreeViewColumnNew()
	combinedColumn.SetTitle("Two fields")
	firstRenderer, _ := gtk.CellRendererTextNew()
	secondRenderer, _ := gtk.CellRendererTextNew()
	combinedColumn.PackStart(firstRenderer, false)
	combinedColumn.PackStart(secondRenderer, false)
	combinedColumn.AddAttribute(firstRenderer, "text", 1)
	combinedColumn.AddAttribute(secondRenderer, "text", 0)
	treeView.AppendColumn(combinedColumn)

	selection, _ := treeView.GetSelection()
	selection.SetMode(gtk.SELECTION_SINGLE)

	return treeView, listStore
}

func updateLanguagesTable(textData *gtk.ListStore) {
	textData.Clear()

	texts, err := givematlib.ListTexts()
	if err != nil {
		panic(err)
	}

	learnablesStatus := givematlib.StatusCacheNew()
	for _, textId := range texts {
		text, err := givematlib.LoadText(textId)
		if err != nil {
			log.Panic("Could not load text:", err)
		}

		knownLearnables, _ := learnablesStatus.ReadLearnableStatus(text.Language)

		iter := textData.Append()
		textData.Set(iter, []int{0, 1, 2}, []interface{}{
			text.Title,
			text.Language,
			len(text.Unknown(knownLearnables)),
		})
	}
}

func createColumn(title string, id int) *gtk.TreeViewColumn {
	cellRenderer, err := gtk.CellRendererTextNew()
	if err != nil {
		log.Fatal("Unable to create text cell renderer:", err)
	}

	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "text", id)
	column.SetSortColumnID(id)
	if err != nil {
		log.Fatal("Unable to create cell column:", err)
	}

	return column
}

func filterByLanguage(model *gtk.TreeModel, iter *gtk.TreeIter) bool {
	if selectedLanguage != "" {
		v, _ := model.GetValue(iter, COLUMN_LANGUAGE)
		gv, _ := v.GoValue()
		return gv.(string) == selectedLanguage
	} else {
		return true
	}
}
