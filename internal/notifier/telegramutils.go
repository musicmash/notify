package notifier

import (
	"fmt"
	"time"

	"github.com/musicmash/notify/internal/config"
	"github.com/musicmash/notify/internal/db"
	"github.com/musicmash/notify/internal/log"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func makeText(release *db.Release) string {
	releaseDate := ""
	state := "released"
	if release.Released.After(time.Now().UTC()) {
		state = "announced"
		releaseDate = fmt.Sprintf("\nRelease date: %s", release.Released.Format(time.RFC850))
	}

	poster := fmt.Sprintf("[‌‌](%s)", release.Poster)
	return fmt.Sprintf("New album %s \n*%s*\n%s%s %s", state, release.ArtistName, release.Title, releaseDate, poster)
}

func makeButtons(release *db.Release) *[][]tgbotapi.InlineKeyboardButton {
	buttons := [][]tgbotapi.InlineKeyboardButton{}
	for _, store := range release.Stores {
		storeDetails, ok := config.Config.Stores[store.StoreName]
		if !ok {
			log.Errorf("Can't make button for store '%s'. Is store exists in the config/stores section?", store.StoreName)
			continue
		}

		buttonLabel := fmt.Sprintf("Open in %s", storeDetails.Name)
		url := fmt.Sprintf(config.Config.Stores[store.StoreName].ReleaseURL, store.StoreID)
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL(buttonLabel, url)))
	}
	return &buttons
}

func MakeMessage(release *db.Release) *tgbotapi.MessageConfig {
	message := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ReplyToMessageID: 0,
			ReplyMarkup:      tgbotapi.InlineKeyboardMarkup{InlineKeyboard: *makeButtons(release)},
		},
		Text:                  makeText(release),
		ParseMode:             "markdown",
		DisableWebPagePreview: false,
	}
	return &message
}
