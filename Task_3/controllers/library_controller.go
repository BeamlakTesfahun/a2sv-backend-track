package controllers

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"library_management/concurrency"
	"library_management/models"
	"library_management/services"
)

type LibraryController struct {
	service *services.Library
	worker  *concurrency.ReservationWorker
	reader  *bufio.Reader
}

func NewLibraryController(lib *services.Library, worker *concurrency.ReservationWorker) *LibraryController {
	return &LibraryController{
		service: lib,
		worker:  worker,
		reader:  bufio.NewReader(os.Stdin),
	}
}

func (lc *LibraryController) SeedData() {
	// seed (member + book)
	_ = lc.service.AddMember(models.Member{ID: 1, Name: "Default Member"})
	lc.service.AddBook(models.Book{ID: 1, Title: "Go Programming", Author: "Alan A. A. Donovan"})
	lc.service.AddBook(models.Book{ID: 2, Title: "The Pragmatic Programmer", Author: "Andrew Hunt"})
}

func (lc *LibraryController) Run() {
	for {
		lc.printMenu()
		choice := lc.readInt("Enter your choice: ")

		switch choice {
		case 1:
			lc.handleAddBook()
		case 2:
			lc.handleRemoveBook()
		case 3:
			lc.handleBorrowBook()
		case 4:
			lc.handleReturnBook()
		case 5:
			lc.handleListAvailableBooks()
		case 6:
			lc.handleListBorrowedBooks()
		case 7:
			lc.handleAddMember()
		case 8:
			lc.handleReserveBook()
		case 9:
			lc.handleSimulateConcurrentReservations()
		case 10:
			fmt.Println("Exiting... Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}

		fmt.Println()
	}
}

func (lc *LibraryController) printMenu() {
	fmt.Println("========== Library Management ==========")
	fmt.Println("1. Add a new book")
	fmt.Println("2. Remove an existing book")
	fmt.Println("3. Borrow a book")
	fmt.Println("4. Return a book")
	fmt.Println("5. List all available books")
	fmt.Println("6. List all borrowed books by a member")
	fmt.Println("7. Add a new member")
	fmt.Println("8. Reserve a book (concurrent)")
	fmt.Println("9. Simulate concurrent reservations")
	fmt.Println("10. Exit")
	fmt.Println("========================================")
}

func (lc *LibraryController) handleAddBook() {
	fmt.Println("---- Add New Book ----")
	id := lc.readInt("Enter book ID: ")
	title := lc.readString("Enter title: ")
	author := lc.readString("Enter author: ")

	book := models.Book{
		ID:     id,
		Title:  title,
		Author: author,
	}

	lc.service.AddBook(book)
	fmt.Println("Book added successfully!")
}

func (lc *LibraryController) handleRemoveBook() {
	fmt.Println("---- Remove Book ----")
	id := lc.readInt("Enter book ID to remove: ")
	lc.service.RemoveBook(id)
	fmt.Println("If the book existed, it has been removed.")
}

func (lc *LibraryController) handleBorrowBook() {
	fmt.Println("---- Borrow Book ----")
	bookID := lc.readInt("Enter book ID: ")
	memberID := lc.readInt("Enter member ID: ")

	err := lc.service.BorrowBook(bookID, memberID)
	if err != nil {
		fmt.Println("Error borrowing book:", err)
		return
	}

	fmt.Println("Book borrowed successfully!")
}

func (lc *LibraryController) handleReturnBook() {
	fmt.Println("---- Return Book ----")
	bookID := lc.readInt("Enter book ID: ")
	memberID := lc.readInt("Enter member ID: ")

	err := lc.service.ReturnBook(bookID, memberID)
	if err != nil {
		fmt.Println("Error returning book:", err)
		return
	}

	fmt.Println("Book returned successfully!")
}

func (lc *LibraryController) handleListAvailableBooks() {
	fmt.Println("---- Available Books ----")
	books := lc.service.ListAvailableBooks()
	if len(books) == 0 {
		fmt.Println("No available books.")
		return
	}

	for _, b := range books {
		fmt.Printf("ID: %d | Title: %s | Author: %s | Status: %s\n",
			b.ID, b.Title, b.Author, b.Status)
	}
}

func (lc *LibraryController) handleListBorrowedBooks() {
	fmt.Println("---- Borrowed Books By Member ----")
	memberID := lc.readInt("Enter member ID: ")

	_, exists := lc.service.GetMember(memberID)
	if !exists {
		fmt.Println("Member not found.")
		return
	}

	books := lc.service.ListBorrowedBooks(memberID)
	if len(books) == 0 {
		fmt.Println("This member has no borrowed books.")
		return
	}

	for _, b := range books {
		fmt.Printf("ID: %d | Title: %s | Author: %s | Status: %s\n",
			b.ID, b.Title, b.Author, b.Status)
	}
}

func (lc *LibraryController) handleAddMember() {
	fmt.Println("---- Add New Member ----")
	id := lc.readInt("Enter member ID: ")
	name := lc.readString("Enter member name: ")

	member := models.Member{
		ID:   id,
		Name: name,
	}

	if err := lc.service.AddMember(member); err != nil {
		fmt.Println("Error adding member:", err)
		return
	}

	fmt.Println("Member added successfully!")
}

func (lc *LibraryController) handleReserveBook() {
	fmt.Println("---- Reserve Book (Concurrent) ----")
	bookID := lc.readInt("Enter book ID: ")
	memberID := lc.readInt("Enter member ID: ")

	errChan := lc.worker.SubmitReservation(bookID, memberID)
	err := <-errChan

	if err != nil {
		fmt.Println("Error reserving book:", err)
		return
	}

	fmt.Println("Reservation request processed successfully!")
}

func (lc *LibraryController) handleSimulateConcurrentReservations() {
	fmt.Println("---- Simulate Concurrent Reservations ----")
	bookID := lc.readInt("Enter book ID for simulation: ")

	memberIDs := []int{1, 2, 3, 4, 5}
	fmt.Printf("Simulating members %v trying to reserve book %d at the same time...\n", memberIDs, bookID)
	lc.worker.SimulateConcurrentReservations(bookID, memberIDs)
}

// input helpers

func (lc *LibraryController) readString(prompt string) string {
	fmt.Print(prompt)
	text, _ := lc.reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func (lc *LibraryController) readInt(prompt string) int {
	for {
		input := lc.readString(prompt)
		value, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid number, please try again.")
			continue
		}
		return value
	}
}
