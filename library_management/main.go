package main

import (
	"library_management/concurrency"
	"library_management/controllers"
	"library_management/services"
)

func main() {
	lib := services.NewLibrary()
	worker := concurrency.NewReservationWorker(lib, 10) // buffered channel size = 10

	controller := controllers.NewLibraryController(lib, worker)
	controller.SeedData()
	controller.Run()
}
