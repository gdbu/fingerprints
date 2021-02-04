# Fingerprints [![GoDoc](https://godoc.org/github.com/gdbu/fingerprints?status.svg)](https://pkg.go.dev/github.com/gdbu/fingerprints) ![Status](https://img.shields.io/badge/status-beta-yellow.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/gdbu/fingerprints)](https://goreportcard.com/report/github.com/gdbu/fingerprints)
Fingerprints is a user detection library intended to aid in finding spam accounts

## Usage
### New
```go
func ExampleNew() {
	var err error
	if testController, err = New("./data"); err != nil {
		log.Fatal(err)
	}
}
```

### Controller.New
```go
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
```

### Controller.GetByIP
```go
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
```

### Controller.GetByUserAgent
```go
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
```

### Controller.GetMatches
```go
func ExampleController_GetMatches() {
	var (
		i   Identifiers
		es  []*Entry
		err error
	)

	i.IPAddress = "64.233.191.255"
	i.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36"
	i.AcceptLanguage = "en-US,en;q=0.9,de-DE;q=0.8,de;q=0.7"

	if es, err = testController.GetMatches(i); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Matching entries", es)
}
```
### Controller.GetMatches (ip only)
```go
func ExampleController_GetMatches_ip_only() {
	var (
		i   Identifiers
		es  []*Entry
		err error
	)

	i.IPAddress = "64.233.191.255"

	if es, err = testController.GetMatches(i); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Matching entries", es)
}
```

### Controller.GetMatches (user agent only)
```go
func ExampleController_GetMatches_user_agent_only() {
	var (
		i   Identifiers
		es  []*Entry
		err error
	)

	i.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36"

	if es, err = testController.GetMatches(i); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Matching entries", es)
}
```

### Controller.GetDuplicates
```go
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
```