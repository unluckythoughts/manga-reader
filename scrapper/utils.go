package scrapper

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

func strAdd(a string, b int) string {
	numA, err := strconv.Atoi(a)
	if err != nil {
		return a + strconv.Itoa(b)
	}
	return strconv.Itoa(numA + b)
}

func getText(s *goquery.Selection, selector string) (string, error) {
	var text string
	if _, attrs, inAttr := hasDataInAttr(selector); inAttr {
		var ok bool
		for _, attr := range attrs {
			text, ok = s.Attr(attr)
			if ok || text != "" {
				break
			}
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

func getTextForSelector(h *colly.HTMLElement, sel string) (string, error) {
	sels := strings.Split(sel, ",")
	var selectorErr error
	for _, s := range sels {
		val, err := getText(h.DOM.Find(s), s)
		if val != "" {
			return val, err
		}

		if err != nil {
			selectorErr = err
		}
	}

	return "", selectorErr
}

func getTextListForSelector(h *colly.HTMLElement, selector string) (texts []string, err error) {
	// fix for <noscript> tags
	h.DOM.Find("noscript").Parent().SetHtml(h.DOM.Find("noscript").Text())

	selectors := strings.Split(selector, ",")
	var selectorErr error
	for _, s := range selectors {
		h.DOM.Find(s).Each(func(i int, sel *goquery.Selection) {
			var text string
			text, err = getText(sel, s)
			texts = append(texts, text)
		})

		if len(texts) > 0 {
			return texts, nil
		}

		if err != nil {
			selectorErr = err
		}
	}

	return texts, selectorErr
}

func hasDataInAttr(selector string) (string, []string, bool) {
	pattern, err := regexp.Compile("\\[[^]]+\\]")
	if err != nil {
		return selector, []string{}, false
	}

	matches := pattern.FindAllString(selector, -1)
	if len(matches) < 1 {
		return selector, []string{}, false
	}

	attr := strings.Trim(matches[len(matches)-1], "[]")
	selector = pattern.ReplaceAllString(selector, "")

	return selector, strings.Split(attr, ","), true
}

func isInternalLink(s string) bool {
	pattern := regexp.MustCompile("^[?&]")

	return pattern.MatchString(s)
}

func getString(num float64) string {
	strFloat := fmt.Sprintf("%f", num)
	strFloat = strings.TrimRight(strFloat, "0")
	strFloat = strings.TrimRight(strFloat, ".")

	return strFloat
}

func GetChapterNumber(text string) string {
	getValue := func(subMatches []string) float64 {
		val := ""
		if len(subMatches) == 1 {
			val = subMatches[0]
		} else if len(subMatches) > 1 {
			val = subMatches[1]
		}

		num, _ := strconv.ParseFloat(val, 64)
		return num
	}

	re := regexp.MustCompile("^[0-9.]+")
	if !re.MatchString(text) {
		re = regexp.MustCompile("(?m)\\b([0-9.]+)\\b")
	}
	matches := re.FindAllStringSubmatch(text, -1)

	if len(matches) == 1 {
		return getString(getValue(matches[0]))
	} else if len(matches) > 1 {
		var num float64
		for _, match := range matches {
			newNum := getValue(match)
			if newNum > num {
				num = newNum
			}
		}
		return getString(num)
	}

	return ""
}
