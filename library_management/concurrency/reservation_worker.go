package concurrency

import (
	"fmt"

	"library_management/services"
)

// ReservationRequest represents a single reservation attempt.
type ReservationRequest struct {
	BookID     int
	MemberID   int
	ResultChan chan error
}

// ReservationWorker processes reservation requests concurrently using
// a channel and goroutines.
type ReservationWorker struct {
	Lib      *services.Library
	Requests chan ReservationRequest
}

// NewReservationWorker creates a new worker and starts its internal loop.
func NewReservationWorker(lib *services.Library, bufferSize int) *ReservationWorker {
	worker := &ReservationWorker{
		Lib:      lib,
		Requests: make(chan ReservationRequest, bufferSize),
	}

	// Start the main processing goroutine
	go worker.start()

	return worker
}

// start listens on the Requests channel and processes each reservation.
// Each reservation is handled in its own goroutine to allow concurrency.
func (w *ReservationWorker) start() {
	for req := range w.Requests {
		go func(r ReservationRequest) {
			err := w.Lib.ReserveBook(r.BookID, r.MemberID)
			r.ResultChan <- err
			close(r.ResultChan)
		}(req)
	}
}

// SubmitReservation enqueues a reservation request and returns a channel
// that will receive the result error (nil if successful).
func (w *ReservationWorker) SubmitReservation(bookID, memberID int) <-chan error {
	result := make(chan error, 1)

	w.Requests <- ReservationRequest{
		BookID:     bookID,
		MemberID:   memberID,
		ResultChan: result,
	}

	return result
}

// helper to demonstrate multiple members trying to reserve the same book at the same time.
func (w *ReservationWorker) SimulateConcurrentReservations(bookID int, memberIDs []int) {
	for _, memberID := range memberIDs {
		go func(mID int) {
			errChan := w.SubmitReservation(bookID, mID)
			if err := <-errChan; err != nil {
				fmt.Printf("[Simulation] Member %d failed to reserve book %d: %v\n", mID, bookID, err)
			} else {
				fmt.Printf("[Simulation] Member %d successfully reserved book %d\n", mID, bookID)
			}
		}(memberID)
	}
}
