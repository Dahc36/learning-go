package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dahc36/learning-go/12-the-context/pkg/identity"
	"github.com/dahc36/learning-go/12-the-context/pkg/tracker"
)

func main() {
	addr := ":8000"
	c := controller{
		logger: tracker.Logger{},
	}
	var h http.Handler = http.HandlerFunc(c.handler)
	h = identity.Middleware(h)
	h = tracker.Middleware(h)
	s := http.Server{
		Addr:         addr,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      h,
	}
	fmt.Println("Listening for requests to", addr)
	// You can do calls like `curl --header "X-User: LokoHanks" localhost:8000/whatsuuuup`
	err := s.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}
}

type logger interface {
	Log(ctx context.Context, message string)
}

type controller struct {
	logger logger
}

func (c controller) handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	c.logger.Log(ctx, r.URL.Path)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"ok": true}`))
}
