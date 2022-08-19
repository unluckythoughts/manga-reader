package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type StrIntList [2]int

func (sil *StrIntList) Scan(value interface{}) error {
	data, ok := value.(string)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal stringList value:", value))
	}

	l := StrIntList{0, 0}
	sl := strings.Split(data, ",")
	for i, s := range sl {
		v, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		l[i] = v
	}

	*sil = l
	return nil
}

func (sil StrIntList) Value() (driver.Value, error) {
	if len(sil) == 0 {
		return "", nil
	}

	l := []string{}
	for _, v := range sil {
		l = append(l, strconv.Itoa(v))
	}

	return strings.Join(l, ","), nil
}

type Favorite struct {
	ID         int        `json:"id" gorm:"column:id;primarykey"`
	UserID     int        `json:"-" gorm:"column:user_id"`
	User       User       `json:"user" gorm:"foreignkey:UserID"`
	MangaID    int        `json:"-" gorm:"column:manga_id"`
	Manga      Manga      `json:"manga" gorm:"foreignkey:MangaID"`
	Progress   StrIntList `json:"progress" gorm:"column:progress"`
	Categories StrList    `json:"categories" gorm:"column:categories"`
}

func (f3556622c9ac203f7800dd88f8efe81126b1bbf8 *Favorite) TableName() string {
	return "favorite"
}

type Category struct {
	ID   int    `json:"-" gorm:"column:id;primarykey"`
	Name string `json:"name" gorm:"column:name"`
}

func (c *Category) TableName() string {
	return "category"
}

type User struct {
	ID   int    `json:"-" gorm:"column:id;primarykey"`
	Name string `json:"name" gorm:"column:name"`
}

func (u *User) TableName() string {
	return "user"
}
