package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/got-many-wheels/dwarf/server/internal/platform/db"
	"github.com/got-many-wheels/dwarf/server/internal/platform/httpserver"
)

func main() {
	addr := flag.String("addr", ":8080", "port of the http server")
	mongoURI := flag.String("mongouri", "mongodb://localhost:27017", "mongodb uri")
	flag.Parse()

	client, err := db.Connect(*mongoURI)
	if err != nil {
		log.Printf("error connecting to db: %v", err)
	}

	defer func() {
		err := client.Disconnect(context.TODO())
		if err != nil {
			log.Println(err)
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/url", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK!"))
	})

	ctx, cancel := context.WithCancel(context.Background())
	s := httpserver.New(mux, *addr)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// spawn new goroutine to block terminate signal. and only call cancel() if
	// the system calls for either SIGTERM or SIGINT
	go func() {
		<-sig
		cancel() // will send ctx.Done() signal to the httpserver to close the server
	}()

	if err := s.Run(ctx); err != nil {
		log.Println(err)
	}
}
