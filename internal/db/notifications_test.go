package db

import (
	"testing"
	"time"

	"github.com/musicmash/notify/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestDB_Notifications_CreateAndGet(t *testing.T) {
	setup()
	defer teardown()

	// action
	now := time.Now().UTC()
	err := DbMgr.CreateNotification(&Notification{UserName: testutil.UserObjque, Date: now, ReleaseID: testutil.StoreIDQ})

	// assert
	assert.NoError(t, err)
	notifications, err := DbMgr.GetNotificationsForUser(testutil.UserObjque)
	assert.NoError(t, err)
	assert.Len(t, notifications, 1)
	assert.Equal(t, testutil.UserObjque, notifications[0].UserName)
	assert.Equal(t, now, notifications[0].Date)
	assert.Equal(t, uint64(testutil.StoreIDQ), notifications[0].ReleaseID)
}

func TestDB_Notifications_FindUsersThatReceivedNotification(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	now := time.Now()
	assert.NoError(t, DbMgr.CreateNotification(&Notification{UserName: testutil.UserObjque, Date: now, ReleaseID: testutil.StoreIDQ}))
	assert.NoError(t, DbMgr.CreateNotification(&Notification{UserName: testutil.UserObjque, Date: now, ReleaseID: 2}))
	assert.NoError(t, DbMgr.CreateNotification(&Notification{UserName: testutil.UserBot, Date: now, ReleaseID: 2}))

	// action
	users, err := DbMgr.FindUsersThatReceivedNotification(testutil.StoreIDQ, []string{testutil.UserObjque, testutil.UserBot})

	// assert
	assert.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, testutil.UserObjque, users[0])
}
