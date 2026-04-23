package view

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/erwindrsno/resi-scanner/internal/constant"
	"github.com/erwindrsno/resi-scanner/pkg/viewmodel"
)

var (
	defaultKeywordPlaceholder = "Filter a word"
	defaultSheetSelector      = "Select a sheet..."
)

type Header struct {
	vm  *viewmodel.ViewModel
	win fyne.Window
}

func NewHeader(vm *viewmodel.ViewModel, win fyne.Window) *Header {
	return &Header{
		vm:  vm,
		win: win,
	}
}

func (h *Header) Render() fyne.CanvasObject {
	return container.NewGridWithColumns(
		3,
		h.createLeftSection(),
		h.createMiddleSection(),
		h.createRightSection())
}

func (h *Header) createLeftSection() *fyne.Container {
	filenameLabel := widget.NewLabelWithData(h.vm.FileName)
	openFileBtn := widget.NewButton("Open File", h.showFilePicker)
	fileRow := container.NewHBox(filenameLabel, openFileBtn)

	//INFO: This doesn't work, because it needs an event listener
	keywordEntry := widget.NewEntryWithData(h.vm.Keyword)
	reg := regexp.MustCompile("[a-z]")
	rawContent := reg.ReplaceAllString(strings.TrimSpace(keywordEntry.Text), "")
	content := strings.ToUpper(rawContent) // <--- This is how you get the content

	keywordEntry.SetPlaceHolder(defaultKeywordPlaceholder)
	keywordFilterBtn := widget.NewButton("Search", func() {
		filename, _ := h.vm.FileName.Get()
		date, _ := h.vm.Date.Get()
		selectedSheet, _ := h.vm.SelectedSheet.Get()
		slog.Info("metadata", "filename", filename, "keyword", keywordEntry.Text, "date", date, "sheet", selectedSheet)
		h.vm.ExecuteSearch(content, date, selectedSheet)
	})
	keywordFilterBtn.Importance = widget.HighImportance
	keywordRow := container.NewBorder(nil, nil, nil, keywordFilterBtn, keywordEntry)

	leftSection := container.NewVBox(fileRow, keywordRow)
	return leftSection
}

func (h *Header) createMiddleSection() *fyne.Container {
	dateLabel := widget.NewLabel("Select date:")
	dateEntry := widget.NewDateEntry()
	dateEntry.SetPlaceHolder("DD/MM/YYYY")
	dateEntry.OnChanged = func(selectedDate *time.Time) {
		if selectedDate != nil {
			h.vm.Date.Set(selectedDate.Format(constant.Layout))
		} else {
			h.vm.Date.Set("")
		}
	}
	// Layout: (top, bottom, left, right, center)
	dateRow := container.NewBorder(nil, nil, dateLabel, nil, dateEntry)

	resiFoundStatusLabel := widget.NewLabel("Status:")
	resiFoundStatusValue := widget.NewLabelWithData(h.vm.ResiFoundStatus)
	resiFoundStatusRow := container.NewHBox(resiFoundStatusLabel, resiFoundStatusValue)
	middleSection := container.NewVBox(dateRow, resiFoundStatusRow)

	return middleSection
}

func (h *Header) createRightSection() *fyne.Container {
	sheetLabel := widget.NewLabel("Select sheet:")
	sheetSelector := h.createSheetSelector()
	sheetRow := container.NewHBox(layout.NewSpacer(), sheetLabel, sheetSelector)

	clearBtn := widget.NewButton("Clear", func() {
		sheetSelector.Refresh()
		h.vm.Clear()
	})
	clearBtn.Importance = widget.DangerImportance
	btnRow := container.NewHBox(layout.NewSpacer(), clearBtn)
	rightSection := container.NewVBox(sheetRow, btnRow)
	return rightSection
}

func (h *Header) showFilePicker() {
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			slog.Error(err.Error())
			return
		}
		if reader == nil {
			slog.Info("File picking cancelled")
			return // User cancelled
		}

		// Get the path and hand it to the ViewModel
		// Capture the full URI
		uri := reader.URI()
		filename := filepath.Base(uri.Path())
		h.vm.FileURI.Set(uri)
		h.vm.FileName.Set(filename)

		// For debugging, you can print the full path or string
		strFileURIDEBUG, _ := h.vm.FileURI.Get()
		strFilePathDEBUG, _ := h.vm.FilePath.Get()
		strFileNameDEBUG, _ := h.vm.FileName.Get()
		slog.Debug("File metadata", "uri", strFileURIDEBUG, "path", strFilePathDEBUG, "name", strFileNameDEBUG)
		h.vm.SelectFile(reader)
		//TODO:Can put loading UI here
		h.showPreviewDialog()
	}, h.win)
}

func (h *Header) showPreviewDialog() {
	filename, err := h.vm.FileName.Get()
	if err != nil {
		slog.Error(err.Error())
	}
	d := dialog.NewConfirm("Import Preview", fmt.Sprintf("You are going to import %s", filename), func(b bool) {
		if b {
			slog.Info("Processing excel", "file", filename)
			h.vm.ImportConfirmed(h.vm.File)
		} else {
			slog.Info("Cancelling excel processing")
		}
	}, h.win)
	d.Resize(fyne.NewSize(400, 300))
	d.Show()
}

func (h *Header) createSheetSelector() fyne.CanvasObject {
	// 1. Create the Select with empty options initially
	selector := widget.NewSelect([]string{}, func(selected string) {
		slog.Info("Sheet selected", "sheet", selected)
		h.vm.SelectedSheet.Set(selected)
	})
	sheetPlaceholder, err := h.vm.SheetPlaceholder.Get()
	if err != nil {
		slog.Error(err.Error())
	}
	selector.PlaceHolder = sheetPlaceholder

	// 2. Attach a listener to the StringList binding
	h.vm.Sheets.AddListener(binding.NewDataListener(func() {
		// This code runs every time vm.Sheets.Set() is called
		newSheets, _ := h.vm.Sheets.Get()

		if len(newSheets) == 0 {
			selector.PlaceHolder = defaultSheetSelector
			selector.ClearSelected()
		} else {
			selector.Options = newSheets
		}
		selector.Refresh()
	}))

	return selector
}
