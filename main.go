package main

import (
	"context"
	"net/http"
	"openpak/wfc-common/common"
	"openpak/wfc-common/logging"
	"openpak/wfc-common/service"
	"openpak/wfc-content/race"
	"openpak/wfc-content/sake"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	if err := common.LoadTokenKey(); err != nil {
		panic(err)
	}
	cfg := common.GetConfig()
	logging.SetLevel(*cfg.LogLevel)
	common.ReadGameList()
	sake.StartServer(false)
	defer sake.Shutdown()
	race.StartServer(false)
	defer race.Shutdown()
	mux := http.NewServeMux()
	sake.RegisterHandlers(mux)
	race.RegisterHandlers(mux)
	health := http.NewServeMux()
	health.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok\n")) })
	guarded := service.Guard(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}
		mux.ServeHTTP(w, r)
	}))
	health.Handle("/", guarded)
	addr := os.Getenv("WFC_HTTP_ADDR")
	if addr == "" {
		addr = ":20122"
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := service.Run(ctx, service.Server(addr, health)); err != nil {
		panic(err)
	}
}
