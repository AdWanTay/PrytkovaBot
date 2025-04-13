package storage

import (
	models "PrytkovaBot/internal/models"
	"database/sql"
	"time"
)

func GetAvailableSlots(db *sql.DB) ([]models.Slot, error) {
	rows, err := db.Query("SELECT id, time FROM slots WHERE is_booked = 0 AND time > datetime('now')")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []models.Slot
	for rows.Next() {
		var s models.Slot
		err := rows.Scan(&s.ID, &s.Time)
		if err != nil {
			return nil, err
		}
		slots = append(slots, s)
	}
	return slots, nil
}

func BookSlot(db *sql.DB, slotID int, userID int64, username string) error {
	_, err := db.Exec(
		"UPDATE slots SET is_booked = 1, user_id = ?, user_name = ? WHERE id = ?",
		userID, username, slotID,
	)
	return err
}

func CreateSlots(db *sql.DB, from, to time.Time) error {
	for day := 0; day < 7; day++ {
		date := from.AddDate(0, 0, day)
		for hour := 12; hour < 16; hour++ {
			slotTime := time.Date(date.Year(), date.Month(), date.Day(), hour, 0, 0, 0, time.Local)
			_, err := db.Exec("INSERT INTO slots (time, is_booked) VALUES (?, 0)", slotTime)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
