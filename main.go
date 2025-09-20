package main

import (
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	tele "gopkg.in/telebot.v4"
)

var (
	albumBuf    = make(map[string][]*tele.Message)
	albumTimers = make(map[string]*time.Timer)
	mu          sync.Mutex
)

func main() {
	chatIdEnv, err := strconv.Atoi(os.Getenv("CHAT_ID"))
	if err != nil {
		log.Fatal(err)
	}

	chatId := tele.ChatID(chatIdEnv)

	pref := tele.Settings{
		Token:  os.Getenv("TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	handler := func(c tele.Context) error {
		m := c.Message()
		if m == nil {
			return c.Send("wtf") // Ну тип реально каким образом
		}

		if m.AlbumID != "" {
			return c.Send("Альбомы (больше одного фото/видео в сообщении) не поддерживаются, т.к. telegram API вместе с telebot говно, пересылка альбомов там максимально заёбная и создателю было впадлу этим заниматься")
		}

		if _, err := b.Copy(chatId, m); err != nil {
			return c.Send(err.Error())
		}
		return c.Send("@cbsle2a")
	}

	b.Handle("/start", func(c tele.Context) error {
		return c.Send("Напишите сообщение, что бы отправить его в @cbsle2a")
	})

	b.Handle(tele.OnText, handler)
	b.Handle(tele.OnPhoto, handler)
	b.Handle(tele.OnVideo, handler)
	b.Handle(tele.OnAudio, handler)
	b.Handle(tele.OnDocument, handler)
	b.Handle(tele.OnVoice, handler)
	b.Handle(tele.OnSticker, handler)
	b.Handle(tele.OnContact, handler)
	b.Handle(tele.OnLocation, handler)

	log.Print("Bot started")

	b.Start()
}
