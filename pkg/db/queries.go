package db

import (
    "database/sql"
    "fmt"
    "log"
    "time"

    "github.com/mariogarzac/tecpool/pkg/models"
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

    _, err = insert.Exec(fname, lname, hashedPassword, phone, email, dob)

    if err != nil {
        log.Println("Error inserting data", err)
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
        log.Println("Wrong username or password", err)
        return err
    }

    return nil
}


func GetUserIDByEmail(email string) (int, error) {
    var userId int

    stmt := "SELECT user_id FROM users WHERE email = ?"
    row := db.QueryRow(stmt, email)
    err = row.Scan(&userId)

    if err != nil {
        return userId, err
    }

    return userId, nil
}

func GetNameByID(id int) (string, error) {
    var name string
    stmt := "SELECT fname FROM users WHERE user_id = ?"
    row := db.QueryRow(stmt, id)
    err = row.Scan(&name)

    if err != nil {
        log.Println("Error getting the user's name ", err)
        return name, err
    }

    return name, nil
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
        log.Println("Error inserting data", err)
        return err
    }

    return nil
}

// Fetch most recent trips
func GetRecentTrips() (*sql.Rows, error) {
    stmt := "SELECT * FROM Trips ORDER BY trip_id DESC LIMIT 4"
    rows, err := db.Query(stmt)
    if err != nil {
        return nil, err
    }

    defer rows.Close()

    return rows, nil
}

func GetTripID() (int, error) {
    var tripId int
    stmt := "SELECT trip_id FROM trips ORDER BY trip_id DESC LIMIT 1"

    row := db.QueryRow(stmt)
    err = row.Scan(&tripId)

    if err != nil {
        return 0, err
    }

    return tripId, nil
}

// SearchTripsByDepartureTime fetches trips from the database that match the given departure time
func SearchTripsByDepartureTime(departureTime time.Time) (*sql.Rows, error) {
    // Convert the departureTime to a string in the format YYYY-MM-DD HH:MM:SS
    formattedTime := departureTime.Format("2006-01-02 15:04:05")

    // Modify the SQL query to format the departure_time column in the same format
    stmt := fmt.Sprintf("SELECT trip_id, car_model, license_plate, departure_time, user_id FROM Trips WHERE DATE_FORMAT(departure_time, '%%Y-%%m-%%d %%H:%%i:%%s') = '%s'", formattedTime)

    rows, err := db.Query(stmt)
    if err != nil {
        return nil, err
    }

    defer rows.Close()

    return rows, nil
}

func DoesTripExist(tripId int) bool {

    stmt := "SELECT trip_id FROM trips WHERE trip_id = ?"
    row := db.QueryRow(stmt, tripId)

    var exists int
    err := row.Scan(&exists)

    if err != nil {
        if err == sql.ErrNoRows {
            log.Println("Trip does not exist", err)
            return false
        }
        log.Println("An error ocurred", err)
        return false
    }

    return true
}

func IsUserInTrip(tripId, userId int) bool {
    stmt := "SELECT user_id FROM trip_participants WHERE trip_id = ? AND user_id = ?"
    row := db.QueryRow(stmt, tripId, userId)

    var inTrip int
    err := row.Scan(&inTrip)

    if err != nil {
        if err == sql.ErrNoRows {
            return false
        }

        log.Println("An error occurred:", err)
        return false
    }

    return true
}
func GetUserTrips(userId int) (map[int]*models.Trip, error){

    stmt := `SELECT t.trip_id, t.car_model, t.departure_time, t.user_id, t.license_plate 
    FROM trips as t, trip_participants as tp
    WHERE tp.trip_id = t.trip_id 
    AND tp.user_id = ?;`

    rows, err := db.Query(stmt, userId)
    if err != nil {
        log.Println("Error querying row: ", err)
        return nil, err
    }

    defer rows.Close()

    tripMap, err := ProcessTrips(rows, userId)
    if err != nil {
        log.Println("An error occured while parsing rows: ", err)
        return tripMap, err
    }

    return tripMap, nil
}

