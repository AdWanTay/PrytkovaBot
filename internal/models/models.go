package models

import "time"

type Slot struct {
	ID       int
	Time     time.Time // в формате HH:MM
	IsBooked bool
	UserID   int
	UserName string
}
