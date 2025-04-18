package main

import (
	"PrytkovaBot/internal/handlers"
	"time"

	tele "gopkg.in/telebot.v4"

	"PrytkovaBot/config"
	"log"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	err := config.Load()
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

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
	b.Start()
}
