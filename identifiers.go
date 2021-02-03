package fingerprints

import "net/http"

func makeIdentifiers(req *http.Request) (i Identifiers) {
	h := req.Header
	// Attempt to get IP address from X-Forwarded-For header value first
	if i.IPAddress = h.Get("X-Forwarded-For"); i.IPAddress == "" {
		// X-Forwarded-For does not exists, get IP from remote address
		i.IPAddress = req.RemoteAddr
	}

	// Get Accept-Language header value
	i.AcceptLanguage = h.Get("Accept-Language")
	// Get User-Agent header value
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
