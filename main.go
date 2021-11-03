package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {

	token := flag.String("token", "", "telegram token")
	userId := flag.String("user-id", "", "user who should receive msg")
	isReplyServer := flag.Bool("reply-server", false, "serve messages and reply user-id")
	oldUsage := flag.Usage
	flag.Usage = func() {
		oldUsage()
		fmt.Printf(`example: 
	%s -token=${TOKEN} -reply-server
	echo this is echo msg | %s -token=${TOKEN} -user-id=${USER}
`, os.Args[0], os.Args[0])
	}

	flag.Parse()

	if *token == "" || (*userId == "" && !*isReplyServer) {
		flag.Usage()
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

	if err := sendStdin(bot, *userId); err != nil {
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

func sendStdin(bot *tgbotapi.BotAPI, userId string) error {
	const maxMessageSize = 3000
	userIdi, err := strconv.Atoi(userId)
	if err != nil {
		return fmt.Errorf("cant parse userId: %w", err)
	}

	//todo потоковое чтение?)
	log.Println("Reading stdin...")
	body, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("cant read stdin: %w", err)
	}

	lines := strings.Split(string(body), "\n")

	builder := strings.Builder{}
	for i := 0; i < len(lines); i++ {
		if builder.Len() + len(lines[i]) + 1 < maxMessageSize {
			builder.WriteString(lines[i])
			builder.WriteString("\n")
			continue
		}
		if len(lines[i]) + 1 >= maxMessageSize {
			chunks := splitStringByChunks(lines[i], maxMessageSize)
			for i, chunk := range chunks {
				if i == len(chunks) {
					chunk += "\n"
				}
				if err := sendMessage(bot, userIdi, chunk); err != nil {
					return err
				}
			}
			builder.Reset()
			continue
		}

		if err := sendMessage(bot, userIdi, builder.String()); err != nil {
			return err
		}
		builder.Reset()
		builder.WriteString(lines[i])
		builder.WriteString("\n")
	}

	if builder.Len() > 0 {
		if err := sendMessage(bot, userIdi, builder.String()); err != nil {
			return err
		}
	}
	return nil
}

func sendMessage(bot *tgbotapi.BotAPI, userId int, msg string) error {
	fmt.Println(msg)
	if _, err := bot.Send(tgbotapi.NewMessage(int64(userId), msg)); err != nil {
		return fmt.Errorf("cant send msg to %d: %w", userId, err)
	}

	return nil
}

func splitStringByChunks(body string, size int) []string {
	chunks := []string{}
	for len(body) != 0 {
		l := size
		if len(body) < l {
			l = len(body)
		}
		chunks = append(chunks, body[:l])
		body = body[l:]
	}

	return chunks
}
