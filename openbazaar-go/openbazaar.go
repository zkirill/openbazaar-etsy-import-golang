package openbazaar

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

const (
	api = "http://localhost:18469/api/v1"
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
	for _, k := range contract.Tags {
		v.Add("keywords", k)
	}
	// Add shipping destinations.
	v.Add("ships_to", "all")
	resp, err := client.PostForm(api+"/contracts", v)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	dec := json.NewDecoder(resp.Body)
	var r response
	dec.Decode(&r)
	if !r.Success {
		return errors.New("Failed to post contract: " + contract.Title)
	}
	return nil
}
