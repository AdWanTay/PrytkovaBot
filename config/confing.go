package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

var BotToken string
var AdminId int64

func Load() error {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}
	adminId, err := strconv.ParseInt(os.Getenv("ADMIN_ID"), 10, 64)
	if err != nil {
		log.Println("Invalid ADMIN_ID")
	}
	AdminId = adminId
	BotToken = os.Getenv("BOT_TOKEN")
	return nil
}
