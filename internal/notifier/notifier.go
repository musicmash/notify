package notifier

import (
	"time"

	artsapi "github.com/musicmash/artists/pkg/api"
	mashapi "github.com/musicmash/musicmash/pkg/api"
	"github.com/musicmash/musicmash/pkg/api/releases"
	"github.com/musicmash/notify/internal/log"
	"github.com/musicmash/notify/internal/notifier/steps"
	"github.com/musicmash/notify/internal/notifier/telegram"
	subsapi "github.com/musicmash/subscriptions/pkg/api"
	"github.com/pkg/errors"
)

type Notifier struct {
	pipe *steps.Pipeline
}

func New(mashClient *mashapi.Provider, artsClient *artsapi.Provider, subsClient *subsapi.Provider) *Notifier {
	return &Notifier{pipe: steps.NewPipeline(mashClient, artsClient, subsClient)}
}

func notify(chatID int64, artistName string, release *releases.Release) error {
	message := makeMessage(artistName, release)
	message.ChatID = chatID
	if err := telegram.SendMessage(message); err != nil {
		return errors.Wrapf(err, "tried to send release to '%d'", chatID)
	}
	return nil
}

func (n *Notifier) Notify(period time.Time) error {
	items, err := n.pipe.Do(period)
	if err != nil {
		return err
	}

	for _, item := range items {
		for _, chat := range item.Chats {
			for _, release := range item.Releases {
				// TODO (m.kalinin): do not nofify if user already recieved an notification
				if err := notify(chat.ID, item.ArtistName, release); err != nil {
					log.Error(err)
				}
			}
		}
	}
	return nil
}
