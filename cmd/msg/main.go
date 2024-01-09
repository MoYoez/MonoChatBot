package msg

import (
	"os"

	"github.com/MoYoez/MonoChatBot/cmd"
	"github.com/MoYoez/MonoChatBot/cmd/request"
	rei "github.com/fumiama/ReiBot"
)

// write in memory, not in database(. || ticket itself will break when time goes by.

func init() {
	// chat should start.
	rei.OnMessageCommand("chat", rei.OnlyToMe).SetBlock(true).Handle(func(ctx *rei.Ctx) {
		// Using command to do something.
		isValid, UserGroup, UserID := cmd.CheckIfTheUserIsValid(ctx)
		if !isValid {
			ctx.SendPlainMessage(true, "未知用户，用户不可为匿名用户~")
			return
		}
		if request.IsReplyModeOn(ctx) {
			ctx.SendPlainMessage(true, "连续对话已经开启了!")
			return
		}
		request.SetReplyModeOn(UserGroup, UserID)
		ctx.SendPlainMessage(true, "连续对话已经开始了~试着给我发信息吧w")
	})

	rei.OnMessageCommand("stop", rei.OnlyToMe).SetBlock(true).Handle(func(ctx *rei.Ctx) {
		isValid, UserGroup, UserID := cmd.CheckIfTheUserIsValid(ctx)
		if !isValid {
			ctx.SendPlainMessage(true, "未知用户，用户不可为匿名用户~")
			return
		}
		request.RemoveChattingModeStatus(UserGroup, UserID)
		ctx.SendPlainMessage(true, "咱已经停止工作了哦~")
	})

	rei.OnMessageCommand("clean").SetBlock(true).Handle(func(ctx *rei.Ctx) {
		isValid, UserGroup, UserID := cmd.CheckIfTheUserIsValid(ctx)
		if !isValid {
			ctx.SendPlainMessage(true, "未知用户，用户不可为匿名用户~")
			return
		}
		request.CleanChattingMemory(UserGroup, UserID)
		ctx.SendPlainMessage(true, "记忆已经清除了~让我们重新开始吧w")
	})

	rei.OnMessageCommand("start", rei.OnlyToMe).SetBlock(true).Handle(func(ctx *rei.Ctx) {
		ctx.SendPlainMessage(true, "MonoChat Bot Here~ \n  我可以做什么? \n\n - 实现连续性对话，来试试 /chat 开始吧~")
	})

	// only on mode it can reply./
	rei.OnMessage(rei.OnlyToMeOrToReply).SetBlock(false).Handle(func(ctx *rei.Ctx) {
		// check first, users should get their ticket.
		if !request.IsReplyModeOn(ctx) {
			return
		}
		// get chat msg.
		getMsg, getKeys := request.Package(ctx)
		sendRequest, err := request.Completions(getMsg, os.Getenv("gptkey"))
		if err != nil {
			panic(err)
		}
		request.SetTicketChatMessageNext(ctx, getMsg, getKeys, sendRequest)
	})
}
