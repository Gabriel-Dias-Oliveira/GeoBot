package main

import (
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	runBot()
}

// Start the bot.
func runBot() {
	geoBot, err := tgbotapi.NewBotAPI(botToken)

	log.Printf("Gettin Contries API...\n")

	countries := getCountriesAPI()

	log.Printf("End of data collection...\n")

	if err != nil {
		log.Panic(err)
	}

	geoBot.Debug = true

	log.Printf("Authorized on account %s", geoBot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := geoBot.GetUpdatesChan(u)

	time.Sleep(time.Millisecond * 500)
	updates.Clear() // Clean old updates.

	for update := range updates {
		if update.Message == nil { // Ignore any non-Message Updates.
			continue
		}

		log.Printf("Message from (%v): %v\n", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

			switch update.Message.Command() {
			case "start":
				msg.Text = startBot()
			case "help":
				msg.Text = helpBot()
			case "info":
				msg.Text = getCountryInfo(countries, update.Message.CommandArguments())
			case "list":
				msg.Text = listCountries(countries)
			case "play":
				msg.Text = countryGame(countries)
			case "answer":
				msg.Text = checkAnswer(update.Message.CommandArguments())
			default:
				msg.Text = defaultAnswer()
			}

			geoBot.Send(msg)
		}
	}
}
