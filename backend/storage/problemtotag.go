package storage

import (
	"context"

	"github.com/uptrace/bun"
)

type ProblemToTag struct {
	bun.BaseModel `bun:"table:problem_to_tag"`

	ProblemID int64    `bun:"problem_id,pk"`
	Problem   *Problem `bun:"rel:belongs-to,join:problem_id=id"`
	TagName   string   `bun:"tag_name,pk"`
	Tag       *Tag     `bun:"rel:belongs-to,join:tag_name=name"`
}

func (db DB) GetProblemToTag(ctx context.Context, problemToTag *ProblemToTag, problemID int64, tagName string) error {
	err := db.NewSelect().Model(problemToTag).Where("problem_id = ?", problemID).Where("tag_name = ?", tagName).Scan(ctx)
	return err
}

func (db DB) CreateProblemToTag(ctx context.Context, problemToTag *ProblemToTag) error {
	_, err := db.NewInsert().Model(problemToTag).Exec(ctx)
	return err
}
