package models

import (
	"errors"
	"time"
)

// top-level data types that our database model will use and return

var (
	ErrNoRecord           = errors.New("models: No matching records found")
	ErrInvalidCredentials = errors.New("model: Invalid Credentials")
	ErrDuplicateEmail     = errors.New("models: Duplicate Email")
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword string
	Created        string
	Active         bool
}
