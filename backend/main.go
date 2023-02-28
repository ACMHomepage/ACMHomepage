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

	// Register many to many model so bun can better recognize m2m relation.
	// This should be done before you use the model for the first time.
	db.RegisterModel((*storage.ProblemToTag)(nil))

	err = migrateDB(ctx, db)
	if err != nil {
		return err
	}

	// Start the HTTP server.
	r := gin.Default()

	// 获取用户统计数据
	r.GET("/api/v1/oj_user/:oj/:handle", func(ctx *gin.Context) {
		ojName, _ := ctx.Params.Get("oj")
		handle, _ := ctx.Params.Get("handle")

		var oj storage.Oj
		if err := db.GetOj(ctx, &oj, ojName); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		var ojUser storage.OjUser
		if err := db.GetOjUser(ctx, &ojUser, ojName, handle); err != nil {
			ojUser.OjName = ojName
			ojUser.Handle = handle
			if err := db.CreateOjUser(ctx, &ojUser); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": err.Error(),
				})
				return
			}
		}

		// 使用爬虫更新用户数据
		crawler.UpdateOjUser(db, ctx, &ojUser)

		ctx.JSON(http.StatusOK, gin.H{
			"user": ojUser,
		})
	})

	// 获取用户提交记录
	r.GET("/api/v1/submission/:oj/:handle", func(ctx *gin.Context) {
		ojName, _ := ctx.Params.Get("oj")
		handle, _ := ctx.Params.Get("handle")

		var oj storage.Oj
		if err := db.GetOj(ctx, &oj, ojName); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		var ojUser storage.OjUser
		if err := db.GetOjUser(ctx, &ojUser, ojName, handle); err != nil {
			ojUser.OjName = ojName
			ojUser.Handle = handle
			if err := db.CreateOjUser(ctx, &ojUser); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": err.Error(),
				})
				return
			}
		}

		crawler.UpdateOjUser(db, ctx, &ojUser)

		if err := db.GetSubmissionList(ctx, &ojUser, ojName, handle); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"submission": ojUser.SubmissionList,
		})
	})

	// TODO: 检索器
	// r.GET("/api/v1/problems", func(ctx *gin.Context) {
	// 	// api/v1/problems?rating=1800-2300&tags=dp,math
	// 	// NOTE: 每天更新一次 UpdateProblem()
	// })

	r.Run()

	return nil
}
