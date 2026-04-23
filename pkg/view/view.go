package view

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/erwindrsno/resi-scanner/internal/constant"
	"github.com/erwindrsno/resi-scanner/pkg/viewmodel"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"

	"fyne.io/fyne/v2/widget"
)

type View struct {
	// Reference to the main window for dialogs/popups
	win fyne.Window
	// Reference to the logic & state
	vm *viewmodel.ViewModel

	dataTable *widget.Table

	h *Header
	f *Footer
}

func NewView(win fyne.Window, vm *viewmodel.ViewModel) *View {
	v := &View{
		win: win,
		vm:  vm,
	}
	v.dataTable = v.CreateTable()
	vm.OnDataChanged = func() {
		v.dataTable.Refresh()     // Fixes the "ghosting"
		v.dataTable.ScrollToTop() // Prevents being stuck at the bottom
	}

	v.h = NewHeader(vm, win)
	v.f = NewFooter(vm, win)
	return v
}

func (v *View) Render() fyne.CanvasObject {
	v.f.OnSearchSuccess = func(rowIdx int) {
		if rowIdx != -1 {
			v.dataTable.ScrollTo(widget.TableCellID{Row: rowIdx, Col: 0})
			// v.dataTable.Select(widget.TableCellID{Row: rowIdx, Col: 0})
			// TODO: Can't use hardcoded value below.
			v.vm.ResiFoundStatus.Set("FOUND.")
		} else {
			v.dataTable.Refresh() // Required for Android stability
			v.vm.ResiFoundStatus.Set("NOT FOUND.")
		}
	}

	return container.NewBorder(
		v.h.Render(), // Top
		v.f.Render(), // Bottom
		nil, nil,     // Left, Right
		v.dataTable, // Center
	)
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
