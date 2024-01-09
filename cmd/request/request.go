package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/FloatTech/ttl"
	"github.com/MoYoez/MonoChatBot/cmd"
	rei "github.com/fumiama/ReiBot"
)

// the ticket will be expired in 120 mins. || reply mode will turn off in 60mins if no reply.
var ticketPackage = ttl.NewCache[TicketCombination, []ChatMessage](time.Minute * 120)
var IsModeOn = ttl.NewCache[TicketCombination, UserOnModelStatus](time.Minute * 60) // set expired time to 60mins.

type UserOnModelStatus struct {
	isModeOn bool
}

type TicketCombination struct {
	// if a group is none ==> set group to 0.
	group int64
	user  int64
}

type ChatGPTResponseBody struct {
	ID      string       `json:"id"`
	Object  string       `json:"object"`
	Created int          `json:"created"`
	Model   string       `json:"model"`
	Choices []ChatChoice `json:"choices"`
}

// ChatGPTRequestBody 请求体
type ChatGPTRequestBody struct {
	Model       string        `json:"model,omitempty"` // default ==> gpt3.5-turbo
	Messages    []ChatMessage `json:"messages,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
	N           int           `json:"n,omitempty"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
}

// ChatMessage 消息
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatChoice struct {
	Index        int `json:"index"`
	Message      ChatMessage
	FinishReason string `json:"finish_reason"`
}

// Completions to Users >> Proxy can use the first one.
func Completions(messages []ChatMessage, apiKey string, isProxylink ...string) (*ChatGPTResponseBody, error) {
	com := ChatGPTRequestBody{
		Messages: messages,
	}
	// default model || user can choose which model they want to use.
	if com.Model == "" {
		com.Model = "gpt-3.5-turbo" // I have no idea if chatgpt can run itself
	}
	body, err := json.Marshal(com)
	if err != nil {
		return nil, err
	}
	// user can set proxy == > in case they cannot use.
	var req *http.Request
	if len(isProxylink) > 0 {
		req, err = http.NewRequest(http.MethodPost, isProxylink[0], bytes.NewReader(body))
	} else {
		req, err = http.NewRequest(http.MethodPost, "https://api.openai.com/v1/chat/completions", bytes.NewReader(body))
	}
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	v := new(ChatGPTResponseBody)
	if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
		return nil, err
	}
	return v, nil
}

func IsReplyModeOn(ctx *rei.Ctx) bool {
	UserGroup, Userid := cmd.ReturnUser(ctx)
	// check cache status.
	getReturn := IsModeOn.Get(TicketCombination{group: UserGroup, user: Userid})
	if getReturn.isModeOn == true {
		return true
	}
	return false
}

// Package Make Msg send quickly. || Should check user first.
func Package(ctx *rei.Ctx) (msg []ChatMessage, keys TicketCombination) {
	UserGroup, Userid := cmd.ReturnUser(ctx)
	key := TicketCombination{
		group: UserGroup,
		user:  Userid,
	}
	messages := ticketPackage.Get(key)
	messages = append(messages, ChatMessage{
		Role:    "user",
		Content: ctx.Message.Text,
	})
	return messages, key
}

// SetTicketChatMessageNext msg here ==> packed old msg || resp ==> use completions resp.
func SetTicketChatMessageNext(ctx *rei.Ctx, msg []ChatMessage, ticket TicketCombination, resp *ChatGPTResponseBody) {
	reply := resp.Choices[0].Message
	reply.Content = strings.TrimSpace(reply.Content)
	msg = append(msg, reply)
	ticketPackage.Set(ticket, msg)
	ctx.SendPlainMessage(true, reply.Content)
}

func SetReplyModeOn(userGroup int64, UserID int64) {
	IsModeOn.Set(TicketCombination{group: userGroup, user: UserID}, UserOnModelStatus{isModeOn: true})
}

// CleanChattingMemory Clean Ticket Memory
func CleanChattingMemory(userGroup int64, userID int64) {
	ticketPackage.Delete(TicketCombination{group: userGroup, user: userID})
}

// RemoveChattingModeStatus Remove ChattingMode Status
func RemoveChattingModeStatus(userGroup int64, userID int64) {
	IsModeOn.Delete(TicketCombination{group: userGroup, user: userID})
}
