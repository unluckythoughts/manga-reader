package models

import (
	"database/sql/driver"
	"strings"
)

const listSeparator = "::;;::"

// List represents a string with multiple values separated by a custom delimiter
type List string

// NewList creates a new List from a slice of strings
func NewList(values []string) List {
	return List(strings.Join(values, listSeparator))
}

// String returns the raw string representation
func (c List) String() string {
	return string(c)
}

// Values returns the List as a slice of strings
func (c List) Values() []string {
	if c == "" {
		return []string{}
	}
	return strings.Split(string(c), listSeparator)
}

// Add appends a new value to the List
func (c *List) Add(value string) {
	if *c == "" {
		*c = List(value)
		return
	}
	*c = List(string(*c) + listSeparator + value)
}

// AddMultiple appends multiple values to the List
func (c *List) AddMultiple(values []string) {
	for _, value := range values {
		c.Add(value)
	}
}

// Count returns the number of values in the List
func (c List) Count() int {
	if c == "" {
		return 0
	}
	return len(c.Values())
}

// IsEmpty checks if the List is empty
func (c List) IsEmpty() bool {
	return c == ""
}

// Scan implements the sql.Scanner interface for database reads
func (c *List) Scan(value interface{}) error {
	if value == nil {
		*c = ""
		return nil
	}
	if sv, ok := value.(string); ok {
		*c = List(sv)
		return nil
	}
	if bv, ok := value.([]byte); ok {
		*c = List(string(bv))
		return nil
	}
	return nil
}

// Value implements the driver.Valuer interface for database writes
func (c List) Value() (driver.Value, error) {
	return string(c), nil
}
