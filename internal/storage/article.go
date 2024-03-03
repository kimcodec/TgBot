package storage

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/kimcodec/TgBot/internal/storage/model"
	"time"
)

type ArticlePostgresStorage struct {
	db *sqlx.DB
}

func NewArticlePostgresStorage(db *sqlx.DB) *ArticlePostgresStorage {
	return &ArticlePostgresStorage{
		db: db,
	}
}

func (s *ArticlePostgresStorage) Store(ctx context.Context, article model.Article) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.ExecContext(
		ctx,
		"INSERT INTO articles(source_id, title, link, summary, published_at) VALUES ($1, $2, $3, $4, $5) "+
			"ON CONFLICT DO NOTHING ",
		article.SourceID, article.Title, article.Link, article.Summary, article.PublishedAt); err != nil {
		return err
	}
	return nil
}

func (s *ArticlePostgresStorage) AllNotPosted(ctx context.Context, limit int64) ([]model.Article, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var dbArticles []dbArticle
	if err := conn.SelectContext(
		ctx,
		&dbArticles,
		"SELECT id, source_id, title,link, summary, published_at, created_at FROM articles "+
			"WHERE posted_at IS NULL ORDER BY published_at DESC LIMIT $1", limit); err != nil {
		return nil, err
	}

	var articles []model.Article
	for _, v := range dbArticles {
		articles = append(articles,
			model.Article{
				SourceID:    v.SourceID,
				ID:          v.ID,
				Title:       v.Title,
				Link:        v.Link,
				Summary:     v.Summary,
				PublishedAt: v.PublishedAt,
				CreatedAt:   v.CreatedAt,
			})
	}
	return articles, nil
}

func (s *ArticlePostgresStorage) MarkPosted(ctx context.Context, id uint64) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.ExecContext(
		ctx,
		"UPDATE articles SET posted_at = $1::timestamp WHERE id = $2",
		time.Now().UTC().Format(time.RFC3339), id); err != nil {
		return err
	}
	return nil
}

type dbArticle struct {
	ID          uint64    `db:"id"`
	SourceID    uint64    `db:"source_id"`
	Title       string    `db:"title"`
	Link        string    `db:"link"`
	Summary     string    `db:"summary"`
	PublishedAt time.Time `db:"published_at"`
	CreatedAt   time.Time `db:"created_at"`
}
