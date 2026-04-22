package location

import "database/sql"

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteLocationRepository(db *sql.DB) LocationRepository {
	return &SQLiteRepository{db: db}
}

// Now implement the interface methods...
func (r *SQLiteRepository) Save(l *Location) error {
	// SQL Insert logic goes here
	return nil
}
