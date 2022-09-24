package models

type Source struct {
	ID        int    `json:"id" gorm:"column:id;primarykey"`
	Name      string `json:"name" gorm:"column:name"`
	Domain    string `json:"domain" gorm:"column:domain"`
	IconURL   string `json:"iconUrl" gorm:"column:icon_url"`
	UpdatedAt string `json:"updatedAt" gorm:"column:updated_at"`
}

func (m Source) TableName() string {
	return "source"
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
