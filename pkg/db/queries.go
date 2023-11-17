package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
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

	stmt := "SELECT password FROM `users` WHERE email = ?"
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

	stmt := "SELECT user_id FROM `users` WHERE email = ?"
	row := db.QueryRow(stmt, email)
	err = row.Scan(&userId)

	if err != nil {
		return userId, err
	}

	return userId, nil
}

func GetNameByID(id int) (string, error) {
	var name string
	stmt := "SELECT fname FROM `users` WHERE user_id = ?"
	row := db.QueryRow(stmt, id)
	err = row.Scan(&name)

	if err != nil {
		log.Println("Error getting the user's name ", err)
		return name, err
	}

	return name, nil
}

func CreateTrip(carModel, licensePlate, departureTime, startLocation string, userId int) error {
	// insert user into table
	var insert *sql.Stmt

	// Actualizar la consulta para incluir start_location
	insert, err := db.Prepare("INSERT into `Trips` (`car_model`, `license_plate`, `departure_time`, `user_id`, `start_location`) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("Error preparing query", err)
		return err
	}
	defer insert.Close()

	// Incluir startLocation en Exec
	_, err = insert.Exec(carModel, licensePlate, departureTime, userId, startLocation)
	if err != nil {
		log.Println("Error inserting data", err)
		return err
	}

	// El resto del código permanece igual
	tripId, err := GetTripID()
	if err != nil {
		return err
	}

	//Create trip chat as well
	insertChat, err := db.Prepare("INSERT INTO `chats` (chat_id, trip_id, chat_name) VALUES (?, ?, ?)")
	if err != nil {
		log.Println("Error preparing chat insert query:", err)
		return err
	}
	defer insertChat.Close()

	tripIDstr := strconv.Itoa(tripId)
	chatName := "Group " + tripIDstr

	_, err = insertChat.Exec(tripId, tripId, chatName)
	if err != nil {
		log.Println("Error inserting chat into db:", err)
		return err
	}

	return nil
}

