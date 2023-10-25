package model

import "time"

type Tasks struct {
	ID          int
	UserID      int64
	Complexity  int
	Deadline    time.Time
	Description string
}
