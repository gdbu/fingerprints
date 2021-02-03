package fingerprints

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
)

// NewHash will generate a new hashed value from an inbound string value and return a pointer to the resulting HAsh
func NewHash(strs ...string) *Hash {
	h := MakeHash(strs...)
	return &h
}

// MakeHash will generate a new hashed value from an inbound string value
func MakeHash(strs ...string) (h Hash) {
	var err error
	enc := sha256.New()
	for _, str := range strs {
		if _, err = enc.Write([]byte(str)); err != nil {
			// Note: This is technically not possible, as the Write func for this type (sha256.digest)
			// does not have any error paths. That being said, it's always best practice to at least
			// check for an error (just in case something changes in the future). Since calling a panic
			// can cause serious problems within running applications, we are going to settle for
			// stdout logging
			log.Printf("error writing to sha256 writer, if you see this please report this on github.\nValue: <%s>\nError: %v\n", str, err)
			return
		}
	}

	enc.Sum(h[:0])
	return
}

// Hash represents a hashed value
type Hash [32]byte

// String returns a hex encoded string representation of the hashed value
func (h *Hash) String() string {
	return hex.EncodeToString((*h)[:])
}
