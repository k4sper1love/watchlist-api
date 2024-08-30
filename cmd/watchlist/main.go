package main

import (
	"github.com/k4sper1love/watchlist-api/internal/watchlist"
	"os"
)

func main() {
	watchlist.Run(os.Args)
}
