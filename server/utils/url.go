package utils

import (
	"regexp"
)

func GetTrucattedURL(url string) string {
	// domainPattern regex pattern for domain and subdomain in urls
	domainPattern := regexp.MustCompile(`^(?i)(https?:\/\/)?(?:[^@\n]+@)?(?:www\.)?[^:\/\n?]+`)
	return domainPattern.ReplaceAllLiteralString(url, "")
}

func GetCompleteURL(domain, url string) string {
	// protocolPattern regex pattern for http protocol in urls
	protocolPattern := regexp.MustCompile(`https?://`)
	// apiPathPattern regex pattern for path data in domain
	apiPathPattern := regexp.MustCompile(`/.*`)
	domain = protocolPattern.ReplaceAllLiteralString(domain, "")
	domain = apiPathPattern.ReplaceAllLiteralString(domain, "")
	url = GetTrucattedURL(url)

	return "https://" + domain + url
}

func GetBaseURL(domain string) string {
	return GetCompleteURL(domain, "/")
}
