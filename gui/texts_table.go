package gui

import (
	"givematerial/givematlib"
	"log"
	"strings"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var selectedLanguage string

const (
	COLUMN_ID = iota
	COLUMN_TITLE
	COLUMN_LANGUAGE
	COLUMN_UNKNOWN_VOCAB_COUNT
	COLUMN_UNKNOWN_VOCAB
	COLUMN_KNOWN_VOCAB_PERCENT
)

type TextTableModel struct {
	rawListStore        *gtk.ListStore
	filterableListStore *gtk.TreeModelFilter
	sortableListStore   *gtk.TreeModelSort
}

type TextTableView struct {
	treeView           *gtk.TreeView
	scrollableTreeView *gtk.ScrolledWindow
}

type TextTableController struct {
	model *TextTableModel
	view  *TextTableView
}

func newTextTableModel() *TextTableModel {
	listStore, _ := gtk.ListStoreNew(
		glib.TYPE_STRING,
		glib.TYPE_STRING,
		glib.TYPE_STRING,
		glib.TYPE_INT,
		glib.TYPE_STRING,
		glib.TYPE_DOUBLE,
	)

	filter, _ := listStore.FilterNew(nil)
	filter.SetVisibleFunc(filterByLanguage)
	sort, _ := gtk.TreeModelSortNew(filter)

	return &TextTableModel{
		rawListStore:        listStore,
		filterableListStore: filter,
		sortableListStore:   sort,
	}
}

func newTextTableView(model *TextTableModel) *TextTableView {
	treeView, _ := gtk.TreeViewNew()

	treeView.SetModel(model.sortableListStore)
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

	selection, _ := treeView.GetSelection()
	selection.SetMode(gtk.SELECTION_SINGLE)

	scrollableTreelist, _ := gtk.ScrolledWindowNew(nil, nil)
	scrollableTreelist.Add(treeView)
	scrollableTreelist.SetVExpand(true)

	return &TextTableView{
		treeView:           treeView,
		scrollableTreeView: scrollableTreelist,
	}
}

func newTextTableController(model *TextTableModel, view *TextTableView) *TextTableController {
	view.treeView.Connect("row-activated", func(tv *gtk.TreeView, path *gtk.TreePath, column *gtk.TreeViewColumn) {
		iter, _ := model.sortableListStore.GetIter(path)
		v, _ := model.sortableListStore.GetValue(iter, COLUMN_ID)
		gv, _ := v.GoValue()
		textId := gv.(string)

		v2, _ := model.sortableListStore.GetValue(iter, COLUMN_UNKNOWN_VOCAB)
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

	return &TextTableController{
		model: model,
		view:  view,
	}
}

func (c *TextTableController) updateTextsTable() {
	listStore := c.model.rawListStore
	listStore.Clear()

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

		iter := listStore.Append()
		numUnknown := len(text.Unknown(knownLearnables))
		numLearnables := len(text.Learnables)

		percentageKnown := 0
		if numLearnables != 0 {
			percentageKnown = 100 * (numLearnables - numUnknown) / numLearnables
		}
		listStore.Set(iter, []int{0, 1, 2, 3, 4, 5}, []interface{}{
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
