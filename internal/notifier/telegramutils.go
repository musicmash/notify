package notifier

import (
	"fmt"
	"time"

	"github.com/musicmash/musicmash/pkg/api/releases"
	"github.com/musicmash/notify/internal/config"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func makeText(artistName string, release *releases.Release) string {
	releaseDate := ""
	state := "released"
	if release.Released.After(time.Now().UTC()) {
		state = "announced"
		releaseDate = fmt.Sprintf("\nRelease date: %s", release.Released.Format(time.RFC850))
	}

	poster := fmt.Sprintf("[‌‌](%s)", release.Poster)
	return fmt.Sprintf("New album %s \n*%s*\n%s%s %s", state, artistName, release.Title, releaseDate, poster)
}

func makeButtons(release *releases.Release) *[][]tgbotapi.InlineKeyboardButton {
	buttons := [][]tgbotapi.InlineKeyboardButton{}
	storeDetails := config.Config.Stores[release.StoreName]
	buttonLabel := fmt.Sprintf("Listen on %s", storeDetails.Name)
	url := fmt.Sprintf(config.Config.Stores[release.StoreName].ReleaseURL, release.StoreID)
	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL(buttonLabel, url)))
	return &buttons
}

func makeMessage(artistName string, release *releases.Release) *tgbotapi.MessageConfig {
	message := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ReplyToMessageID: 0,
			ReplyMarkup:      tgbotapi.InlineKeyboardMarkup{InlineKeyboard: *makeButtons(release)},
		},
		Text:                  makeText(artistName, release),
		ParseMode:             "markdown",
		DisableWebPagePreview: false,
	}
	return &message
}
