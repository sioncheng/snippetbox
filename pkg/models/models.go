package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models:no matching record found")

type Snippet struct {
	Id      int64
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}
