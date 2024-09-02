/*
Package main initializes and starts the Watchlist API application.

The `main` function calls `watchlist.Run` with command-line arguments to start the application.
*/

package main

import (
	"github.com/k4sper1love/watchlist-api/internal/watchlist"
	"os"
)

func main() {
	// Start the application with command-line arguments.
	watchlist.Run(os.Args)
}
