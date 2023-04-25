package Models

type BookingDetail struct {
	BookingId int    `json: "bookingId"`
	SeatId    int    `json: "seatId"`
	Name      string `json:"name"`
	Contact   int    `json:"contact"`
}

type SeatDetails struct {
	Id             int    `json:"id"`
	SeatIdentifier string `json:"seat_identifier"`
	SeatClass      string `json:"seat_class"`
	IsBooked       bool   `json:"is_booked"`
}

type SeatListModel struct {
	Id             int    `json:"id"`
	SeatIdentifier string `json:"seat_identifier"`
	SeatClass      string `json:"seat_class"`
	Price          string `json:"price"`
}
