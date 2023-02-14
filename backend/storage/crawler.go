package storage

import (
	"context"

	"github.com/uptrace/bun"
)

type OjUser struct {
	bun.BaseModel `bun:"table:oj_user"`

	ID            int64  `bun:"id,pk,autoincrement"`
	OjID          int64  `bun:"oj_id"`
	Handle        string `bun:"handle"`
	AcceptCount   int64  `bun:"accept_count"`
	MaxRating     int64  `bun:"max_rating"`
	CurrentRating int64  `bun:"current_rating"`
	Link          string `bun:"link"`
}

type Problem struct {
	bun.BaseModel `bun:"table:problem"`

	ID          int64  `bun:"id,pk,autoincrement"`
	OjID        int64  `bun:"oj_id"`
	ProblemName string `bun:"problem_name"`
	Rating      int64  `bun:"rating"`
	Link        string `bun:"link"`
}

type Submission struct {
	bun.BaseModel `bun:"table:submission"`

	ID         int64  `bun:"id,pk,autoincrement"`
	OjUserID   int64  `bun:"oj_user_id"`
	ProblemID  int64  `bun:"problem_id"`
	VerdictID  int64  `bun:"verdict_id"`
	SubmitTime int64  `bun:"submit_time"`
	Link       string `bun:"link"`
}

type Tag struct {
	bun.BaseModel `bun:"table:tag"`

	ID      int64  `bun:"id,pk,autoincrement"`
	TagName string `bun:"tag_name"`
}

type ProblemTag struct {
	bun.BaseModel `bun:"table:problem_tag"`

	ID        int64 `bun:"id,pk,autoincrement"`
	ProblemID int64 `bun:"problem_id"`
	TagID     int64 `bun:"tag_id"`
}

type Oj struct {
	bun.BaseModel `bun:"table:oj"`

	ID     int64  `bun:"id,pk,autoincrement"`
	OjName string `bun:"oj_name"`
}

type Verdict struct {
	bun.BaseModel `bun:"table:verdict"`

	ID          int64  `bun:"id,pk,autoincrement"`
	VerdictName string `bun:"verdict_name"`
}

func (db DB) GetOjByName(ctx context.Context, oj *Oj) error {
	err := db.NewSelect().Model(oj).Where("oj_name = ?", oj.OjName).Scan(ctx)
	return err
}

func (db DB) GetOjByID(ctx context.Context, oj *Oj) error {
	err := db.NewSelect().Model(oj).Where("id = ?", oj.ID).Scan(ctx)
	return err
}

func (db DB) GetOjUser(ctx context.Context, ojUser *OjUser) error {
	err := db.NewSelect().Model(ojUser).Where("oj_id = ? handle = ?", ojUser.OjID, ojUser.Handle).Scan(ctx)
	return err
}

func (db DB) SetOjUser(ctx context.Context, ojUser *OjUser) error {
	_, err := db.NewUpdate().Model(ojUser).WherePK().Exec(ctx)
	return err
}

func (db DB) GetVerdictID(ctx context.Context, verdict *Verdict) error {
	err := db.NewSelect().Model(verdict).Where("verdict_name = ?", verdict.VerdictName).Scan(ctx)
	return err
}

func (db DB) GetAcceptCount(ctx context.Context, ojUser *OjUser) error {
	acceptVerdictID := 1
	count, err := db.NewSelect().Model((*Submission)(nil)).Where("oj_user_id = ? verdict_id = ?", ojUser.ID, acceptVerdictID).Count(ctx)
	ojUser.AcceptCount = int64(count)
	return err
}

func (db DB) GetProblemByLink(ctx context.Context, problem *Problem) (bool, error) {
	exists, err := db.NewSelect().Model(problem).Where("link = ?", problem.Link).Exists(ctx)
	return exists, err
}
