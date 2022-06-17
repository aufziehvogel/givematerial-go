package gui

import (
	"givematerial/givematlib"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func ReadTextViewContent(textView *gtk.TextView) (string, error) {
	buffer, err := textView.GetBuffer()
	if err != nil {
		return "", err
	}

	start, end := buffer.GetBounds()

	return buffer.GetText(start, end, false)
}

func newLanguagesListStore(languages map[givematlib.Language]struct{}) (*gtk.ListStore, error) {
	listStore, err := gtk.ListStoreNew(
		glib.TYPE_STRING,
	)
	if err != nil {
		return nil, err
	}

	for language, _ := range languages {
		iter := listStore.Append()
		listStore.Set(iter, []int{0}, []interface{}{language.ShortCode()})
	}

	return listStore, nil
}
