package migrations

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/skogkatt-org/ACMHomepage/backend/storage"
)

func init() {
    hello := new(storage.Hello)

    up := func(ctx context.Context, db *bun.DB) error {
        _, err := db.NewCreateTable().Model(hello).Exec(ctx)
        if err != nil {
            return err
        }

        _, err = db.NewInsert().Model(&storage.Hello{
            ID: 1,
            Data: "world",
        }).Exec(ctx)
        if err != nil {
            return err
        }

        return nil
    }
    down := func(ctx context.Context, db *bun.DB) error {
        db.NewDropTable().Model(hello).Exec(ctx)
        return nil
    }

    Migrations.MustRegister(up, down)
}
