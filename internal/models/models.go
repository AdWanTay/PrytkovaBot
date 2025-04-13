package models

import "time"

type Slot struct {
	ID       int
	Date     string    // в формате YYYY-MM-DD
	Time     time.Time // в формате HH:MM
	IsBooked bool
	BookedBy int64     // Telegram ID пользователя
	BookedAt time.Time // Когда было забронировано
}
