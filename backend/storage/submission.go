package storage

import (
	"context"

	"github.com/uptrace/bun"
)

type Submission struct {
	bun.BaseModel `bun:"table:submission"`

	ID          int64   `bun:"id,pk,autoincrement"`
	OjUserID    int64   `bun:"oj_user_id"`
	ProblemID   int64   `bun:"problem_id"`
	VerdictName string  `bun:"verdict_name"`
	SubmitTime  float64 `bun:"submit_time"`
	Link        string  `bun:"link"`
}

func (db DB) GetSubmissionByLink(ctx context.Context, submission *Submission, link string) error {
	err := db.NewSelect().Model(submission).Where("link = ?", link).Scan(ctx)
	return err
}

func (db DB) CreateSubmission(ctx context.Context, submission *Submission) error {
	_, err := db.NewInsert().Model(submission).Exec(ctx)
	return err
}

func (db DB) GetAcceptCount(ctx context.Context, ojUserID int64) (int64, error) {
	count, err := db.NewSelect().Model((*Submission)(nil)).Where("oj_user_id = ?", ojUserID).Where("verdict_name = ?", "OK").DistinctOn("problem_id").Count(ctx)
	return int64(count), err
}
