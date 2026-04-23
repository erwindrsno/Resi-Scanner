package service

import (
	"database/sql"
	"log/slog"
	"time"

	"github.com/erwindrsno/resi-scanner/internal/constant"
	"github.com/erwindrsno/resi-scanner/internal/shipment"
	"github.com/erwindrsno/resi-scanner/internal/util"
	"github.com/xuri/excelize/v2"
)

type Service struct {
	Repo shipment.ShipmentRepository
}

func (s *Service) ProcessExcel(f *excelize.File) error {
	rawDate, err := f.GetCellValue(constant.Sheets[0], "B1")
	if err != nil {
		slog.Error(err.Error())
	}
	date, err := util.ParseDate(rawDate, constant.Layout)
	if err != nil {
		slog.Error(err.Error())
	}

	return s.Repo.WithTransaction(func(tx *sql.Tx) error {
		for _, sheet := range constant.Sheets {
			rows, err := f.Rows(sheet)

			if err != nil {
				slog.Error(err.Error())
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
					slog.Error(err.Error())
				}
				if row[0] == "TOTAL" {
					break
				}

				var locationId int
				switch sheet {
				case constant.Sheets[0]:
					locationId = 1
				case constant.Sheets[1]:
					locationId = 2
				case constant.Sheets[2]:
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

func (s *Service) RunSearch(keyword, date, selectedSheet string) []shipment.Shipment {
	parsedDate, err := time.Parse(constant.Layout, date)
	if err != nil {
		slog.Error(err.Error())
	}
	filter := &shipment.ShipmentFilter{
		Keyword:  keyword,
		Date:     &parsedDate,
		Location: selectedSheet,
	}

	results, err := s.Repo.Get(filter)
	if err != nil {
		slog.Error(err.Error())
	}
	return results
}

func (s *Service) RunHighlight(resiNumber string) bool {
	err := s.Repo.UpdateIsScannedByResiNumber(resiNumber)
	if err != nil {
		slog.Error(err.Error())
		return false
	}
	return true
}
