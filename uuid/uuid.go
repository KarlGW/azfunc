package uuid

import (
	"crypto/rand"
	"fmt"
)

// New creates a new UUID.
func New() (string, error) {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		return "", err
	}

	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}
