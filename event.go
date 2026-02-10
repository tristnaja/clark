package main

import (
	"context"
	"strings"

	"github.com/gen2brain/beeep"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

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
			sender := v.Info.Sender.String()
			relation, isVIP := ast.VIP[sender]
			userMsg := v.Message.GetConversation()
			if userMsg == "" {
				userMsg = v.Message.GetExtendedTextMessage().GetText()
			}

			if !ast.Status || !isVIP || v.Info.IsFromMe || v.Info.IsGroup {
				return
			}

			if strings.Contains(strings.ToLower(userMsg), "get him to me") {
				beeep.Alert("Attention Master!", relation+"needs you!", "")
				sendSelf(waClient, "🚨 Attention Master!"+relation+"needs you!")
				reply(waClient, v, "I've alerted him. One Moment.")
				return
			}

			aiResp := ast.GetAIResponse(sender, userMsg)
			reply(waClient, v, aiResp)
		}
	}
}
