package main

import (
	"sync"

	"github.com/zkirill/openbazaar-etsy-import-golang/web"

	"github.com/skratchdot/open-golang/open"
)

const (
	listingsFile = "listings.csv"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	// Start the web server.
	go func() {
		web.Start()
	}()

	// Open website in browser.
	// NOTE: Possible race condition if web server doesn't come up quickly enough.
	open.Run("http://localhost:8000")

	wg.Wait()
}
