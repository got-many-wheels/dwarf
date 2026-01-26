package url

import (
	"context"

	"github.com/got-many-wheels/dwarf/server/internal/core"
	"github.com/got-many-wheels/dwarf/server/internal/service/sequence"
	"github.com/got-many-wheels/dwarf/server/utils"
)

type Store interface {
	InsertBatch(ctx context.Context, items []core.URL) error
	Get(ctx context.Context, code string) (core.URL, error)
	Delete(ctx context.Context, code string) error
}

type Service struct {
	store Store
	seq   sequence.Store
}

func New(store Store, seq sequence.Store) *Service {
	return &Service{store: store, seq: seq}
}

func (s *Service) InsertBatch(ctx context.Context, items []core.URL) error {
	for i := range items {
		seq, err := s.seq.Next(ctx, "url")
		if err != nil {
			return err
		}
		// generate code with base62 coming from the record count
		// of urls if no custom code provided
		if len(items[i].Code) == 0 {
			items[i].Code = utils.DecimalToBase62(seq)
		}
	}
	return s.store.InsertBatch(ctx, items)
}

func (s *Service) Get(ctx context.Context, code string) (core.URL, error) {
	return s.store.Get(ctx, code)
}

func (s *Service) Delete(ctx context.Context, code string) error {
	return s.store.Delete(ctx, code)
}
