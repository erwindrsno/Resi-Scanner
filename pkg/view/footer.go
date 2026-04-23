package view

import (
	"log/slog"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/erwindrsno/resi-scanner/pkg/viewmodel"
)

var (
	defaultScanPlaceholderMessage = "Enter or Scan barcode..."
)

type Footer struct {
	vm  *viewmodel.ViewModel
	win fyne.Window

	OnSearchSuccess func(rowIdx int)
}

func NewFooter(vm *viewmodel.ViewModel, win fyne.Window) *Footer {
	return &Footer{
		vm:  vm,
		win: win,
	}
}

func (f *Footer) Render() fyne.CanvasObject {
	scanEntry := widget.NewEntry()
	scanEntry.ActionItem = nil
	scanEntry.SetPlaceHolder(defaultScanPlaceholderMessage)

	// 1. Create a shared helper function to avoid repeating code
	runHighlight := func() {
		slog.Info("Process highlighting")
		content := strings.TrimSpace(scanEntry.Text)

		rowIdx := f.vm.ExecuteHighlight(content)

		f.OnSearchSuccess(rowIdx)

		// Clear and refocus for the next scan
		scanEntry.SetText("")
		f.win.Canvas().Focus(scanEntry)
	}

	// 2. Trigger via Enter Key
	scanEntry.OnSubmitted = func(content string) {
		runHighlight()
	}

	// 3. Trigger via Save Button
	saveBtn := widget.NewButton("Save", func() {
		runHighlight()
	})

	footer := container.NewBorder(nil, nil, nil, saveBtn, scanEntry)

	return footer
}
