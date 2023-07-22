package models

type MangaFavorite struct {
	ID         uint         `json:"id" gorm:"column:id;primarykey"`
	UserID     int          `json:"-" gorm:"column:user_id"`
	User       User         `json:"user" gorm:"foreignkey:UserID"`
	MangaID    uint         `json:"-" gorm:"column:manga_id"`
	Manga      Manga        `json:"manga" gorm:"foreignkey:MangaID"`
	Progress   StrFloatList `json:"progress" gorm:"column:progress"`
	Categories StrList      `json:"categories" gorm:"column:categories"`
	UpdatedAt  StrTimeStamp `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}

func (f MangaFavorite) TableName() string {
	return "favorite"
}

type NovelFavorite struct {
	ID         uint         `json:"id" gorm:"column:id;primarykey"`
	UserID     int          `json:"-" gorm:"column:user_id"`
	User       User         `json:"user" gorm:"foreignkey:UserID"`
	NovelID    uint         `json:"-" gorm:"column:novel_id"`
	Novel      Novel        `json:"novel" gorm:"foreignkey:NovelID"`
	Progress   StrFloatList `json:"progress" gorm:"column:progress"`
	Categories StrList      `json:"categories" gorm:"column:categories"`
	UpdatedAt  StrTimeStamp `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}

func (f NovelFavorite) TableName() string {
	return "novel_favorite"
}

type Category struct {
	ID        int          `json:"-" gorm:"column:id;primarykey"`
	Name      string       `json:"name" gorm:"column:name"`
	UpdatedAt StrTimeStamp `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}

func (c Category) TableName() string {
	return "category"
}

type User struct {
	ID        int          `json:"-" gorm:"column:id;primarykey"`
	Name      string       `json:"name" gorm:"column:name"`
	UpdatedAt StrTimeStamp `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}

func (u User) TableName() string {
	return "user"
}
