package migrations

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/skogkatt-org/ACMHomepage/backend/storage"
)

func init() {
	news := new(storage.News)

	up := func(ctx context.Context, db *bun.DB) error {
		_, err := db.NewCreateTable().Model(news).Exec(ctx)
		if err != nil {
			return err
		}

		_, err = db.NewInsert().Model(&storage.News{
			ID:      1,
			Title:   "test1",
			Image:   "test-image1",
			Content: "test-content1",
		}).Exec(ctx)
		if err != nil {
			return err
		}

		_, err = db.NewInsert().Model(&storage.News{
			ID:      2,
			Title:   "test2",
			Image:   "test-image2",
			Content: "test-content2",
		}).Exec(ctx)
		if err != nil {
			return err
		}

		return nil
	}
	down := func(ctx context.Context, db *bun.DB) error {
		db.NewDropTable().Model(news).Exec(ctx)
		return nil
	}

	Migrations.MustRegister(up, down)
}
