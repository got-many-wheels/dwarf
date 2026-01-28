package url

import (
	"context"
	"fmt"

	coreurl "github.com/got-many-wheels/dwarf/server/internal/core/url"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repo struct {
	col *mongo.Collection
}

func New(db *mongo.Database) *Repo {
	r := &Repo{}
	col := db.Collection("urls")
	r.col = col
	return r
}

func (r *Repo) InsertBatch(ctx context.Context, items []coreurl.URL) error {
	_, err := r.col.InsertMany(context.TODO(), items)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) Get(ctx context.Context, code string) (coreurl.URL, error) {
	var doc coreurl.URL
	filter := bson.M{"code": bson.M{"$eq": code}}
	err := r.col.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		// TODO: should use error factory that holds map of errors from the store
		// instead of explicitly state what's the error is
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("could not find url with code %q", code)
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
