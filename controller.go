package fingerprints

import (
	"context"
	"net/http"

	"github.com/gdbu/stringset"
	"github.com/mojura/mojura"
	"github.com/mojura/mojura/filters"
)

// Relationship key const block
const (
	RelationshipUsers           = "users"
	RelationshipIPAddresses     = "ipAddresses"
	RelationshipUserAgents      = "userAgents"
	RelationshipAcceptLanguages = "acceptLanguages"
	RelationshipSignatures      = "signatures"
)

// relationships is a collection of all the supported relationship keys
var relationships = []string{
	RelationshipUsers,
	RelationshipIPAddresses,
	RelationshipUserAgents,
	RelationshipAcceptLanguages,
	RelationshipSignatures,
}

// New will return a new instance of the Controller
func New(dir string) (cc *Controller, err error) {
	var c Controller
	if c.m, err = mojura.New("fingerprints", dir, &Entry{}, relationships...); err != nil {
		return
	}

	// Assign pointer reference to our controller
	cc = &c
	return
}

// Controller represents a management layer to facilitate the retrieval and modification of Entries
type Controller struct {
	// Core will manage the data layer and will utilize the underlying back-end
	m *mojura.Mojura
}

// New will insert a new Entry to the back-end
func (c *Controller) New(ctx context.Context, userID string, i Identifiers) (err error) {
	// Make new entry from provided userID and identifiers
	e := makeEntry(userID, i)

	// Validate entry
	if err = e.Validate(); err != nil {
		return
	}

	// Set users match as the primary filter
	fs := []mojura.Filter{filters.Match(RelationshipUsers, userID)}
	// Append remaining match filters as secondary filters
	fs = appendMatchFilters(fs, i)

	var exists bool
	// Open a new read-only transaction
	if err = c.m.ReadTransaction(ctx, func(txn *mojura.Transaction) (err error) {
		// Attempt to see if entry exists
		exists, err = c.entryExists(txn, fs)
		return
	}); err != nil {
		return
	}

	// Check to see if matching entry for this user already exist
	if exists {
		// Exact match entry already exists, bail out
		return
	}

	// Open a batched read/write transaction
	err = c.m.Batch(ctx, func(txn *mojura.Transaction) (err error) {
		// Now that we've opened a write transaction, ensure the entry still does not exists
		if exists, err = c.entryExists(txn, fs); exists || err != nil {
			return
		}

		// Insert new entry into DB
		return c.new(txn, &e)
	})

	return
}

// NewFromHTTPRequst will create and insert a new entry from a given http request
func (c *Controller) NewFromHTTPRequst(ctx context.Context, userID string, req *http.Request) (err error) {
	return c.New(ctx, userID, makeIdentifiers(req))
}

// Get will retrieve an Entry which has the same ID as the provided entryID
func (c *Controller) Get(entryID string) (entry *Entry, err error) {
	var e Entry
	// Attempt to get Entry with the provided ID, pass reference to entry for which values to be applied
	if err = c.m.Get(entryID, &e); err != nil {
		return
	}

	// Assign reference to retrieved Entry
	entry = &e
	return
}

// GetByUser will retrieve all Entries associated with given user
func (c *Controller) GetByUser(userID string) (es []*Entry, err error) {
	filter := filters.Match(RelationshipUsers, userID)
	opts := mojura.NewFilteringOpts(filter)
	_, err = c.m.GetFiltered(&es, opts)
	return
}

// GetByIP will retrieve all Entries associated with given IP address
func (c *Controller) GetByIP(ipAddress string) (es []*Entry, err error) {
	filter := filters.Match(RelationshipIPAddresses, ipAddress)
	opts := mojura.NewFilteringOpts(filter)
	_, err = c.m.GetFiltered(&es, opts)
	return
}

// GetByUserAgent will retrieve all Entries associated with given user agent
func (c *Controller) GetByUserAgent(userAgent string) (es []*Entry, err error) {
	hashed := NewHash(userAgent).String()
	filter := filters.Match(RelationshipUserAgents, hashed)
	opts := mojura.NewFilteringOpts(filter)
	_, err = c.m.GetFiltered(&es, opts)
	return
}

