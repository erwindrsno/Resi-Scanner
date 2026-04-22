package shipment

import (
	"database/sql"
	"fmt"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteShipmentRepository(db *sql.DB) ShipmentRepository {
	return &SQLiteRepository{db: db}
}

// Now implement the interface methods...
func (r *SQLiteRepository) Insert(s *Shipment) error {
	// SQL Insert logic goes here
	return nil
}

func (r *SQLiteRepository) InsertByTx(tx *sql.Tx, s *Shipment) error {
	fmt.Println(s.ResiNumber)
	query := `INSERT INTO shipments (id, resi_number, location_id, destination, weight, koli, notes, date, is_scanned) 
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := tx.Exec(query, s.Id.String(), s.ResiNumber, s.LocationId, s.Destination, s.Weight, s.Koli, s.Notes, s.Date.String(), s.IsScanned)
	return err
}

func (r *SQLiteRepository) Get(filter *ShipmentFilter) ([]Shipment, error) {
	query := `
		SELECT s.resi_number, s.destination, s.weight, s.koli, s.notes, s.is_scanned
		FROM shipments s 
		JOIN locations l ON s.location_id = l.id
		WHERE 1=1`
	var args []interface{}

	if filter.Date != nil {
		query += " AND s.date = ?"
		fmt.Printf("Date query: %s\n", filter.Date)
		args = append(args, filter.Date.String())
	}

	if filter.Keyword != "" {
		query += " AND s.resi_number LIKE ?"
		fmt.Printf("Keyword query: %s\n", filter.Keyword)
		args = append(args, "%"+filter.Keyword+"%")
	}

	if filter.Location != "" {
		query += " AND l.code LIKE ?"
		fmt.Printf("Location query: %s\n", filter.Location)
		args = append(args, "%"+filter.Location+"%")
	}
	fmt.Println(query)
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var shipments = make([]Shipment, 0)
	for rows.Next() {
		var s Shipment
		err := rows.Scan(&s.ResiNumber, &s.Destination, &s.Weight, &s.Koli, &s.Notes, &s.IsScanned)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		shipments = append(shipments, s)
	}
	return shipments, nil
}

func (r *SQLiteRepository) GetByResi(resi string, date string) (*Shipment, error) {
	// SQL Select logic goes here
	return nil, nil
}

func (r *SQLiteRepository) UpdateIsScannedByResiNumber(resiNumber string) error {
	query := `UPDATE shipments 
    SET is_scanned = CASE 
        WHEN is_scanned >= 2 THEN 0 
        ELSE is_scanned + 1 
    END 
    WHERE resi_number = ?`

	// 2. Execute the query
	result, err := r.db.Exec(query, resiNumber)
	if err != nil {
		return fmt.Errorf("failed to update shipment %s: %w", resiNumber, err)
	}

	// 3. Optional: Check if any row was actually updated
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no shipment found with resi number: %s", resiNumber)
	}
	return nil
}

func (r *SQLiteRepository) ImportFromExcel(filePath string) error {
	// Excelize logic goes here
	return nil
}

func (r *SQLiteRepository) WithTransaction(fn func(tx *sql.Tx) error) error {
	// 1. Start the transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// 2. Use defer to rollback if something goes wrong
	defer tx.Rollback()

	// 3. Execute the logic passed from the service
	if err := fn(tx); err != nil {
		return err
	}

	// 4. If no error, commit to disk
	return tx.Commit()
}
