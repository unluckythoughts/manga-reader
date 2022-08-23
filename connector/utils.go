package connector

import (
	"github.com/unluckythoughts/manga-reader/models"
)

type strFloat string

func (sf *strFloat) UnmarshalJSON(data []byte) error {
	*sf = strFloat(data)

	return nil
}

func uniqChapters(chapters []models.Chapter) []models.Chapter {
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
