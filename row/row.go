package row

import "strings"

// Field represents the CSV field.
type Field int

const (
	// FieldTitle is the listing title.
	FieldTitle = iota
	// FieldDescription is the listing description.
	FieldDescription
	// FieldPrice is the listing price.
	FieldPrice
	// FieldCurrencyCode is the listing currency code.
	FieldCurrencyCode
	// FieldQuantity is the listing quantity.
	FieldQuantity
	// FieldTags is the listing tags.
	FieldTags
	// FieldMaterials is the listing materials.
	FieldMaterials
	// FieldImage1 is the listing image.
	FieldImage1
	// FieldImage2 is the listing image.
	FieldImage2
	// FieldImage3 is the listing image.
	FieldImage3
	// FieldImage4 is the listing image.
	FieldImage4
	// FieldImage5 is the listing image.
	FieldImage5
	// FieldVariation1Type is the listing variation type.
	FieldVariation1Type
	// FieldVariation1Name is the listing variation name.
	FieldVariation1Name
	// FieldVariation1Values is the listing variation values.
	FieldVariation1Values
	// FieldVariation2Type is the listing variation type.
	FieldVariation2Type
	// FieldVariation2Name is the listing variation name.
	FieldVariation2Name
	// FieldVariation2Values is the listing variation values.
	FieldVariation2Values
)

// Row represents a CSV row containing a listing.
type Row struct {
	// Title is the title of the listing.
	Title string
	// Description is the description of the listing.
	Description string
	// Price is the price of the listing.
	Price string
	// Image is the image of the listing.
	Image string
	// Tags is the tags (keywords).
	Tags []string
	// CurrencyCode is the currency code.
	CurrencyCode string
}

// Parse parses the record and returns a row.
func Parse(record []string) (row Row, err error) {
	r := Row{
		Title:        record[FieldTitle],
		Description:  record[FieldDescription],
		Price:        record[FieldPrice],
		Image:        record[FieldImage1],
		Tags:         strings.Split(record[FieldTags], ","),
		CurrencyCode: record[FieldCurrencyCode],
	}
	return r, nil
}
