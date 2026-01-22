package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/got-many-wheels/dwarf/server/internal/core"
	"github.com/got-many-wheels/dwarf/server/internal/platform/db"
	"github.com/got-many-wheels/dwarf/server/internal/platform/httpserver"
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
	urlsColl := client.Database("dwarf").Collection("urls")
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

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{code}", func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("code")
		if len(code) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("url code needed, got an empty string"))
			return
		}
		filter := bson.M{"code": bson.M{"$eq": code}}
		var doc core.URL
		err = urlsColl.FindOne(context.TODO(), filter).Decode(&doc)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				http.Error(w, fmt.Sprintf("could not find url with code %q", code), http.StatusNotFound)
				return
			} else {
				log.Printf("urlsColl.FindOne: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		http.Redirect(w, r, doc.Long, http.StatusSeeOther)
	})

	mux.HandleFunc("DELETE /url/{code}", func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("code")
		if len(code) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("url code needed, got an empty string"))
			return
		}
		filter := bson.M{"code": bson.M{"$eq": code}}
		res, err := urlsColl.DeleteOne(context.TODO(), filter)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				http.Error(w, fmt.Sprintf("could not find url with code %q", code), http.StatusNotFound)
				return
			} else {
				log.Printf("urlsColl.DeleteOne: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		if res.DeletedCount > 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		http.Error(w, "not found", http.StatusNotFound)
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
		result, err := urlsColl.InsertOne(context.TODO(), doc)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				http.Error(w, fmt.Sprintf("url with code %q already exist", doc.Code), http.StatusConflict)
				return
			}
			log.Printf("coll.InsertOne: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		buf := []byte(doc.String())
		buf = fmt.Appendf(buf, "with _id of %v", result.InsertedID)
		w.WriteHeader(http.StatusCreated)
		w.Write(buf)
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
