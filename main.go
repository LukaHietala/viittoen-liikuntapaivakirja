package main

import (
	"context"
	"embed"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/LukaHietala/viittoen-liikuntapaivakirja/api"
	"github.com/LukaHietala/viittoen-liikuntapaivakirja/db"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed web
var contentFS embed.FS

func main() {
	api.InitAuth("salaisuus")
	server := &http.Server{
		Addr:    "0.0.0.0:3000",
		Handler: api.Serve(contentFS),
	}

	d, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer d.Close()
	db.DB = d

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal(err)
	}
}
