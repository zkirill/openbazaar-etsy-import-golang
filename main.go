package main

import (
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/zkirill/openbazaar-etsy-import-golang/openbazaar-go"
	"github.com/zkirill/openbazaar-etsy-import-golang/row"
)

const (
	listingsFile = "listings.csv"
)

func main() {
	fmt.Println("Welcome to OpenBazaar import script!")

	var username, password, shipOrigin string

	// Get username.
	for len(username) == 0 {
		fmt.Println("Please enter your OpenBazaar server username: ")
		_, err := fmt.Scanln(&username)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Get password.
	for len(password) == 0 {
		fmt.Println("Thanks! Please enter your OpenBazaar server password: ")
		_, err := fmt.Scanln(&password)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Get shipping origin.
	for len(shipOrigin) == 0 {
		fmt.Println("Thanks! Please enter a shipping origin in the form of a country code, or enter ALL: ")
		_, err := fmt.Scanln(&shipOrigin)
		if err != nil {
			log.Fatal(err)
		}
	}

	client, err := openbazaar.Client(username, password)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Openning file " + listingsFile + "...")
	contents, err := ioutil.ReadFile(listingsFile)
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(bytes.NewReader(contents))

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
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// Check if this is the first row, which we skip.
		if !parsedHeaders {
			parsedHeaders = true
			continue
		}
		row, err := row.Parse(record)
		if err != nil {
			log.Println("Error parsing CSV row: " + err.Error())
			break
		}
		rows = append(rows, row)
	}

	// Create OB contracts.
	for _, row := range rows {
		c := openbazaar.Contract{
			Title:          row.Title,
			Image:          row.Image,
			Price:          row.Price,
			Description:    row.Description,
			Tags:           row.Tags,
			CurrencyCode:   row.CurrencyCode,
			ShippingOrigin: shipOrigin,
		}
		// Get the image.
		resp, err := http.Get(row.Image)
		if err != nil {
			log.Println("Error downloading listing image: " + err.Error())
			break
		}
		buffer, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error downloading listing image: " + err.Error())
			break
		}
		err = resp.Body.Close()

		if err != nil {
			log.Println("Error downloading listing image: " + err.Error())
			break
		}

		// Encode image to base64.
		img := base64.StdEncoding.EncodeToString(buffer)
		// Upload image to OB and get image hash back.
		c.Image, err = openbazaar.UploadImage(client, img)
		if err != nil {
			log.Println("Error downloading listing image: " + err.Error())
			break
		}
		err = openbazaar.PostContract(client, c)
		if err != nil {
			log.Println("Error uploading contract: " + err.Error())
			break
		}
	}
}
