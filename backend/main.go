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

	// Register many to many model so bun can better recognize m2m relation.
	// This should be done before you use the model for the first time.
	db.RegisterModel((*storage.NewsToTag)(nil))

	// Start the HTTP server.
	r := gin.Default()
	registerAPI(r, db)
	return nil
}

func registerAPI(r *gin.Engine, db storage.DB) {
	// CreateNews
	r.POST("/api/v1/news", func(ctx *gin.Context) {
		var news storage.News
		if err := ctx.ShouldBind(&news); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		if err := db.CreateNews(ctx, &news); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"data": news,
		})
	})
	r.GET("/api/v1/news", func(ctx *gin.Context) {
		tagName := ctx.Query("tag_name")
		if tagName == "" {
			// ListNews
			var newsList []storage.News
			if err := db.ListNews(ctx, &newsList); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": err.Error(),
				})
				return
			}
			ctx.JSON(http.StatusOK, gin.H{
				"data": newsList,
			})
		} else {
			// ListNewsByTag
			var tag storage.Tag
			if err := db.GetTag(ctx, &tag, tagName); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": err.Error(),
				})
				return
			}
			ctx.JSON(http.StatusOK, gin.H{
				"data": tag,
			})
		}
	})
	// GetNews
	r.GET("/api/v1/news/:id", func(ctx *gin.Context) {
		tmp, _ := ctx.Params.Get("id")
		id, _ := strconv.Atoi(tmp)
		news := new(storage.News)
		if err := db.GetNews(ctx, news, id); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"data": news,
		})
	})
	// UpdateNews
	r.PUT("/api/v1/news/:id", func(ctx *gin.Context) {
		tmp, _ := ctx.Params.Get("id")
		id, _ := strconv.Atoi(tmp)
		var news storage.News
		if err := ctx.ShouldBind(&news); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		if err := db.UpdateNews(ctx, &news, id); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"data": news,
		})
	})
	// DeleteNews
	r.DELETE("/api/v1/news/:id", func(ctx *gin.Context) {
		tmp, _ := ctx.Params.Get("id")
		id, _ := strconv.Atoi(tmp)
		news := new(storage.News)
		if err := db.DeleteNews(ctx, news, id); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		// delete news_to_tag with such news
		if err := db.DeleteNewsToTagByNews(ctx, &[]storage.NewsToTag{}, int(news.ID)); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"data": news,
		})
	})
	// AddTag
	r.POST("/api/v1/news/:id/tags", func(ctx *gin.Context) {
		var tag storage.Tag
		if err := ctx.ShouldBind(&tag); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		// search tag with tag_name in table tag
		if err := db.GetTag(ctx, &storage.Tag{}, tag.Name); err != nil {
			// cannot find such tag, insert it
			if err := db.CreateTag(ctx, &tag); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": err.Error(),
				})
				return
			}
		}
		// create news_to_tag in table news_to_tag
		tmp, _ := ctx.Params.Get("id")
		newsID, _ := strconv.Atoi(tmp)
		newsToTag := storage.NewsToTag{
			NewsID:  int64(newsID),
			TagName: tag.Name,
		}
		if err := db.CreateNewsToTag(ctx, &newsToTag); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"data": newsToTag,
		})
	})
	// DeleteTag
	r.DELETE("/api/v1/news/:id/tags/:name", func(ctx *gin.Context) {
		tmp, _ := ctx.Params.Get("id")
		newsID, _ := strconv.Atoi(tmp)
		tagName, _ := ctx.Params.Get("name")
		var news_to_tag storage.NewsToTag
		if err := db.DeleteNewsToTag(ctx, &news_to_tag, newsID, tagName); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"data": news_to_tag,
		})
	})
	r.Run()
}
