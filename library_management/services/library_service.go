package services

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"library_management/models"
)

const (
	StatusAvailable = "Available"
	StatusBorrowed  = "Borrowed"
	StatusReserved  = "Reserved" // for reservations
)

// library operations
type LibraryManager interface {
	AddBook(book models.Book)
	RemoveBook(bookID int)
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) []models.Book

	// concurrent reservations
	ReserveBook(bookID int, memberID int) error
}

// Library implements LibraryManager
type Library struct {
	mu           sync.Mutex
	books        map[int]models.Book
	members      map[int]models.Member
	reservations map[int]int // bookID -> memberID
}

// new Library instance
func NewLibrary() *Library {
	return &Library{
		books:        make(map[int]models.Book),
		members:      make(map[int]models.Member),
		reservations: make(map[int]int),
	}
}

// member-related helpers

func (l *Library) AddMember(member models.Member) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, exists := l.members[member.ID]; exists {
		return fmt.Errorf("member with ID %d already exists", member.ID)
	}
	member.BorrowedBooks = []models.Book{}
	l.members[member.ID] = member
	return nil
}

func (l *Library) GetMember(id int) (models.Member, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	m, ok := l.members[id]
	return m, ok
}

// library manager implementation

func (l *Library) AddBook(book models.Book) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if book.Status == "" {
		book.Status = StatusAvailable
	}
	l.books[book.ID] = book
}

func (l *Library) RemoveBook(bookID int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	delete(l.books, bookID)
	delete(l.reservations, bookID)
}

func (l *Library) BorrowBook(bookID int, memberID int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	book, ok := l.books[bookID]
	if !ok {
		return fmt.Errorf("book with ID %d not found", bookID)
	}

	member, ok := l.members[memberID]
	if !ok {
		return fmt.Errorf("member with ID %d not found", memberID)
	}

	if book.Status == StatusBorrowed {
		return errors.New("book is already borrowed")
	}

	// only the reserving member can borrow it if reserved
	if book.Status == StatusReserved {
		reservingMemberID, reserved := l.reservations[bookID]
		if !reserved {
			return errors.New("book is reserved but reservation data is missing")
		}
		if reservingMemberID != memberID {
			return errors.New("book is reserved by another member")
		}
		delete(l.reservations, bookID)
	}
	
	book.Status = StatusBorrowed
	l.books[bookID] = book

	// Add book to member's borrowed list
	member.BorrowedBooks = append(member.BorrowedBooks, book)
	l.members[memberID] = member

	return nil
}

func (l *Library) ReturnBook(bookID int, memberID int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	book, ok := l.books[bookID]
	if !ok {
		return fmt.Errorf("book with ID %d not found", bookID)
	}

	if book.Status != StatusBorrowed {
		return errors.New("book is not currently borrowed")
	}

	member, ok := l.members[memberID]
	if !ok {
		return fmt.Errorf("member with ID %d not found", memberID)
	}

	// Find and remove book from member's borrowed list.
	index := -1
	for i, b := range member.BorrowedBooks {
		if b.ID == bookID {
			index = i
			break
		}
	}

	if index == -1 {
		return errors.New("this member has not borrowed the specified book")
	}

	member.BorrowedBooks = append(member.BorrowedBooks[:index], member.BorrowedBooks[index+1:]...)
	l.members[memberID] = member

	// Update book status back to available.
	book.Status = StatusAvailable
	l.books[bookID] = book

	return nil
}

func (l *Library) ListAvailableBooks() []models.Book {
	l.mu.Lock()
	defer l.mu.Unlock()

	var available []models.Book
	for _, book := range l.books {
		if book.Status == StatusAvailable {
			available = append(available, book)
		}
	}
	return available
}

func (l *Library) ListBorrowedBooks(memberID int) []models.Book {
	l.mu.Lock()
	defer l.mu.Unlock()

	member, ok := l.members[memberID]
	if !ok {
		return []models.Book{}
	}
	// avoid external mutation issues
	result := make([]models.Book, len(member.BorrowedBooks))
	copy(result, member.BorrowedBooks)
	return result
}

// ReserveBook handles reservation logic.
// - If the book is available, it becomes "Reserved" for that member.
// - If already reserved or borrowed, it returns an error.
// - A goroutine automatically cancels the reservation after 5 seconds
//   if the book hasn't been borrowed.
func (l *Library) ReserveBook(bookID int, memberID int) error {
	l.mu.Lock()

	book, ok := l.books[bookID]
	if !ok {
		l.mu.Unlock()
		return fmt.Errorf("book with ID %d not found", bookID)
	}

	if book.Status == StatusBorrowed {
		l.mu.Unlock()
		return errors.New("book is already borrowed")
	}

	if book.Status == StatusReserved {
		l.mu.Unlock()
		return errors.New("book is already reserved")
	}

	if _, ok := l.members[memberID]; !ok {
		l.mu.Unlock()
		return fmt.Errorf("member with ID %d not found", memberID)
	}

	book.Status = StatusReserved
	l.books[bookID] = book
	l.reservations[bookID] = memberID

	fmt.Printf("Book %d reserved for member %d\n", bookID, memberID)

	go l.autoCancelReservation(bookID, memberID, 5*time.Second)

	l.mu.Unlock()
	return nil
}

// autoCancelReservation runs in a separate goroutine and unreserves
// the book if it hasn't been borrowed by the same member within the given duration.
func (l *Library) autoCancelReservation(bookID, memberID int, d time.Duration) {
	time.Sleep(d)

	l.mu.Lock()
	defer l.mu.Unlock()

	book, ok := l.books[bookID]
	if !ok {
		return
	}

	currentReservingMemberID, reserved := l.reservations[bookID]
	if !reserved {
		return
	}

	// Only cancel if:
	// - still reserved
	// - reserved for the same member
	if book.Status == StatusReserved && currentReservingMemberID == memberID {
		book.Status = StatusAvailable
		l.books[bookID] = book
		delete(l.reservations, bookID)
		fmt.Printf("Reservation for book %d by member %d has been auto-cancelled\n", bookID, memberID)
	}
}
