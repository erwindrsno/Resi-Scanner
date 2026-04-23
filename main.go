package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/erwindrsno/resi-scanner/internal/database"
	"github.com/erwindrsno/resi-scanner/pkg/view"
	"github.com/erwindrsno/resi-scanner/pkg/viewmodel"

	// "github.com/erwindrsno/resi-scanner/location"
	"github.com/erwindrsno/resi-scanner/internal/service"
	"github.com/erwindrsno/resi-scanner/internal/shipment"
)

func main() {
	myapp := app.NewWithID("com.resiscanner.tool")

	dbPath := database.GetDatabasePath(myapp)
	db, err := database.InitDB(dbPath)

	if err != nil {
		fmt.Println(err)
		// panic(err)
	}
	defer db.Close()

	shipmentRepo := shipment.NewSQLiteShipmentRepository(db)
	// locationRepo := location.NewSQLiteLocationRepository(db)
	service := service.Service{Repo: shipmentRepo}

	win := myapp.NewWindow("Resi Scanner")
	win.Resize(fyne.NewSize(1000, 800))

	vm := viewmodel.NewViewModel(&service)
	v := view.NewView(win, vm)

	win.SetContent(v.Render())
	win.ShowAndRun()
}
