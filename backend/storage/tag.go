package storage

import (
	"context"

	"github.com/uptrace/bun"
)

type Tag struct {
	bun.BaseModel `bun:"table:tag"`

	Name        string    `bun:"name,pk"`
	ProblemList []Problem `bun:"m2m:problem_to_tag,join:Tag=Problem"`
}

func (db DB) GetTag(ctx context.Context, tag *Tag, name string) error {
	err := db.NewSelect().Model(tag).Where("name = ?", name).Scan(ctx)
	return err
}

func (db DB) CreateTag(ctx context.Context, tag *Tag) error {
	_, err := db.NewInsert().Model(tag).Exec(ctx)
	return err
}
