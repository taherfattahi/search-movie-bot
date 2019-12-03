package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/imroc/req"
	"github.com/tidwall/gjson"
	"log"
)

func initialTeleBot() (*tgbotapi.BotAPI, tgbotapi.UpdateConfig) {
	bot, err := tgbotapi.NewBotAPI("MyAwesomeBotToken")

	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return bot, u
}

func main() {

	bot, u := initialTeleBot()

	var yourApiKey string

	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		//if update.Message == nil { // ignore any non-Message Updates
		//	continue
		//}

		//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		//msg.ReplyToMessageID = update.Message.MessageID
		//
		//bot.Send(msg)

		if update.CallbackQuery == nil {

			authHeader := req.Header{
				"Accept": "application/json",
			}

			response, err := req.Get("http://www.omdbapi.com/?apikey=" + yourApiKey + "&s=" + update.Message.Text, authHeader)

			if err != nil {
				log.Panic(err)
			}

			resString, _ := response.ToString()

			result := gjson.Get(resString, "Search")

			var nestResult []gjson.Result
			var nestInlineKeyboard []tgbotapi.InlineKeyboardButton

			//msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

			i := 0
			j := 1
			result.ForEach(func(key, value gjson.Result) bool {
				nestResult = append(nestResult, gjson.GetMany(value.String(), "Title", "imdbID")...)

				nestInlineKeyboard = append(nestInlineKeyboard, tgbotapi.NewInlineKeyboardButtonData(nestResult[i].Str, nestResult[j].Str))

				//msg.Text  = nestResult[i].Str
				//bot.Send(msg)

				i = i + 2
				j = j + 2
				return true // keep iterating
			})

			//fmt.Println(nestResult)

			var inlineKey [][]tgbotapi.InlineKeyboardButton

			for i := 0; i < len(nestInlineKeyboard); i++ {
				inlineKey = append(inlineKey, nestInlineKeyboard[i:i+1])

			}

			var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
				inlineKey...,
			)

			if update.Message != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

				msg.ReplyMarkup = numericKeyboard
				bot.Send(msg)
			}

		} else {

			fmt.Print(update)

			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))

			//bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data))

		}


		//var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		//	tgbotapi.NewInlineKeyboardRow(
		//		tgbotapi.NewInlineKeyboardButtonURL("1.com","http://1.com"),
		//		tgbotapi.NewInlineKeyboardButtonSwitch("2sw","open 2"),
		//		tgbotapi.NewInlineKeyboardButtonData("3","3"),
		//	),
		//	tgbotapi.NewInlineKeyboardRow(
		//		tgbotapi.NewInlineKeyboardButtonData("4","4"),
		//		tgbotapi.NewInlineKeyboardButtonData("5","5"),
		//		tgbotapi.NewInlineKeyboardButtonData("6","6"),
		//	),
		//)

		//if update.CallbackQuery != nil{
		//	fmt.Print(update)
		//
		//	bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID,update.CallbackQuery.Data))
		//
		//	bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,update.CallbackQuery.Data))
		//}

		//if update.Message != nil {
		//	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		//
		//	msg.ReplyMarkup = numericKeyboard
		//	bot.Send(msg)
		//}

		//if update.Message != nil {
		//	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		//
		//	switch update.Message.Text {
		//	case "open":
		//		msg.ReplyMarkup = numericKeyboard
		//	}
		//
		//	bot.Send(msg)
		//}

		//if update.Message.IsCommand() {
		//	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		//	switch update.Message.Command() {
		//	case "help":
		//		msg.Text = "please type your movie name"
		//	case "sayhi":
		//		msg.Text = "Hi :)"
		//	case "status":
		//		msg.Text = "I'm ok."
		//	default:
		//		msg.Text = "I don't know that command"
		//	}
		//
		//	bot.Send(msg)
		//}
	}
}
