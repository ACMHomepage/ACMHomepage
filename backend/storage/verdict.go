package storage

import (
	"context"

	"github.com/uptrace/bun"
)

type Verdict struct {
	bun.BaseModel `bun:"table:verdict"`

	Name string `bun:"name,pk"`
}

func (db DB) GetVerdict(ctx context.Context, verdict *Verdict, name string) error {
	err := db.NewSelect().Model(verdict).Where("name = ?", name).Scan(ctx)
	return err
}

func (db DB) CreateVerdict(ctx context.Context, verdict *Verdict) error {
	_, err := db.NewInsert().Model(verdict).Exec(ctx)
	return err
}
