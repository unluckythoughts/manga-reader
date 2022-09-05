package scrapper

import (
	"html"
	"regexp"
	"strconv"
	"strings"

	"github.com/unluckythoughts/manga-reader/models"

	"github.com/PuerkitoBio/goquery"
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
		text = s.Text()
	}

	text = html.UnescapeString(text)
	text = strings.ReplaceAll(strings.TrimSpace(text), "</br>", "\n")

	return text, nil
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
		if !strings.HasPrefix(s, "[") {
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

func GetChapterTitle(title string) string {
	pattern := regexp.MustCompile(`(?i).*chapter`)
	if pattern.MatchString(title) {
		return pattern.ReplaceAllString(title, "Chapter")
	}

	return title
}

func GetChapterFromInfoSelectors(s *goquery.Selection, info models.MangaInfo) (models.Chapter, error) {
	chapter := models.Chapter{}
	title, err := GetTextForSelector(s, info.ChapterTitle)
	if err != nil {
		return chapter, errors.Wrapf(err, "could not get chapter title with selector %s", info.ChapterTitle)
	}
	title = GetChapterTitle(title)

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

func getSelectors(selector string) []string {
	pattern, err := regexp.Compile(`(?im)[^\n,(]*(?:\([^)]*\))?`)
	if err != nil {
		return []string{selector}
	}

	selector = strings.TrimSpace(selector)
	return pattern.FindAllString(selector, -1)
}

func hasDataInAttr(selector string) (string, []string, bool) {
	multiAttrPattern := regexp.MustCompile(`(?i):(?:is|matches)\([^)]*\)`)
	pattern := regexp.MustCompile(`(?i)\[[^]]+\]`)

	matches := pattern.FindAllString(selector, -1)
	if len(matches) < 1 {
		return selector, []string{}, false
	}

	attrs := []string{}
	for _, match := range matches {
		attrs = append(attrs, strings.Trim(match, "[]"))
	}
	selector = pattern.ReplaceAllString(selector, "")
	selector = multiAttrPattern.ReplaceAllString(selector, "")

	return selector, attrs, true
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
		re = regexp.MustCompile("[Cc]hapter[ :-]*([0-9.]+)")
		if !re.MatchString(text) {
			re = regexp.MustCompile("(?m)\\b([0-9.]+)\\b")
		}
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
