package url

import (
	"context"
	"fmt"
	"log"

	"github.com/got-many-wheels/dwarf/server/internal/core"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repo struct {
	col     *mongo.Collection
	colname string
}

func New(db *mongo.Database) *Repo {
	r := &Repo{}
	col := db.Collection("urls")
	r.col = col
	return r
}

func (r *Repo) InsertBatch(ctx context.Context, items []core.URL) error {
	_, err := r.col.InsertMany(context.TODO(), items)
	log.Println(items)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) Get(ctx context.Context, code string) (core.URL, error) {
	var doc core.URL
	filter := bson.M{"code": bson.M{"$eq": code}}
	err := r.col.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("Could not find url with code %q", code)
		} else {
			return doc, err
		}
	}
	return doc, nil
}

func (r *Repo) Delete(ctx context.Context, code string) error {
	filter := bson.M{"code": bson.M{"$eq": code}}
	_, err := r.col.DeleteOne(context.TODO(), filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("could not find url with code %q", code)
		} else {
			return err
		}
	}
	return nil
}
