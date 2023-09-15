package db

import (
	"database/sql"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)


func RegisterUser(fname, lname, password, email, phone, dob string) error {

    // check if user already exists
    stmt := "SELECT email FROM Users WHERE email = ?"
    row := db.QueryRow(stmt, email)

    var isDuplicate string
    err = row.Scan(&isDuplicate)

    if err != sql.ErrNoRows {
        log.Println("Error user already exists", err)
        return err
    }

    // insert user into table
    var insert *sql.Stmt

    insert, err := db.Prepare("INSERT into `Users` (`fname`, `lname`, `password`, `email`, `phone_number`, `dob`) VALUES (?, ?, ?, ?, ?, ?)")
    if err != nil {
        log.Println("Error preparing query", err)
    }
    defer insert.Close()

    var hashedPassword []byte
    hashedPassword, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        log.Println("Error creating the hash ", err)
        return err
    }

    _, err = insert.Exec(fname,lname, hashedPassword, phone, email, dob)

    if err != nil {
        log.Println("Error inserting data" , err)
        return err
    }

    return nil
}

// ValidateUserInfo takes username and password and checks with the db to see if
// there is a match
func ValidateUserInfo(email, password string) error {
    var hash string

    stmt := "SELECT password FROM users WHERE email = ?"
    row := db.QueryRow(stmt, string(email))
    err = row.Scan(&hash)

    if err != nil {
        log.Println("Error selecting hash in db by username", err) 
        return err
    }

    err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    if err != nil {
        log.Println("Wrong things", err)
        return err
    }

    fmt.Println("Login success")

    return nil
}

//Fetch most recent trips
func GetRecentTrips() (*sql.Rows, error) {
    stmt := "SELECT * FROM Trips ORDER BY trip_id DESC LIMIT 4"
    rows, err := db.Query(stmt)
    if err != nil {
        return nil, err
    }
    return rows, nil
}

