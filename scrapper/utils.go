package scrapper

import (
	"fmt"
	"html"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/pkg/errors"
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

	text = html.UnescapeString(text)
	text = strings.ReplaceAll(strings.TrimSpace(text), "</br>", "\n")

	return text, nil
}

func parseHumanReadableFormat(date string) (string, error) {
	t := strings.Split(strings.TrimSpace(date), " ")

	if len(t) < 3 {
		t = append([]string{"1"}, t...)
	}

	if t[2] != "ago" {
		return date, errors.New("human readable date does not end with 'ago'")
	}

	var err error
	var num int64
	var timeSpan int64

	if t[0] == "a" || t[0] == "an" || t[0] == "1" || t[0] == "few" {
		num = -1
	} else {
		num, err = strconv.ParseInt(t[0], 10, 64)
		if err != nil {
			return date, errors.Wrapf(err, "human readable date could not parse num of quantity %s", t[0])
		}

		num = -1 * num
	}

	switch {
	case strings.HasPrefix(t[1], "second"):
		timeSpan = int64(time.Second)
	case strings.HasPrefix(t[1], "minute"):
		timeSpan = int64(time.Minute)
	case strings.HasPrefix(t[1], "hour"):
		timeSpan = int64(time.Hour)
	case strings.HasPrefix(t[1], "day"):
		timeSpan = int64(24 * time.Hour)
	case strings.HasPrefix(t[1], "week"):
		timeSpan = int64(7 * 24 * time.Hour)
	case strings.HasPrefix(t[1], "month"):
		timeSpan = int64(30 * 24 * time.Hour)
	case strings.HasPrefix(t[1], "year"):
		timeSpan = int64(365 * 24 * time.Hour)
	default:
		return date, errors.Errorf("human readable date could not parse time span %s", t[1])
	}

	newDate := time.Now().Add(time.Duration(num * timeSpan)).Format("2006-01-02")

	return newDate, nil
}

func ParseDate(date string, format string) (string, error) {
	if strings.HasSuffix(date, "ago") || format == HUMAN_READABLE_DATE_FORMAT {
		return parseHumanReadableFormat(date)
	}

	if format == "" {
		return date, nil
	}

	t, err := time.Parse(format, date)
	if err != nil {
		return date, err
	}

	return t.Format("2006-01-02"), nil
}

func getTextForSelector(h *colly.HTMLElement, sel string) (string, error) {
	sels := getSelectors(sel)
	var selectorErr error
	for _, s := range sels {
		// fmt.Println(h.DOM.Html())
		val, err := getText(h.DOM.Find(getSelector(s)), s)
		if val != "" {
			return val, err
		}

		if err != nil {
			selectorErr = err
		}
	}

	return "", selectorErr
}

func getTextListForSelector(h *colly.HTMLElement, selector string, includeNoScript bool) (texts []string, err error) {
	// fix for <noscript> tags
	if includeNoScript {
		h.DOM.Find("noscript").Parent().SetHtml(h.DOM.Find("noscript").Text())
	}

	selectors := getSelectors(selector)
	var selectorErr error
	for _, s := range selectors {
		h.DOM.Find(getSelector(s)).Each(func(i int, sel *goquery.Selection) {
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

func getSelectors(selector string) []string {
	pattern, err := regexp.Compile("[^[,]+(\\[[^]]+\\])?")
	if err != nil {
		return []string{selector}
	}

	selector = strings.TrimSpace(selector)
	selectors := []string{}
	for _, m := range pattern.FindAllStringSubmatch(selector, -1) {
		selectors = append(selectors, m[0])
	}

	return selectors
}

func getSelector(selector string) string {
	pattern, err := regexp.Compile("\\[[^]]+\\]")
	if err != nil {
		return selector
	}

	return pattern.ReplaceAllString(selector, "")
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
