package openbazaar

import (
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/zkirill/openbazaar-etsy-import-golang/row"
)

var (
	// ImportListingsParsed is the number of parsed listings.
	ImportListingsParsed int
	// ImportFinished is true when import has finished.
	ImportFinished bool
	// ImportFailed is true when when import has failed.
	ImportFailed bool
)

const (
	api = "http://localhost:18469/api/v1"
	// Maximum number of keywords per listing.
	maxKeywords = 10
)

// Contract represents an OpenBazaar contract.
type Contract struct {
	Title          string
	Image          string
	Tags           []string
	Price          string
	Description    string
	CurrencyCode   string
	ShippingOrigin string
	ShipsTo        []string
}

// Response represents the response from OpenBazaar server.
type response struct {
	Success bool `json:"success"`
}

// Client sets up the client logs in.
func Client(username string, password string) (client *http.Client, err error) {
	cjar, _ := cookiejar.New(nil)
	client = &http.Client{
		Jar: cjar,
	}
	err = Login(client, username, password)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// Login logs in with the username and password provided.
func Login(client *http.Client, username string, password string) error {
	resp, err := client.PostForm(api+"/login",
		url.Values{"username": {username}, "password": {password}})
	if err != nil {
		log.Fatal("Failed to login to OpenBazaar.")
		return err
	}
	dec := json.NewDecoder(resp.Body)
	var r response
	dec.Decode(&r)
	if !r.Success {
		return errors.New("Failed to login to OpenBazaar.")
	}
	return nil
}

// UploadImage uploads an image.
func UploadImage(client *http.Client, img string) (hash string, err error) {
	// Represents the response from OpenBazaar server.
	type imageUploadResponse struct {
		// Success represents the result of the upload.
		Success bool `json:"success"`
		// Hashes is the list of image hashes.
		Hashes []string `json:"image_hashes"`
	}

	resp, err := client.PostForm(api+"/upload_image",
		url.Values{
			"image": {img},
		})
	if err != nil {
		log.Fatal(err.Error())
		return "", err
	}
	dec := json.NewDecoder(resp.Body)
	var r imageUploadResponse
	err = dec.Decode(&r)
	if err != nil {
		log.Fatal(err.Error())
		return "", err
	}
	if !r.Success {
		return "", errors.New("Failed to upload image.")
	}
	return r.Hashes[0], nil
}

// PostContract posts a contract to OpenBazaar local server.
func PostContract(client *http.Client, contract Contract) error {
	// Represents the contract post response from OpenBazaar server.
	type response struct {
		// Success represents the result of the upload.
		Success bool `json:"success"`
		// Reason contains the error message in the event of an error.
		Reason string `json:"reason"`
	}
	v := url.Values{
		"title":                  {contract.Title},
		"expiration_date":        {""},
		"metadata_category":      {"physical good"},
		"description":            {contract.Description},
		"currency_code":          {contract.CurrencyCode},
		"price":                  {contract.Price},
		"process_time":           {"TBD"},
		"terms_conditions":       {""},
		"returns":                {""},
		"shipping_currency_code": {contract.CurrencyCode},
		"shipping_domestic":      {""},
		"shipping_international": {""},
		"category":               {""},
		"condition":              {""},
		"sku":                    {""},
		"images":                 {contract.Image},
		"free_shipping":          {"false"},
		"nsfw":                   {"false"},
		"shipping_origin":        {contract.ShippingOrigin},
	}
	// Add keywords.
	for idx, kw := range contract.Tags {
		if idx >= maxKeywords {
			break
		}
		v.Add("keywords", kw)
	}
	// Add shipping destinations.
	for _, sd := range contract.ShipsTo {
		v.Add("ships_to", sd)
	}
	resp, err := client.PostForm(api+"/contracts", v)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	dec := json.NewDecoder(resp.Body)
	var r response
	dec.Decode(&r)
	if !r.Success {
		return errors.New("Failed to post contract: " + r.Reason)
	}
	return nil
}

// Import imports data into OpenBazaar.
func Import(username string, password string, shipOrigin string, listings []byte) {
	client, err := Client(username, password)

	if err != nil {
		log.Fatal(err)
		ImportFailed = true
	}

	r := csv.NewReader(bytes.NewReader(listings))

	// Rows created from parsing the CSV file.
	var rows []row.Row

	// Have the headers been parsed?
	var parsedHeaders bool

	// Parse each CSV row. Skip the headers.
	for {
		record, err := r.Read()
		if err == io.EOF {
			msg := fmt.Sprintf("Parsed %d rows.", len(rows))
			log.Println(msg)
			ImportFailed = true
			break
		}
		if err != nil {
			log.Fatal(err)
			ImportFailed = true
		}

		// Check if this is the first row, which we skip.
		if !parsedHeaders {
			parsedHeaders = true
			continue
		}
		row, err := row.Parse(record)
		if err != nil {
			log.Println("Error parsing CSV row: " + err.Error())
			ImportFailed = true
			break
		}
		rows = append(rows, row)
	}

	// Create OB contracts.
	for _, row := range rows {
		c := Contract{
			Title:          row.Title,
			Image:          row.Image,
			Price:          row.Price,
			Description:    row.Description,
			Tags:           row.Tags,
			CurrencyCode:   row.CurrencyCode,
			ShippingOrigin: shipOrigin,
			ShipsTo:        []string{"all"},
		}
		// Get the image.
		resp, err := http.Get(row.Image)
		if err != nil {
			log.Println("Error downloading listing image: " + err.Error())
			ImportFailed = true
			break
		}
		buffer, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error downloading listing image: " + err.Error())
			ImportFailed = true
			break
		}
		err = resp.Body.Close()

		if err != nil {
			log.Println("Error downloading listing image: " + err.Error())
			ImportFailed = true
			break
		}

		// Encode image to base64.
		img := base64.StdEncoding.EncodeToString(buffer)
		// Upload image to OB and get image hash back.
		c.Image, err = UploadImage(client, img)
		if err != nil {
			log.Println("Error downloading listing image: " + err.Error())
			ImportFailed = true
			break
		}
		err = PostContract(client, c)
		if err != nil {
			log.Println("Error uploading contract: " + err.Error())
			ImportFailed = true
			break
		}
		ImportListingsParsed++
	}
	ImportFinished = true
}
