package gui

import (
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func displayStatusMessages(
	messageChan <-chan string,
	statusBar *gtk.Statusbar,
) {
	ticker := time.NewTicker(1 * time.Second)
	done := false

	for !done {
		select {
		case statusText, ok := <-messageChan:
			if ok {
				glib.IdleAdd(func() {
					statusBar.Push(0, statusText)
				})
			} else {
				done = true
			}
		case <-ticker.C:
		}
	}
}
