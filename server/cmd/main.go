package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/got-many-wheels/dwarf/server/internal/core"
	"github.com/got-many-wheels/dwarf/server/internal/platform/db"
	"github.com/got-many-wheels/dwarf/server/internal/platform/httpserver"
	services "github.com/got-many-wheels/dwarf/server/internal/service"
	urlrepo "github.com/got-many-wheels/dwarf/server/internal/store/url"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

	// build the necessary indexes for the collection
	db := client.Database("dwarf")
	urlsColl := db.Collection("urls")
	models := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "code", Value: 1},
			},
			Options: options.Index().SetUnique(true).SetName("code_unique"),
		},
	}
	_, err = urlsColl.Indexes().CreateMany(context.TODO(), models)
	if err != nil {
		log.Printf("error while creating collection indexes: %v", err)
		return
	}

	services := buildServices(db)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{code}", func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("code")
		u, err := services.URL.Get(context.Background(), code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, u.Long, http.StatusSeeOther)
	})

	mux.HandleFunc("DELETE /url/{code}", func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("code")
		err := services.URL.Delete(context.Background(), code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("POST /url", func(w http.ResponseWriter, r *http.Request) {
		var doc core.URL
		err := json.NewDecoder(r.Body).Decode(&doc)
		if err != nil {
			log.Printf("could not decode request body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		doc.CreatedAt = time.Now().UTC()
		err = services.URL.InsertBatch(context.Background(), []core.URL{doc})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(doc.String()))
	})

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
