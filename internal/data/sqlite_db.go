package data

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"

	"fyne.io/fyne/v2"
	_ "github.com/mattn/go-sqlite3"
)

func InitDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	// 2. Create the tables (The "Schema")
	query := `
    CREATE TABLE IF NOT EXISTS locations (
        id TEXT PRIMARY KEY,
        code TEXT UNIQUE
    );

    CREATE TABLE IF NOT EXISTS shipments (
        id TEXT PRIMARY KEY,
        resi_number TEXT,
        location_id TEXT,
				destination TEXT,
        weight REAL,
        koli INTEGER,
        notes TEXT,
        date TEXT, -- YYYY-MM-DD
        is_scanned INTEGER DEFAULT 0,
        FOREIGN KEY(location_id) REFERENCES locations(id),
        UNIQUE(resi_number, destination, date)
    );

		INSERT OR IGNORE INTO locations (id, code) VALUES ('1', 'BTH');
    INSERT OR IGNORE INTO locations (id, code) VALUES ('2', 'TBK');
    INSERT OR IGNORE INTO locations (id, code) VALUES ('3', 'TNJ');
	`

	_, err = db.Exec(query)

	db.SetMaxOpenConns(1)
	if err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}

// func InitDB(dbPath string) (*sql.DB, error) {
//     // 1. Force create the directory structure on the phone
//     dir := filepath.Dir(dbPath)
//     err := os.MkdirAll(dir, 0755)
//     if err != nil {
//         return nil, fmt.Errorf("failed to create directory: %w", err)
//     }
//
//     // 2. Open the connection using the 'sqlite3' driver
//     db, err := sql.Open("sqlite3", dbPath)
//     if err != nil {
//         return nil, err
//     }
//
//     // 3. Ping the database to ensure the CGO connection is actually alive
//     // This is the best way to catch path/permission errors in Go
//     // before the C layer triggers a native crash.
//     err = db.Ping()
//     if err != nil {
//         return nil, fmt.Errorf("database ping failed: %w", err)
//     }
//
//     return db, nil
// }

func GetDatabasePath(myApp fyne.App) string {
	if runtime.GOOS == "android" {
		root := myApp.Storage().RootURI().Path()
		return filepath.Join(root, "resi-scanner.db")
	}
	return "resi-scanner.db"
}
