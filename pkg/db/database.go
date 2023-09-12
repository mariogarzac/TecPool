package db

import (
    "database/sql"
    "log"

    _ "github.com/go-sql-driver/mysql"
)

var (
    db  *sql.DB
    err error
)

// InitDB initializes the database connection pool.
func InitDB(connectionString string) {
    db, err = sql.Open("mysql", connectionString)
    if err != nil {
        log.Fatal(err)
    }

    // Test the database connection
    err = db.Ping()
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Connected to the database")
}

// GetDB returns the database connection pool.
func GetDB() *sql.DB {
    return db
}

// CloseDB closes the database connection.
func CloseDB() {
    if db != nil {
        db.Close()
        log.Println("Database connection closed")
    }
}
