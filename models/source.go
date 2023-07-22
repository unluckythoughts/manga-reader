package models

type Source struct {
	ID        uint         `json:"id" gorm:"column:id;primarykey"`
	Name      string       `json:"name" gorm:"column:name"`
	Domain    string       `json:"domain" gorm:"column:domain"`
	IconURL   string       `json:"iconUrl" gorm:"column:icon_url"`
	UpdatedAt StrTimeStamp `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}

func (m Source) TableName() string {
	return "source"
}

type NovelSource struct {
	ID        uint         `json:"id" gorm:"column:id;primarykey"`
	Name      string       `json:"name" gorm:"column:name"`
	Domain    string       `json:"domain" gorm:"column:domain"`
	IconURL   string       `json:"iconUrl" gorm:"column:icon_url"`
	UpdatedAt StrTimeStamp `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}

func (m NovelSource) TableName() string {
	return "novel_source"
}
