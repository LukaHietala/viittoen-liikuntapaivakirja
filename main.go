package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lukahietala/viittoen-liikuntapaivakirja/api"
	"github.com/lukahietala/viittoen-liikuntapaivakirja/auth"
	"github.com/lukahietala/viittoen-liikuntapaivakirja/db"
	"github.com/lukahietala/viittoen-liikuntapaivakirja/services"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := services.LoadEnv()
	if err != nil {
		log.Fatal(err)
	}

	conn, err := db.Connect("viittoen.db?_journal_mode=WAL&_foreign_keys=on&_busy_timeout=5000")
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer conn.Close()

	store := db.NewStore(conn)
	mailService := services.NewMailService(cfg.STMP_SERVER, cfg.STMP_DOMAIN, cfg.STMP_PASSWORD)
	authService := auth.New(cfg.JWTSecret, store)

	handler := api.NewHandler(store, authService, mailService)

	httpServer := &http.Server{
		Addr:    ":3000",
		Handler: handler,
	}

	go func() {
		log.Printf("listening on %s", httpServer.Addr)
		err := httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Println(err)
	}

}
