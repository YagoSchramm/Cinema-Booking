package main

import (
	"log"
	"net/http"

	"github.com/YagoSchramm/Cinema-Booking/internal/adapters"
	"github.com/YagoSchramm/Cinema-Booking/internal/booking"
	"github.com/YagoSchramm/Cinema-Booking/internal/utils"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /movies", listMovies)
	mux.Handle("GET /", http.FileServer(http.Dir("static")))
	client := adapters.NewClient("localhost:6379")
	srv := booking.NewService(booking.NewRedisStore(client))
	bookingHandler := booking.NewHandler(*srv)
	mux.HandleFunc("GET /movies/{movieID}/seats", bookingHandler.ListSeats)
	mux.HandleFunc("POST /movies/{movieID}/seats/{seatID}/hold", bookingHandler.HoldSeat)
	mux.HandleFunc("PUT /sessions/{sessionID}/confirm", bookingHandler.ConfirmSession)
	mux.HandleFunc("DELETE /sessions/{sessionID}", bookingHandler.ReleaseSession)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

var movies = []movieResponse{
	{ID: "inception", Title: "Inception", Rows: 5, SeatsPerRow: 8},
	{ID: "dune", Title: "Dune: Part Two", Rows: 4, SeatsPerRow: 6},
}

func listMovies(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, movies)
}

type movieResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Rows        int    `json:"rows"`
	SeatsPerRow int    `json:"seats_per_row"`
}
