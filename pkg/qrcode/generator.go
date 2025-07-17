package qrcode

import (
	goqrcode "github.com/skip2/go-qrcode"
)

// Generate creates a QR code image from data
func Generate(data string) ([]byte, error) {
	return goqrcode.Encode(data, goqrcode.Medium, 256)
}

// GenerateWithSize creates a QR code image with custom size
func GenerateWithSize(data string, size int) ([]byte, error) {
	return goqrcode.Encode(data, goqrcode.Medium, size)
}

// GenerateWithLevel creates a QR code image with custom error correction level
func GenerateWithLevel(data string, level goqrcode.RecoveryLevel, size int) ([]byte, error) {
	return goqrcode.Encode(data, level, size)
}
