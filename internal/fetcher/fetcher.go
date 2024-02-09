package fetcher

import (
	"context"
	"github.com/kimcodec/TgBot/internal/source"
	"github.com/kimcodec/TgBot/internal/storage/model"
	"log"
	"strings"
	"sync"
	"time"
)

type ArticleStorage interface {
	Store(ctx context.Context, article model.Article) error
}

type SourceProvider interface {
	Sources(ctx context.Context) ([]model.Source, error)
}

type Source interface {
	ID() uint64
	Name() string
	Fetch(ctx context.Context) ([]model.Item, error)
}

type Fetcher struct {
	articles ArticleStorage
	sources  SourceProvider

	fetchInterval    time.Duration
	filteredKeywords []string // Подразумевается, что фильтр отсеивает по неинтересу
}

func NewFetcher(articleStorage ArticleStorage,
	sourceProvider SourceProvider,
	fetchInterval time.Duration,
	filteredKeywords []string) *Fetcher {
	return &Fetcher{
		articles:         articleStorage,
		sources:          sourceProvider,
		fetchInterval:    fetchInterval,
		filteredKeywords: filteredKeywords,
	}
}

func (f *Fetcher) Start(ctx context.Context) error {
	ticker := time.NewTicker(f.fetchInterval)

	if err := f.fetch(ctx); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := f.fetch(ctx); err != nil {
				return err
			}
		}
	}
}

func (f *Fetcher) fetch(ctx context.Context) error {
	sources, err := f.sources.Sources(ctx)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, src := range sources {
		wg.Add(1)

		rssSource := source.NewRSSSourceFromModel(src)

		go func(source Source) {
			defer wg.Done()

			items, err := source.Fetch(ctx)
			if err != nil {
				log.Printf("[ERROR] Fetching items from source %s : %v", source.Name(), source.ID())
				return
			}

			if err := f.processItem(ctx, source, items); err != nil {
				log.Printf("[ERROR] Processing items from source %s : %v", source.Name(), source.ID())
				return
			}
		}(rssSource)

		wg.Wait()
	}
	return nil
}

func (f *Fetcher) processItem(ctx context.Context, source Source, items []model.Item) error {
	for _, item := range items {
		item.Date = item.Date.UTC()

		if f.itemShouldBeSkipped(item) {
			continue
		}

		if err := f.articles.Store(ctx, model.Article{
			SourceID:    source.ID(),
			Title:       item.Title,
			Link:        item.Link,
			Summary:     item.Summary,
			PublishedAt: item.Date,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (f *Fetcher) itemShouldBeSkipped(item model.Item) bool {
	categoriesSet := make(map[string]struct{})
	for _, v := range item.Categories {
		categoriesSet[v] = struct{}{}
	}

	for _, keyword := range f.filteredKeywords {
		isTitleContainsKeyword := strings.Contains(strings.ToLower(item.Title), keyword)
		_, isCategoriesContainsKeyword := categoriesSet[keyword]
		if isCategoriesContainsKeyword || isTitleContainsKeyword {
			return true
		}
	}
	return false
}
