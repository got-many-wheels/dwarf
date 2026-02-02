package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/got-many-wheels/dwarf/server/internal/platform/config"
	"github.com/got-many-wheels/dwarf/server/internal/platform/database"
	"github.com/got-many-wheels/dwarf/server/internal/platform/httpserver"
	services "github.com/got-many-wheels/dwarf/server/internal/service"
	urlrepo "github.com/got-many-wheels/dwarf/server/internal/store/url"
	"github.com/got-many-wheels/dwarf/server/internal/transport/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

// @title Dwarf URL API
// @version 1.0
// @description URL shortener service
// @BasePath /
func main() {
	cfg, err := config.Init()
	if err != nil {
		panic(err)
	}

	pool, err := database.Connect(cfg.DatabaseURI)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	services := buildServices(pool)
	mux := mux.New(services)

	ctx, cancel := context.WithCancel(context.Background())
	s := httpserver.New(mux, cfg.Port)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// spawn new goroutine to block terminate signal, and only call cancel() if
	// the system calls for either SIGTERM or SIGINT
	go func() {
		<-sig
		cancel() // will send ctx.Done() signal to the httpserver to close the server
	}()

	if err := s.Run(ctx); err != nil {
		log.Println(err)
	}
}

func buildServices(pool *pgxpool.Pool) services.Services {
	stores := defaultStoreFactory(pool)
	return services.New(stores)
}

func defaultStoreFactory(pool *pgxpool.Pool) services.Stores {
	return services.Stores{
		URL: urlrepo.New(pool),
	}
}
