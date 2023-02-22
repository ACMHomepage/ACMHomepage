package storage

import (
	"context"

	"github.com/uptrace/bun"
)

type Tag struct {
	bun.BaseModel `bun:"table:tag"`

	Name     string `bun:"name,pk"`
	NewsList []News `bun:"m2m:news_to_tag,join:Tag=News"`
}

func (db DB) CreateTag(ctx context.Context, tag *Tag) error {
	_, err := db.NewInsert().Model(tag).Exec(ctx)
	return err
}

func (db DB) GetTag(ctx context.Context, tag *Tag, name string) error {
	err := db.NewSelect().Model(tag).Where("name = ?", name).Relation("NewsList").Order("tag.name ASC").Scan(ctx)
	return err
}
