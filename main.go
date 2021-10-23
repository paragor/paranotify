package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {

	token := flag.String("token", "", "telegram token")
	userId := flag.String("user-id", "", "user who should receive msg")
	isReplyServer := flag.Bool("reply-server", false, "serve messages and reply user-id")
	flag.Parse()

	if *token == "" || (*userId == "" && !*isReplyServer) {
		flag.PrintDefaults()
		os.Exit(1)
	}

	bot, err := tgbotapi.NewBotAPI(*token)
	if err != nil {
		log.Fatal(err)
	}

	if err := printBotInfo(bot); err != nil {
		log.Fatal(err)
	}

	if *isReplyServer {
		if err := replyServer(bot); err != nil {
			log.Fatal(err)
		}
		return
	}

	if err := sendMessage(bot, *userId); err != nil {
		log.Fatal(err)
	}
}

func printBotInfo(bot *tgbotapi.BotAPI) error {
	me, err := bot.GetMe()
	if err != nil {
		return fmt.Errorf("cant get info about bot self: %w", err)
	}
	log.Printf("bot name: %s", me.String())

	return nil
}

func replyServer(bot *tgbotapi.BotAPI) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(tgbotapi.UpdateConfig{})

	if err != nil {
		return fmt.Errorf("cant create updates chan: %w", err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprint(update.Message.Chat.ID))
		msg.ReplyToMessageID = update.Message.MessageID

		if _, err := bot.Send(msg); err != nil {
			return fmt.Errorf("cant send msg to %d: %w", update.Message.Chat.ID, err)
		}
	}

	return nil
}

func sendMessage(bot *tgbotapi.BotAPI, userId string) error {
	userIdi, err := strconv.Atoi(userId)
	if err != nil {
		return fmt.Errorf("cant parse userId: %w", err)
	}

	log.Println("Reading stdin...")
	body, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("cant read stdin: %w", err)
	}

	msg := tgbotapi.NewMessage(int64(userIdi), string(body))
	if _, err := bot.Send(msg); err != nil {
		return fmt.Errorf("cant send msg to %d: %w", userIdi, err)
	}
	return nil
}
