package models

import (
	"database/sql/driver"
	"strings"
)

const ContentSeparator = "::;;::"

// Content represents a string with multiple values separated by a custom delimiter
type Content string

// NewContent creates a new Content from a slice of strings
func NewContent(values []string) Content {
	return Content(strings.Join(values, ContentSeparator))
}

// String returns the raw string representation
func (c Content) String() string {
	return string(c)
}

// Values returns the content as a slice of strings
func (c Content) Values() []string {
	if c == "" {
		return []string{}
	}
	return strings.Split(string(c), ContentSeparator)
}

// Add appends a new value to the content
func (c *Content) Add(value string) {
	if *c == "" {
		*c = Content(value)
		return
	}
	*c = Content(string(*c) + ContentSeparator + value)
}

// AddMultiple appends multiple values to the content
func (c *Content) AddMultiple(values []string) {
	for _, value := range values {
		c.Add(value)
	}
}

// Count returns the number of values in the content
func (c Content) Count() int {
	if c == "" {
		return 0
	}
	return len(c.Values())
}

// IsEmpty checks if the content is empty
func (c Content) IsEmpty() bool {
	return c == ""
}

// Scan implements the sql.Scanner interface for database reads
func (c *Content) Scan(value interface{}) error {
	if value == nil {
		*c = ""
		return nil
	}
	if sv, ok := value.(string); ok {
		*c = Content(sv)
		return nil
	}
	if bv, ok := value.([]byte); ok {
		*c = Content(string(bv))
		return nil
	}
	return nil
}

// Value implements the driver.Valuer interface for database writes
func (c Content) Value() (driver.Value, error) {
	return string(c), nil
}
