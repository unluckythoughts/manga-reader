package scrapper

import (
	"html"
	"regexp"
	"strconv"
	"strings"

	"github.com/unluckythoughts/manga-reader/models"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/pkg/errors"
	"github.com/unluckythoughts/manga-reader/utils"
)

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

func getTextForSelector(h *colly.HTMLElement, sel string) (string, error) {
	sels := getSelectors(sel)
	var selectorErr error
	for _, s := range sels {
		// fmt.Println(h.DOM.Html())
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

func GetTextForSelector(selection *goquery.Selection, sel string) (string, error) {
	if sel == "" {
		return "", nil
	}
	// x, e := selection.Html()
	// fmt.Println(x, e)

	sels := getSelectors(sel)
	var selectorErr error
	for _, s := range sels {
		selForText := selection
		if s != "" {
			selForText = selection.Find(s)
		}

		val, err := getText(selForText, s)
		if val != "" {
			return val, err
		}

		if err != nil {
			selectorErr = err
		}
	}

	return "", selectorErr
}

func GetElementForSelector(selection *goquery.Selection, sel string) (*goquery.Selection, bool) {
	sels := getSelectors(sel)
	for _, s := range sels {
		if selection.Find(s).Length() > 0 {
			return selection.Find(s).First(), true
		}
	}

	return nil, false
}

func getMangaFromSelectors(
	s *goquery.Selection,
	titleSelector,
	uRLSelector,
	imageURLSelector,
	slugSelector,
	synopsisSelector,
	otherIDSelector string,
) (models.Manga, error) {
	manga := models.Manga{}

	title, err := GetTextForSelector(s, titleSelector)
	if err != nil {
		return manga, errors.Wrapf(err, "could not get manga title with selector %s", titleSelector)
	}

	mangaURL, err := GetTextForSelector(s, uRLSelector)
	if err != nil {
		return manga, errors.Wrapf(err, "could not get manga url with selector %s", uRLSelector)
	}

	imageURL, err := GetTextForSelector(s, imageURLSelector)
	if err != nil {
		return manga, errors.Wrapf(err, "could not get manga image url with selector %s", imageURLSelector)
	}

	slug, err := GetTextForSelector(s, slugSelector)
	if err != nil {
		return manga, errors.Wrapf(err, "could not get manga slug with selector %s", slugSelector)
	}

	synopsis, err := GetTextForSelector(s, synopsisSelector)
	if err != nil {
		return manga, errors.Wrapf(err, "could not get manga synopsis with selector %s", synopsisSelector)
	}

	otherID, err := GetTextForSelector(s, otherIDSelector)
	if err != nil {
		return manga, errors.Wrapf(err, "could not get manga other id with selector %s", otherIDSelector)
	}

	manga.Title = title
	manga.URL = mangaURL
	manga.ImageURL = imageURL
	manga.Slug = slug
	manga.Synopsis = synopsis
	manga.OtherID = otherID

	return manga, nil
}

func GetMangaFromListSelectors(s *goquery.Selection, list models.MangaList) (models.Manga, error) {
	return getMangaFromSelectors(s,
		list.MangaTitle,
		list.MangaURL,
		list.MangaImageURL,
		list.MangaSlug,
		"",
		list.MangaOtherID,
	)
}

func GetMangaFromInfoSelectors(s *goquery.Selection, info models.MangaInfo) (models.Manga, error) {
	return getMangaFromSelectors(s,
		info.Title,
		"",
		info.ImageURL,
		info.Slug,
		info.Synopsis,
		info.OtherID,
	)
}

func UniqChapters(chapters []models.Chapter) []models.Chapter {
	chapterMap := map[string]models.Chapter{}
	for _, c := range chapters {
		chapterMap[c.Number] = c
	}

	chapters = []models.Chapter{}
	for _, c := range chapterMap {
		chapters = append(chapters, c)
	}

	return chapters
}

func GetChapterFromInfoSelectors(s *goquery.Selection, info models.MangaInfo) (models.Chapter, error) {
	chapter := models.Chapter{}
	title, err := GetTextForSelector(s, info.ChapterTitle)
	if err != nil {
		return chapter, errors.Wrapf(err, "could not get chapter title with selector %s", info.ChapterTitle)
	}

	number, err := GetTextForSelector(s, info.ChapterNumber)
	if err != nil {
		return chapter, errors.Wrapf(err, "could not get chapter number with selector %s", info.ChapterNumber)
	}
	number = GetChapterNumber(number)

	chapterURL, err := GetTextForSelector(s, info.ChapterURL)
	if err != nil {
		return chapter, errors.Wrapf(err, "could not get chapter url with selector %s", info.ChapterURL)
	}

	uploadDate, err := GetTextForSelector(s, info.ChapterUploadDate)
	if err != nil {
		return chapter, errors.Wrapf(err, "could not get chapter upload date with selector %s", info.ChapterUploadDate)
	}
	parsedDate, err := utils.ParseDate(uploadDate, info.ChapterUploadDateFormat)
	if err != nil {
		return chapter, errors.Wrapf(err, "could not parse chapter upload date %s with format %s", uploadDate, info.ChapterUploadDateFormat)
	}

	chapter.Title = title
	chapter.Number = number
	chapter.URL = chapterURL
	chapter.UploadDate = parsedDate

	return chapter, nil
}

func GetImagesListForSelector(selection *goquery.Selection, selector string, includeNoScript bool) (images []string, err error) {
	// fix for <noscript> tags
	if includeNoScript {
		selection.Find("noscript").Parent().SetHtml(selection.Find("noscript").Text())
	}

	selectors := getSelectors(selector)
	var selectorErr error
	for _, s := range selectors {
		selection.Find(s).Each(func(i int, sel *goquery.Selection) {
			var text string
			text, selectorErr = getText(sel, s)
			images = append(images, text)
		})

		if len(images) > 0 {
			return images, nil
		}

		if selectorErr != nil {
			return images, errors.Wrapf(selectorErr, "could not get chapter image url for selector %s", selector)
		}
	}

	return images, selectorErr
}

func getTextListForSelector(h *colly.HTMLElement, selector string, includeNoScript bool) (texts []string, err error) {
	// fix for <noscript> tags
	if includeNoScript {
		h.DOM.Find("noscript").Parent().SetHtml(h.DOM.Find("noscript").Text())
	}

	selectors := getSelectors(selector)
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

func getSelectors(selector string) []string {
	pattern, err := regexp.Compile("[^[,]*(\\[[^]]+\\])?")
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
		return utils.GetString(getValue(matches[0]))
	} else if len(matches) > 1 {
		var num float64
		for _, match := range matches {
			newNum := getValue(match)
			if newNum > num {
				num = newNum
			}
		}
		return utils.GetString(num)
	}

	return ""
}
