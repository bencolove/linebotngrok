package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

// var lineBot

var channelSecret string = ""
var channelAccessToken string = ""
var bot *linebot.Client
var err error

func init() {

	channelSecret, err = GetEnvString("CHANNEL_SECRET")
	channelAccessToken, err = GetEnvString("CHANNEL_TOKEN")

	httpClient := &http.Client{}
	bot, err = linebot.New(channelSecret, channelAccessToken,
		linebot.WithHTTPClient(httpClient))
}

func BuildLinebotHandler() (http.HandlerFunc, error) {
	if err != nil {
		return nil, err
	}
	return lineCallbackHandler, nil
}

func lineCallbackHandler(w http.ResponseWriter, r *http.Request) {
	// parse line message
	events, err := bot.ParseRequest(r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			// 400
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			// text
			handleTextMessage(event.Message.(*linebot.TextMessage), event.ReplyToken, event.Source)
		}
	}
}

type TextHandler = func(*linebot.TextMessage, string, *linebot.EventSource)

var textHandlers map[string]TextHandler

func init() {
	// add init text handlers
	textHandlers = make(map[string]TextHandler)

	textHandlers["/help"] = handleHelp
	textHandlers["/ngrok"] = handleNgrokRequest
}

func handleTextMessage(message *linebot.TextMessage, replyToken string, source *linebot.EventSource) {
	content := message.Text
	if handler, ok := textHandlers[content]; ok {
		// text handler exists
		handler(message, replyToken, source)
	}
}

func handleHelp(message *linebot.TextMessage, replyToken string, source *linebot.EventSource) {
	if _, err := bot.ReplyMessage(replyToken,
		linebot.NewTextMessage("Menus").WithQuickReplies(
			linebot.NewQuickReplyItems(
				linebot.NewQuickReplyButton(
					"",
					linebot.NewMessageAction("ngroks", "/ngrok"),
				),
			),
		)).Do(); err != nil {
		fmt.Printf("error: %+v\n", err)
	}
}

func handleNgrokRequest(message *linebot.TextMessage, replyToken string, source *linebot.EventSource) {
	//urls, err := ListTunnels()
	msg := ""
	if err != nil {
		msg = err.Error()
	} else {
		urlResText := []string{}
		for data := range GetNgrokTunnels() {
			content := ""
			if data.Err != nil {
				content = data.Err.Error()
			} else {
				content = formatNgrokTunnels(data.Val)
			}
			urlResText = append(urlResText, content)
		}
		if len(urlResText) > 0 {
			msg = strings.Join(urlResText, "\n")
		}
	}

	bot.ReplyMessage(replyToken,
		linebot.NewTextMessage(msg),
	).Do()

}

func formatNgrokTunnels(t *NgrokTunnel) string {
	if t != nil {
		b := strings.Builder{}
		if t.Protocal == "tcp" {
			b.WriteString(t.Protocal)
			b.WriteString("://")
			b.WriteString(t.PublicURL)
			b.WriteString(":")
			b.WriteString(strconv.Itoa(t.PublicPort))
			b.WriteString(" -> localhost:")
			b.WriteString(strconv.Itoa(t.LocalPort))
		} else {
			b.WriteString(t.Protocal)
			b.WriteString("://")
			b.WriteString(t.PublicURL)
		}
		return b.String()
	}
	return ""
}
