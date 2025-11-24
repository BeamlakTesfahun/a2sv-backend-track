package main

import "library_management/controllers"

func main() {
	controller := controllers.NewLibraryController()
	controller.Run()
}
