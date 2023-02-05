package storage

import (
	"context"

	"github.com/uptrace/bun"
)

type Hello struct {
	bun.BaseModel `bun:"table:hello"`

	ID   int64  `bun:"id,pk,autoincrement"`
	Data string `bun:"data"`
}

func (db DB) GetHello(ctx context.Context, hello *Hello) error {
	err := db.NewSelect().Model(hello).Where("id = ?", 1).Scan(ctx)
	return err
}

func (db DB) SetHello(ctx context.Context, hello *Hello) error {
	_, err := db.NewUpdate().Model(hello).WherePK().Exec(ctx)
	return err
}
