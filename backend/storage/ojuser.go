package storage

import (
	"context"

	"github.com/uptrace/bun"
)

type OjUser struct {
	bun.BaseModel `bun:"table:oj_user"`

	ID            int64  `bun:"id,pk,autoincrement"`
	Handle        string `bun:"handle"`
	AcceptCount   int64  `bun:"accept_count"`
	MaxRating     int64  `bun:"max_rating"`
	CurrentRating int64  `bun:"current_rating"`
	Link          string `bun:"link"`

	OjName         string        `bun:"oj_name"`
	SubmissionList []*Submission `bun:"rel:has-many,join:id=oj_user_id"`
}

func (db DB) GetOjUser(ctx context.Context, ojUser *OjUser, ojName string, handle string) error {
	err := db.NewSelect().Model(ojUser).Where("oj_name = ?", ojName).Where("handle = ?", handle).Scan(ctx)
	return err
}

func (db DB) GetSubmissionList(ctx context.Context, ojUser *OjUser, ojName string, handle string) error {
	err := db.NewSelect().Model(ojUser).Where("oj_name = ?", ojName).Where("handle = ?", handle).Relation("SubmissionList").Scan(ctx)
	return err
}

func (db DB) CreateOjUser(ctx context.Context, ojUser *OjUser) error {
	_, err := db.NewInsert().Model(ojUser).Exec(ctx)
	return err
}

func (db DB) UpdateOjUser(ctx context.Context, ojUser *OjUser) error {
	_, err := db.NewUpdate().Model(ojUser).WherePK().Exec(ctx)
	return err
}
