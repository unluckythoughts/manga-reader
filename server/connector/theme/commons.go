package theme

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/unluckythoughts/book-reader/server/models"
	"github.com/unluckythoughts/book-reader/server/utils"
	"github.com/unluckythoughts/go-scraper"
)

func getBook(data string, conn models.Connector) (models.Book, error) {
	book := models.Book{}
	url, err := scraper.GetTextSingle(data, conn.Selectors.Book.URL)
	if err != nil {
		return book, err
	}

	title, err := scraper.GetTextSingle(data, conn.Selectors.Book.Title)
	if err != nil {
		return book, err
	}

	imageURL, err := scraper.GetTextSingle(data, conn.Selectors.Book.CoverImage)
	if err != nil {
		return book, err
	}

	synopsis, err := scraper.GetTextSingle(data, conn.Selectors.Book.Synopsis)
	if err != nil {
		return book, err
	}

	book = models.Book{
		URL:      utils.GetRelativeURL(url, conn.Domain),
		Title:    title,
		ImageURL: imageURL,
		Synopsis: synopsis,
		Type:     conn.Type,
	}
	return book, nil
}

func getChapter(data string, conn models.Connector) (models.Chapter, error) {
	chapter := models.Chapter{}
	sels := conn.Selectors.Book.Chapter

	url, err := scraper.GetTextSingle(data, sels.URL)
	if err != nil {
		return chapter, err
	}

	title := ""
	if sels.Title != "" {
		title, err = scraper.GetTextSingle(data, sels.Title)
		if err != nil {
			return chapter, err
		}

		// Clean the text - remove chapter number from title if present
		cleanPattern := regexp.MustCompile(`[Cc]hapter[ -]?[0-9.]+[- :]*`)
		title = cleanPattern.ReplaceAllString(title, "")
	}

	number := ""
	if sels.Number != "" {
		floatNum, err := scraper.GetFloat(data, sels.Number)
		if err != nil {
			return chapter, err
		}
		number = strconv.FormatFloat(floatNum, 'f', -1, 64)
	}

	uploadDateText, err := scraper.GetTextSingle(data, sels.UploadDate)
	if err != nil {
		return chapter, err
	}

	chapter = models.Chapter{
		URL:    utils.GetRelativeURL(url, conn.Domain),
		Title:  title,
		Number: number,
	}

	if uploadDateText != "" && sels.DateFormat != "" {
		uploadDate, err := scraper.GetTime(data, sels.UploadDate, sels.DateFormat)
		if err != nil {
			return chapter, err
		}

		chapter.UploadDate = *uploadDate
	}

	return chapter, nil
}

func cleanData(data string, patterns []models.Pattern) string {
	for _, pattern := range patterns {
		if pattern.IsRegex() {
			if !pattern.IsCaseSensitive() {
				data = strings.ToLower(data)
				pattern.Match = strings.ToLower(pattern.Match)
			}
			re, err := regexp.Compile(pattern.Match)
			if err != nil {
				continue
			}
			data = re.ReplaceAllString(data, pattern.ReplaceWith)
		}

		data = strings.ReplaceAll(data, pattern.Match, pattern.ReplaceWith)
	}

	return data
}
