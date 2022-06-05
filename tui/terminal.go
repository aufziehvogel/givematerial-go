package tui

import (
	"fmt"
	"givematerial/givematlib"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func updateStatus(messageChan <-chan string, statusView *tview.TextView, app *tview.Application) {
	ticker := time.NewTicker(1 * time.Second)
	done := false

	for !done {
		select {
		case statusText, ok := <-messageChan:
			if ok {
				app.QueueUpdateDraw(func() {
					statusView.SetText(statusText)
				})
			} else {
				done = true
			}
		case <-ticker.C:
		}
	}
}

func updateProviders(
	config *givematlib.ApplicationConfig,
	messageChan chan<- string,
) error {
	defer close(messageChan)

	err := updateAnki(config.AnkiFile, config.AnkiDecksForLanguage, messageChan)
	if err != nil {
		return err
	}

	return updateWanikani(config.WanikaniApiToken, messageChan)
}

func updateWanikani(apiToken string, messageChan chan<- string) error {
	w := givematlib.NewWanikaniProvider(apiToken)

	messageChan <- "Reading status from Wanikani"
	learnables, err := w.FetchLearnables()
	if err != nil {
		return err
	}

	messageChan <- "Saving Wanikani status"
	err = givematlib.SaveLearnableStatus("wanikani", "ja", learnables)
	messageChan <- "Wanikani status update complete"
	return err
}

func updateAnki(
	ankiFile string,
	decksForLanguage map[string][]string,
	messageChan chan<- string,
) error {
	for language, decks := range decksForLanguage {
		p := givematlib.AnkiProvider{
			Decks:    decks,
			AnkiFile: ankiFile,
		}
		messageChan <- fmt.Sprintf(
			"Reading status from Anki for %s (decks: %v)",
			language,
			decks,
		)
		learnables, err := p.FetchLearnables()
		if err != nil {
			return err
		}

		messageChan <- fmt.Sprintf(
			"Saving Anki status for %s (%d words)",
			language,
			len(learnables),
		)
		err = givematlib.SaveLearnableStatus("anki", language, learnables)
		if err != nil {
			return err
		}
		messageChan <- "Anki status update complete"
	}

	return nil
}

func newEventHandler(
	config *givematlib.ApplicationConfig,
	statusView *tview.TextView,
	app *tview.Application,
) func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		k := event.Key()

		if k == tcell.KeyTAB {
			currentPrimitive =
				(currentPrimitive + 1) % len(selectablePrimitives)
			app.SetFocus(selectablePrimitives[currentPrimitive])
		} else if k == tcell.KeyRune && event.Rune() == 'u' {
			statusMessages := make(chan string)
			// TODO: Error handling
			go updateProviders(config, statusMessages)
			go updateStatus(statusMessages, statusView, app)
			return nil
		}
		return event
	}
}

var selectablePrimitives []tview.Primitive
var currentPrimitive int

// Show a navigable tree view of the current directory.
func Init(config *givematlib.ApplicationConfig) {
	rootDir := "."
	root := tview.NewTreeNode(rootDir).
		SetColor(tcell.ColorRed)
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	texts, err := givematlib.ListTexts()
	if err != nil {
		panic(err)
	}
	// TODO: Build tree in givematlib/texts.go
	folders := make(map[string]*tview.TreeNode)
	for _, textID := range texts {
		// TODO: Error handling
		text, _ := givematlib.LoadText(textID)
		hierarchies := strings.SplitN(text.Title, " - ", 2)

		switch len(hierarchies) {
		case 2:
			treeNode, ok := folders[hierarchies[0]]

			if !ok {
				treeNode = tview.NewTreeNode(hierarchies[0]).
					SetColor(tcell.ColorGreen)
				treeNode.Collapse()
				root.AddChild(treeNode)
				folders[hierarchies[0]] = treeNode
			}
			element := tview.NewTreeNode(hierarchies[1]).
				SetReference(textID)
			treeNode.AddChild(element)
		case 1:
			element := tview.NewTreeNode(hierarchies[0]).
				SetReference(textID)
			root.AddChild(element)
		}
	}

	app := tview.NewApplication()

	status := tview.NewTextView().
		SetText("Status not set.")
	content := tview.NewTextView()
	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		if len(node.GetChildren()) > 0 {
			if node.IsExpanded() {
				node.Collapse()
			} else {
				node.Expand()
			}
		}

		reference := node.GetReference()
		if reference == nil {
			return
		}

		textID := reference.(string)
		text, err := givematlib.LoadText(textID)
		if err != nil {
			status.SetText(
				fmt.Sprintf("Error loading text %q", textID),
			)
			return
		}

		knownLearnables, err := givematlib.ReadLearnableStatus(text.Language)
		if err != nil {
			status.SetText(fmt.Sprintf("Error loading learnables: %v", err))
			return
		}
		unknownLearnables := text.Unknown(knownLearnables)
		content.SetText(fmt.Sprintf(
			"Unknown: %s\n\n%s",
			strings.Join(unknownLearnables, ", "),
			text.Fulltext,
		))
		content.ScrollToBeginning()
	})

	app.SetInputCapture(newEventHandler(config, status, app))

	grid := tview.NewGrid().
		SetRows(0, 1).
		SetColumns(30, 0).
		SetBorders(true).
		AddItem(tree, 0, 0, 1, 1, 0, 0, true).
		AddItem(content, 0, 1, 1, 1, 0, 0, true).
		AddItem(status, 1, 0, 1, 2, 0, 0, false)

	selectablePrimitives = []tview.Primitive{tree, content}

	if err := app.SetRoot(grid, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
