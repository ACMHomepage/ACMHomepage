package storage

import (
	"context"

	"github.com/uptrace/bun"
)

type Problem struct {
	bun.BaseModel `bun:"table:problem"`

	ID     int64  `bun:"id,pk,autoincrement"`
	Name   string `bun:"name"`
	Rating int64  `bun:"rating"`
	Link   string `bun:"link"`

	OjName  string `bun:"oj_name"`
	TagList []Tag  `bun:"m2m:problem_to_tag,join:Problem=Tag"`
}

func (db DB) GetProblemByLink(ctx context.Context, problem *Problem, link string) error {
	err := db.NewSelect().Model(problem).Where("link = ?", link).Scan(ctx)
	return err
}

func (db DB) CreateProblem(ctx context.Context, problem *Problem) error {
	_, err := db.NewInsert().Model(problem).Exec(ctx)
	return err
}

func (db DB) GetProblemByRating(ctx context.Context, problems *[]Problem, rating_from int64, rating_to int64) error {
	err := db.NewSelect().Model(problems).Where("rating >= ?", rating_from).Where("rating <= ?", rating_to).Scan(ctx)
	return err
}
