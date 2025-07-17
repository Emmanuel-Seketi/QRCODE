package shorturl

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// Generate creates a short URL slug
func Generate() (string, error) {
	bytes := make([]byte, 4) // 4 bytes will produce an 8-character hex string
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	slug := hex.EncodeToString(bytes)
	fmt.Printf("Generated short URL slug: %s\n", slug)
	return slug, nil
}