func ProcessTrips(rows *sql.Rows, uid int) (map[int]*models.Trip, error) {
    // close rows when done
    defer rows.Close() 

    // loop through all the rows and store the results in a map
    tripMap := make(map[int]*models.Trip)

    for rows.Next() {
        var (
            tripId int
            carModel string
            departureTime []uint8
            userId int
            licensePlate string
        )

        if err := rows.Scan(&tripId, &carModel, &departureTime, &userId, &licensePlate); err != nil {
            log.Println("Error scanning rows: ", err)
            return tripMap, err
        }

        parsedTime, err := time.Parse("2006-01-02 15:04:05", string(departureTime))

        date := parsedTime.Format("2006-01-02")
        time := parsedTime.Format("15:04:05")    

        if err != nil {
            log.Println("Error parsing time: ", err)
            return tripMap, err
        }

        trip := models.Trip{
            TripID: tripId,
            CarModel: carModel,
            Date: date,
            Time: time,
            UserID: uid,
            LicensePlate: licensePlate,
        }

        // add trip to map
        tripMap[tripId] = &trip

    }

    if err := rows.Err(); err != nil {
        log.Println("Error while exiting loop: ", err)
        return tripMap, err
    }

    return tripMap, nil
}

func AddUserToTrip(userId, tripId uint64) error {

    insert, err := db.Prepare("INSERT into `trip_participants` (user_id, trip_id) VALUES (?, ?)")

    if err != nil {
        log.Println("Error preparing query", err)
        return err
    }
    defer insert.Close()

    _, err = insert.Exec(userId, tripId)

    if err != nil {
        log.Println("Error inserting data", err)
        return err
    }

    return nil
}


func SaveMessage(tripId, userId int, msg string) error {
    // Check if the chat exists
    stmt := "SELECT * FROM chats WHERE chat_id = ?"
    row := db.QueryRow(stmt, tripId)

    var chatID int
    var tripID int
    if err := row.Scan(&chatID, &tripID); err != nil {
        if err == sql.ErrNoRows{
            // Chat doesn't exist, create it
            insertChat, err := db.Prepare("INSERT INTO `chats` (chat_id, trip_id) VALUES (?, ?)")
            if err != nil {
                log.Println("Error preparing chat insert query:", err)
                return err
            }
            defer insertChat.Close()

            _, err = insertChat.Exec(tripId, tripId)
            if err != nil {
                log.Println("Error inserting chat into db:", err)
                return err
            }
        } else {
            log.Println("Error checking chat existence:", err)
            return err
        }
    }

    // Insert the message
    insertMessage, err := db.Prepare("INSERT INTO `messages` (chat_id, user_id, message) VALUES (?, ?, ?)")
    if err != nil {
        log.Println("Error preparing message insert query:", err)
        return err
    }
    defer insertMessage.Close()

    _, err = insertMessage.Exec(tripId, userId, msg)
    if err != nil {
        log.Println("Error inserting message into db:", err)
        return err
    }

    return nil
}

func LoadMessages(tripId, uid int) ([]*models.Message, error) {

    loadedMessages := make([]*models.Message ,0)

    stmt := `SELECT m.chat_id, m.user_id, m.message, m.time FROM messages as m
    LEFT JOIN chats as c ON m.chat_id = c.trip_id
    WHERE m.chat_id = ?
    ORDER BY m.time DESC`

    rows, err := db.Query(stmt, tripId)
    if err != nil {
        log.Println("Error querying row: ", err)
        return nil, err
    }

    for rows.Next() {
        var (
            chatId int 
            userId int 
            message string 
            time []uint8
        )

        if err := rows.Scan(&chatId, &userId, &message, &time); err != nil {
            log.Println("Error scanning rows: ", err)
            return loadedMessages, err
        }

        if err != nil {
            log.Println("Error parsing time: ", err)
            return loadedMessages, err
        }

        msg := models.Message{
            Message: message,
            UserID: userId,
            ChatID: chatId,
            Time: time,
            Self: userId == uid,
        }

        loadedMessages= append(loadedMessages, &msg)

    }

    if err := rows.Err(); err != nil {
        log.Println("Error while exiting loop: ", err)
        return loadedMessages, err
    }

    return loadedMessages, nil
}
