package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Load(t *testing.T) {
	err := Load([]byte(`
db:
    type:  'mysql'
    host:  'mariadb'
    name:  'notify'
    login: 'notify'
    pass:  'notify'
    log: false

log:
    level: DEBUG
    file: 'notify.log'

notify:
    count_of_skipped_hours: 8
    telegram_token: "12340255:BBBZZZJJJJJAAAEEEEE"

sentry:
  enabled: true
  key: "https://xxxxx:yyyyy@sentry.io/123456"


stores:
  itunes:
    url: https://api.music.apple.com
    fetch_workers: 25
    name: "Apple Music"
    fetch: true
    meta:
      region: "us"


musicmash: "http://musicmash"

subscriptions: "http://subscriptions"
`))

	assert.NoError(t, err)

	assert.Equal(t, "mysql", Config.DB.Type)
	assert.Equal(t, "mariadb", Config.DB.Host)
	assert.Equal(t, "notify", Config.DB.Name)
	assert.Equal(t, "notify", Config.DB.Login)
	assert.Equal(t, "notify", Config.DB.Pass)
	assert.False(t, Config.DB.Log)

	assert.Equal(t, "DEBUG", Config.Log.Level)
	assert.Equal(t, "notify.log", Config.Log.File)

	assert.Equal(t, float64(8), Config.Notifier.CountOfSkippedHours)
	assert.Equal(t, "12340255:BBBZZZJJJJJAAAEEEEE", Config.Notifier.TelegramToken)

	assert.True(t, Config.Sentry.Enabled)
	assert.Equal(t, "https://xxxxx:yyyyy@sentry.io/123456", Config.Sentry.Key)

	assert.Equal(t, "https://api.music.apple.com", Config.Stores["itunes"].URL)
	assert.Equal(t, 25, Config.Stores["itunes"].FetchWorkers)
	assert.Len(t, Config.Stores["itunes"].Meta, 1)
	assert.Equal(t, "us", Config.Stores["itunes"].Meta["region"])
	assert.Equal(t, "Apple Music", Config.Stores["itunes"].Name)
	assert.True(t, Config.Stores["itunes"].Fetch)

	assert.Equal(t, "http://musicmash", Config.Musicmash)
	assert.Equal(t, "http://subscriptions", Config.Subscriptions)
	assert.Equal(t, "http://artists", Config.Artists)
}
