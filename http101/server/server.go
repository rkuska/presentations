package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

const (
	ContentType           = "Content-Type"
	ContentTypeJSON       = "application/json; charset=utf-8"
	DefaultTimeoutHeaders = 3 * time.Second
	DefaultTimeoutHandler = 3 * time.Second
	DefaultIdleTimeout    = 6 * time.Minute
)

type response struct {
	Msg string
}

func SleepyHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	timer := time.NewTimer(10 * time.Millisecond)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("hell broke lose"))
	case <-timer.C:
		EncodeJSONResponse(response{strings.Repeat("x", 502)}, w)
	}
}

func EncodeJSONResponse(response interface{}, w http.ResponseWriter) error {
	w.Header().Set(ContentType, ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	return encoder.Encode(response)
}

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.Handle(
		"/sleepyget",
		http.TimeoutHandler(http.HandlerFunc(SleepyHandler), DefaultTimeoutHandler, "request handling timeout"),
	).Methods("GET")
	return r
}

type Server struct {
	*http.Server
}

func NewServer(addr string, h http.Handler, lw io.Writer) Server {
	s := http.Server{
		Addr:              addr,
		Handler:           h,
		ReadHeaderTimeout: DefaultTimeoutHeaders,
		IdleTimeout:       DefaultIdleTimeout,
		ErrorLog:          log.New(lw, "serverLogger", 0),
	}
	return Server{Server: &s}
}
