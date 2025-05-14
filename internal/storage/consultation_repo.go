package storage

import (
	"PrytkovaBot/internal/models"
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

func GetTimeBySlotId(db *sql.DB, slotId int64) (time.Time, error) {
	var t time.Time
	err := db.QueryRow("SELECT time FROM slots WHERE id = $1", slotId).Scan(&t)
	if err != nil {
		return t, err
	}
	return t, nil
}

func BookSlot(db *sql.DB, slotID int, userID int64, username string) error {
	_, err := db.Exec(
		"UPDATE slots SET is_booked = 1, user_id = ?, user_name = ? WHERE id = ?",
		userID, username, slotID,
	)
	return err
}

func GetBookedSlots(db *sql.DB) ([]models.Slot, error) {
	rows, err := db.Query(`
		SELECT id, time, user_id, user_name 
		FROM slots 
		WHERE is_booked = 1 AND time >= CURRENT_TIMESTAMP
		ORDER BY time ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []models.Slot
	for rows.Next() {
		var s models.Slot
		if err := rows.Scan(&s.ID, &s.Time, &s.UserID, &s.UserName); err != nil {
			return nil, err
		}
		slots = append(slots, s)
	}

	return slots, nil
}
func CreateSlots(db *sql.DB) error {
	// Считаем количество свободных (не забронированных) слотов в будущем
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM slots WHERE time > ? AND is_booked = 0", time.Now()).Scan(&count)
	if err != nil {
		return err
	}

	// Если свободных слотов >= 10 — ничего не делаем
	if count >= 10 {
		return nil
	}

	// Создаем недостающие слоты на 14 дней вперед (если таких еще нет)
	for day := 0; day < 14; day++ {
		date := time.Now().AddDate(0, 0, day)
		for hour := 12; hour < 16; hour++ {
			slotTime := time.Date(date.Year(), date.Month(), date.Day(), hour, 0, 0, 0, time.Local)

			// Проверим, существует ли уже такой слот
			var exists bool
			err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM slots WHERE time = ?)", slotTime).Scan(&exists)
			if err != nil {
				return err
			}
			if exists {
				continue
			}

			// Вставляем слот
			_, err = db.Exec("INSERT INTO slots (time, is_booked) VALUES (?, 0)", slotTime)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
