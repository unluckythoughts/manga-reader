package models

import (
	"github.com/unluckythoughts/go-microservice/tools/web"
)

type INovelConnector interface {
	GetSource() Source
	GetNovelList(ctx web.Context) ([]Novel, error)
	GetNovelInfo(ctx web.Context, novelURL string) (Novel, error)
	GetNovelChapter(ctx web.Context, chapterURL string) ([]string, error)
}

type NovelSource struct {
	ID        int    `json:"id" gorm:"column:id;primarykey"`
	Name      string `json:"name" gorm:"column:name"`
	Domain    string `json:"domain" gorm:"column:domain"`
	IconURL   string `json:"iconUrl" gorm:"column:icon_url"`
	UpdatedAt string `json:"updatedAt" gorm:"column:updated_at"`
}

func (m NovelSource) TableName() string {
	return "novel_source"
}
