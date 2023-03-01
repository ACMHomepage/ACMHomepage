package migrations

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/skogkatt-org/ACMHomepage/backend/storage"
)

func init() {
	news := new(storage.News)
	tag := new(storage.Tag)
	newsToTag := new(storage.NewsToTag)

	up := func(ctx context.Context, db *bun.DB) error {
		_, err := db.NewCreateTable().Model(news).Exec(ctx)
		if err != nil {
			return err
		}
		_, err = db.NewCreateTable().Model(tag).Exec(ctx)
		if err != nil {
			return err
		}
		_, err = db.NewCreateTable().Model(newsToTag).Exec(ctx)
		if err != nil {
			return nil
		}
		return nil
	}
	down := func(ctx context.Context, db *bun.DB) error {
		db.NewDropTable().Model(news).Exec(ctx)
		db.NewDropTable().Model(tag).Exec(ctx)
		db.NewDropTable().Model(newsToTag).Exec(ctx)
		return nil
	}

	Migrations.MustRegister(up, down)
}
