package migrations

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/skogkatt-org/ACMHomepage/backend/storage"
)

func init() {
	verdict := new(storage.Verdict)
	submission := new(storage.Submission)
	ojUser := new(storage.OjUser)
	problem := new(storage.Problem)
	tag := new(storage.Tag)
	problemToTag := new(storage.ProblemToTag)
	oj := new(storage.Oj)

	up := func(ctx context.Context, db *bun.DB) error {
		_, err := db.NewCreateTable().Model(verdict).Exec(ctx)
		if err != nil {
			return err
		}

		_, err = db.NewCreateTable().Model(submission).Exec(ctx)
		if err != nil {
			return err
		}

		_, err = db.NewCreateTable().Model(ojUser).Exec(ctx)
		if err != nil {
			return err
		}

		_, err = db.NewCreateTable().Model(problem).Exec(ctx)
		if err != nil {
			return err
		}

		_, err = db.NewCreateTable().Model(tag).Exec(ctx)
		if err != nil {
			return err
		}

		_, err = db.NewCreateTable().Model(problemToTag).Exec(ctx)
		if err != nil {
			return err
		}

		_, err = db.NewCreateTable().Model(oj).Exec(ctx)
		if err != nil {
			return err
		}

		_, err = db.NewInsert().Model(&[]storage.Oj{
			storage.Oj{
				Name: "codeforces",
			},
			storage.Oj{
				Name: "atcoder",
			},
		}).Exec(ctx)
		if err != nil {
			panic(err)
		}

		// Enum: FAILED, OK, PARTIAL, COMPILATION_ERROR, RUNTIME_ERROR, WRONG_ANSWER, PRESENTATION_ERROR, TIME_LIMIT_EXCEEDED, MEMORY_LIMIT_EXCEEDED, IDLENESS_LIMIT_EXCEEDED, SECURITY_VIOLATED, CRASHED, INPUT_PREPARATION_CRASHED, CHALLENGED, SKIPPED, TESTING, REJECTED.
		// TODO: verdict如何兼容不同平台

		return nil
	}
	down := func(ctx context.Context, db *bun.DB) error {
		db.NewDropTable().Model(verdict).Exec(ctx)
		db.NewDropTable().Model(submission).Exec(ctx)
		db.NewDropTable().Model(ojUser).Exec(ctx)
		db.NewDropTable().Model(problem).Exec(ctx)
		db.NewDropTable().Model(tag).Exec(ctx)
		db.NewDropTable().Model(problemToTag).Exec(ctx)
		db.NewDropTable().Model(oj).Exec(ctx)

		return nil
	}

	Migrations.MustRegister(up, down)
}
