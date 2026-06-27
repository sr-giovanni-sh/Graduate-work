package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT "",
	title VARCHAR(128) NOT NULL DEFAULT "",
	comment TEXT NOT NULL DEFAULT "",
	repeat VARCHAR(128) NOT NULL DEFAULT ""
)
`

/*
Init initializes the SQLite database connection and creates the scheduler table
if it doesn’t exist or the file is empty.

Behavior:
  - Determines the database file path: uses the TODO_DBFILE environment variable
    if set; otherwise, uses the provided dbFile argument.
  - Checks if the database file exists and whether it’s empty. If the file doesn’t
    exist or has zero size, the install flag is set to true to trigger schema creation.
  - Opens a connection to the SQLite database.
  - Pings the database to verify connectivity.
  - If installation is required, executes the schema SQL to create the scheduler table.

Returns:
- nil on successful initialization.
- An error describing the failure if any step (opening, pinging, or schema execution) fails.
*/
func Init(dbFile string) error {
	if envDBFile := os.Getenv("TODO_DBFILE"); envDBFile != "" {
		dbFile = envDBFile
	}

	stat, err := os.Stat(dbFile)
	var install bool
	if err != nil {
		install = true
	} else if stat.Size() == 0 {
		install = true
	}

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("error open database: %v", err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("error pinging database: %v", err)
	}

	if install {
		_, err = db.Exec(schema)
		if err != nil {
			return fmt.Errorf("error executing schema: %v", err)
		}
	}

	return nil
}

/*
CloseDB gracefully closes the active database connection.

Behavior:
- Checks whether the global db variable is initialized.
- Calls Close() on the database handle if it’s not nil.
- Does nothing if no database connection is currently open.
*/
func CloseDB() {
	if db != nil {
		db.Close()
	}
}
