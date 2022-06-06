package gui

import "github.com/gotk3/gotk3/gtk"

func ReadTextViewContent(textView *gtk.TextView) (string, error) {
	buffer, err := textView.GetBuffer()
	if err != nil {
		return "", err
	}

	start, end := buffer.GetBounds()

	return buffer.GetText(start, end, false)
}
