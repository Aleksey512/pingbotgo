package tgbot

import (
	"context"
	"fmt"
	"log"
	"misbotgo/internal/app/ping"
	"misbotgo/internal/config"
	"misbotgo/internal/database/storage"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotHandler struct {
	bot     *tgbotapi.BotAPI
	cfg     *config.Settings
	storage storage.Storage
}

func NewBotHandler(
	bot *tgbotapi.BotAPI, cfg *config.Settings, storage storage.Storage,
) *BotHandler {
	return &BotHandler{
		bot:     bot,
		cfg:     cfg,
		storage: storage,
	}
}

func (h *BotHandler) Start(ctx context.Context) error {
	h.bot.Debug = true
	log.Printf("Authorized on account %s", h.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	updates := h.bot.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			return nil
		case update := <-updates:
			h.handleUpdate(ctx, update)
		}
	}
}

func (h *BotHandler) handleUpdate(ctx context.Context, update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	if update.Message.IsCommand() {
		msg.Text = h.handleCommand(ctx, update.Message)
	} else {
		msg.Text = h.handleMessage(update.Message)
	}
	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func (h *BotHandler) handleMessage(msg *tgbotapi.Message) string {
	return msg.Text
}

func (h *BotHandler) handleCommand(ctx context.Context, msg *tgbotapi.Message) string {
	switch msg.Command() {
	case "start":
		return h.handleStartCommand()
	case "subscribe":
		return h.handleSubscribeCommand(ctx, msg)
	case "unsubscribe":
		return h.handleUnsubscribeCommand(ctx, msg)
	case "ping_now":
		return h.handlePingNowCommand()
	case "config":
		return h.handleConfigCommand()
	default:
		return "Unknown command"
	}
}

func (h *BotHandler) handleStartCommand() string {
	return "Бот для мониторинга серверов.\n\n" +
		"Команды:\n" +
		"/subscribe - подписаться на уведомления\n" +
		"/unsubscribe - отписаться от уведомлений\n" +
		"/ping_now - проверить серверы сейчас\n" +
		"/config - посмотреть список пингуемых серверов\n"
}

func (h *BotHandler) handlePingNowCommand() string {
	return ping.GeneratePingReport(h.cfg.Servers)
}

func (h *BotHandler) handleSubscribeCommand(ctx context.Context, msg *tgbotapi.Message) string {
	chatIDStr := strconv.FormatInt(msg.Chat.ID, 10)

	redisCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := h.storage.AddChatID(redisCtx, chatIDStr); err != nil {
		log.Printf("Failed to subscribe chat %d: %v", msg.Chat.ID, err)
		return "Failed to subscribe. Please try again later."
	}
	return "You've been successfully subscribed to notifications!"
}

func (h *BotHandler) handleUnsubscribeCommand(ctx context.Context, msg *tgbotapi.Message) string {
	chatIDStr := strconv.FormatInt(msg.Chat.ID, 10)

	redisCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := h.storage.RemoveChatID(redisCtx, chatIDStr); err != nil {
		log.Printf("Failed to unsubscribe chat %d: %v", msg.Chat.ID, err)
		return "Failed to unsubscribe. Please try again later."
	}
	return "You've been successfully unsubscribed from notifications!"
}

func (h *BotHandler) handleConfigCommand() string {
	var report strings.Builder
	for name, ip := range h.cfg.Servers {
		report.WriteString(fmt.Sprintf("%s - %s\n", name, ip))
	}

	return report.String()
}
