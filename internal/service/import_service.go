package service

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/erwindrsno/resi-scanner/internal/constant"
	"github.com/erwindrsno/resi-scanner/internal/shipment"
	"github.com/erwindrsno/resi-scanner/internal/util"
	"github.com/xuri/excelize/v2"
)

type ImportService struct {
	Repo shipment.ShipmentRepository
}

func (s *ImportService) ProcessExcel(f *excelize.File) error {
	rawDate, err := f.GetCellValue(constant.Sheets[0], "B1")
	if err != nil {
		fmt.Println(err)
	}
	date, err := util.ParseDate(rawDate, constant.Layout)
	if err != nil {
		fmt.Println(err)
	}

	return s.Repo.WithTransaction(func(tx *sql.Tx) error {
		for _, sheet := range constant.Sheets {
			rows, err := f.Rows(sheet)

			if err != nil {
				fmt.Println(err)
			}

			counter := 0
			for rows.Next() {
				//skip first 3 rows
				counter++
				if counter <= 3 {
					continue
				}

				row, err := rows.Columns()
				if err != nil {
					fmt.Println(err)
				}
				if row[0] == "TOTAL" {
					break
				}

				var locationId int
				switch sheet {
				case "BTH":
					locationId = 1
				case "TBK":
					locationId = 2
				case "TNJ":
					locationId = 3
				default:
					locationId = 0
				}

				shipment := &shipment.Shipment{
					Id:          util.GenerateUUID(),
					ResiNumber:  row[1],
					LocationId:  locationId,
					Destination: row[2],
					Weight:      util.ParseWeight(row[3]),
					Koli:        util.ParseKoli(row[4]),
					Notes:       "",
					Date:        date,
				}
				if err := s.Repo.InsertByTx(tx, shipment); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (s *ImportService) RunSearch(keyword, date, selectedSheet string) []shipment.Shipment {
	parsedDate, err := time.Parse(constant.Layout, date)
	if err != nil {
		fmt.Println(err)
	}
	// parsedDate, err := time.Parse(constant.Layout, date)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	filter := &shipment.ShipmentFilter{
		Keyword:  keyword,
		Date:     &parsedDate,
		Location: selectedSheet,
	}
	fmt.Println(filter)

	results, err := s.Repo.Get(filter)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(results)
	return results
}

func (s *ImportService) RunHighlight(resiNumber string) bool {
	err := s.Repo.UpdateIsScannedByResiNumber(resiNumber)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
