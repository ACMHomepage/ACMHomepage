package storage

import (
	"context"

	"github.com/uptrace/bun"
)

type Oj struct {
	bun.BaseModel `bun:"table:oj"`

	Name        string     `bun:"name,pk"`
	OjUserList  []*OjUser  `bun:"rel:has-many,join:name=oj_name"`
	ProblemList []*Problem `bun:"rel:has-many,join:name=oj_name"`
}

func (db DB) GetOj(ctx context.Context, oj *Oj, name string) error {
	err := db.NewSelect().Model(oj).Where("name = ?", name).Scan(ctx)
	return err
}
