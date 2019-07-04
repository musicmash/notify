package cron

import (
	"testing"
	"time"

	"github.com/musicmash/notify/internal/db"
	"github.com/stretchr/testify/assert"
)

func setup() {
	db.DbMgr = db.NewFakeDatabaseMgr()
}

func teardown() {
	_ = db.DbMgr.DropAllTables()
	_ = db.DbMgr.Close()
}

func TestCron_IsMustFetch_FirstRun(t *testing.T) {
	// first run means that no records in last_fetches
	setup()
	defer teardown()
	c := cron{ActionName: db.ActionNotify}

	// action
	must := c.IsMustFetch()

	// assert
	assert.True(t, must)
}

func TestCron_IsMustFetch_ReloadApp_AfterFetching(t *testing.T) {
	// fetch was successful and someone restart the app
	setup()
	defer teardown()
	c := cron{ActionName: db.ActionNotify, CountOfSkippedHoursToRun: 8}

	// arrange
	assert.NoError(t, db.DbMgr.SetLastActionDate(db.ActionNotify, time.Now().UTC()))

	// action
	must := c.IsMustFetch()

	// assert
	assert.False(t, must)
}

func TestCron_IsMustFetch_ReloadApp_AfterOldestFetching(t *testing.T) {
	// fetch was successful some times ago and someone restart the app
	setup()
	defer teardown()
	c := cron{ActionName: db.ActionNotify}

	// arrange
	assert.NoError(t, db.DbMgr.SetLastActionDate(db.ActionNotify, time.Now().UTC().Add(-time.Hour*48)))

	// action
	must := c.IsMustFetch()

	// assert
	assert.True(t, must)
}
