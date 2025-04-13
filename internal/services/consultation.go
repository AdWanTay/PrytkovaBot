package services

import (
	"PrytkovaBot/internal/storage"
	"database/sql"
	"fmt"
	"strings"
)

func FormatAvailableSlots(db *sql.DB) (string, error) {
	slots, err := storage.GetAvailableSlots(db)
	if err != nil {
		return "", err
	}
	if len(slots) == 0 {
		return "Нет доступных слотов 🫠", nil
	}

	var b strings.Builder
	b.WriteString("Свободные слоты:\n")
	for _, s := range slots {
		b.WriteString(fmt.Sprintf("🕒 %s /book_%d\n", s.Time.Format("02 Jan 15:04"), s.ID))
	}
	return b.String(), nil
}
