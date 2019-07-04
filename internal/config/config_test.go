package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Load(t *testing.T) {
	err := Load([]byte(`
http:
    port: 9933

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

notifier:
    count_of_skipped_hours: 8
    telegram_token: "12340255:BBBZZZJJJJJAAAEEEEE"

sentry:
  enabled: true
  key: "https://xxxxx:yyyyy@sentry.io/123456"
`))

	assert.NoError(t, err)

	assert.Equal(t, "", Config.HTTP.IP)
	assert.Equal(t, 9933, Config.HTTP.Port)

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
}
