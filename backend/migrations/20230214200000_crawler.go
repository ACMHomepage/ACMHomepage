package migrations

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/skogkatt-org/ACMHomepage/backend/storage"
)

func init() {
	ojUser := new(storage.OjUser)
	problem := new(storage.Problem)
	submission := new(storage.Submission)
	tag := new(storage.Tag)
	problemTag := new(storage.ProblemTag)
	oj := new(storage.Oj)
	verdict := new(storage.Verdict)

	up := func(ctx context.Context, db *bun.DB) error {
		_, err := db.NewCreateTable().Model(ojUser).Exec(ctx)
		if err != nil {
			return err
		}

		_, err = db.NewCreateTable().Model(problem).Exec(ctx)
		if err != nil {
			return err
		}

		_, err = db.NewCreateTable().Model(submission).Exec(ctx)
		if err != nil {
			return err
		}

		_, err = db.NewCreateTable().Model(tag).Exec(ctx)
		if err != nil {
			return err
		}

		_, err = db.NewCreateTable().Model(problemTag).Exec(ctx)
		if err != nil {
			return err
		}

		_, err = db.NewCreateTable().Model(oj).Exec(ctx)
		if err != nil {
			return err
		}

		_, err = db.NewCreateTable().Model(verdict).Exec(ctx)
		if err != nil {
			return err
		}

		_, err = db.NewInsert().Model(&[]storage.Oj{storage.Oj{
			ID:     1,
			OjName: "codeforces",
		}, storage.Oj{
			ID:     2,
			OjName: "atcoder",
		}}).Exec(ctx)
		if err != nil {
			panic(err)
		}

		// Enum: FAILED, OK, PARTIAL, COMPILATION_ERROR, RUNTIME_ERROR, WRONG_ANSWER, PRESENTATION_ERROR, TIME_LIMIT_EXCEEDED, MEMORY_LIMIT_EXCEEDED, IDLENESS_LIMIT_EXCEEDED, SECURITY_VIOLATED, CRASHED, INPUT_PREPARATION_CRASHED, CHALLENGED, SKIPPED, TESTING, REJECTED.
		_, err = db.NewInsert().Model(&[]storage.Verdict{storage.Verdict{
			ID:          1,
			VerdictName: "OK",
		}, storage.Verdict{
			ID:          2,
			VerdictName: "WRONG_ANSWER",
		}, storage.Verdict{
			ID:          3,
			VerdictName: "TIME_LIMIT_EXCEEDED",
		}, storage.Verdict{
			ID:          4,
			VerdictName: "MEMORY_LIMIT_EXCEEDED",
		}, storage.Verdict{
			ID:          5,
			VerdictName: "RUNTIME_ERROR",
		}, storage.Verdict{
			ID:          6,
			VerdictName: "COMPILATION_ERROR",
		}}).Exec(ctx)
		if err != nil {
			panic(err)
		}

		return nil
	}
	down := func(ctx context.Context, db *bun.DB) error {
		db.NewDropTable().Model(ojUser).Exec(ctx)
		db.NewDropTable().Model(problem).Exec(ctx)
		db.NewDropTable().Model(submission).Exec(ctx)
		db.NewDropTable().Model(tag).Exec(ctx)
		db.NewDropTable().Model(problemTag).Exec(ctx)
		db.NewDropTable().Model(oj).Exec(ctx)
		db.NewDropTable().Model(verdict).Exec(ctx)
		return nil
	}

	Migrations.MustRegister(up, down)
}
