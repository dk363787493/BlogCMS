package dto

import (
	"BlogCMS/model"
	"errors"
	"time"
)

type ArticleDTO struct {
	ID               uint
	CategoryLevelOne int
	CategoryLevelTwo int
	Title            string
	CoverPath        string
	Tags             string
	Content          string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type BatchArticleDTO struct {
	ArticleDTOs []*ArticleDTO
	Page        uint
	Size        uint
	Total       uint
}

func (dto *ArticleDTO) Transfer(article *model.Article) error {
	if dto == nil {
		return errors.New("article is nil")
	}
	dto.ID = article.ID
	dto.CategoryLevelOne = article.CategoryLevelOne
	dto.CategoryLevelTwo = article.CategoryLevelTwo
	dto.Title = article.Title
	dto.CoverPath = article.CoverPath
	dto.Tags = article.Tags
	dto.CreatedAt = time.Now()
	dto.UpdatedAt = time.Now()
	dto.Content = string(article.Content)
	return nil
}
