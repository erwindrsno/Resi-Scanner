package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/erwindrsno/resi-scanner/internal/data"
	"github.com/erwindrsno/resi-scanner/pkg/ui"

	// "github.com/erwindrsno/resi-scanner/location"
	"github.com/erwindrsno/resi-scanner/internal/service"
	"github.com/erwindrsno/resi-scanner/internal/shipment"
)

func main() {
	myapp := app.NewWithID("com.resiscanner.tool")

	dbPath := data.GetDatabasePath(myapp)
	db, err := data.InitDB(dbPath)

	if err != nil {
		fmt.Println(err)
		// panic(err)
	}
	defer db.Close()

	shipmentRepo := shipment.NewSQLiteShipmentRepository(db)
	// locationRepo := location.NewSQLiteLocationRepository(db)
	importService := service.ImportService{Repo: shipmentRepo}
	// err = importService.ProcessExcel("a")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// shipmentRepo.GetByResi("a", "B")

	// myapp.Storage().RootURI().Path()
	win := myapp.NewWindow("Resi Scanner")
	win.Resize(fyne.NewSize(1000, 800))

	vm := ui.NewViewModel(&importService)
	v := ui.NewView(win, vm)

	win.SetContent(v.Render())
	win.ShowAndRun()
}
