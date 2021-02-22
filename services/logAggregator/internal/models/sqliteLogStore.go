package models

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"

	// sqlite drivers
	_ "github.com/mattn/go-sqlite3"
)

const (
	driverName = "sqlite3"
)

// SQLiteLogStore is a SQLite backed implementation
// of LogStore
type SQLiteLogStore struct {
	LogStore
	*sql.DB
}

// NewSQLiteLogStore retuns a new SQLiteLogStore
// for a given dsn.
// if initDB, the database schema will be initialized.
// Logs fatally if fails.
func NewSQLiteLogStore(dsn string, initDB bool) *SQLiteLogStore {
	// open database file
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		log.Fatalf("Failed to open SQLite DB: %v", err)
	}
	// create store
	store := &SQLiteLogStore{DB: db}
	if initDB {
		// initialize DB schema
		store.initializeDB()
	}

	return store
}

// Log stores a new LogEntry in the database
func (s *SQLiteLogStore) Log(newEntry LogEntry) *LogError {
	// prepare insert
	ins, err := s.DB.Prepare("INSERT INTO Logs (correlationID, timeUnix, service, statusCode, notes) VALUES (?,?,?,?,?)")
	if err != nil {
		return &LogError{Err: err, Code: 500}
	}

	// execute insert
	_, err = ins.Exec(
		newEntry.CorrelationID.String(),
		newEntry.TimeUnix,
		newEntry.Service,
		newEntry.StatusCode,
		newEntry.Notes,
	)
	if err != nil {
		return &LogError{Err: err, Code: 500}
	}

	return nil
}

// Query returns LogEntries matching a LogQuery.
// TODO: UNIMPLEMENTED
func (s *SQLiteLogStore) Query(query LogQuery) ([]LogEntry, *QueryError) {

	// generate SQL based on the LogQuery
	stmt, args, err := s.generateSQL(query)
	if err != nil {
		return nil, &QueryError{Err: err, Code: http.StatusInternalServerError}
	}

	// Execute the query
	results, err := stmt.Query(args...)
	if err != nil {
		return nil, &QueryError{Err: err, Code: http.StatusInternalServerError}
	}

	entries := make([]LogEntry, 0)

	// Scan the query results into a slice of LogEntry
	for results.Next() {
		var entry LogEntry
		err = results.Scan(
			&entry.CorrelationID,
			&entry.TimeUnix,
			&entry.Service,
			&entry.StatusCode,
			&entry.Notes,
		)
		if err != nil {
			return nil, &QueryError{Err: err, Code: http.StatusInternalServerError}
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

// generateSQL returns a prepared SQL statement and a list of arguments
// to execute against the Logs table based on a LogQuery
func (s *SQLiteLogStore) generateSQL(query LogQuery) (*sql.Stmt, []interface{}, error) {

	var argsList []interface{}
	wheres := make([]string, 0)

	// select all columns except logID
	selectClause := "SELECT correlationID, timeUnix, service, statusCode, notes FROM Logs "

	// Add predicates and args for all LogQuery fields (that aren't default/nil)
	if query.CorrelationID != uuid.Nil {
		wheres = append(wheres, "correlationID = ?")
		argsList = append(argsList, query.CorrelationID.String())
	}
	if query.TimeUnixStart != 0 {
		wheres = append(wheres, "timeUnix >= ?")
		argsList = append(argsList, query.TimeUnixStart)
	}
	if query.TimeUnixStop != 0 {
		wheres = append(wheres, "timeUnix <= ?")
		argsList = append(argsList, query.TimeUnixStop)
	}
	if query.Service != "" {
		wheres = append(wheres, "service = ?")
		argsList = append(argsList, query.Service)
	}
	if query.StatusCodeStart != 0 {
		wheres = append(wheres, "statusCode >= ?")
		argsList = append(argsList, query.StatusCodeStart)
	}
	if query.StatusCodeStop != 0 {
		wheres = append(wheres, "statusCode <= ?")
		argsList = append(argsList, query.StatusCodeStop)
	}

	// any ordering or limit clauses
	trailers := " ORDER BY timeUnix"

	// Build where clauses
	var whereClause string
	if len(wheres) > 0 {
		whereClause = "WHERE " + strings.Join(wheres, " AND ")
	}

	// Build the full query
	fullquery := selectClause + whereClause + trailers

	// prepare the statment and return
	stmt, err := s.Prepare(fullquery)
	if err != nil {
		return nil, nil, err
	}

	return stmt, argsList, nil
}

// Initializes the Users table in SQLite syntax, for testing mostly
func (s *SQLiteLogStore) initializeDB() {

	log.Println("Initializing database...")
	stmt, err := s.DB.Prepare(`
		CREATE TABLE IF NOT EXISTS Logs (
			logID INTEGER PRIMARY KEY AUTOINCREMENT,
			correlationID TEXT NOT NULL,
			timeUnix INTEGER NOT NULL,
			service TEXT NOT NULL,
			statusCode INTEGER NOT NULL,
			notes TEXT
		)`)

	if err != nil {
		log.Fatalf("Error initializing DB: %v", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Fatalf("Error initializing DB: %v", err)
	}
	log.Println("Database initialized!")
}
