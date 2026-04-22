package shipment

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Shipment struct {
	Id          uuid.UUID
	ResiNumber  string
	LocationId  int
	Destination string
	Weight      float64
	Koli        int
	Notes       string
	Date        time.Time
	IsScanned   int
}

type ShipmentFilter struct {
	Keyword  string
	Date     *time.Time
	Location string
}

type ShipmentRepository interface {
	Get(filter *ShipmentFilter) ([]Shipment, error)
	GetByResi(resi string, date string) (*Shipment, error)
	Insert(s *Shipment) error
	InsertByTx(tx *sql.Tx, s *Shipment) error
	UpdateIsScannedByResiNumber(resiNumber string) error
	ImportFromExcel(filePath string) error
	WithTransaction(fn func(tx *sql.Tx) error) error
}
