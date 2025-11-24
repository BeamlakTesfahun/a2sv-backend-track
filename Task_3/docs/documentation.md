# Console-Based Library Management System (Go)

## 1. Overview

This project is a simple console-based Library Management System written in Go.

It demonstrates:

-   Structs (Book, Member)
-   Interfaces (LibraryManager)
-   Methods on structs
-   Use of slices and maps to manage in-memory data
-   Basic console input/output for user interaction

---

## 2. Features

-   Add a new book
-   Remove an existing book
-   Register a new library member
-   Borrow a book (if available)
-   Return a borrowed book
-   List all available books
-   List all books borrowed by a specific member

---

## 3. Architecture & Folder Structure

library_management/
├── main.go
├── controllers/
│ └── library_controller.go
├── models/
│ ├── book.go
│ └── member.go
├── services/
│ └── library_service.go
├── docs/
│ └── documentation.md
└── go.mod

### Models

Book:

-   ID (int)
-   Title (string)
-   Author (string)
-   Status (string) — "Available" or "Borrowed"

Member:

-   ID (int)
-   Name (string)
-   BorrowedBooks ([]Book)

### Service Layer

LibraryManager interface:

AddBook(book Book)
RemoveBook(bookID int)
BorrowBook(bookID int, memberID int) error
ReturnBook(bookID int, memberID int) error
ListAvailableBooks() []Book
ListBorrowedBooks(memberID int) []Book

The Library struct implements this interface and stores:

-   books in map[int]Book
-   members in map[int]Member

It contains the core logic for:

-   adding and removing books
-   borrowing and returning books
-   listing available and borrowed books

### Controller Layer

The LibraryController handles:

-   printing menu options
-   reading user input
-   invoking service methods

---

## 4. How to Run

From the project root (where go.mod is):

go run .

Then follow the menu instructions in the terminal.

---

## 5. Error Handling

The service layer includes error handling for:

-   Book not found
-   Member not found
-   Borrowing a book already borrowed
-   Returning a book not borrowed
-   Returning a book not borrowed by that member

Errors are returned to the controller and printed to the user.

---

## 6. Possible Future Improvements

-   Save data to JSON or a database
-   Add book search functionality
-   Add due dates or overdue fees
-   Add user roles (admin vs member)

---

## 7. Concurrency & Reservations

The system has been extended to support concurrent book reservations using:

-   Goroutines
-   Channels
-   Mutexes (`sync.Mutex`)
-   Timers (`time.Sleep`)

### Reservation Flow

1. The `LibraryManager` interface now includes:

    - `ReserveBook(bookID int, memberID int) error`

2. The `Library` struct uses:

    - `sync.Mutex` to guard shared state
    - `map[int]Book` for books
    - `map[int]Member` for members
    - `map[int]int` for active reservations (`bookID -> memberID`)

3. When `ReserveBook` is called:
    - If the book is **available**, it is marked as `"Reserved"` for the given member.
    - If the book is already **reserved** or **borrowed**, an error is returned.
    - A goroutine starts a **5-second timer**. If the book is still reserved for the same member after 5 seconds and has not been borrowed, the reservation is **automatically cancelled** and the status returns to `"Available"`.

### ReservationWorker

The `concurrency/reservation_worker.go` file defines:

-   `ReservationWorker`, which owns a `Requests` channel.
-   `ReservationRequest`, which includes:
    -   `BookID`
    -   `MemberID`
    -   `ResultChan` (channel used to send back the error result)

The worker:

-   Listens on the `Requests` channel.
-   For each request, starts a goroutine that calls `Library.ReserveBook`.
-   Sends the result (error or nil) back on `ResultChan`.

### Using the Worker

The controller:

-   Calls `worker.SubmitReservation(bookID, memberID)` to enqueue a new reservation.
-   Waits on the returned channel for the result.

A menu option also demonstrates **simulated concurrent reservations**, where multiple members try to reserve the same book at the same time. The `sync.Mutex` in `Library` ensures data consistency, and double reservations are prevented.
