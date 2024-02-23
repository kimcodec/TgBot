package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kimcodec/TgBot/internal/botkit"
	"strconv"
)

type SourceDeleter interface {
	Delete(ctx context.Context, id uint64) error
}

func DeleteSource(deleter SourceDeleter) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		id, err := strconv.Atoi(update.Message.CommandArguments())
		if err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID,
				"Введен некорректный ID\\! "+
					"ID должен быть целочисленным значением\\!")
			msg.ParseMode = tgbotapi.ModeMarkdownV2
			bot.Send(msg)
			return err
		}

		if err := deleter.Delete(ctx, uint64(id)); err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось удалить источник\\!")
			msg.ParseMode = tgbotapi.ModeMarkdownV2
			bot.Send(msg)
			return err
		}

		reply := tgbotapi.NewMessage(update.Message.Chat.ID, "Источник успешно удалено\\!")
		reply.ParseMode = tgbotapi.ModeMarkdownV2
		if _, err := bot.Send(reply); err != nil {
			return err
		}

		return nil
	}
}
