package steps

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/musicmash/notify/internal/testutil"
	subsapi "github.com/musicmash/subscriptions/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestNotifierSteps_Subscription(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	items := []*Item{{ArtistID: testutil.StoreIDQ}, {ArtistID: testutil.StoreIDW}}
	client := subsapi.NewProvider(server.URL, 1)
	step := NewFindSubscribersStep(client)
	mux.HandleFunc("/v1/subscriptions", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		_, _ = fmt.Fprintf(w,
			`[{"artist_id":%d,"user_name":"%s"},
					 {"artist_id":%d,"user_name":"%s"},
					 {"artist_id":%d,"user_name":"%s"}]`,
			testutil.StoreIDQ, testutil.UserObjque,
			testutil.StoreIDW, testutil.UserObjque,
			testutil.StoreIDW, testutil.UserBot,
		)
	})

	// action
	items, err := step.Do(time.Now(), items)

	// assert
	assert.NoError(t, err)
	assert.Len(t, items, 2)
	// item 1
	assert.Equal(t, int64(testutil.StoreIDQ), items[0].ArtistID)
	assert.Empty(t, items[0].ArtistName)
	assert.Empty(t, items[0].Chats)
	assert.Len(t, items[0].Subscribers, 1)
	assert.Equal(t, []string{testutil.UserObjque}, items[0].Subscribers)
	// item 2
	assert.Equal(t, int64(testutil.StoreIDW), items[1].ArtistID)
	assert.Empty(t, items[1].ArtistName)
	assert.Empty(t, items[1].Chats)
	assert.Len(t, items[1].Subscribers, 2)
	assert.Equal(t, []string{testutil.UserObjque, testutil.UserBot}, items[1].Subscribers)
}
