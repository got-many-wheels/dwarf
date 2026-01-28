package coresequence

type Sequence struct {
	ID  string `bson:"_id"`
	Seq int64  `bson:"seq"`
}
