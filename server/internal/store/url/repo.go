package url

import (
	"context"

	coreurl "github.com/got-many-wheels/dwarf/server/internal/core/url"
	"github.com/got-many-wheels/dwarf/server/internal/store/database/sqlc"
	"github.com/got-many-wheels/dwarf/server/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	q    *sqlc.Queries
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repo {
	return &Repo{q: sqlc.New(pool), pool: pool}
}

func (r *Repo) InsertBatch(ctx context.Context, items []coreurl.URL) error {
	if len(items) == 0 {
		return nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	qtx := r.q.WithTx(tx)

	for idx, item := range items {
		row, err := qtx.CreateURL(ctx, sqlc.CreateURLParams{
			Long: item.Long,
			Code: item.Code,
		})
		if err != nil {
			return err
		}

		// generate url code if not provided
		if row.Code == "" {
			code := utils.DecimalToBase62(row.ID)
			qtx.UpdateURL(ctx, sqlc.UpdateURLParams{
				ID:   row.ID,
				Code: code,
				Long: row.Long,
			})
			items[idx].Code = code
		}
	}

	return tx.Commit(ctx)
}

func (r *Repo) Get(ctx context.Context, code string) (coreurl.URL, error) {
	row, err := r.q.GetURLByCode(ctx, code)
	if err != nil {
		return coreurl.URL{}, err
	}
	return mapUrl(row), nil
}

func (r *Repo) Delete(ctx context.Context, code string) error {
	return r.q.DeleteURLByCode(ctx, code)
}

func mapUrl(row sqlc.Url) coreurl.URL {
	return coreurl.URL{
		Id:        int(row.ID),
		Code:      row.Code,
		Long:      row.Long,
		CreatedAt: row.CreatedAt.Time,
	}
}
