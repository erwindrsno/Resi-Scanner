package ui

import (
	"fmt"
	"image/color"
	"path/filepath"
	"regexp"

	// "regexp"
	"strconv"
	"strings"
	"time"

	"github.com/erwindrsno/resi-scanner/internal/constant"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"

	"fyne.io/fyne/v2/widget"
)

type View struct {
	// Reference to the main window for dialogs/popups
	win fyne.Window
	// Reference to the logic & state
	vm *ViewModel

	dataTable *widget.Table
}

func NewView(win fyne.Window, vm *ViewModel) *View {
	v := &View{
		win: win,
		vm:  vm,
	}
	v.dataTable = v.CreateTable()
	vm.OnDataChanged = func() {
		v.dataTable.Refresh()     // Fixes the "ghosting"
		v.dataTable.ScrollToTop() // Prevents being stuck at the bottom
	}
	return v
}

func (v *View) Render() fyne.CanvasObject {
	return container.NewBorder(
		v.CreateHeader(), // Top
		v.CreateFooter(), // Bottom
		nil, nil,         // Left, Right
		v.dataTable, // Center
	)
}

func (v *View) CreateHeader() fyne.CanvasObject {
	leftSection := v.createLeftSection()
	middleSection := v.createMiddleSection()
	rightSection := v.createRightSection()

	header := container.NewGridWithColumns(3, leftSection, middleSection, rightSection)
	return header
}

