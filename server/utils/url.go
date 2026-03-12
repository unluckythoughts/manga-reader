package utils

import (
	"regexp"
	"strings"
)

func GetTrucattedURL(url string) string {
	// domainPattern regex pattern for domain and subdomain in urls
	domainPattern := regexp.MustCompile(`^(?i)(https?:\/\/)?(?:[^@\n]+@)?(?:www\.)?[^:\/\n?]+`)
	return domainPattern.ReplaceAllLiteralString(url, "")
}

func GetRelativeURL(url, domain string) string {
	if strings.Contains(url, domain) {
		if !strings.HasPrefix(domain, "http") {
			pattern := regexp.MustCompile(`^(?i)(https?:\/\/)?(?:[^@\n]+@)?(?:www\.)?`)
			url = pattern.ReplaceAllLiteralString(url, "")
		}

		return strings.TrimPrefix(url, domain)
	}

	return url
}

func GetCompleteURL(domain, url string) string {
	if strings.HasPrefix(url, "http") {
		return url
	}

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
