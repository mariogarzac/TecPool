package db

import (
	"database/sql"
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

    return nil
}

func GetNameByEmail(email string) (string, error) {
    var name string
    stmt := "SELECT fname FROM users WHERE email = ?"
    row := db.QueryRow(stmt, email)
    err = row.Scan(&name)

    if err != nil {
        log.Println("Error getting the user's name ", err)
        return name, err
    }

    return name, nil
}

func GetUserId(email string) (int, error) {
    var userId int

    stmt := "SELECT user_id FROM users WHERE email = ?"
    row := db.QueryRow(stmt, email)
    err = row.Scan(&userId)

    
    if err != nil {
        return userId, err
    }

    return userId, nil
}

func CreateTrip(carModel, licensePlate, departureTime string, userId int) error {
    
    // insert user into table
    var insert *sql.Stmt

    insert, err := db.Prepare("INSERT into `Trips` (`car_model`, `license_plate`, `departure_time`, `user_id`) VALUE (?, ?, ?, ?)")
    if err != nil {
        log.Println("Error preparing query", err)
    }
    defer insert.Close()

    if err != nil {
        log.Println("Error creating ride", err)
        return err
    }

    _, err = insert.Exec(carModel, licensePlate, departureTime, userId)

    if err != nil {
        log.Println("Error inserting data" , err)
        return err
    }

    return nil
}
