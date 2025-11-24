package services

import (
	"errors"
	"fmt"

	"library_management/models"
)

const (
	StatusAvailable = "Available"
	StatusBorrowed  = "Borrowed"
)

// library operations
type LibraryManager interface {
	AddBook(book models.Book)
	RemoveBook(bookID int)
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) []models.Book
}

// Library implements LibraryManager
type Library struct {
	books   map[int]models.Book
	members map[int]models.Member
}

// new Library instance
func NewLibrary() *Library {
	return &Library{
		books:   make(map[int]models.Book),
		members: make(map[int]models.Member),
	}
}

// member-related helpers

func (l *Library) AddMember(member models.Member) error {
	if _, exists := l.members[member.ID]; exists {
		return fmt.Errorf("member with ID %d already exists", member.ID)
	}
	member.BorrowedBooks = []models.Book{}
	l.members[member.ID] = member
	return nil
}

func (l *Library) GetMember(id int) (models.Member, bool) {
	m, ok := l.members[id]
	return m, ok
}

// library manager implementation

func (l *Library) AddBook(book models.Book) {
	if book.Status == "" {
		book.Status = StatusAvailable
	}
	l.books[book.ID] = book
}

func (l *Library) RemoveBook(bookID int) {
	delete(l.books, bookID)
}

func (l *Library) BorrowBook(bookID int, memberID int) error {
	book, ok := l.books[bookID]
	if !ok {
		return fmt.Errorf("book with ID %d not found", bookID)
	}

	if book.Status != StatusAvailable {
		return errors.New("book is not available for borrowing")
	}

	member, ok := l.members[memberID]
	if !ok {
		return fmt.Errorf("member with ID %d not found", memberID)
	}

	book.Status = StatusBorrowed
	l.books[bookID] = book

	member.BorrowedBooks = append(member.BorrowedBooks, book)
	l.members[memberID] = member

	return nil
}

func (l *Library) ReturnBook(bookID int, memberID int) error {
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

	// Find the book in member's borrowed list
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

	// remove book from member
	member.BorrowedBooks = append(member.BorrowedBooks[:index], member.BorrowedBooks[index+1:]...)
	l.members[memberID] = member

	// mark book as available
	book.Status = StatusAvailable
	l.books[bookID] = book

	return nil
}

func (l *Library) ListAvailableBooks() []models.Book {
	var available []models.Book
	for _, book := range l.books {
		if book.Status == StatusAvailable {
			available = append(available, book)
		}
	}
	return available
}

func (l *Library) ListBorrowedBooks(memberID int) []models.Book {
	member, ok := l.members[memberID]
	if !ok {
		return []models.Book{}
	}
	return member.BorrowedBooks
}
