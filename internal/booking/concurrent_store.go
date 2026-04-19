package booking

import "sync"

type ConcurrentStore struct {
	bookings map[string]Booking
	sync.RWMutex
}

func NewConcurrentStore() *ConcurrentStore {
	return &ConcurrentStore{bookings: map[string]Booking{}}
}
func (ms *ConcurrentStore) Book(b Booking) error {
	ms.Lock()
	defer ms.Unlock()
	if _, exists := ms.bookings[b.SeatID]; exists {
		return OcuppedSeatError
	}
	ms.bookings[b.ID] = b
	return nil
}
func (ms *ConcurrentStore) ListBookings(movieID string) []Booking {
	ms.RLock()
	defer ms.RUnlock()
	var bookings []Booking
	for _, book := range ms.bookings {
		if movieID == book.MovieID {
			bookings = append(bookings, book)
		}
	}
	return bookings
}
