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

	"github.com/skogkatt-org/ACMHomepage/backend/crawler"
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

	r.GET("/api/v1/oj_user", func(ctx *gin.Context) {
		// example: /api/v1/oj_user?oj=codeforces&handle=jiangly
		ojName := ctx.DefaultQuery("oj", "codeforces")
		handle := ctx.DefaultQuery("handle", "jiangly")

		oj := new(storage.Oj)
		oj.OjName = ojName
		db.GetOj(ctx, oj)

		ojUser := new(storage.OjUser)
		ojUser.OjID = oj.ID
		ojUser.Handle = handle
		crawler.UpdateOjUser(db, ctx, ojUser)
		ctx.JSON(http.StatusOK, gin.H{
			"handle":         ojUser.Handle,
			"current_rating": ojUser.CurrentRating,
			"max_rating":     ojUser.MaxRating,
			"accept_count":   ojUser.AcceptCount,
		})
	})

	// TODO
	// 检索器
	// r.GET("/api/v1/problems", func(ctx *gin.Context) {
	// 	// /api/v1/problems?rating=1800-2300&tags=dp,math
	// 	// rating := range
	// 	// tags := slice
	// 	// NOTE: 每天更新一次 UpdateProblem()
	// 	// db.GetProblem() // storage/storage.go
	// 	// return json
	// })

	// 提交记录
	// r.GET("/api/v1/submission", func(ctx *gin.Context) {
	// 	// /api/v1/submission?handle=flowrays&submit_time=20230206000000-20230207000000
	// 	// handle := xxx
	// 	// submit_time := range
	// 	// UpdateOjUser(oj, handle, db)
	// 	// db.GetOjuser()
	// 	// return json
	// })

	r.Run()

	return nil
}
