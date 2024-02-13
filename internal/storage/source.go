package storage

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/kimcodec/TgBot/internal/storage/model"
	"time"
)

type SourcePostgresStorage struct {
	db *sqlx.DB
}

func NewSourcePostgresStorage(db *sqlx.DB) *SourcePostgresStorage {
	return &SourcePostgresStorage{
		db: db,
	}
}

func (s *SourcePostgresStorage) Sources(ctx context.Context) ([]model.Source, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var sources []dbSource
	if err := conn.SelectContext(ctx, &sources, "SELECT * FROM source"); err != nil {
		return nil, err
	}

	var sourcesModels []model.Source
	for _, v := range sources {
		source := model.Source{
			ID:        v.ID,
			Name:      v.Name,
			FeedURL:   v.FeedUrl,
			CreatedAt: v.CreatedAt,
		}
		sourcesModels = append(sourcesModels, source)
	}

	return sourcesModels, nil
}

func (s *SourcePostgresStorage) SourceByID(ctx context.Context, id uint64) (*model.Source, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var source dbSource
	if err := conn.GetContext(ctx, &source, "SELECT * FROM source WHERE id = $1", id); err != nil {
		return nil, err
	}

	return &model.Source{
		ID:        source.ID,
		Name:      source.Name,
		FeedURL:   source.FeedUrl,
		CreatedAt: source.CreatedAt,
	}, nil
}

func (s *SourcePostgresStorage) Add(ctx context.Context, source model.Source) (uint64, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	var id uint64

	row := conn.QueryRowxContext(
		ctx,
		"INSERT INTO source(name, feed_url, created_at) VALUES ($1, $2, $3) RETURNING id",
		source.Name,
		source.FeedURL,
		source.CreatedAt)
	if row.Err() != nil {
		return 0, err
	}

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *SourcePostgresStorage) Delete(ctx context.Context, id uint64) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.ExecContext(ctx, "DELETE FROM source WHERE id = $1", id); err != nil {
		return err
	}
	return nil
}

type dbSource struct {
	ID        uint64    `db:"id"`
	Name      string    `db:"name"`
	FeedUrl   string    `db:"feed_url"`
	CreatedAt time.Time `db:"created_at"`
}
