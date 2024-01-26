package main

import (
	"errors"
	"net/http"
	"time"
)

type pressureGauge struct {
	ch chan struct{}
}

func new(limit int) *pressureGauge {
	ch := make(chan struct{}, limit)
	for i := 0; i < limit; i++ {
		ch <- struct{}{}
	}

	return &pressureGauge{ch: ch}
}

func (pg *pressureGauge) Process(f func()) error {
	select {
	case <-pg.ch:
		f()
		pg.ch <- struct{}{}
		return nil
	default:
		return errors.New("no more capacity")
	}
}

func doThingThatShouldBeLimited() string {
	time.Sleep(time.Second * 2)
	return "done"
}

func main() {
	pg := new(3)
	http.HandleFunc("/request", func(w http.ResponseWriter, r *http.Request) {
		err := pg.Process(func() {
			w.Write([]byte(doThingThatShouldBeLimited()))
		})
		if err != nil {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("Too many requests"))
		}
	})
	http.ListenAndServe(":8080", nil)
}
