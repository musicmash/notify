package db

import (
	"testing"

	"github.com/musicmash/notify/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestDB_Chat_EnsureExists(t *testing.T) {
	setup()
	defer teardown()

	// action
	err := DbMgr.EnsureChatExists(&Chat{ID: 10000420, UserName: testutil.UserObjque})

	// assert
	assert.NoError(t, err)
	chatID, err := DbMgr.FindChatByUserName(testutil.UserObjque)
	assert.NoError(t, err)
	assert.Equal(t, int64(10000420), *chatID)
}
