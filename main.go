package main

import (
	"errors"
	"fmt"
	"log"
	"main/config"
	"net/http"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
	"github.com/yanun0323/pkg/logs"
)

func main() {
	conf := config.LoadConfig()
	logs.Infof("conf: %+v", conf)
	bot, err := messaging_api.NewMessagingApiAPI(
		conf.ChannelAccessToken,
	)
	if err != nil {
		logs.Fatalf("create bot, err: %+v", err)
	}

	// Setup HTTP Server for receiving requests from LINE platform
	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		logs.Info("/callback called...")

		cb, err := webhook.ParseRequest(conf.ChannelSecret, req)
		if err != nil {
			logs.Errorf("parse request, err: %+v", err)
			if errors.Is(err, webhook.ErrInvalidSignature) {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		logs.Info("Handling events...")
		for _, event := range cb.Events {
			log.Printf("/callback called%+v...\n", event)

			switch e := event.(type) {
			case webhook.MessageEvent:
				switch message := e.Message.(type) {
				case webhook.TextMessageContent:
					if _, err = bot.ReplyMessage(
						&messaging_api.ReplyMessageRequest{
							ReplyToken: e.ReplyToken,
							Messages: []messaging_api.MessageInterface{
								messaging_api.TextMessage{
									Text: message.Text,
								},
							},
						},
					); err != nil {
						logs.Errorf("reply message, err: %+v", err)
					} else {
						logs.Info("Sent text reply.")
					}
				case webhook.StickerMessageContent:
					replyMessage := fmt.Sprintf(
						"sticker id is %s, stickerResourceType is %s", message.StickerId, message.StickerResourceType)
					if _, err = bot.ReplyMessage(
						&messaging_api.ReplyMessageRequest{
							ReplyToken: e.ReplyToken,
							Messages: []messaging_api.MessageInterface{
								messaging_api.TextMessage{
									Text: replyMessage,
								},
							},
						}); err != nil {
						logs.Errorf("reply message, err: %+v", err)
					} else {
						logs.Info("Sent sticker reply.")
					}
				default:
					logs.Errorf("Unsupported message content: %T\n", message)
				}
			default:
				logs.Errorf("Unsupported event: %T\n", event)
			}
		}
	})

	// This is just sample code.
	// For actual use, you must support HTTPS by using `ListenAndServeTLS`, a reverse proxy or something else.
	port := conf.Port
	if port == "" {
		port = "5000"
	}
	logs.Infof("Listening on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logs.Fatalf("listen and serve, err: %+v", err)
	}
}
