package models

type Novel struct {
	ID        uint           `json:"id" gorm:"column:id;primarykey"`
	SourceID  uint           `json:"-" gorm:"column:source_id"`
	Source    NovelSource    `json:"source" gorm:"foreignKey:SourceID"`
	URL       string         `json:"url" gorm:"column:url;unique"`
	Title     string         `json:"title" gorm:"column:title"`
	Slug      string         `json:"slug" gorm:"column:slug"`
	ImageURL  string         `json:"imageUrl" gorm:"column:image_url"`
	Synopsis  string         `json:"synopsis" gorm:"column:synopsis"`
	OtherID   string         `json:"otherId" gorm:"column:other_id"`
	Chapters  []NovelChapter `json:"chapters" gorm:"foreignkey:NovelID"`
	UpdatedAt StrTimeStamp   `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}

func (m Novel) TableName() string {
	return "novel"
}

type NovelChapter struct {
	ID         uint         `json:"id" gorm:"column:id;primarykey"`
	URL        string       `json:"url" gorm:"column:url;unique"`
	Title      string       `json:"title" gorm:"column:title"`
	NovelID    uint         `json:"-" gorm:"column:novel_id"`
	Number     string       `json:"number" gorm:"column:number"`
	UploadDate string       `json:"uploadDate" gorm:"column:upload_date"`
	OtherID    string       `json:"otherId" gorm:"column:other_id"`
	Completed  bool         `json:"completed" gorm:"column:completed"`
	Downloaded bool         `json:"downloaded" gorm:"column:downloaded"`
	UpdatedAt  StrTimeStamp `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}

func (c NovelChapter) TableName() string {
	return "novel_chapter"
}
