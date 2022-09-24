package scrapper

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/unluckythoughts/manga-reader/models"
	"github.com/unluckythoughts/manga-reader/utils"
)

func getNovelFromSelectors(
	s *goquery.Selection,
	titleSelector,
	uRLSelector,
	imageURLSelector,
	slugSelector,
	synopsisSelector,
	otherIDSelector string,
) (models.Novel, error) {
	novel := models.Novel{}

	title, err := GetTextForSelector(s, titleSelector)
	if err != nil {
		return novel, errors.Wrapf(err, "could not get manga title with selector %s", titleSelector)
	}

	novelURL, err := GetTextForSelector(s, uRLSelector)
	if err != nil {
		return novel, errors.Wrapf(err, "could not get manga url with selector %s", uRLSelector)
	}

	imageURL, err := GetTextForSelector(s, imageURLSelector)
	if err != nil {
		return novel, errors.Wrapf(err, "could not get manga image url with selector %s", imageURLSelector)
	}

	slug, err := GetTextForSelector(s, slugSelector)
	if err != nil {
		return novel, errors.Wrapf(err, "could not get manga slug with selector %s", slugSelector)
	}

	synopsis, err := GetAllTextForSelector(s, synopsisSelector)
	if err != nil {
		return novel, errors.Wrapf(err, "could not get manga synopsis with selector %s", synopsisSelector)
	}

	otherID, err := GetTextForSelector(s, otherIDSelector)
	if err != nil {
		return novel, errors.Wrapf(err, "could not get manga other id with selector %s", otherIDSelector)
	}

	novel.Title = title
	novel.URL = novelURL
	novel.ImageURL = imageURL
	novel.Slug = slug
	novel.Synopsis = strings.Join(synopsis, "\n")
	novel.OtherID = otherID

	return novel, nil
}

func GetNovelFromListSelectors(s *goquery.Selection, list models.NovelList) (models.Novel, error) {
	return getNovelFromSelectors(s,
		list.NovelTitle,
		list.NovelURL,
		list.NovelImageURL,
		list.NovelSlug,
		"",
		list.NovelOtherID,
	)
}

func GetNovelFromInfoSelectors(s *goquery.Selection, info models.NovelInfo) (models.Novel, error) {
	return getNovelFromSelectors(s,
		info.Title,
		"",
		info.ImageURL,
		info.Slug,
		info.Synopsis,
		info.OtherID,
	)
}

func GetNovelChapterFromInfoSelectors(s *goquery.Selection, info models.NovelInfo) (models.NovelChapter, error) {
	chapter := models.NovelChapter{}
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
