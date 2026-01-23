package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/got-many-wheels/dwarf/server/internal/platform/database"
	"github.com/got-many-wheels/dwarf/server/internal/platform/httpserver"
	services "github.com/got-many-wheels/dwarf/server/internal/service"
	urlrepo "github.com/got-many-wheels/dwarf/server/internal/store/url"
	"github.com/got-many-wheels/dwarf/server/internal/transport/mux"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func main() {
	addr := flag.String("addr", ":8080", "port of the http server")
	mongoURI := flag.String("mongouri", "mongodb://localhost:27017", "mongodb uri")
	flag.Parse()

	db, err := database.Init(*mongoURI, "dwarf")
	if err != nil {
		panic(err)
	}
	defer func() {
		err := db.Client.Disconnect(context.TODO())
		if err != nil {
			log.Println(err)
		}
	}()

	services := buildServices(db.DB)
	mux := mux.New(services)

	ctx, cancel := context.WithCancel(context.Background())
	s := httpserver.New(mux, *addr)

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

func buildServices(db *mongo.Database) services.Services {
	stores := defaultStoreFactory(db)
	return services.New(stores)
}

func defaultStoreFactory(db *mongo.Database) services.Stores {
	return services.Stores{
		URL: urlrepo.New(db),
	}
}
