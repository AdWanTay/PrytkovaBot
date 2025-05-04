package main

import (
	"PrytkovaBot/internal/handlers"
	"PrytkovaBot/internal/services"
	"PrytkovaBot/internal/storage"
	"fmt"
	"time"

	tele "gopkg.in/telebot.v4"

	"PrytkovaBot/config"
	"log"
)

func main() {
	err := config.Load()
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	db, err := storage.InitDB("data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	services.CreateSlotsPerPeriod(db, 24*time.Hour)

	pref := tele.Settings{
		Token:  config.BotToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
		return
	}

	handlers.RegisterHandlers(b, config.AdminId)
	fmt.Println("BOT STARTED")

	b.Start()
}
