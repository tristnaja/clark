package internal

import (
	"context"
	_ "embed"
	"fmt"
	"strings"

	"github.com/gen2brain/beeep"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

//go:embed utils/clark.png
var clarkIcon []byte

func sendSelf(cli *whatsmeow.Client, msg string) {
	targetJID := cli.Store.ID.ToNonAD()

	cli.SendMessage(context.Background(), targetJID, &waE2E.Message{
		Conversation: proto.String(msg),
	})
}

func reply(cli *whatsmeow.Client, v *events.Message, msg string) {
	cli.SendMessage(context.Background(), v.Info.Sender, &waE2E.Message{
		Conversation: proto.String(msg),
	})
}

func EventHandler(waClient *whatsmeow.Client, ast *Assistant) whatsmeow.EventHandler {
	return func(evt any) {
		switch v := evt.(type) {
		case *events.Message:
			if v == nil || v.Info.Chat.IsEmpty() || v.Info.Sender.IsEmpty() || v.Message == nil {
				fmt.Println("Warning: Received nil message data")
				return
			}
			sender := v.Info.Sender.String()
			relation, isVIP := ast.VIP.CheckVIP(sender)

			var userMsg string
			if conversation := v.Message.GetConversation(); conversation != "" {
				userMsg = conversation
			} else if userMsg == "" {
				if extendedMessage := v.Message.GetExtendedTextMessage(); extendedMessage != nil {
					userMsg = extendedMessage.GetText()
				}
			} else {
				fmt.Println("Warning: Message has no recognizable content")
			}

			if !ast.Status || !isVIP || v.Info.IsFromMe || v.Info.IsGroup {
				return
			}

			if strings.Contains(strings.ToLower(userMsg), "get him to me") {
				beeep.Notify("Attention Sir!", relation+" needs you!", clarkIcon)
				sendSelf(waClient, "🚨 Attention Master!\n"+relation+" needs you!")
				reply(waClient, v, "I've alerted him. One Moment.")
				return
			}

			aiResp, err := ast.GetAIResponse(sender, userMsg)
			if err != nil {
				fmt.Printf("AI response error: %v\n", err)
				reply(waClient, v, "I apologize, but I'm experiencing technical difficulties. Please try again later.")
				return
			}
			reply(waClient, v, aiResp)
		}
	}
}
