package sequence

import "context"

type Store interface {
	// Increment the sequence `name` in the collection
	Next(ctx context.Context, name string) (int64, error)
}

type Service struct {
	store Store
}

func New(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) Next(ctx context.Context, name string) (int64, error) {
	return s.Next(ctx, name)
}