// GetMatches will retrieve all Entries associated with given user agent and ip pair
func (c *Controller) GetMatches(i Identifiers) (es []*Entry, err error) {
	err = c.m.ReadTransaction(context.Background(), func(txn *mojura.Transaction) (err error) {
		es, err = c.getMatches(txn, "", i)
		return
	})

	return
}

// GetDuplicates will get all exact signature duplicates
func (c *Controller) GetDuplicates() (dups map[string]stringset.Map, err error) {
	fn := func(txn *mojura.Transaction) (err error) {
		dups, err = c.getDuplicates(txn)
		return
	}

	err = c.m.ReadTransaction(context.Background(), fn)
	return
}

// ForEach will iterate through all Entries
// Note: The error constant mojura.Break can returned by the iterating func to end the iteration early
func (c *Controller) ForEach(fn func(*Entry) error, opts *mojura.IteratingOpts) (err error) {
	onIterate := func(key string, val mojura.Value) (err error) {
		var e *Entry
		if e, err = asEntry(val); err != nil {
			return
		}

		// Pass iterating Entry to iterating function
		return fn(e)
	}

	// Iterate through all entries
	err = c.m.ForEach(onIterate, opts)
	return
}

// Delete will remove an Entry for a given user
func (c *Controller) Delete(ctx context.Context, entryID string) (removed *Entry, err error) {
	err = c.m.Transaction(ctx, func(txn *mojura.Transaction) (err error) {
		removed, err = c.delete(txn, entryID)
		return
	})

	return
}

// Close will close the controller and it's underlying dependencies
func (c *Controller) Close() (err error) {
	// Since we only have one dependency, we can just call this func directly
	return c.m.Close()
}

func (c *Controller) new(txn *mojura.Transaction, e *Entry) (err error) {
	// Attempt to validate Entry
	if err = e.Validate(); err != nil {
		// Entry is not valid, return validation error
		return
	}

	// Entry is valid!

	// Insert Entry into mojura.Core and return the results
	_, err = txn.New(e)
	return
}

// get will retrieve an entry by User ID
func (c *Controller) get(txn *mojura.Transaction, entryID string) (e *Entry, err error) {
	var entry Entry
	if err = txn.Get(entryID, &entry); err != nil {
		return
	}

	e = &entry
	return
}

// Delete will remove an Entry for a given userID
func (c *Controller) delete(txn *mojura.Transaction, userID string) (removed *Entry, err error) {
	var e *Entry
	if e, err = c.get(txn, userID); err != nil {
		return
	}

	// Remove Entry from mojura.Core
	if err = txn.Remove(e.ID); err != nil {
		return
	}

	removed = e
	return
}

func (c *Controller) getMatches(txn *mojura.Transaction, lastID string, i Identifiers) (es []*Entry, err error) {
	fs := appendMatchFilters(nil, i)
	opts := mojura.NewFilteringOpts(fs...)
	opts.LastID = lastID
	_, err = txn.GetFiltered(&es, opts)
	return
}

// New will insert a new Entry to the back-end
func (c *Controller) entryExists(txn *mojura.Transaction, fs []mojura.Filter) (ok bool, err error) {
	var cur mojura.Cursor
	// Initialize a new cursor with the provided filters
	if cur, err = txn.Cursor(fs...); err != nil {
		return
	}

	// Attempt to get the very first matching entry
	_, err = cur.First()
	switch err {
	case nil:
		// No error encountered, entry exists
		ok = true
	case mojura.Break:
		// Break error encountered, no entry exists
		err = nil
	default:
		// Unexpected error encountered, oof
	}

	return
}

func (c *Controller) getDuplicates(txn *mojura.Transaction) (dups map[string]stringset.Map, err error) {
	// Initialize new duplicates seeker
	d := newDuplicates()
	// Populate duplicates map
	if err = d.populate(txn); err != nil {
		return
	}

	// Reference duplicates map
	dups = d.m
	return
}
