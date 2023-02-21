package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"

	"github.com/skogkatt-org/ACMHomepage/backend/migrations"
	"github.com/skogkatt-org/ACMHomepage/backend/storage"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		panic(fmt.Errorf("unexcepeted error: %w", err))
	}
}

type config struct {
	// PostgresURL is the URL to connect the Postgres server.
	// It should not be not empty.
	PostgresURL string
}

func getConfig() (*config, error) {
	config := &config{
		PostgresURL: os.Getenv("POSTGRES_URL"),
	}

	if config.PostgresURL == "" {
		return nil, errors.New("Cannot got PostgresURL, failed")
	}

	return config, nil
}

func migrateDB(ctx context.Context, db storage.DB) error {
	migrator := migrate.NewMigrator(db.DB, migrations.Migrations)
	err := migrator.Init(ctx)
	if err != nil {
		return err
	}
	err = migrator.Lock(ctx)
	if err != nil {
		return err
	}
	_, err = migrator.Migrate(ctx)
	if err != nil {
		return err
	}
	err = migrator.Unlock(ctx)
	if err != nil {
		return err
	}

	return nil
}

func run(ctx context.Context) error {
	// Get the config.
	config, err := getConfig()
	if err != nil {
		return err
	}

	// Connect and init the database.
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(config.PostgresURL)))
	db := storage.DB{DB: bun.NewDB(sqldb, pgdialect.New())}
	err = migrateDB(ctx, db)
	if err != nil {
		return err
	}

	// Start the HTTP server.
	r := gin.Default()
	registerAPI(r, db)
	return nil
}

func registerAPI(r *gin.Engine, db storage.DB) {
	r.GET("/api/v1/news", func(ctx *gin.Context) {
		var newsList []storage.News
		err := db.ListNews(ctx, &newsList)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"data": newsList,
		})
	})
	r.Run()
}
