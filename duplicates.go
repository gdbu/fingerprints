package fingerprints

import (
	"github.com/gdbu/stringset"
	"github.com/mojura/mojura"
	"github.com/mojura/mojura/filters"
)

func newDuplicates() *duplicates {
	var d duplicates
	d.m = make(map[string]stringset.Map)
	return &d
}

type duplicates struct {
	lastSignature string

	m map[string]stringset.Map
}

func (d *duplicates) clean() {
	// We've rotated to a new signature, check to see if the last one has a
	// sufficient amount of entries to stay within the duplicates map
	if len(d.m[d.lastSignature]) < 2 {
		// Map has less than two entries, delete signature key
		delete(d.m, d.lastSignature)
	}
}

func (d *duplicates) compare(signature string) (ok bool, err error) {
	// Return true if signature is set
	ok = signature != ""

	// Check to see if signature has changed
	if signature == d.lastSignature {
		// Current signature matches the most recently set signature, return
		return
	}

	// Clean previous signature
	d.clean()

	// Update last signature to be current signature
	d.lastSignature = signature
	return
}

func (d *duplicates) iterate(entryID string, val mojura.Value) (err error) {
	e := val.(*Entry)
	// Attempt to get stringset for current signature
	ss, ok := d.m[d.lastSignature]
	if !ok {
		// No stringset exists for this signature, initialize a new one
		ss = stringset.MakeMap()
		// Associated created stringset with the signature
		d.m[d.lastSignature] = ss
	}

	// Set user ID to stringset
	ss.Set(e.UserID)
	return
}

func (d *duplicates) populate(txn *mojura.Transaction) (err error) {
	f := filters.Comparison(RelationshipSignatures, d.compare)
	// Initialize iterating options
	opts := mojura.NewIteratingOpts(f)
	// Iterate through entries
	err = txn.ForEach(d.iterate, opts)
	// Clean last signature
	d.clean()
	return
}
