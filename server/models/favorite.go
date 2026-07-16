package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/unluckythoughts/go-microservice/v2/tools/auth"
	"gorm.io/gorm"
)

type Progress string

func (p Progress) GetValue() (string, uint, error) {
	strProgress := string(p)
	var chapter string
	var page uint
	_, err := fmt.Sscanf(strProgress, "%s:%d", &chapter, &page)
	if err != nil {
		return "", 0, fmt.Errorf("failed to parse Progress: %v", err)
	}
	return chapter, page, nil
}

func (p *Progress) ScanValue(chapter string, page int) {
	*p = Progress(fmt.Sprintf("%s:%d", chapter, page))
}

func (p Progress) Value() (driver.Value, error) {
	return string(p), nil
}

func (p *Progress) Scan(value any) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("failed to scan Progress: expected string, got %T", value)
	}
	*p = Progress(str)
	return nil
}

type FavoriteData struct {
	Sources  []uint            `json:"sources,omitempty"`
	Books    map[uint]Progress `json:"books,omitempty"`
	Settings map[string]string `json:"settings,omitempty"`
}

func (f FavoriteData) Value() (driver.Value, error) {
	return json.Marshal(f)
}

func (f *FavoriteData) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan FavoriteData: expected []byte, got %T", value)
	}
	return json.Unmarshal(bytes, f)
}

// Favorite represents a user's favorite book with progress tracking
type Favorite struct {
	gorm.Model
	UserID uint         `gorm:"column:user_id" json:"user_id,omitempty"`
	Data   FavoriteData `gorm:"column:data;type:jsonb" json:"data,omitempty"`

	// Relationships
	User *auth.User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name for Favorite model
func (Favorite) TableName() string {
	return "favorites"
}
