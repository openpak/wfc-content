package main

import (
	"context"
	"net/http"
	"openpak/wfc-common/service"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	key := service.Key()
	race, e := service.Proxy(os.Getenv("WFC_RACE_URL"), key)
	if e != nil {
		panic(e)
	}
	sake, e := service.Proxy(os.Getenv("WFC_SAKE_URL"), key)
	if e != nil {
		panic(e)
	}
	h := http.NewServeMux()
	h.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok\n")) })
	h.Handle("/", service.Guard(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/RaceService/") {
			race.ServeHTTP(w, r)
			return
		}
		sake.ServeHTTP(w, r)
	})))
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	addr := os.Getenv("WFC_HTTP_ADDR")
	if addr == "" {
		addr = ":20122"
	}
	if e := service.Run(ctx, service.Server(addr, h)); e != nil {
		panic(e)
	}
}
