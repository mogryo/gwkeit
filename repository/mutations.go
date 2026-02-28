package repository

import (
	"context"
	"database/sql"

	"github.com/gwkeit/dto"
	"github.com/gwkeit/gwkeitdb"
	"github.com/gwkeit/slicelib"
)

func (r *Repository) SaveSnippet(
	ctx context.Context,
	snippetInput *dto.Snippet,
) (int64, error) {
	r.checkVersion(ctx)
	tx, err := r.db.Begin()
	if err != nil {
		return -1, err
	}
	defer tx.Rollback()
	qtx := r.queries.WithTx(tx)

	snippetId, _ := qtx.InsertSnippet(
		ctx, gwkeitdb.InsertSnippetParams{
			Title:       snippetInput.Title,
			Body:        snippetInput.Body,
			Description: snippetInput.Description,
			Url:         snippetInput.UrlText,
			Language:    sql.NullString{String: snippetInput.Language, Valid: snippetInput.Language != ""},
		},
	)

	existingTags, _ := qtx.FindTagsByTag(ctx, snippetInput.Tags)
	existingTagNames := slicelib.Map(existingTags, func(tag gwkeitdb.Tag) string {
		return tag.Tag
	})
	existingTagIds := slicelib.Map(existingTags, func(tag gwkeitdb.Tag) int64 {
		return tag.ID
	})

	nonExistingTags := slicelib.Difference(snippetInput.Tags, existingTagNames)
	if len(nonExistingTags) > 0 {
		for _, tag := range nonExistingTags {
			id, _ := qtx.InsertTag(ctx, tag)
			existingTagIds = append(existingTagIds, id)
		}
	}

	urlIds := make([]int64, 0)
	if len(snippetInput.UrlList) > 0 {
		for _, url := range snippetInput.UrlList {
			id, _ := qtx.InsertUrl(ctx, gwkeitdb.InsertUrlParams{Url: url, SnippetID: snippetId})
			urlIds = append(urlIds, id)
		}
	}

	for _, tagId := range existingTagIds {
		_ = qtx.InsertSnippetTag(ctx, gwkeitdb.InsertSnippetTagParams{SnippetID: snippetId, TagID: tagId})
	}

	return snippetId, tx.Commit()
}

func (r *Repository) UpdateSnippet(
	ctx context.Context,
	snippetId int64,
	newSnippetData *dto.Snippet,
) error {
	r.checkVersion(ctx)
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := r.queries.WithTx(tx)

	err = qtx.UpdateSnippet(
		ctx,
		gwkeitdb.UpdateSnippetParams{
			ID:          snippetId,
			Title:       newSnippetData.Title,
			Body:        newSnippetData.Body,
			Description: newSnippetData.Description,
			Language:    sql.NullString{String: newSnippetData.Language, Valid: newSnippetData.Language != ""},
		},
	)
	if err != nil {
		return err
	}

	err = r.updateSnippetTags(ctx, qtx, snippetId, newSnippetData.Tags)
	if err != nil {
		return err
	}

	err = r.updateSnippetUrls(ctx, qtx, snippetId, newSnippetData.UrlList)
	if err != nil {
		return err
	}

	return tx.Commit()
}
func (r *Repository) updateSnippetTags(
	ctx context.Context,
	qtx *gwkeitdb.Queries,
	snippetId int64,
	newTagNames []string,
) error {
	r.checkVersion(ctx)
	existingSnippetTags, err := qtx.FindTagsBySnippetId(ctx, snippetId)
	if err != nil {
		return err
	}

	unusedTags := slicelib.DifferenceGetA(existingSnippetTags, newTagNames, func(tag gwkeitdb.Tag) string { return tag.Tag })
	err = qtx.DeleteSnippetTags(
		ctx,
		gwkeitdb.DeleteSnippetTagsParams{
			SnippetID: snippetId,
			TagIds:    slicelib.Map(unusedTags, func(tag gwkeitdb.Tag) int64 { return tag.ID }),
		},
	)
	if err != nil {
		return err
	}

	newSnippetTagNames := slicelib.Difference(newTagNames, slicelib.Map(existingSnippetTags, func(tag gwkeitdb.Tag) string { return tag.Tag }))
	existingTagNames, _ := qtx.FindTagsByTag(ctx, newSnippetTagNames)
	existingTagIds := slicelib.Map(existingTagNames, func(tag gwkeitdb.Tag) int64 { return tag.ID })

	tagNamesToAdd := slicelib.DifferenceGetB(newSnippetTagNames, existingTagNames, func(tag gwkeitdb.Tag) string { return tag.Tag })
	for _, tagName := range tagNamesToAdd {
		id, err := qtx.InsertTag(ctx, tagName)
		if err != nil {
			return err
		}

		existingTagIds = append(existingTagIds, id)
	}

	for _, tagId := range existingTagIds {
		_ = qtx.InsertSnippetTag(ctx, gwkeitdb.InsertSnippetTagParams{SnippetID: snippetId, TagID: tagId})
	}

	err = r.deleteOrphanTags(ctx, qtx, unusedTags)
	return err
}

func (r *Repository) updateSnippetUrls(
	ctx context.Context,
	qtx *gwkeitdb.Queries,
	snippetId int64,
	newUrls []string,
) error {
	r.checkVersion(ctx)
	existingSnippetUrls, err := qtx.FindUrlsBySnippetId(ctx, snippetId)

	unusedUrls := slicelib.DifferenceGetA(existingSnippetUrls, newUrls, func(url gwkeitdb.Url) string { return url.Url })
	err = qtx.DeleteSnippetUrls(ctx, slicelib.Map(unusedUrls, func(url gwkeitdb.Url) int64 { return url.ID }))
	if err != nil {
		return err
	}

	newSnippetUrls := slicelib.Difference(newUrls, slicelib.Map(existingSnippetUrls, func(url gwkeitdb.Url) string { return url.Url }))
	for _, url := range newSnippetUrls {
		_, err = qtx.InsertUrl(ctx, gwkeitdb.InsertUrlParams{Url: url, SnippetID: snippetId})
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) deleteOrphanTags(ctx context.Context, qtx *gwkeitdb.Queries, tags []gwkeitdb.Tag) error {
	r.checkVersion(ctx)
	for _, tag := range tags {
		exists, err := qtx.SnippetTagExists(ctx, tag.Tag)
		if err != nil {
			return err
		}
		if exists == 0 {
			err = qtx.DeleteTagByTag(ctx, tag.Tag)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
