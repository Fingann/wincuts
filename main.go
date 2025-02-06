// main delegates application startup to the app package to maintain a clear separation of concerns.
// This ensures that main remains minimal and focused solely on bootstrapping the application.
package main

import (
	"log"
	"wincuts/app"
)

// main is the entry point of the application.
func main() {
	log.SetFlags(0)
	log.SetPrefix("error: ")

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
