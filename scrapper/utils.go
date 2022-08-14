package scrapper

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

func getText(s *goquery.Selection, selector string) (string, error) {
	var text string
	if attr, inAttr := hasDataInAttr(selector); inAttr {
		var ok bool
		text, ok = s.Attr(attr)
		if !ok {
			return "", nil
		}
	} else {
		var err error
		text, err = s.Html()
		if err != nil {
			return text, err
		}
	}

	text = strings.ReplaceAll(strings.TrimSpace(text), "</br>", "\n")

	return text, nil
}

func getTextListForSelector(h *colly.HTMLElement, s string) (texts []string, err error) {
	h.DOM.Find(s).Each(func(i int, sel *goquery.Selection) {
		var text string
		text, err = getText(sel, s)
		texts = append(texts, text)
	})

	return texts, err
}

func hasDataInAttr(selector string) (string, bool) {
	pattern, err := regexp.Compile("\\[[^]]+\\]")
	if err != nil {
		return "", false
	}

	matches := pattern.FindAllString(selector, -1)
	if len(matches) < 1 {
		return "", false
	}

	attr := strings.Trim(matches[len(matches)-1], "[]")

	return attr, true
}

func isInternalLink(s string) bool {
	pattern := regexp.MustCompile("^[/?&]")

	return pattern.MatchString(s)
}
