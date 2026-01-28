package sequence

import (
	"context"

	coresequence "github.com/got-many-wheels/dwarf/server/internal/core/sequence"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Repo struct {
	col *mongo.Collection
}

func New(db *mongo.Database) *Repo {
	r := &Repo{}
	col := db.Collection("sequences")
	r.col = col
	return r
}

func (r *Repo) Next(ctx context.Context, name string) (int64, error) {
	opts := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)
	var seq coresequence.Sequence
	err := r.col.FindOneAndUpdate(
		ctx,
		bson.M{"_id": name},
		bson.M{"$inc": bson.M{"seq": 1}},
		opts,
	).Decode(&seq)
	return seq.Seq, err
}
