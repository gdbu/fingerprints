package fingerprints

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/gdbu/stringset"
	"github.com/hatchify/errors"
)

var (
	testCtx        = context.Background()
	testController *Controller
)

func TestController_GetDuplicates(t *testing.T) {
	var (
		c   *Controller
		err error
	)

	if c, err = testInit(); err != nil {
		t.Fatal(err)
	}
	defer testTeardown(c)

	var i Identifiers
	i.IPAddress = "64.233.191.255"
	i.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36"
	i.AcceptLanguage = "en-US,en;q=0.9,de-DE;q=0.8,de;q=0.7"

	if err = c.New(context.Background(), "user_0", i); err != nil {
		t.Fatal(err)
	}
	if err = c.New(context.Background(), "user_1", i); err != nil {
		t.Fatal(err)
	}
	if err = c.New(context.Background(), "user_1", i); err != nil {
		t.Fatal(err)
	}
	if err = c.New(context.Background(), "user_0", i); err != nil {
		t.Fatal(err)
	}

	if err = c.New(context.Background(), "user_3", i); err != nil {
		t.Fatal(err)
	}

	i.AcceptLanguage = "en-US,en;q=0.5"

	if err = c.New(context.Background(), "user_2", i); err != nil {
		t.Fatal(err)
	}

	var dups map[string]stringset.Map
	if dups, err = c.GetDuplicates(); err != nil {
		t.Fatal(err)
	}

	if len(dups) != 1 {
		t.Fatalf("invalid number of duplicate maps, expected %d and received %d", 1, len(dups))
	}

	var count int
	for _, dup := range dups {
		switch count {
		case 1:
			if len(dup) != 3 {
				t.Fatalf("invalid number of duplicate entries, expected %d and received %d", 3, len(dup))
			}

			if !dup.Has("user_0") {
				t.Fatalf("expected duplicate to contain <%s> and does not", "user_0")
			}

			if !dup.Has("user_1") {
				t.Fatalf("expected duplicate to contain <%s> and does not", "user_1")
			}

			if dup.Has("user_2") {
				t.Fatalf("expected duplicate to not contain <%s> and does", "user_2")
			}

			if !dup.Has("user_3") {
				t.Fatalf("expected duplicate to contain <%s> and does not", "user_3")
			}

		default:
		}
	}
}

func ExampleNew() {
	var err error
	if testController, err = New("./data"); err != nil {
		log.Fatal(err)
	}
}

func ExampleController_New() {
	var (
		i   Identifiers
		err error
	)

	i.IPAddress = "64.233.191.255"
	i.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36"
	i.AcceptLanguage = "en-US,en;q=0.9,de-DE;q=0.8,de;q=0.7"

	if err = testController.New(context.Background(), "user_0", i); err != nil {
		log.Fatal(err)
	}
}

func ExampleController_GetByIP() {
	var (
		es  []*Entry
		err error
	)
	if es, err = testController.GetByIP("[IP Address]"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Matching entries", es)
}

func ExampleController_GetByUserAgent() {
	var (
		es  []*Entry
		err error
	)
	if es, err = testController.GetByUserAgent("[User Agent]"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Matching entries", es)
}

func ExampleController_GetDuplicates() {
	var (
		dups map[string]stringset.Map
		err  error
	)
	if dups, err = testController.GetDuplicates(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Duplicates", dups)
}

func testInit() (c *Controller, err error) {
	if err = os.Mkdir("./test_dir", 0744); err != nil {
		return
	}

	return New("./test_dir")
}

func testTeardown(c *Controller) (err error) {
	var errs errors.ErrorList
	errs.Push(c.Close())
	errs.Push(os.RemoveAll("./test_dir"))
	return errs.Err()
}
