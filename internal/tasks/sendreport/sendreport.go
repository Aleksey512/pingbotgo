package sendreport

import (
	"context"
	"log"
	"misbotgo/internal/app/ping"
	"misbotgo/internal/config"
	"misbotgo/internal/database/storage"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendReportToSubscribersTask(
	ctx context.Context,
	cfg *config.Settings,
	storage storage.Storage,
	bot *tgbotapi.BotAPI,
) error {
	storageCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	report := ping.GeneratePingReport(cfg.Servers)
	chatIds, err := storage.GetChatIDs(storageCtx)
	if err != nil {
		log.Printf("Error in get chat ids %v", err)
		return err
	}
	for _, val := range chatIds {
		idInt, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			log.Printf("Error convert id: %v", err)
			continue
		}

		msg := tgbotapi.NewMessage(idInt, report)
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Error sending message: %v", err)
		}
		time.Sleep(time.Second / 10) // 0.1 second
	}
	return nil
}
