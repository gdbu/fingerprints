package fingerprints

import (
	"github.com/hatchify/errors"
	"github.com/mojura/mojura"
)

const (
	// ErrEmptyUserID is returned when the User ID for an Entry is empty
	ErrEmptyUserID = errors.Error("invalid user ID, cannot be empty")
	// ErrEmptyIPAddress is returned when the IP address for an Entry is empty
	ErrEmptyIPAddress = errors.Error("invalid IP address, cannot be empty")
	// ErrEmptyUserAgent is returned when the User Agent for an Entry is empty
	ErrEmptyUserAgent = errors.Error("invalid user agent, cannot be empty")
)

func makeEntry(userID string) (e Entry) {
	e.UserID = userID
	return
}

// Entry represents a stored entry within the Controller
type Entry struct {
	// Include mojura.Entry to auto-populate fields/methods needed to match the
	mojura.Entry

	// UserID which Entry is related to
	UserID string `json:"userID"`
	// Include identifiers
	Identifiers
}

// Signature will return the signature of the entry
func (e *Entry) Signature() string {
	h := NewHash(e.IPAddress, e.UserAgent, e.AcceptLanguage)
	return h.String()
}

// GetRelationships will return the relationship IDs associated with the Entry
func (e *Entry) GetRelationships() (r mojura.Relationships) {
	r.Append(e.UserID)
	r.Append(e.IPAddress)
	r.Append(NewHash(e.UserAgent).String())
	r.Append(e.AcceptLanguage)
	r.Append(e.Signature())
	return
}

// Validate will ensure an Entry is valid
func (e *Entry) Validate() (err error) {
	// An error list allows us to collect all the errors and return them as a group
	var errs errors.ErrorList
	// Check to see if User ID is set
	if len(e.UserID) == 0 {
		errs.Push(ErrEmptyUserID)
	}

	if len(e.IPAddress) == 0 {
		errs.Push(ErrEmptyIPAddress)
	}

	if len(e.UserAgent) == 0 {
		errs.Push(ErrEmptyUserAgent)
	}

	// Note: AcceptLanguage is not required

	// Note: If error list is empty, a nil value is returned
	return errs.Err()
}
