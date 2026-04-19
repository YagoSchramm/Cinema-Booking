package booking

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/YagoSchramm/Cinema-Booking/internal/utils"
)

type Handler struct {
	srv Service
}

func NewHandler(srv Service) *Handler {
	return &Handler{srv: srv}
}

type holdSeatRequest struct {
	UserID string `json:"user_id"`
}

func (h *Handler) ConfirmSession(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("sessionID")

	var req holdSeatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return
	}

	if req.UserID == "" {
		return
	}

	session, err := h.srv.ConfirmSeat(r.Context(), sessionID, req.UserID)
	if err != nil {
		return
	}

	utils.WriteJSON(w, http.StatusOK, sessionResponse{
		SessionID: session.ID,
		MovieID:   session.MovieID,
		SeatID:    session.SeatID,
		UserID:    req.UserID,
		Status:    session.Status,
	})
}

type sessionResponse struct {
	SessionID string `json:"session_id"`
	MovieID   string `json:"movie_id"`
	SeatID    string `json:"seat_id"`
	UserID    string `json:"user_id"`
	Status    string `json:"status"`
	ExpiresAt string `json:"expires_at,omitempty"`
}

func (h *Handler) ReleaseSession(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("sessionID")

	var req holdSeatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println(err)
		return
	}
	if req.UserID == "" {
		return
	}

	err := h.srv.ReleaseSeat(r.Context(), sessionID, req.UserID)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) HoldSeat(w http.ResponseWriter, r *http.Request) {
	movieID := r.PathValue("movieID")
	seatID := r.PathValue("seatID")

	var req holdSeatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println(err)
		return
	}

	data := Booking{
		UserID:  req.UserID,
		SeatID:  seatID,
		MovieID: movieID,
	}

	session, err := h.srv.Book(data)
	if err != nil {
		log.Println(err)
		return
	}

	type holdResponse struct {
		SessionID string `json:"session_id"`
		MovieID   string `json:"movieID"`
		SeatID    string `json:"seat_id"`
		ExpiresAt string `json:"expires_at"`
	}

	utils.WriteJSON(w, http.StatusCreated, holdResponse{
		SeatID:    seatID,
		MovieID:   session.MovieID,
		SessionID: session.ID,
		ExpiresAt: session.ExpiresAt.Format(time.RFC3339),
	})
}

func (h *Handler) ListSeats(w http.ResponseWriter, r *http.Request) {
	movieID := r.PathValue("movieID")

	bookings := h.srv.ListBooking(movieID)

	seats := make([]seatInfo, 0, len(bookings))
	for _, b := range bookings {
		seats = append(seats, seatInfo{
			SeatID:    b.SeatID,
			UserID:    b.UserID,
			Booked:    true,
			Confirmed: b.Status == "confirmed",
		})
	}

	utils.WriteJSON(w, http.StatusOK, seats)
}

type seatInfo struct {
	SeatID    string `json:"seat_id"`
	UserID    string `json:"user_id"`
	Booked    bool   `json:"booked"`
	Confirmed bool   `json:"confirmed"`
}
