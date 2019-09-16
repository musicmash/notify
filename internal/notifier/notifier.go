package notifier

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	artsapi "github.com/musicmash/artists/pkg/api"
	mashapi "github.com/musicmash/musicmash/pkg/api"
	"github.com/musicmash/musicmash/pkg/api/releases"
	"github.com/musicmash/notify/internal/db"
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

func markReleaseAsDeliveredTo(userName string, releaseID uint64) error {
	return db.DbMgr.CreateNotification(&db.Notification{
		Date:      time.Now().UTC(),
		UserName:  userName,
		ReleaseID: releaseID,
	})
}

func isComing(release *releases.Release) bool {
	// if release day tomorrow or later, than that means coming release is here
	return release.Released.After(time.Now().UTC().Truncate(24 * time.Hour))
}

func (n *Notifier) Notify(period time.Time) error {
	items, err := n.pipe.Do(period)
	if err != nil {
		return err
	}

	for _, item := range items {
		for _, chat := range item.Chats {
			for _, release := range item.Releases {
				_, err := db.DbMgr.IsUserNotified(chat.UserName, release.ID, isComing(release))
				switch err {
				case nil:
					log.Debugln(fmt.Sprintf("user '%s' already notified about '%d'", chat.UserName, release.ID))
					continue
				case gorm.ErrRecordNotFound:
					break
				default:
					log.Error(err)
					continue
				}

				if err := notify(chat.ID, item.ArtistName, release); err != nil {
					log.Error(err)
					continue
				}

				_ = markReleaseAsDeliveredTo(chat.UserName, release.ID)
			}
		}
	}
	return nil
}
