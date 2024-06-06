package model

import (
	"BlogCMS/db"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Article struct {
	ID               uint      `gorm:"column:id;primary_key;autoIncrement"`
	CategoryLevelOne uint      `gorm:"column:category_level_one"`
	CategoryLevelTwo uint      `gorm:"column:category_level_two"`
	Title            string    `gorm:"column:title"`
	CoverPath        string    `gorm:"column:cover_path"`
	Tags             string    `gorm:"column:tags"`
	Content          []byte    `gorm:"column:content"`
	CreatedAt        time.Time `gorm:"column:created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at"`
}

func (Article) TableName() string {
	return "articles"
}

func GetArticlesByCategory(categoryLevelOne uint8, categoryLevelTwo uint8, tags string, page uint, limit uint) ([]*Article, error) {
	var articles []*Article
	if page < 1 {
		return nil, errors.New("page must be greater than or equal 1")
	}
	params := make([]any, 0)
	params = append(params, categoryLevelOne)
	params = append(params, categoryLevelTwo)
	whereSql := "category_level_one = ? and category_level_two=?"
	if tags != "" {
		whereSql = fmt.Sprintf("%s and tags=? ", whereSql)
		params = append(params, tags)
	}
	from := (page - 1) * limit
	params = append(params, from)
	params = append(params, limit)
	whereSql = fmt.Sprintf("%s limit ?,?", whereSql)
	result := db.MysqlDB.Where(whereSql, params...).Find(&articles)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return articles, nil
}

func InsertArticle(article *Article) (int64, error) {
	result := db.MysqlDB.Create(article)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
