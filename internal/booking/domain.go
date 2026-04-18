package booking

import "errors"

var (
	OcuppedSeatError = errors.New("The seat is already taken")
)

type Booking struct {
	ID      string
	MovieID string
	SeatID  string
	UserID  string
	Status  string
}
type BookingStore interface {
	Book(b Booking) error
	ListBookings(movieID string) []Booking
}
