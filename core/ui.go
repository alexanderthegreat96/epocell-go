package core

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"os"
)

func (t *TableModel) RowCount() int {
	return len(t.Rows)
}

func (t *TableModel) CreateRow(i int) fyne.CanvasObject {
	return container.NewGridWithColumns(len(t.Columns),
		widget.NewLabel(t.Rows[i][0]),
		widget.NewLabel(t.Rows[i][1]),
		widget.NewLabel(t.Rows[i][2]),
	)
}

func PickFile(win fyne.Window, filePathEntry *widget.Entry) {
	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return
	}

	// Create a file URI for the current directory
	currentDirURI := storage.NewFileURI(currentDir)

	lister, _ := storage.ListerForURI(currentDirURI)

	dialogOpen := dialog.NewFileOpen(func(f fyne.URIReadCloser, err error) {
		if err != nil || f == nil {
			return
		}
		// Get the selected file path without the "file://" prefix
		filePath := f.URI().Path()
		err = f.Close()
		if err != nil {
			return
		}
		filePathEntry.SetText(filePath)
	}, win)
	dialogOpen.SetLocation(lister)

	dialogOpen.SetFilter(storage.NewExtensionFileFilter([]string{".xls", ".xlsx"})) // Set file filters if needed
	dialogOpen.Show()
}
func PickFileCsv(win fyne.Window, filePathEntry *widget.Entry) {
	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return
	}

	// Create a file URI for the current directory
	currentDirURI := storage.NewFileURI(currentDir)

	lister, _ := storage.ListerForURI(currentDirURI)

	dialogOpen := dialog.NewFileOpen(func(f fyne.URIReadCloser, err error) {
		if err != nil || f == nil {
			return
		}
		// Get the selected file path without the "file://" prefix
		filePath := f.URI().Path()
		err = f.Close()
		if err != nil {
			return
		}
		filePathEntry.SetText(filePath)
	}, win)
	dialogOpen.SetLocation(lister)

	dialogOpen.SetFilter(storage.NewExtensionFileFilter([]string{".csv"})) // Set file filters if needed
	dialogOpen.Show()
}

func newClickableLabel(text string, onClicked func()) *clickableLabel {
	label := &clickableLabel{
		Label:     widget.NewLabel(text),
		onClicked: onClicked,
	}

	return label
}

func (c *clickableLabel) Tapped(*fyne.PointEvent) {
	if c.onClicked != nil {
		c.onClicked()
	}
}

func createForm(storeCache StoreCache) fyne.CanvasObject {
	storeNameEntry := widget.NewEntry()
	salesUnitsEntry := widget.NewEntry()

	storeNameEntry.SetText(storeCache.StoreName)
	salesUnitsEntry.SetText(storeCache.SalesUnits)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Store Name", Widget: storeNameEntry},
			{Text: "Sales Units", Widget: salesUnitsEntry},
			// Add more form fields as needed
		},
	}
	return form
}
