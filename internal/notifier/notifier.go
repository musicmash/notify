package notifier

import (
	"time"

	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/log"
	"github.com/musicmash/musicmash/internal/notifier/telegram"
	"github.com/pkg/errors"
)

func notify(chatID int64, releases []*db.Release) {
	for i := range releases {
		message := MakeMessage(releases[i])
		message.ChatID = chatID
		if err := telegram.SendMessage(chatID, message); err != nil {
			log.Error(errors.Wrapf(err, "tried to send release to '%d'", chatID))
		}
	}
}

func Notify() {
	users, err := db.DbMgr.GetUsersWithReleases(time.Now().UTC())
	if err != nil {
		log.Error(errors.Wrap(err, "tried to get users with releases for notify stage"))
		return
	}

	last, err := db.DbMgr.GetLastActionDate(db.ActionNotify)
	if err != nil {
		log.Error(errors.Wrap(err, "tried to get last_action for notify stage"))
		return
	}

	for _, user := range users {
		chat, err := db.DbMgr.FindChatByUserName(user)
		if err != nil {
			log.Error(errors.Wrapf(err, "tried to get chat for '%s' for notify stage", user))
			continue
		}

		feed, err := db.DbMgr.GetUserFeedSince(user, last.Date)
		if err != nil {
			log.Error(errors.Wrapf(err, "tried to get feed for '%s' for notify stage", user))
			return
		}

		notify(*chat, feed.Announced)
		notify(*chat, feed.Released)
	}
}
