package notifier

import (
	"github.com/jinzhu/gorm"
	artsapi "github.com/musicmash/artists/pkg/api"
	"github.com/musicmash/artists/pkg/api/artists"
	mashapi "github.com/musicmash/musicmash/pkg/api"
	"github.com/musicmash/musicmash/pkg/api/releases"
	"github.com/musicmash/notify/internal/config"
	"github.com/musicmash/notify/internal/db"
	"github.com/musicmash/notify/internal/log"
	"github.com/musicmash/notify/internal/notifier/telegram"
	subsapi "github.com/musicmash/subscriptions/pkg/api"
	"github.com/musicmash/subscriptions/pkg/api/subscriptions"
	"github.com/pkg/errors"
)

type Notifier struct {
	mashClient *mashapi.Provider
	subsClient *subsapi.Provider
	artsClient *artsapi.Provider
}

func getUniqueArtists(groupedArtists map[int64][]*releases.Release) []int64 {
	artists := []int64{}
	for artistID, _ := range groupedArtists {
		artists = append(artists, artistID)
	}
	return artists
}

func groupArtistReleases(newReleases []*releases.Release) map[int64][]*releases.Release {
	uniqueArtistsWithReleases := map[int64][]*releases.Release{}
	for _, release := range newReleases {
		if _, exist := uniqueArtistsWithReleases[release.ArtistID]; !exist {
			uniqueArtistsWithReleases[release.ArtistID] = []*releases.Release{}
		}

		uniqueArtistsWithReleases[release.ArtistID] = append(uniqueArtistsWithReleases[release.ArtistID], release)
	}
	return uniqueArtistsWithReleases
}

func groupUserSubscriptions(subs []*subscriptions.Subscription) map[string][]*subscriptions.Subscription {
	groupedSubscriptions := map[string][]*subscriptions.Subscription{}
	for _, subscription := range subs {
		if _, exist := groupedSubscriptions[subscription.UserName]; !exist {
			groupedSubscriptions[subscription.UserName] = []*subscriptions.Subscription{}
		}

		groupedSubscriptions[subscription.UserName] = append(groupedSubscriptions[subscription.UserName], subscription)
	}
	return groupedSubscriptions
}

func groupArtistInfo(arts []*artists.Artist) map[int64]*artists.Artist {
	groupedArtists := make(map[int64]*artists.Artist, len(arts))
	for _, artist := range arts {
		groupedArtists[artist.ID] = artist
	}
	return groupedArtists
}

func notify(chatID int64, artistName string, release *releases.Release) error {
	message := makeMessage(artistName, release)
	message.ChatID = chatID
	if err := telegram.SendMessage(message); err != nil {
		return errors.Wrapf(err, "tried to send release to '%d'", chatID)
	}
	return nil
}

func (n *Notifier) Notify() error {
	last, err := db.DbMgr.GetLastActionDate(db.ActionNotify)
	if err != nil {
		return errors.Wrap(err, "tried to get last_action for notify stage")
	}

	newReleases, err := releases.Get(n.mashClient, last.Date)
	if err != nil {
		return errors.Wrapf(err, "tried to get releases since %s", last.Date.String())
	}

	if len(newReleases) == 0 {
		log.Info("no new releases")
		return nil
	}

	artistWithReleases := groupArtistReleases(newReleases)
	ids := getUniqueArtists(artistWithReleases)
	subs, err := subscriptions.GetArtistsSubscriptions(n.subsClient, ids)
	if err != nil {
		return errors.Wrapf(err, "tried to get releases since %s", last.Date.String())
	}

	if len(subs) == 0 {
		log.Infof("no one subscribed on %d artists with newest releases")
		return nil
	}

	artistsInfo, err := artists.GetFullInfo(n.artsClient, ids)
	if err != nil {
		return errors.Wrap(err, "tried to get artist details")
	}
	groupedArtistsInfo := groupArtistInfo(artistsInfo)

	chats := map[string]int64{}
	for _, subscription := range subs {
		if _, exists := chats[subscription.UserName]; !exists {
			chat, err := db.DbMgr.FindChatByUserName(subscription.UserName)
			if err != nil {
				if gorm.IsRecordNotFoundError(err) {
					log.Debugf("user '%s' doesn't have a telegram chat")
					continue
				}

				log.Error(errors.Wrapf(err, "tried to get chat for '%s' for notify stage", subscription.UserName))
				continue
			}

			chats[subscription.UserName] = *chat
		}

		for _, release := range artistWithReleases[subscription.ArtistID] {
			if _, exist := config.Config.Stores[release.StoreName]; !exist {
				log.Warnf("found release in '%s', but details about store not provided in config", release.StoreName)
				continue
			}

			err := notify(chats[subscription.UserName], groupedArtistsInfo[release.ArtistID].Name, release)
			if err != nil {
				log.Error(err)
				continue
			}
		}
		//db.DbMgr.MarkReleasesAsDelivered(subscription, releases)
	}
	return nil
}
