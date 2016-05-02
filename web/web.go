package web

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/zkirill/openbazaar-etsy-import-golang/openbazaar-go"
)

var (
	homeHTML   []byte
	importHTML []byte
)

// Handle main page.
func homeHandler(w http.ResponseWriter, req *http.Request) {
	w.Write(homeHTML)
}

// Handle status requests.
func statusHandler(w http.ResponseWriter, req *http.Request) {
	p := struct {
		Count    int  `json:"count"`
		Success  bool `json:"success"`
		Finished bool `json:"finished"`
	}{
		Finished: openbazaar.ImportFinished,
		Success:  openbazaar.ImportFailed,
		Count:    openbazaar.ImportListingsParsed,
	}

	// Marshal into JSON.
	b, err := json.Marshal(p)

	if err != nil {
		log.WithField("error", err.Error()).
			Errorln("Failed to marshal values into JSON.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Setup file attachment.
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

// Handle import request.
func importHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseMultipartForm(0)
	file, _, err := req.FormFile("listings")
	if err != nil {
		log.WithField("error", err.Error()).
			Errorln("Failed to parse form.")
		w.Write([]byte(err.Error()))
		return
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.WithField("error", err.Error()).
			Errorln("Failed to read uploaded listings.")
		w.Write([]byte(err.Error()))
		return
	}
	username := req.FormValue("username")
	password := req.FormValue("password")
	origin := req.FormValue("origin")

	if len(origin) == 0 {
		origin = "ALL"
	}

	go openbazaar.Import(username, password, origin, data)

	w.Write(importHTML)
}

// Load HTML page into memory. Panics if file is not found.
func loadHTML(filename string, destination *[]byte) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.WithField("error", err.Error()).
			Panicln("Failed to read HTML file.")
	}
	*destination = data
}

// Preload assets needed by the web server.
func preloadAssets() {
	loadHTML("home.html", &homeHTML)
	loadHTML("import.html", &importHTML)
}

// Start creates a new router and begins handling web requests.
func Start() {
	preloadAssets()

	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/import", importHandler)
	r.HandleFunc("/status", statusHandler)
	http.Handle("/", r)
	log.Infoln("Accepting requests.")
	err := http.ListenAndServe(
		"localhost:8000",
		r,
	)
	if err != nil {
		panic("Error: " + err.Error())
	}
}
