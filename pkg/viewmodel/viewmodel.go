package viewmodel

import (
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"github.com/erwindrsno/resi-scanner/internal/constant"
	"github.com/erwindrsno/resi-scanner/internal/service"
	"github.com/erwindrsno/resi-scanner/internal/shipment"
	"github.com/xuri/excelize/v2"
)

var (
	defaultFileMessage        = "No file selected."
	defaultSheetSelector      = "Select a sheet..."
	defaultKeywordPlaceholder = "Filter a word"
	defaultResiFoundStatus    = "-"
)

type ViewModel struct {
	//file section
	FileName binding.String
	FilePath binding.String
	FileURI  binding.URI
	File     *excelize.File

	//date section
	Date binding.String

	//sheet section
	Sheets           binding.StringList
	SelectedSheet    binding.String
	SheetPlaceholder binding.String

	//Keyword section
	Keyword            binding.String
	KeywordPlaceholder binding.String

	ResiFoundStatus binding.String
	Service         *service.Service

	//footer
	Resi binding.String

	//data
	Shipments []shipment.Shipment

	OnDataChanged func()

	// Pvm              *PreviewViewModel
}

func NewViewModel(service *service.Service) *ViewModel {
	vm := &ViewModel{
		FileName:           binding.NewString(),
		FilePath:           binding.NewString(),
		FileURI:            binding.NewURI(),
		File:               nil,
		Date:               binding.NewString(),
		Sheets:             binding.NewStringList(),
		SelectedSheet:      binding.NewString(),
		SheetPlaceholder:   binding.NewString(),
		Keyword:            binding.NewString(),
		KeywordPlaceholder: binding.NewString(),
		ResiFoundStatus:    binding.NewString(),
		Service:            service,
		Resi:               binding.NewString(),
		Shipments:          make([]shipment.Shipment, 0),

		// Pvm:              pvm,
	}

	vm.FileName.Set(defaultFileMessage)
	vm.ResiFoundStatus.Set(defaultResiFoundStatus)
	vm.SheetPlaceholder.Set(defaultSheetSelector)
	vm.KeywordPlaceholder.Set(defaultKeywordPlaceholder)
	vm.Sheets.Set(constant.Sheets)
	return vm
}

func (vm *ViewModel) SelectFile(reader fyne.URIReadCloser) {
	defer reader.Close()
	f, err := excelize.OpenReader(reader)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	vm.File = f
}

func (vm *ViewModel) Clear() {
	vm.FileName.Set(defaultFileMessage)
	vm.FilePath.Set("")
	vm.FileURI.Set(nil)
	vm.File = nil
	vm.Sheets.Set(nil)
	vm.SelectedSheet.Set("")
	vm.SheetPlaceholder.Set(defaultSheetSelector)
	vm.Keyword.Set("")
	vm.ResiFoundStatus.Set(defaultResiFoundStatus)
	vm.Resi.Set("")
	vm.Shipments = []shipment.Shipment{}

}

func (vm *ViewModel) ExecuteSearch(keyword, date, sheet string) {
	// 1. Fetch data
	if sheet == "" {
		return
	}
	results := vm.Service.RunSearch(keyword, date, sheet)

	// 2. Update the raw data
	vm.Shipments = results

	// 3. POKE THE UI: "Hey, I'm done!"
	if vm.OnDataChanged != nil {
		vm.OnDataChanged()
	}
}

func (vm *ViewModel) ExecuteHighlight(resiNumber string) int {
	vm.Service.RunHighlight(resiNumber)

	keyword, err := vm.Keyword.Get()
	if err != nil {
		slog.Error(err.Error())
	}
	date, err := vm.Date.Get()
	if err != nil {
		slog.Error(err.Error())
	}
	sheet, err := vm.SelectedSheet.Get()
	if err != nil {
		slog.Error(err.Error())
	}

	// 2. Find the row index of the scanned Resi
	targetRow := -1
	for i, s := range vm.Shipments {
		if s.ResiNumber == resiNumber {
			targetRow = i
			break
		}
	}

	vm.ExecuteSearch(keyword, date, sheet)

	return targetRow
}

func (vm *ViewModel) ImportConfirmed(f *excelize.File) {
	err := vm.Service.ProcessExcel(f)
	if err != nil {
		slog.Error(err.Error())
	}
}
