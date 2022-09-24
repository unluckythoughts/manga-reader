package models

type Novel struct {
	ID       int            `json:"id" gorm:"column:id;primarykey"`
	SourceID int            `json:"-" gorm:"column:source_id"`
	Source   NovelSource    `json:"source" gorm:"foreignKey:SourceID"`
	URL      string         `json:"url" gorm:"column:url;unique"`
	Title    string         `json:"title" gorm:"column:title"`
	Slug     string         `json:"slug" gorm:"column:slug"`
	ImageURL string         `json:"imageUrl" gorm:"column:image_url"`
	Synopsis string         `json:"synopsis" gorm:"column:synopsis"`
	OtherID  string         `json:"otherId" gorm:"column:other_id"`
	Chapters []NovelChapter `json:"chapters" gorm:"foreignkey:NovelID"`
}

func (m Novel) TableName() string {
	return "novel"
}

type NovelChapter struct {
	ID         int    `json:"id" gorm:"column:id;primarykey"`
	URL        string `json:"url" gorm:"column:url;unique"`
	Title      string `json:"title" gorm:"column:title"`
	NovelID    int    `json:"-" gorm:"column:novel_id"`
	Number     string `json:"number" gorm:"column:number"`
	UploadDate string `json:"uploadDate" gorm:"column:upload_date"`
	OtherID    string `json:"otherId" gorm:"column:other_id"`
	Completed  bool   `json:"completed" gorm:"column:completed"`
	Downloaded bool   `json:"downloaded" gorm:"column:downloaded"`
}

func (c NovelChapter) TableName() string {
	return "novel_chapter"
}
