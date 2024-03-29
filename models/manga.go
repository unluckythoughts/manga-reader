package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
)

type Manga struct {
	ID        uint           `json:"id" gorm:"column:id;primarykey"`
	SourceID  uint           `json:"-" gorm:"column:source_id"`
	Source    Source         `json:"source" gorm:"foreignKey:SourceID"`
	URL       string         `json:"url" gorm:"column:url;unique"`
	Title     string         `json:"title" gorm:"column:title"`
	Slug      string         `json:"slug" gorm:"column:slug"`
	ImageURL  string         `json:"imageUrl" gorm:"column:image_url"`
	Synopsis  string         `json:"synopsis" gorm:"column:synopsis"`
	OtherID   string         `json:"otherId" gorm:"column:other_id"`
	Chapters  []MangaChapter `json:"chapters" gorm:"foreignkey:MangaID"`
	UpdatedAt StrTimeStamp   `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}

func (m Manga) TableName() string {
	return "manga"
}

type MangaChapter struct {
	ID         uint         `json:"id" gorm:"column:id;primarykey"`
	URL        string       `json:"url" gorm:"column:url;unique"`
	Title      string       `json:"title" gorm:"column:title"`
	MangaID    uint         `json:"mangaId" gorm:"column:manga_id"`
	Number     string       `json:"number" gorm:"column:number"`
	UploadDate string       `json:"uploadDate" gorm:"column:upload_date"`
	ImageURLs  StrList      `json:"imageUrls" gorm:"column:image_urls"`
	OtherID    string       `json:"otherId" gorm:"column:other_id"`
	Completed  bool         `json:"completed" gorm:"column:completed"`
	Downloaded bool         `json:"downloaded" gorm:"column:downloaded"`
	UpdatedAt  StrTimeStamp `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}

func (c MangaChapter) TableName() string {
	return "chapter"
}

type StrList []string

func (l *StrList) Scan(value interface{}) error {
	data, ok := value.(string)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal stringList value:", value))
	}

	*l = strings.Split(data, ",")
	return nil
}

func (l StrList) Value() (driver.Value, error) {
	if len(l) == 0 {
		return "", nil
	}
	return strings.Join(l, ","), nil
}
