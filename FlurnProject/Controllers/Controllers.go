package Controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"FlurnProject/Models"
	"FlurnProject/db"
)

func GetAllSeats(c *gin.Context) {
	fmt.Println("IN: GetAllSeats")
	defer fmt.Println("IN: GetAllSeats")

	db, dberr := db.GetSQLConnection()
	if dberr != nil {
		fmt.Println("ERROR in DB ", dberr)
		// logginghelper.LogError("ERROR in DB ", dberr)
		return
	}

	var seatDetails []Models.SeatDetails
	query := `SELECT sl.id, sl.seat_identifier, sl.seat_class, 
	CASE 
	WHEN bookingId IS NOT NULL THEN TRUE 
	ELSE FALSE
	END AS is_booked
	FROM project.seatList sl 
	LEFT JOIN project.booking b 
	ON sl.id=b.seatId 
	ORDER BY seat_class;`
	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Error")
	}

	for rows.Next() {
		var record Models.SeatDetails
		err := rows.Scan(&record.Id, &record.SeatIdentifier, &record.SeatClass, &record.IsBooked)
		if err != nil {
			fmt.Println("Error while executing the query")
			return
		}
		seatDetails = append(seatDetails, record)
	}

	c.JSON(http.StatusOK, gin.H{"result": seatDetails})
}

func GetSeatDetailsById(c *gin.Context) {
	fmt.Println("IN: GetSeatDetailsById")
	defer fmt.Println("OUT: GetSeatDetailsById")

	seatId := c.Param("id")

	db, dberr := db.GetSQLConnection()
	if dberr != nil {
		fmt.Println("ERROR in DB ", dberr)
		// logginghelper.LogError("ERROR in DB ", dberr)
		return
	}

	var seatListData Models.SeatListModel
	// Query for fetching data from seatlist table of SQL
	query := "SELECT id, seat_identifier, seat_class FROM seatlist WHERE id=?;"
	// Here we check if seat is already booked
	readErr := db.QueryRow(query, seatId).Scan(&seatListData.Id, &seatListData.SeatIdentifier, &seatListData.SeatClass)
	if readErr != nil {
		fmt.Println(readErr)
		c.JSON(http.StatusInternalServerError, "Error")
	}
	// fmt.Println("seatListData: ", seatListData)

	// Fetching the price of the seatClass
	query = `SELECT
	CASE
		WHEN temp.percent < 40 THEN (SELECT CASE WHEN min_price = '' THEN normal_price ELSE min_price END AS price FROM project.seatpricing WHERE seat_class=?)
		WHEN temp.percent >= 40 && temp.percent <=60 THEN (SELECT CASE WHEN normal_price = '' THEN max_price ELSE normal_price END AS price FROM project.seatpricing WHERE seat_class=?)
		ELSE (SELECT CASE WHEN max_price = '' THEN normal_price ELSE max_price END AS price FROM project.seatpricing WHERE seat_class=?) 
	END AS price
	FROM (select ( (select COUNT('seat_class') 
				FROM project.booking AS b
				LEFT JOIN project.seatlist AS sl
				ON b.seatId=sl.id 
				WHERE sl.seat_class=? group by sl.seat_class)*100 / COUNT('seat_class') ) AS percent, sl.seat_class AS seat_class
			FROM project.seatlist AS sl WHERE sl.seat_class=? group by sl.seat_class) AS temp;`

	readErr = db.QueryRow(query, seatListData.SeatClass, seatListData.SeatClass, seatListData.SeatClass, seatListData.SeatClass, seatListData.SeatClass).Scan(&seatListData.Price)
	if readErr != nil {
		fmt.Println(readErr)
		c.JSON(http.StatusInternalServerError, "Error")
	}
	fmt.Println("seatListData: ", seatListData)

	c.JSON(http.StatusOK, gin.H{"result": seatListData})

}

func CreateBooking(c *gin.Context) {
	fmt.Println("IN: CreateBooking")
	defer fmt.Println("OUT: CreateBooking")

	var bookingDetail []Models.BookingDetail

	err := c.BindJSON(&bookingDetail)
	if err != nil {
		fmt.Println("Error while Binding the data")
		c.JSON(http.StatusInternalServerError, "Error")
		return
	}
	fmt.Println("bookingDetail: ", bookingDetail)

	db, dberr := db.GetSQLConnection()
	if dberr != nil {
		fmt.Println("ERROR in DB ", dberr)
		// logginghelper.LogError("ERROR in DB ", dberr)
		return
	}

	var seatArray []string
	for _, record := range bookingDetail {
		seatArray = append(seatArray, fmt.Sprint(record.SeatId))
	}

	var noOfSeatsNotAvailable int
	query := "SELECT COUNT(*) FROM booking WHERE seatId in (" + strings.Join(seatArray, ",") + ")"
	// Here we check if seat is already booked
	readErr := db.QueryRow(query).Scan(&noOfSeatsNotAvailable)
	if readErr != nil {
		fmt.Println(readErr)
		c.JSON(http.StatusInternalServerError, "Error")
	}
	fmt.Println("noOfSeatsNotAvailable: ", noOfSeatsNotAvailable)

	if noOfSeatsNotAvailable > 0 {
		c.JSON(http.StatusInternalServerError, "SEATS_NOT_AVAILABLE")
		return
	}

	sqlStr := "INSERT INTO project.booking (seatId, name, contact) VALUES "
	vals := []interface{}{}

	for _, row := range bookingDetail {
		sqlStr += "(?, ?, ?),"
		vals = append(vals, row.SeatId, row.Name, row.Contact)
	}

	sqlStr = sqlStr[:len(sqlStr)-1]
	//prepare the statement
	stmt, _ := db.Prepare(sqlStr)

	//format all vals at once
	_, _ = stmt.Exec(vals...)
	// fmt.Println("res: ", res)

	c.JSON(http.StatusOK, "Data saved successfully")

}

func GetBookingDetails(c *gin.Context) {
	fmt.Println("IN: CreateBooking")
	defer fmt.Println("OUT: CreateBooking")
	contactNo, ok := c.GetQuery("contact")
	if !ok {
		c.JSON(http.StatusInternalServerError, "Error")
	}

	db, dberr := db.GetSQLConnection()
	if dberr != nil {
		fmt.Println("ERROR in DB ", dberr)
		// logginghelper.LogError("ERROR in DB ", dberr)
		return
	}

	var bookingDetail []Models.BookingDetail

	rows, err := db.Query("SELECT bookingId, seatId, name, contact FROM booking WHERE contact=?;", contactNo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Error")
	}

	for rows.Next() {
		var record Models.BookingDetail
		err := rows.Scan(&record.BookingId, &record.SeatId, &record.Name, &record.Contact)
		if err != nil {
			fmt.Println("Error while executing the query")
			return
		}
		bookingDetail = append(bookingDetail, record)
	}

	c.JSON(http.StatusOK, gin.H{"result": bookingDetail})
}
