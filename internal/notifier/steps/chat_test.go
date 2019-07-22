package steps

import (
	"testing"
	"time"

	"github.com/musicmash/notify/internal/db"
	"github.com/musicmash/notify/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestNotifierSteps_Chat(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, db.DbMgr.EnsureChatExists(&db.Chat{UserName: testutil.UserObjque, ID: 29}))
	subscribers := []string{testutil.UserObjque, testutil.UserBot}
	items := []*Item{{ArtistID: testutil.StoreIDQ, Subscribers: subscribers}}
	step := NewFindSubscriberChatsStep()

	// action
	items, err := step.Do(time.Now(), items)

	// assert
	assert.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Equal(t, int64(testutil.StoreIDQ), items[0].ArtistID)
	assert.Len(t, items[0].Subscribers, 1)
	assert.Equal(t, testutil.UserObjque, items[0].Subscribers[0])
	assert.Len(t, items[0].Chats, 1)
	assert.Equal(t, testutil.UserObjque, items[0].Chats[0].UserName)
	assert.Equal(t, int64(29), items[0].Chats[0].ID)
	assert.Empty(t, items[0].Releases)
	assert.Empty(t, items[0].ArtistName)
}
