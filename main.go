package main

import (
	"os"

	rei "github.com/fumiama/ReiBot"
	"github.com/joho/godotenv"

	_ "github.com/MoYoez/MonoChatBot/cmd/msg"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	rei.Run(rei.Bot{
		Token:  os.Getenv("bot"),
		Buffer: 256,
		UpdateConfig: tgba.UpdateConfig{
			Offset:  0,
			Limit:   0,
			Timeout: 60,
		},
		Debug: false,
	})
}
