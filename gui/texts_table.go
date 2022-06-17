package gui

import (
	"givematerial/givematlib"
	"log"
	"strings"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var selectedLanguage string
var languageTableFilter *gtk.TreeModelFilter
var textTableModel *gtk.ListStore

const (
	COLUMN_ID = iota
	COLUMN_TITLE
	COLUMN_LANGUAGE
	COLUMN_UNKNOWN_VOCAB_COUNT
	COLUMN_UNKNOWN_VOCAB
	COLUMN_KNOWN_VOCAB_PERCENT
)

func createTextsTable() (*gtk.TreeView, *gtk.ListStore) {
	treeView, _ := gtk.TreeViewNew()

	listStore, _ := gtk.ListStoreNew(
		glib.TYPE_STRING,
		glib.TYPE_STRING,
		glib.TYPE_STRING,
		glib.TYPE_INT,
		glib.TYPE_STRING,
		glib.TYPE_DOUBLE,
	)
	textTableModel = listStore

	filter, _ := listStore.FilterNew(nil)
	filter.SetVisibleFunc(filterByLanguage)
	sort, _ := gtk.TreeModelSortNew(filter)
	languageTableFilter = filter

	treeView.SetModel(sort)
	treeView.AppendColumn(createColumn("Title", COLUMN_TITLE))
	treeView.AppendColumn(createColumn("Language", COLUMN_LANGUAGE))
	treeView.AppendColumn(createColumn("Unknown Vocabulary", COLUMN_UNKNOWN_VOCAB_COUNT))

	progressRenderer, _ := gtk.CellRendererProgressNew()
	progressColumn, _ := gtk.TreeViewColumnNewWithAttribute(
		"Vocabulary Known",
		progressRenderer,
		"value",
		COLUMN_KNOWN_VOCAB_PERCENT,
	)
	progressColumn.SetSortColumnID(COLUMN_KNOWN_VOCAB_PERCENT)
	treeView.AppendColumn(progressColumn)

	treeView.Connect("row-activated", func(tv *gtk.TreeView, path *gtk.TreePath, column *gtk.TreeViewColumn) {
		iter, _ := sort.GetIter(path)
		v, _ := sort.GetValue(iter, COLUMN_ID)
		gv, _ := v.GoValue()
		textId := gv.(string)

		v2, _ := sort.GetValue(iter, COLUMN_UNKNOWN_VOCAB)
		gv2, _ := v2.GoValue()
		unknownVocab := gv2.(string)

		win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
		win.SetTitle("Popup")

		box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
		win.Add(box)

		unknownVocabLabel, _ := gtk.LabelNew("foo")
		box.Add(unknownVocabLabel)
		unknownVocabView, _ := gtk.TextViewNew()
		unknownVocabView.SetWrapMode(gtk.WRAP_WORD)
		bufferVocab, _ := unknownVocabView.GetBuffer()
		bufferVocab.SetText(unknownVocab)
		box.Add(unknownVocabView)

		scrollView, _ := gtk.ScrolledWindowNew(nil, nil)
		box.PackStart(scrollView, true, true, 10)

		textView, _ := gtk.TextViewNew()
		textView.SetEditable(false)
		buffer, _ := textView.GetBuffer()
		text, _ := givematlib.LoadText(textId)
		buffer.SetText(text.Fulltext)

		scrollView.Add(textView)

		win.ShowAll()
	})

	selection, _ := treeView.GetSelection()
	selection.SetMode(gtk.SELECTION_SINGLE)

	return treeView, listStore
}

func updateTextsTable(textData *gtk.ListStore) {
	textData.Clear()

	texts, err := givematlib.ListTexts()
	if err != nil {
		panic(err)
	}

	learnablesStatus := givematlib.NewStatusCache()
	for _, textId := range texts {
		text, err := givematlib.LoadText(textId)
		if err != nil {
			log.Panic("Could not load text:", err)
		}

		knownLearnables, _ := learnablesStatus.ReadLearnableStatus(text.Language)

		iter := textData.Append()
		numUnknown := len(text.Unknown(knownLearnables))
		numLearnables := len(text.Learnables)

		percentageKnown := 0
		if numLearnables != 0 {
			percentageKnown = 100 * (numLearnables - numUnknown) / numLearnables
		}
		textData.Set(iter, []int{0, 1, 2, 3, 4, 5}, []interface{}{
			textId,
			text.Title,
			text.Language,
			len(text.Unknown(knownLearnables)),
			strings.Join(text.Unknown(knownLearnables), ", "),
			percentageKnown,
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
