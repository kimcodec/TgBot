package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kimcodec/TgBot/internal/botkit"
)

func ViewCmdStart() botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		var (
			msg string = "*Основной список команда*\\:\n\n" +
				"*\\/add* \\- добавляет источник\\. В качестве аргумента принимает модель источника в формате JSON\\.\n" +
				"*\\/delete* \\- удаляет источник\\. В качестве аргумента принимает ID источника\\.\n" +
				"*\\/listsources* \\- выводит список всех источников\\.\n" +
				"*\\/source* \\- выводит источник\\. В качестве аргумента принимает ID источника\\.\n" +
				"*\\/update* \\- обновляет информацию об источник\\. В качестве аргумента принимает модель " +
				"источника в формате JSON\\(Включая ID\\)\\."
			reply = tgbotapi.NewMessage(update.Message.Chat.ID, msg)
		)
		reply.ParseMode = tgbotapi.ModeMarkdownV2

		if _, err := bot.Send(reply); err != nil {
			return err
		}
		return nil
	}
}
