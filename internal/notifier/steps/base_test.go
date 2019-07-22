package steps

import (
	"net/http"
	"net/http/httptest"

	"github.com/musicmash/notify/internal/db"
)

var (
	server *httptest.Server
	mux    *http.ServeMux
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	db.DbMgr = db.NewFakeDatabaseMgr()
}

func teardown() {
	_ = db.DbMgr.DropAllTables()
	_ = db.DbMgr.Close()
	server.Close()
}
