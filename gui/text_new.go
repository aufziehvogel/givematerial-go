package gui

import (
	"givematerial/givematlib"

	"github.com/gotk3/gotk3/gtk"
)

func createNewTextWindow() (*gtk.Window, error) {
	w, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		return nil, err
	}

	grid, err := gtk.GridNew()
	grid.SetColumnSpacing(10)
	grid.SetRowSpacing(10)
	if err != nil {
		return nil, err
	}

	titleLabel, err := gtk.LabelNew("Title")
	if err != nil {
		return nil, err
	}

	languageLabel, err := gtk.LabelNew("Language")
	if err != nil {
		return nil, err
	}

	titleField, err := gtk.EntryNew()
	if err != nil {
		return nil, err
	}

	languages, err := loadLanguages()
	if err != nil {
		return nil, err
	}
	languagesListStore, err := newLanguagesListStore(languages)
	if err != nil {
		return nil, err
	}

	languageField, err := gtk.ComboBoxNewWithModelAndEntry(languagesListStore)
	languageField.SetEntryTextColumn(0)

	textField, err := gtk.TextViewNew()
	if err != nil {
		return nil, err
	}
	textField.SetHExpand(true)
	textField.SetVExpand(true)

	saveButton, err := gtk.ButtonNewWithLabel("Save")
	if err != nil {
		return nil, err
	}
	saveButton.Connect("clicked", func(obj *gtk.Button) {
		title, err := titleField.GetText()
		if err != nil {
			return
		}
		text, err := ReadTextViewContent(textField)
		if err != nil {
			return
		}

		iter, err := languageField.GetActiveIter()
		model, err := languageField.GetModel()
		v, _ := model.ToTreeModel().GetValue(iter, 0)
		gv, _ := v.GoValue()
		language := givematlib.MakeLanguageFromShortCode(gv.(string))

		extractor := givematlib.GetExtractorForLanguage(language)
		learnables := extractor.ExtractLearnables(text)

		t := givematlib.Text{
			Title:      title,
			Fulltext:   text,
			Language:   language.ShortCode(),
			Learnables: learnables,
		}
		givematlib.SaveText(t)
		// TODO: Emit a textAdded event

		w.Close()
	})

	grid.Attach(titleLabel, 0, 0, 1, 1)
	grid.AttachNextTo(titleField, titleLabel, gtk.POS_RIGHT, 1, 1)
	grid.AttachNextTo(languageLabel, titleLabel, gtk.POS_BOTTOM, 1, 1)
	grid.AttachNextTo(languageField, languageLabel, gtk.POS_RIGHT, 1, 1)
	grid.AttachNextTo(textField, languageLabel, gtk.POS_BOTTOM, 2, 1)
	grid.AttachNextTo(saveButton, textField, gtk.POS_BOTTOM, 2, 1)

	w.Add(grid)

	return w, nil
}
