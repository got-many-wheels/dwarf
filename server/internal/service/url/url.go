package url

import (
	"context"

	coreurl "github.com/got-many-wheels/dwarf/server/internal/core/url"
)

type Store interface {
	InsertBatch(ctx context.Context, items []coreurl.URL) error
	Get(ctx context.Context, code string) (coreurl.URL, error)
	Delete(ctx context.Context, code string) error
}

type Service struct {
	store Store
}

func New(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) InsertBatch(ctx context.Context, items []coreurl.URL) error {
	return s.store.InsertBatch(ctx, items)
}

func (s *Service) Get(ctx context.Context, code string) (coreurl.URL, error) {
	return s.store.Get(ctx, code)
}

func (s *Service) Delete(ctx context.Context, code string) error {
	return s.store.Delete(ctx, code)
}
