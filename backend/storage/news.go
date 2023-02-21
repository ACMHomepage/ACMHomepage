package storage

import (
	"context"

	"github.com/uptrace/bun"
)

type News struct {
	bun.BaseModel `bun:"table:news"`

	ID      int64  `bun:"id,pk,autoincrement"`
	Title   string `bun:"title"`
	Image   string `bun:"image"`
	Content string `bun:"content"`
}

func (db DB) CreateNews(ctx context.Context, news *News) error {
	_, err := db.NewInsert().Model(news).Exec(ctx)
	return err
}

func (db DB) ListNews(ctx context.Context, newsList *[]News) error {
	err := db.NewSelect().Model(newsList).Scan(ctx)
	return err
}

func (db DB) GetNews(ctx context.Context, news *News, id int) error {
	err := db.NewSelect().Model(news).Where("id = ?", id).Scan(ctx)
	return err
}

func (db DB) UpdateNews(ctx context.Context, news *News, id int) error {
	_, err := db.NewUpdate().Model(news).Where("id = ?", id).Exec(ctx)
	return err
}

func (db DB) DeleteNews(ctx context.Context, news *News, id int) error {
	_, err := db.NewDelete().Model(news).Where("id = ?", id).Exec(ctx)
	return err
}
