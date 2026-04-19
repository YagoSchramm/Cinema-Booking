package booking

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var defaultHodlTTL = time.Minute * 2

type RedisStore struct {
	rdb *redis.Client
}

func NewRedisStore(rdb *redis.Client) *RedisStore {
	return &RedisStore{rdb: rdb}
}

func (rs *RedisStore) Book(b Booking) (Booking, error) {
	session, err := rs.hold(b)
	if err != nil {
		return Booking{}, err
	}
	log.Printf("Session Booked %v", session)
	return session, nil
}
func (rs *RedisStore) ListBookings(movieID string) []Booking {
	pattern := fmt.Sprintf("seat:%s:*", movieID)
	var sessions []Booking

	ctx := context.Background()

	iter := rs.rdb.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		val, err := rs.rdb.Get(ctx, iter.Val()).Result()
		if err != nil {
			continue
		}
		session, err := parseSession(val)
		if err != nil {
			continue
		}
		sessions = append(sessions, session)
	}

	return sessions
}
func (rs *RedisStore) Confirm(ctx context.Context, sessionID string, userID string) (Booking, error) {
	session, sk, err := rs.getSession(ctx, sessionID, userID)
	if err != nil {
		return Booking{}, err
	}
	rs.rdb.Persist(ctx, sk)
	rs.rdb.Persist(ctx, sessionKey(sessionID))
	session.Status = "confirmed"
	data := Booking{
		ID:      session.ID,
		MovieID: string(session.MovieID),
		SeatID:  session.SeatID,
		UserID:  session.UserID,
		Status:  "confirmed",
	}
	val, _ := json.Marshal(data)
	rs.rdb.Set(ctx, sk, val, 0)
	return session, err
}
func (rs *RedisStore) Release(ctx context.Context, sessionID string, userID string) error {
	_, sk, err := rs.getSession(ctx, sessionID, userID)
	if err != nil {
		return err
	}
	rs.rdb.Del(ctx, sk, sessionKey(sessionID))
	return nil
}
func sessionKey(id string) string {
	return fmt.Sprintf("session:%s", id)
}
func (rs *RedisStore) getSession(ctx context.Context, sessionID string, userID string) (Booking, string, error) {
	sk, err := rs.rdb.Get(ctx, sessionKey(sessionID)).Result()
	if err != nil {
		return Booking{}, "", err
	}
	val, err := rs.rdb.Get(ctx, sk).Result()
	if err != nil {
		return Booking{}, "", err
	}
	session, err := parseSession(val)
	if err != nil {
		return Booking{}, "", err
	}
	return session, sk, err

}
func (rs *RedisStore) hold(b Booking) (Booking, error) {
	id := uuid.NewString()
	now := time.Now()
	ctx := context.Background()
	key := fmt.Sprintf("seat:%s:%s", b.MovieID, b.SeatID)
	b.ID = id
	val, _ := json.Marshal(b)
	res := rs.rdb.SetArgs(ctx, key, val, redis.SetArgs{
		Mode: "NX",
		TTL:  defaultHodlTTL,
	})
	ok := res.Val() == "OK"
	if !ok {
		return Booking{}, OcuppedSeatError
	}
	rs.rdb.Set(ctx, sessionKey(id), key, defaultHodlTTL)
	return Booking{
		ID:        id,
		MovieID:   b.MovieID,
		SeatID:    b.SeatID,
		UserID:    b.UserID,
		Status:    "held",
		ExpiresAt: now.Add(defaultHodlTTL),
	}, nil
}
func parseSession(val string) (Booking, error) {
	var data Booking
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return Booking{}, err
	}
	return Booking{
		ID:      data.ID,
		MovieID: data.MovieID,
		SeatID:  data.SeatID,
		UserID:  data.UserID,
		Status:  data.Status,
	}, nil
}
