package fingerprints

import "net/http"

func makeIdentifiers(req *http.Request) (i Identifiers) {
	h := req.Header
	if i.IPAddress = h.Get("X-Forwarded-For"); i.IPAddress == "" {
		i.IPAddress = req.RemoteAddr
	}

	i.AcceptLanguage = h.Get("Accept-Language")
	i.UserAgent = h.Get("User-Agent")
	return
}

// Identifiers are the identifiers used for fingerprinting
type Identifiers struct {
	// IPAddress is the IP address of the user
	IPAddress string `json:"ipAddress"`
	// UserAgent is the user agent of the user
	UserAgent string `json:"userAgent"`
	// AcceptLanguage is the acceptLanguage value for the user
	AcceptLanguage string `json:"acceptLanguage"`
}
