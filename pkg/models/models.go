package models

import (
	"errors"
	"time"
)

// top-level data types that our database model will use and return

var ErrNoRecord = errors.New("models: No matching records found")

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}
