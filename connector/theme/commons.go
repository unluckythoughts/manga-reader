package theme

import (
	"regexp"
)

var (
	// protocolPattern regex pattern for http protocol in urls
	protocolPattern = regexp.MustCompile(`https?://`)
	// apiPathPattern regex pattern for path data in domain
	apiPathPattern = regexp.MustCompile(`/.*`)
	// domainPattern regex pattern for domain and subdomain in urls
	domainPattern = regexp.MustCompile(`^(?i)(https?:\/\/)?(?:[^@\n]+@)?(?:www\.)?[^:\/\n?]+`)
)

func GetCompleteURL(url, domain string) string {
	domain = protocolPattern.ReplaceAllLiteralString(domain, "")
	domain = apiPathPattern.ReplaceAllLiteralString(domain, "")
	url = domainPattern.ReplaceAllLiteralString(url, "")

	return "https://" + domain + url
}

func GetTrucattedURL(url string) string {
	url = domainPattern.ReplaceAllLiteralString(url, "")

	return url
}
