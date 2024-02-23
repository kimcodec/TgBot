package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kimcodec/TgBot/internal/botkit"
	"github.com/kimcodec/TgBot/internal/storage/model"
)

type SourceUpdater interface {
	Update(ctx context.Context, source model.Source) error
}

func ViewCmdUpdateSource(updater SourceUpdater) botkit.ViewFunc {
	type updateSourceArgs struct {
		Id   uint64 `json:"id"`
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		args, err := botkit.ParseJSON[updateSourceArgs](update.Message.CommandArguments())
		if err != nil {
			return err
		}

		source := model.Source{
			ID:      args.Id,
			Name:    args.Name,
			FeedURL: args.URL,
		}

		if err := updater.Update(ctx, source); err != nil {
			return err
		}

		reply := tgbotapi.NewMessage(update.Message.Chat.ID, "Источник успешно обновлен\\!")
		reply.ParseMode = tgbotapi.ModeMarkdownV2
		if _, err := bot.Send(reply); err != nil {
			return err
		}
		return nil
	}
}
