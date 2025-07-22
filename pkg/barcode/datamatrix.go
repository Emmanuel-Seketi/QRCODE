package barcode

import (
	"bytes"
	"image/png"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/datamatrix"
)

// DataMatrixOptions holds configuration for Data Matrix generation
type DataMatrixOptions struct {
	Data string
	Size string // "auto", "10x10", "12x12", etc.
}

// GenerateDataMatrix creates a Data Matrix barcode image
func GenerateDataMatrix(options DataMatrixOptions) ([]byte, error) {
	// Create the barcode
	code, err := datamatrix.Encode(options.Data)
	if err != nil {
		return nil, err
	}

	// Scale the barcode (default to 200x200 pixels)
	scaledCode, err := barcode.Scale(code, 200, 200)
	if err != nil {
		return nil, err
	}

	// Convert to PNG bytes
	var buf bytes.Buffer
	err = png.Encode(&buf, scaledCode)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// GenerateDataMatrixWithSize creates a Data Matrix barcode with custom size
func GenerateDataMatrixWithSize(data string, width, height int) ([]byte, error) {
	// Create the barcode
	code, err := datamatrix.Encode(data)
	if err != nil {
		return nil, err
	}

	// Scale the barcode to specified size
	scaledCode, err := barcode.Scale(code, width, height)
	if err != nil {
		return nil, err
	}

	// Convert to PNG bytes
	var buf bytes.Buffer
	err = png.Encode(&buf, scaledCode)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}