// Fetch most recent trips
func GetRecentTrips() ([]*models.Trip, error) {

	stmt := `SELECT trip_id, car_model, departure_time, license_plate, start_location
    FROM Trips ORDER BY trip_id DESC LIMIT 4`

	var trips []*models.Trip
	rows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}

	for rows.Next() {

		var (
			tripID        int
			carModel      string
			departureTime []uint8
			licensePlate  string
			startLocation string
		)

		if err := rows.Scan(&tripID, &carModel, &departureTime, &licensePlate, &startLocation); err != nil {
			log.Println("Error parsing recent trips", err)
			return nil, err
		}

		trip := models.Trip{
			TripID:        tripID,
			CarModel:      carModel,
			Date:          string(departureTime),
			LicensePlate:  licensePlate,
			StartLocation: startLocation,
		}

		trips = append(trips, &trip)
	}

	defer rows.Close()

	return trips, nil
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
func SearchTripsByDepartureTime(departureTime time.Time) ([]*models.Trip, error) {
	// Convert the departureTime to a string in the format YYYY-MM-DD HH:MM:SS
	formattedTime := departureTime.Format("2006-01-02 15:04:05")

	// Modify the SQL query to format the departure_time column in the same format
	stmt := fmt.Sprintf("SELECT trip_id, car_model, license_plate, departure_time, user_id, start_location FROM Trips WHERE DATE_FORMAT(departure_time, '%%Y-%%m-%%d %%H:%%i:%%s') = '%s'", formattedTime)

	rows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var trips []*models.Trip

	for rows.Next() {
		var (
			tripID        int
			carModel      string
			licensePlate  string
			departureTime []uint8
			userID        int
			startLocation string
		)

		if err := rows.Scan(&tripID, &carModel, &licensePlate, &departureTime, &userID, &startLocation); err != nil {
			log.Println("Error fetching trips by time", err)
		}
		trip := models.Trip{
			TripID:        tripID,
			CarModel:      carModel,
			LicensePlate:  licensePlate,
			Date:          string(departureTime),
			UserID:        userID,
			StartLocation: startLocation,
		}

		trips = append(trips, &trip)
	}

	return trips, nil
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
func GetUserTrips(userId int) (map[int]*models.Trip, error) {

	stmt := `SELECT t.trip_id, t.car_model, t.departure_time, t.user_id, t.license_plate, t.image_url, t.start_location
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
			tripId        int
			carModel      string
			departureTime []uint8
			userId        int
			licensePlate  string
			image         sql.NullString
			startLocation string
		)

		if err := rows.Scan(&tripId, &carModel, &departureTime, &userId, &licensePlate, &image, &startLocation); err != nil {
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

		var imageURL string

		if !image.Valid {
			imageURL = GetCarImage(carModel)

			if err := SaveImage(imageURL, tripId); err != nil {
				log.Fatal("Error saving image", err)
			}
		} else {
			imageURL = GetImageFromDB(tripId)

			if err != nil {
				log.Fatal("Error getting image from db", err)
			}
		}

		chatName, err := GetTripTitle(tripId)

		trip := models.Trip{
			TripID:        tripId,
			CarModel:      carModel,
			Date:          date,
			Time:          time,
			UserID:        uid,
			LicensePlate:  licensePlate,
			ChatName:      chatName,
			Image:         imageURL,
			StartLocation: startLocation,
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

func GetImageFromDB(tripID int) string {
	stmt := "select image_url from trips where trip_id = ?"

	row := db.QueryRow(stmt, tripID)
	if err != nil {
		log.Println("Error getting image from db:", err)
		return ""
	}

	var imageURL string
	err := row.Scan(&imageURL)

	if err != nil {
		log.Println("Error occured", err)
		return ""
	}

	return imageURL
}

type Image struct {
	Photos []struct {
		Src struct {
			Tiny string `json:"tiny"`
		} `json:"src"`
	} `json:"photos"`
}

func GetCarImage(carModel string) string {
	apiURL := "https://api.pexels.com/v1/search?query=" + carModel
	apiKey := "fqT1IUpPiOY3OY0edw8Gbmb8xpDy1eE9XCZsmgvRvKcIqzrkcBP3Z9CF"

	// Create a new HTTP client
	client := &http.Client{}

	// Create a new GET request
	request, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Println("Error creating request:", err)
		return ""
	}

	// Add the Authorization header
	request.Header.Add("Authorization", apiKey)

	// Make the request
	response, err := client.Do(request)
	if err != nil {
		log.Println("Error making the request:", err)
		return ""
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Error reading the response body:", err)
		return ""
	}

	// Print the API response
	var image Image

	// Unmarshal JSON
	err = json.Unmarshal(body, &image)
	if err != nil {
		log.Println("Error unmarshaling JSON:", err)
		return ""
	}

	// Verifica si hay al menos una foto en el slice
	if len(image.Photos) > 0 {
		return image.Photos[0].Src.Tiny
	} else {
		log.Println("No photos found for car model:", carModel)
		return ""
	}
}

func SaveImage(imageUrl string, tripID int) error {

	stmt := "UPDATE `trips` SET image_url = ? WHERE trip_id = ?"

	_, err = db.Exec(stmt, imageUrl, tripID)

	if err != nil {
		log.Println("Error inserting data", err)
		return err
	}

	return nil
}

func GetParticipants(tripId int) ([]*models.Users, error) {

	var users []*models.Users

	stmt := `SELECT u.fname, u.lname, u.phone_number
    FROM users u JOIN trip_participants tp 
    ON u.user_id = tp.user_id 
    WHERE tp.trip_id = ?;`

	rows, err := db.Query(stmt, tripId)
	if err != nil {
		log.Println("Error loading trip participants ", err)
		return users, err
	}

	for rows.Next() {
		var (
			fname       string
			lname       string
			phoneNumber string
		)

		if err := rows.Scan(&fname, &lname, &phoneNumber); err != nil {
			log.Println("Error saving name and phone number", err)
			return users, err
		}

		user := models.Users{
			Fname: fname,
			Lname: lname,
			Phone: phoneNumber,
		}

		users = append(users, &user)
	}

	return users, nil

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

func SaveMessage(chatId, userId int, msg string) error {

	// Insert the message
	insertMessage, err := db.Prepare("INSERT INTO `messages` (chat_id, user_id, message) VALUES (?, ?, ?)")
	if err != nil {
		log.Println("Error preparing message insert query:", err)
		return err
	}
	defer insertMessage.Close()

	_, err = insertMessage.Exec(chatId, userId, msg)
	if err != nil {
		log.Println("Error inserting message into db:", err)
		return err
	}

	return nil
}

func LoadMessages(tripId, uid int) ([]*models.Message, error) {

	loadedMessages := make([]*models.Message, 0)

	stmt := `SELECT m.chat_id, m.user_id, m.message, m.time FROM messages as m
    LEFT JOIN chats as c ON m.chat_id = c.trip_id
    WHERE m.chat_id = ?
    ORDER BY m.time ASC`

	rows, err := db.Query(stmt, tripId)
	if err != nil {
		log.Println("Error querying row: ", err)
		return nil, err
	}

	for rows.Next() {
		var (
			chatId  int
			userId  int
			message string
			time    []uint8
		)

		if err := rows.Scan(&chatId, &userId, &message, &time); err != nil {
			log.Println("Error scanning rows: ", err)
			return loadedMessages, err
		}

		msg := models.Message{
			Message: message,
			UserID:  userId,
			ChatID:  chatId,
			Time:    time,
			Self:    userId == uid,
		}

		loadedMessages = append(loadedMessages, &msg)

	}

	if err := rows.Err(); err != nil {
		log.Println("Error while exiting loop: ", err)
		return loadedMessages, err
	}

	return loadedMessages, nil
}

func GetTripTitle(tripId int) (string, error) {

	stmt := "SELECT `chat_name` FROM chats WHERE trip_id = ?"

	row := db.QueryRow(stmt, tripId)

	var chatName string
	if err := row.Scan(&chatName); err != nil {
		return chatName, err
	}

	return chatName, nil
}

func UpdateGroupName(name string, chatID int) error {

	stmt := "UPDATE chats SET chat_name = ? WHERE chat_id = ?"
	_, err := db.Exec(stmt, name, chatID)

	if err != nil {
		log.Println("Error updating chat name:", err)
		return err
	}

	return nil
}

func GetUserInfo(userID int) []*models.Users {

	stmt := `select fname, lname, email, phone_number
    from users where user_id = ?`

	row := db.QueryRow(stmt, userID)

	var (
		fname string
		lname string
		email string
		phone string
	)

	if err := row.Scan(&fname, &lname, &email, &phone); err != nil {
		log.Println("Error getting user info", err)
		return nil
	}

	user := models.Users{
		Fname: fname,
		Lname: lname,
		Email: email,
		Phone: phone,
	}

	var users []*models.Users
	users = append(users, &user)

	return users
}

func ChangePassword(userID int, oldPassword, newPassword string) error {

	var hash string

	stmt := "SELECT password FROM `users` WHERE user_id = ?"
	row := db.QueryRow(stmt, userID)
	err = row.Scan(&hash)

	if err != nil {
		log.Println("Error selecting hash in db by username", err)
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(oldPassword))
	if err != nil {
		log.Println("Passwords do not match. Not changing password", err)
		return err
	}

	log.Println("Passwords matched!")

	var hashedPassword []byte
	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error creating the hash ", err)
		return err
	}

	stmt = "UPDATE `users` SET password = ? WHERE user_id = ?"

	_, err = db.Exec(stmt, hashedPassword, userID)

	if err != nil {
		log.Println("Error inserting data", err)
		return err
	}

	return nil
}
