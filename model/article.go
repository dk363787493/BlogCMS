package model

import (
	"BlogCMS/db"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Article struct {
	ID               uint      `gorm:"column:id;primary_key;autoIncrement" json:"id"`
	CategoryLevelOne int       `gorm:"column:category_level_one" json:"category_level_one"`
	CategoryLevelTwo int       `gorm:"column:category_level_two" json:"category_level_two"`
	Title            string    `gorm:"column:title" json:"title"`
	CoverPath        string    `gorm:"column:cover_path" json:"cover_path"`
	Tags             string    `gorm:"column:tags" json:"tags"`
	Description      string    `gorm:"column:description" json:"description"`
	Content          string    `gorm:"column:content" json:"content"`
	CreatedAt        time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Article) TableName() string {
	return "articles"
}

func GetArticles(isShort uint8, categoryLevelOne uint8, categoryLevelTwo uint8, tags string, page int, limit int) ([]*Article, error) {
	var articles []*Article
	if page < 1 {
		return nil, errors.New("page must be greater than or equal 1")
	}
	params := make([]any, 0)
	whereSql := "1=1"
	if categoryLevelOne != 0 {
		params = append(params, categoryLevelOne)
		whereSql = fmt.Sprintf("%s and category_level_one =? ", whereSql)
	}
	if categoryLevelTwo != 0 {
		params = append(params, categoryLevelTwo)
		whereSql = fmt.Sprintf("%s and category_level_two =? ", whereSql)
	}
	if tags != "" {
		params = append(params, tags)
		whereSql = fmt.Sprintf("%s and tags=? ", whereSql)
	}
	from := (page - 1) * limit
	params = append(params, from)
	params = append(params, limit)
	whereSql = fmt.Sprintf("%s limit ?,?", whereSql)
	var result *gorm.DB
	if isShort == 1 {
		result = db.MysqlDB.Select("id,category_level_one,category_level_two,title,cover_path,tags,description").Where(whereSql, params...).Find(&articles)
	} else {
		result = db.MysqlDB.Where(whereSql, params...).Find(&articles)
	}
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			fmt.Println("data not found")
			return nil, nil
		}
		return nil, result.Error
	}
	return articles, nil
}

func GetTotalArticlesCnt(categoryLevelOne uint8, categoryLevelTwo uint8, tags string) (int64, error) {
	var totalCnt int64

	params := make([]any, 0)
	params = append(params, categoryLevelOne)
	params = append(params, categoryLevelTwo)
	whereSql := "category_level_one = ? and category_level_two=?"
	if tags != "" {
		whereSql = fmt.Sprintf("%s and tags=? ", whereSql)
		params = append(params, tags)
	}
	result := db.MysqlDB.Model(&Article{}).Where(whereSql, params...).Count(&totalCnt)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, result.Error
	}
	return totalCnt, nil
}

func GetArticlesById(id uint) (Article, error) {
	var article Article
	tx := db.MysqlDB.Where("id=?", id).First(&article)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return Article{}, nil
	}
	if tx.Error != nil {
		return Article{}, tx.Error
	}
	return article, nil
}

func InsertArticle(article *Article) (int64, error) {
	result := db.MysqlDB.Create(article)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
