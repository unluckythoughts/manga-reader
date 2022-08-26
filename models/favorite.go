package models

type Favorite struct {
	ID         int          `json:"id" gorm:"column:id;primarykey"`
	UserID     int          `json:"-" gorm:"column:user_id"`
	User       User         `json:"user" gorm:"foreignkey:UserID"`
	MangaID    int          `json:"-" gorm:"column:manga_id"`
	Manga      Manga        `json:"manga" gorm:"foreignkey:MangaID"`
	Progress   StrFloatList `json:"progress" gorm:"column:progress"`
	Categories StrList      `json:"categories" gorm:"column:categories"`
}

func (f Favorite) TableName() string {
	return "favorite"
}

type Category struct {
	ID   int    `json:"-" gorm:"column:id;primarykey"`
	Name string `json:"name" gorm:"column:name"`
}

func (c Category) TableName() string {
	return "category"
}

type User struct {
	ID   int    `json:"-" gorm:"column:id;primarykey"`
	Name string `json:"name" gorm:"column:name"`
}

func (u User) TableName() string {
	return "user"
}
