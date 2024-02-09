package source

import (
	"context"
	"github.com/SlyMarbo/rss"
	"github.com/kimcodec/TgBot/internal/storage/model"
)

type RSSSource struct {
	URL        string
	SourceID   uint64
	SourceName string
}

func (r *RSSSource) ID() uint64 {
	return r.SourceID
}

func (r *RSSSource) Name() string {
	return r.SourceName
}

func NewRSSSourceFromModel(m model.Source) *RSSSource {
	return &RSSSource{
		URL:        m.FeedURL,
		SourceID:   m.ID,
		SourceName: m.Name,
	}
}

func (r *RSSSource) Fetch(ctx context.Context) ([]model.Item, error) {
	feed, err := r.loadFeed(ctx)
	if err != nil {
		return nil, err
	}

	var items []model.Item
	for _, v := range feed.Items {
		item := model.Item{
			Title:      v.Title,
			Categories: v.Categories,
			Link:       v.Link,
			Date:       v.Date,
			Summary:    v.Summary,
			SourceName: r.SourceName,
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *RSSSource) loadFeed(ctx context.Context) (*rss.Feed, error) {
	var (
		feedCh = make(chan *rss.Feed)
		errCh  = make(chan error)
	)

	go func() {
		feed, err := rss.Fetch(r.URL)
		if err != nil {
			errCh <- err
			return
		}
		feedCh <- feed
	}()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case err := <-errCh:
			return nil, err
		case feed := <-feedCh:
			return feed, nil
		}
	}
}
