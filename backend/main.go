package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
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
	r.GET("/api/v1/hello", func(ctx *gin.Context) {
		hello := new(storage.Hello)
		db.GetHello(ctx, hello)

		ctx.JSON(http.StatusOK, gin.H{
			"hello": hello.Data,
		})
	})
	r.PUT("/api/v1/hello", func(ctx *gin.Context) {
		// Get the body of request.
		buf, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		// Update the data in database.
		err = db.SetHello(ctx, &storage.Hello{
			ID:   1,
			Data: string(buf),
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		// Return response.
		ctx.JSON(http.StatusOK, gin.H{
			"hello": string(buf),
		})
	})

	// test for news
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
	// r.GET("/api/v1/news/:id", func(ctx *gin.Context) {
	// 	id, _ := ctx.Params.Get("id")
	// 	news := new(storage.News)
	// 	err := db.GetNews(ctx, news, id)
	// 	if err != nil {
	// 		ctx.JSON(http.StatusInternalServerError, gin.H{
	// 			"message": err.Error(),
	// 		})
	// 		return
	// 	}
	// 	ctx.JSON(http.StatusOK, gin.H{
	// 		"data": news,
	// 	})
	// })
	r.Run()

	return nil
}
