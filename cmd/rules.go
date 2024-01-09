package cmd

import (
	rei "github.com/fumiama/ReiBot"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// CheckIfTheUserIsValid Check User is valid and ready to go.
func CheckIfTheUserIsValid(ctx *rei.Ctx) (IsValid bool, UserGroup int64, UserID int64) {
	// Check User Status, including their group status, if this user is unknown.
	getUserChannelStatus := ctx.Event.Value.(*tgba.Message).From.FirstName
	if getUserChannelStatus == "Group" || getUserChannelStatus == "Channel" || ctx.Message.From.ID == 777000 { // unknownUser.
		// ignore Group | Channel || ==> Maybe Someone really set their name to Channel | Group lol.
		return false, 0, 0
	}
	if !ctx.Message.Chat.IsGroup() && !ctx.Message.Chat.IsSuperGroup() {
		// group setted to none, private chat.
		return true, 0, ctx.Message.From.ID
	}
	return true, ctx.Message.Chat.ID, ctx.Message.From.ID
}

// ReturnUser ===> User is Valid.
func ReturnUser(ctx *rei.Ctx) (UserGroup int64, UserID int64) {
	if !ctx.Message.Chat.IsGroup() && !ctx.Message.Chat.IsSuperGroup() {
		// group setted to none, private chat.
		return 0, ctx.Message.From.ID
	}
	return ctx.Message.Chat.ID, ctx.Message.From.ID
}
