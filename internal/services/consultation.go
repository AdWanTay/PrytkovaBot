package services

import (
	"PrytkovaBot/internal/storage"
	"database/sql"
	"fmt"
	"gopkg.in/telebot.v4"
	"log"
	"strings"
	"time"
)

func FormatAvailableSlots(db *sql.DB) ([][]telebot.InlineButton, error) {
	buttonsInRow := 3
	maxRows := 6
	slots, err := storage.GetAvailableSlots(db)
	if err != nil {
		return nil, err
	}

	if len(slots) == 0 {
		return nil, fmt.Errorf("нет доступных слотов 🫠")
	}
	var buttons [][]telebot.InlineButton

	for j := 0; j < len(slots)/buttonsInRow && j < maxRows; j++ {
		var row []telebot.InlineButton
		for i := 0; i < buttonsInRow && j*buttonsInRow+i < len(slots); i++ {
			row = append(row, telebot.InlineButton{
				Text:   fmt.Sprintf("%s", slots[j*buttonsInRow+i].Time.Format("02.01 • 15:04")),
				Unique: fmt.Sprintf("book@%d", slots[j*buttonsInRow+i].ID),
			})
		}
		buttons = append(buttons, row)
	}

	return buttons, nil
}

var months = [...]string{
	"января", "февраля", "марта", "апреля", "мая", "июня",
	"июля", "августа", "сентября", "октября", "ноября", "декабря",
}

var weekdays = [...]string{
	"воскресенье", "понедельник", "вторник", "среда", "четверг", "пятница", "суббота",
}

func FormatBookedSlots(db *sql.DB) (string, error) {
	slots, err := storage.GetBookedSlots(db)
	if err != nil {
		return "", err
	}

	if len(slots) == 0 {
		return "Нет забронированных слотов.", nil
	}

	var sb strings.Builder
	sb.WriteString("Записи на консультацию:\n\n")

	var prevDate string
	for _, s := range slots {
		currentDate := formatDateRuWithWeekday(s.Time)
		if prevDate != currentDate {
			if prevDate != "" {
				sb.WriteString("\n")
			}
			sb.WriteString(fmt.Sprintf("📅 %s\n", currentDate))
			prevDate = currentDate
		}
		line := fmt.Sprintf("• %s — @%s (id %d)\n",
			s.Time.Format("15:04"),
			s.UserName,
			s.UserID,
		)
		sb.WriteString(line)
	}

	return sb.String(), nil
}

// formatDateRuWithWeekday форматирует дату в стиле "17 апреля (среда)"
func formatDateRuWithWeekday(t time.Time) string {
	day := t.Day()
	month := months[t.Month()-1]
	weekday := weekdays[t.Weekday()] // Sunday == 0
	return fmt.Sprintf("%d %s (%s)", day, month, weekday)
}

func CreateSlotsPerPeriod(db *sql.DB, period time.Duration) {
	// Запускаем горутину для периодического создания слотов
	go func() {
		ticker := time.NewTicker(period)
		defer ticker.Stop()

		// Первоначальный запуск создания слотов
		err := storage.CreateSlots(db)
		if err != nil {
			log.Printf("Ошибка при первоначальном создании слотов: %v", err)
		}

		for {
			select {
			case <-ticker.C:
				// Периодическое создание слотов
				err := storage.CreateSlots(db)
				if err != nil {
					log.Printf("Ошибка при создании слотов: %v", err)
				} else {
					log.Println("Слоты для консультаций успешно созданы")
				}
			}
		}
	}()
}
