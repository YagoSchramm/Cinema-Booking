package booking

type MemoryStore struct {
	bookings map[string]Booking
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{bookings: map[string]Booking{}}
}
func (ms *MemoryStore) Book(b Booking) error {
	if _, exists := ms.bookings[b.SeatID]; exists {
		return OcuppedSeatError
	}
	ms.bookings[b.ID] = b
	return nil
}
func (ms *MemoryStore) ListBookings(movieID string) []Booking {
	var bookings []Booking
	for _, book := range ms.bookings {
		if movieID == book.MovieID {
			bookings = append(bookings, book)
		}
	}
	return bookings
}
