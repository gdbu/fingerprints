package fingerprints

import (
	"fmt"

	"github.com/mojura/mojura"
	"github.com/mojura/mojura/filters"
)

func asEntry(val mojura.Value) (e *Entry, err error) {
	var ok bool
	// Attempt to assert the value as an *Entry
	if e, ok = val.(*Entry); !ok {
		// Invalid type provided, return error
		err = fmt.Errorf("invalid entry type, expected %T and received %T", e, val)
		return
	}

	return
}

func appendMatchFilters(in []mojura.Filter, i Identifiers) (out []mojura.Filter) {
	out = in
	if len(i.IPAddress) != 0 {
		ipFilter := filters.Match(RelationshipIPAddresses, i.IPAddress)
		out = append(out, ipFilter)
	}

	if len(i.UserAgent) != 0 {
		hashed := NewHash(i.UserAgent).String()
		uaFilter := filters.Match(RelationshipUserAgents, hashed)
		out = append(out, uaFilter)
	}

	if len(i.AcceptLanguage) != 0 {
		lngFilter := filters.Match(RelationshipAcceptLanguages, i.AcceptLanguage)
		out = append(out, lngFilter)
	}

	return
}
