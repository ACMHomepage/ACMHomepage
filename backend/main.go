package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

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
	registerAPI(r, db)
	return nil
}

func registerAPI(r *gin.Engine, db storage.DB) {
	// Get user statistical data.
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

		// Use crawler to update user data.
		crawler.UpdateOjUser(db, ctx, &ojUser)

		ctx.JSON(http.StatusOK, gin.H{
			"user": ojUser,
		})
	})

	// Get user submissions.
	r.GET("/api/v1/submissions/:oj/:handle", func(ctx *gin.Context) {
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

	// Search problems with rating range.
	r.GET("/api/v1/problems", func(ctx *gin.Context) {
		rating_from, err := strconv.ParseInt(ctx.DefaultQuery("rating_from", "800"), 10, 64)
		if err != nil {
			panic(err)
		}
		rating_to, err := strconv.ParseInt(ctx.DefaultQuery("rating_to", "3500"), 10, 64)
		if err != nil {
			panic(err)
		}
		problems := make([]storage.Problem, 0)
		db.GetProblemByRating(ctx, &problems, rating_from, rating_to)
		ctx.JSON(http.StatusOK, gin.H{
			"problems": problems,
		})
	})

	r.Run()
}
