package storage

import (
	"context"

	"github.com/uptrace/bun"
)

type NewsToTag struct {
	bun.BaseModel `bun:"table:news_to_tag"`

	NewsID  int64  `bun:"news_id,pk"`
	News    *News  `bun:"rel:belongs-to,join:news_id=id"`
	TagName string `bun:"tag_name,pk"`
	Tag     *Tag   `bun:"rel:belongs-to,join:tag_name=name"`
}

func (db DB) CreateNewsToTag(ctx context.Context, newsToTag *NewsToTag) error {
	_, err := db.NewInsert().Model(newsToTag).Exec(ctx)
	return err
}

func (db DB) DeleteNewsToTag(ctx context.Context, news_to_tag *NewsToTag, newsID int, tagName string) error {
	_, err := db.NewDelete().Model(news_to_tag).Where("news_id = ?", newsID).Where("tag_name = ?", tagName).Returning("*").Exec(ctx)
	return err
}

func (db DB) DeleteNewsToTagByNews(ctx context.Context, newsToTagList *[]NewsToTag, newsID int) error {
	_, err := db.NewDelete().Model(newsToTagList).Where("news_id = ?", newsID).Returning("*").Exec(ctx)
	return err
}
