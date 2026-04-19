package booking

import "sync"

type ConcurrentStore struct {
	bookings map[string]Booking
	sync.RWMutex
}

func NewConcurrentStore() *ConcurrentStore {
	return &ConcurrentStore{bookings: map[string]Booking{}}
}
func (cs *ConcurrentStore) Book(b Booking) error {
	cs.Lock()
	defer cs.Unlock()
	if _, exists := cs.bookings[b.SeatID]; exists {
		return OcuppedSeatError
	}
	cs.bookings[b.ID] = b
	return nil
}
func (cs *ConcurrentStore) ListBookings(movieID string) []Booking {
	cs.RLock()
	defer cs.RUnlock()
	var bookings []Booking
	for _, book := range cs.bookings {
		if movieID == book.MovieID {
			bookings = append(bookings, book)
		}
	}
	return bookings
}
