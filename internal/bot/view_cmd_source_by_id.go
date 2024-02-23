package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kimcodec/TgBot/internal/botkit"
	"github.com/kimcodec/TgBot/internal/storage/model"
	"strconv"
)

type SourceGetter interface {
	SourceByID(ctx context.Context, id uint64) (*model.Source, error)
}

func ViewCmdSourceByID(getter SourceGetter) botkit.ViewFunc {
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

		src, err := getter.SourceByID(ctx, uint64(id))
		if err != nil {
			return err
		}

		reply := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			fmt.Sprintf("Источник с id %d\\:\n\n%s", id, formatSource(*src)),
		)
		reply.ParseMode = tgbotapi.ModeMarkdownV2
		if _, err := bot.Send(reply); err != nil {
			return err
		}
		return nil
	}
}
