package main

import (
	"context"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	"github.com/kimcodec/TgBot/internal/bot"
	"github.com/kimcodec/TgBot/internal/botkit"
	"github.com/kimcodec/TgBot/internal/summary"
	_ "github.com/lib/pq"
	"os"
	"os/signal"
	"syscall"

	"github.com/kimcodec/TgBot/internal/config"
	"github.com/kimcodec/TgBot/internal/fetcher"
	"github.com/kimcodec/TgBot/internal/notifier"
	"github.com/kimcodec/TgBot/internal/storage"

	"log"
)

func main() {
	cfg := config.Get()
	botAPI, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		log.Printf("failed to create bot: %v", err)
		return
	}

	db, err := sqlx.Connect("postgres", cfg.DatabaseDSN)
	if err != nil {
		log.Printf("failed to connect to database: %v", err)
		return
	}
	defer db.Close()

	var (
		articleStorage = storage.NewArticlePostgresStorage(db)
		sourceStorage  = storage.NewSourcePostgresStorage(db)
		fetcher        = fetcher.NewFetcher(
			articleStorage,
			sourceStorage,
			cfg.FetchInterval,
			cfg.FilterKeywords,
		)
		notifier = notifier.New(
			articleStorage,
			summary.NewOpenAISummarizer(cfg.OpenAIKey, cfg.OpenAIPrompt),
			botAPI,
			cfg.NotificationInterval,
			2*cfg.FetchInterval,
			cfg.TelegramChannelID,
		)
	)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	newsBot := botkit.NewBot(botAPI)
	newsBot.RegisterCmdView("start", bot.ViewCmdStart())

	go func(ctx context.Context) {
		if err := fetcher.Start(ctx); err != nil {
			if errors.Is(err, context.Canceled) {
				log.Printf("[ERROR] failed to start fetcher: %v", err)
				return
			}
		}
	}(ctx)

	go func(ctx context.Context) {
		if err := notifier.Start(ctx); err != nil {
			if errors.Is(err, context.Canceled) {
				log.Printf("[ERROR] failed to start notifier: %v", err)
				return
			}
		}
	}(ctx)

	if err := newsBot.Run(ctx); err != nil {
		if errors.Is(err, context.Canceled) {
			log.Printf("[ERROR] failed to start bot: %v", err)
			return
		}
		log.Println("Bot stopped")
	}
}
