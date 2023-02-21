package storage

import (
	"context"

	"github.com/uptrace/bun"
)

type Tag struct {
	bun.BaseModel `bun:"table:tag"`

	ID   int64  `bun:"id,pk,autoincrement"`
	Name string `bun:"name"`
}

func (db DB) CreateTag(ctx context.Context, tag *Tag) error {
	_, err := db.NewInsert().Model(tag).Exec(ctx)
	return err
}

func (db DB) ListTag(ctx context.Context, tagList *[]Tag) error {
	err := db.NewSelect().Model(tagList).Scan(ctx)
	return err
}

func (db DB) GetTag(ctx context.Context, tag *Tag, id int) error {
	err := db.NewSelect().Model(tag).Where("id = ?", id).Scan(ctx)
	return err
}

func (db DB) UpdateTag(ctx context.Context, tag *Tag, id int) error {
	_, err := db.NewUpdate().Model(tag).Where("id = ?", id).Exec(ctx)
	return err
}

func (db DB) DeleteTag(ctx context.Context, tag *Tag, id int) error {
	_, err := db.NewDelete().Model(tag).Where("id = ?", id).Exec(ctx)
	return err
}