func (v *View) CreateFooter() fyne.CanvasObject {
	scanEntry := widget.NewEntry()
	scanEntry.ActionItem = nil
	scanEntry.SetPlaceHolder(defaultScanPlaceholderMessage)

	// 1. Create a shared helper function to avoid repeating code
	runHighlight := func() {
		reg := regexp.MustCompile("[a-z]")
		rawContent := reg.ReplaceAllString(strings.TrimSpace(scanEntry.Text), "")
		content := strings.ToUpper(rawContent) // <--- This is how you get the content
		// content := strings.ToUpper(strings.TrimSpace(scanEntry.Text))
		if content == "" {
			return
		}

		fmt.Printf("Triggering Highlight for: %s\n", content)
		rowIdx := v.vm.ExecuteHighlight(content)

		if rowIdx != -1 {
			v.dataTable.ScrollTo(widget.TableCellID{Row: rowIdx, Col: 0})
			// v.dataTable.Select(widget.TableCellID{Row: rowIdx, Col: 0})
			v.vm.ResiFoundStatus.Set("FOUND.")
			v.dataTable.Refresh() // Required for Android stability
		} else {
			v.vm.ResiFoundStatus.Set("NOT FOUND.")
		}

		// Clear and refocus for the next scan
		scanEntry.SetText("")
		v.win.Canvas().Focus(scanEntry)
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

func (v *View) CreateTable() *widget.Table {
	t := widget.NewTable(
		// 1. Length: Rows = slice length, Cols = 6
		func() (int, int) {
			if len(v.vm.Shipments) == 0 {
				return 0, 0
			}
			return len(v.vm.Shipments), 6
		},
		// 2. Create: Keep your clever Stack logic
		func() fyne.CanvasObject {
			bg := canvas.NewRectangle(color.Transparent)
			lbl := widget.NewLabel("")
			lbl.Alignment = fyne.TextAlignCenter
			return container.NewStack(bg, lbl)
		},
		// 3. Update: Map column index to Shipment fields
		func(id widget.TableCellID, o fyne.CanvasObject) {
			stack := o.(*fyne.Container)
			bg := stack.Objects[0].(*canvas.Rectangle)
			lbl := stack.Objects[1].(*widget.Label)

			// Safety check: Prevent index out of bounds during clear/refresh
			if id.Row >= len(v.vm.Shipments) {
				return
			}

			ship := v.vm.Shipments[id.Row]

			// Map the Column ID to the specific field in your object
			var value string
			switch id.Col {
			case 0:
				value = fmt.Sprintf("%d", id.Row+1) // Row number/No.
			case 1:
				value = ship.ResiNumber
			case 2:
				// Handle the time formatting here
				value = ship.Destination // Or whatever your field is
			case 3:
				value = fmt.Sprintf("%.2f", ship.Weight)
			case 4:
				value = strconv.Itoa(ship.Koli)
			default:
				value = "-"
			}

			switch {
			case ship.IsScanned == 1:
				bg.FillColor = constant.Yellow
				lbl.TextStyle = fyne.TextStyle{Bold: true}
			case ship.IsScanned >= 2: // Uses >= in case it gets scanned more than twice
				bg.FillColor = constant.Green
				lbl.TextStyle = fyne.TextStyle{Bold: true}
			default:
				// Reset for rows with 0 scans (Crucial for recycling!)
				bg.FillColor = color.Transparent
				lbl.TextStyle = fyne.TextStyle{Bold: false}
			}

			lbl.SetText(value)

			bg.Refresh()
		},
	)

	// Set the column widths
	widths := []float32{40, 150, 150, 80, 80, 200}
	for i, w := range widths {
		t.SetColumnWidth(i, w)
	}

	return t
}

func (v *View) CreateSheetSelector() fyne.CanvasObject {
	// 1. Create the Select with empty options initially
	selector := widget.NewSelect([]string{}, func(selected string) {
		fmt.Println("Selected sheet:", selected)
		v.vm.SelectedSheet.Set(selected)
	})
	sheetPlaceholder, err := v.vm.SheetPlaceholder.Get()
	if err != nil {
		fmt.Println(err)
	}
	selector.PlaceHolder = sheetPlaceholder

	// 2. Attach a listener to the StringList binding
	v.vm.Sheets.AddListener(binding.NewDataListener(func() {
		// This code runs every time vm.Sheets.Set() is called
		newSheets, _ := v.vm.Sheets.Get()

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

func (v *View) ShowFilePicker() {
	// Dialogs are UI-specific, so they stay in the View
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			// Handle UI error (popup)
			fmt.Println(err)
			return
		}
		if reader == nil {
			fmt.Println(err)
			return // User cancelled
		}

		// Get the path and hand it to the ViewModel
		// Capture the full URI
		uri := reader.URI()
		filename := filepath.Base(uri.Path())
		v.vm.FileURI.Set(uri)
		v.vm.FileName.Set(filename)

		// For debugging, you can print the full path or string
		strFileURIDEBUG, _ := v.vm.FileURI.Get()
		strFilePathDEBUG, _ := v.vm.FilePath.Get()
		strFileNameDEBUG, _ := v.vm.FileName.Get()
		fmt.Println("Full URI String:", strFileURIDEBUG)
		fmt.Println("Path (if available):", strFilePathDEBUG)
		fmt.Println("FIlename (if available):", strFileNameDEBUG)
		v.vm.SelectFile(reader)
		//TODO:Can put loading UI here
		v.ShowPreviewDialog()
	}, v.win)
}

func (v *View) ShowPreviewDialog() {
	filename, err := v.vm.FileName.Get()
	if err != nil {
		fmt.Println(err)
	}
	d := dialog.NewConfirm("Import Preview", fmt.Sprintf("You are going to import %s", filename), func(b bool) {
		if b {
			fmt.Println("You clicked confirm!!!")
			err := v.vm.ImportService.ProcessExcel(v.vm.File)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("You clicked cancel :(")
		}
	}, v.win)
	d.Resize(fyne.NewSize(400, 300))
	d.Show()
}

func (v *View) createLeftSection() *fyne.Container {

	// scanEntry := widget.NewEntry()
	// scanEntry.ActionItem = nil
	// scanEntry.SetPlaceHolder(defaultScanPlaceholderMessage)
	//
	// // 1. Create a shared helper function to avoid repeating code
	// runHighlight := func() {
	// 	reg := regexp.MustCompile("[a-z]")
	// 	rawContent := reg.ReplaceAllString(strings.TrimSpace(scanEntry.Text), "")
	// 	content := strings.ToUpper(rawContent) // <--- This is how you get the content

	filenameLabel := widget.NewLabelWithData(v.vm.FileName)
	openFileBtn := widget.NewButton("Open File", v.ShowFilePicker)
	fileRow := container.NewHBox(filenameLabel, openFileBtn)

	keywordEntry := widget.NewEntryWithData(v.vm.Keyword)
	reg := regexp.MustCompile("[a-z]")
	rawContent := reg.ReplaceAllString(strings.TrimSpace(keywordEntry.Text), "")
	content := strings.ToUpper(rawContent) // <--- This is how you get the content

	keywordEntry.SetPlaceHolder(defaultKeywordPlaceholder)
	keywordFilterBtn := widget.NewButton("Search", func() {
		filename, _ := v.vm.FileName.Get()
		date, _ := v.vm.Date.Get()
		selectedSheet, _ := v.vm.SelectedSheet.Get()
		fmt.Printf("You filenaem is: %s\n", filename)
		fmt.Printf("You searched for: %s\n", keywordEntry.Text)
		fmt.Printf("Your selected date is: %s\n", date)
		fmt.Printf("Your selected sheet is: %s\n", selectedSheet)
		v.vm.ExecuteSearch(content, date, selectedSheet)
	})
	keywordFilterBtn.Importance = widget.HighImportance
	keywordRow := container.NewBorder(nil, nil, nil, keywordFilterBtn, keywordEntry)

	leftSection := container.NewVBox(fileRow, keywordRow)
	return leftSection
}

func (v *View) createMiddleSection() *fyne.Container {
	dateLabel := widget.NewLabel("Select date:")
	dateEntry := widget.NewDateEntry()
	dateEntry.SetPlaceHolder("DD/MM/YYYY")
	dateEntry.OnChanged = func(selectedDate *time.Time) {
		if selectedDate != nil {
			fmt.Printf("Date is: %s\n", selectedDate)
			v.vm.Date.Set(selectedDate.Format(constant.Layout))
		} else {
			v.vm.Date.Set("")
		}
	}
	// Layout: (top, bottom, left, right, center)
	dateRow := container.NewBorder(nil, nil, dateLabel, nil, dateEntry)

	resiFoundStatusLabel := widget.NewLabel("Status:")
	resiFoundStatusValue := widget.NewLabelWithData(v.vm.ResiFoundStatus)
	resiFoundStatusRow := container.NewHBox(resiFoundStatusLabel, resiFoundStatusValue)
	middleSection := container.NewVBox(dateRow, resiFoundStatusRow)

	return middleSection
}

func (v *View) createRightSection() *fyne.Container {
	sheetLabel := widget.NewLabel("Select sheet:")
	sheetSelector := v.CreateSheetSelector()
	sheetRow := container.NewHBox(layout.NewSpacer(), sheetLabel, sheetSelector)

	clearBtn := widget.NewButton("Clear", func() {
		sheetSelector.Refresh()
		v.vm.Clear()
	})
	clearBtn.Importance = widget.DangerImportance
	btnRow := container.NewHBox(layout.NewSpacer(), clearBtn)
	rightSection := container.NewVBox(sheetRow, btnRow)
	return rightSection
}
