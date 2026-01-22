package url

import (
	"context"

	"github.com/got-many-wheels/dwarf/server/internal/core"
)

type URLStore interface {
	InsertBatch(ctx context.Context, items []core.URL) error
	Get(ctx context.Context, code string) (core.URL, error)
	Delete(ctx context.Context, code string) error
}

type URLService struct {
	store URLStore
}

func NewURLService(store URLStore) *URLService {
	return &URLService{store: store}
}

func (s *URLService) InsertBatch(ctx context.Context, items []core.URL) error {
	return s.store.InsertBatch(ctx, items)
}

func (s *URLService) Get(ctx context.Context, code string) (core.URL, error) {
	return s.store.Get(ctx, code)
}

func (s *URLService) Delete(ctx context.Context, code string) error {
	return s.store.Delete(ctx, code)
}
