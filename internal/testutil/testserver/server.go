package testserver

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Env struct {
	Mux    *http.ServeMux
	Server *httptest.Server
}

func Setup() *Env {
	mux := http.NewServeMux()
	return &Env{
		Mux:    mux,
		Server: httptest.NewServer(mux),
	}
}

func (t *Env) TearDown() {
	t.Server.Close()
}

type HandleReqWithoutBodyOpts struct {
	Mux         *http.ServeMux
	URL         string
	RawResponse string
	Method      string
	HTTPStatus  int
	CallFlag    *bool
}

func HandleReqWithoutBody(t *testing.T, opts HandleReqWithoutBodyOpts) {
	opts.Mux.HandleFunc(opts.URL, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(opts.HTTPStatus)
		_, _ = fmt.Fprintf(w, opts.RawResponse)

		if r.Method != opts.Method {
			t.Fatalf("expected %s method but got %s", opts.Method, r.Method)
		}

		*opts.CallFlag = true
	})
}
