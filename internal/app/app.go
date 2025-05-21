package app

import (
	"context"
	"log"
	"misbotgo/internal/config"
	"misbotgo/internal/database/sqlite"
	"misbotgo/internal/database/storage"
	"misbotgo/internal/scheduler"
	"misbotgo/internal/tasks/sendreport"
	"misbotgo/internal/tgbot"
	"os"
	"os/signal"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Main() int {
	cfg, bot, storage, sch, err := initialize()
	if err != nil {
		log.Panicf("Initialization failed: %v", err)
	}
	defer closeResources(storage)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	addTasks(sch, cfg, bot, storage)

	errChan := make(chan error, 1)
	go startBotHandler(ctx, bot, cfg, storage, errChan)
	sch.Start(ctx)

	waitForShutdown(cancel, errChan)

	sch.Stop()

	log.Println("Application shutdown complete")
	return 0
}

func initialize() (
	*config.Settings,
	*tgbotapi.BotAPI,
	*sqlite.SQLiteStorage,
	*scheduler.Scheduler,
	error,
) {
	cfg, err := config.NewSettings()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TgBotToken)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	sqliteStorage, err := sqlite.NewSQLiteStorage(cfg)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	schdler := scheduler.NewScheduler()

	return cfg, bot, sqliteStorage, schdler, nil
}

func closeResources(storage storage.Storage) {
	if err := storage.Close(); err != nil {
		log.Printf("Error closing storage connection: %v", err)
	}
}

func startBotHandler(
	ctx context.Context,
	bot *tgbotapi.BotAPI,
	cfg *config.Settings,
	storage storage.Storage,
	errChan chan<- error,
) {
	botHandler := tgbot.NewBotHandler(bot, cfg, storage)
	log.Println("Starting bot")
	if err := botHandler.Start(ctx); err != nil {
		errChan <- err
	}
}

func addTasks(
	sch *scheduler.Scheduler,
	cfg *config.Settings,
	bot *tgbotapi.BotAPI,
	storage storage.Storage,
) {
	sch.AddTask("send_report_to_subs", 5*time.Minute, func(ctx context.Context) error {
		return sendreport.SendReportToSubscribersTask(ctx, cfg, storage, bot)
	})
}

func waitForShutdown(cancel context.CancelFunc, errChan <-chan error) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigChan:
		log.Printf("Received signal: %v. Shutting down...", sig)
		cancel()
		select {
		case <-time.After(5 * time.Second):
			log.Println("Shutdown timeout exceeded, forcing exit")
		case <-errChan:
			log.Println("App stopped gracefully")
		}
	case err := <-errChan:
		log.Printf("App stopped with error: %v", err)
		cancel()
	}
}
