package util

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

// GenerateULID generates an ULID to use as primary key
func GenerateULID() string {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	ms := ulid.Timestamp(time.Now())
	generatedUlid, err := ulid.New(ms, entropy)
	if err != nil {
		panic(err)
	}

	return generatedUlid.String()
}
