package main

import (
	"BlogCMS/model"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func main() {
	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20
	r.Use(CORSMiddleware())

	r.GET("/article/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		articles, err := model.GetArticlesById(uint(id))
		if err != nil {
			fmt.Println("err:", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "data": articles})
	})

	r.GET("/article", func(c *gin.Context) {
		isShort, _ := strconv.Atoi(c.Query("is_short"))
		categoryLevelOne := c.Query("category_level_one")
		categoryOne, _ := strconv.Atoi(categoryLevelOne)
		categoryLevelTwo := c.Query("category_level_two")
		categoryTwo, _ := strconv.Atoi(categoryLevelTwo)
		tags := c.Query("tags")
		limitStr := c.Query("limit")
		limit := 10
		if limitStr != "" {
			limit, _ = strconv.Atoi(limitStr)
		}
		pageStr := c.Query("page")
		page := 1
		if pageStr != "" {
			page, _ = strconv.Atoi(pageStr)
		}
		articles, err := model.GetArticles(uint8(isShort), uint8(categoryOne), uint8(categoryTwo), tags, page, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		cnt, _ := model.GetTotalArticlesCnt(uint8(categoryOne), uint8(categoryTwo), tags)
		c.JSON(http.StatusOK, gin.H{"success": true, "data": articles, "totalCnt": cnt})
	})
	r.GET("/photo/:id", func(c *gin.Context) {
		// 从请求中获取图片ID
		photoID := c.Param("id")
		// 从数据库中检索图片
		photo, err := model.FindImageByName(photoID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// 没有找到图片
				c.String(http.StatusNotFound, "Photo not found")
				return
			}
			// 发生其他错误
			log.Printf("Error retrieving photo from database: %v", err)
			c.String(http.StatusInternalServerError, "Internal server error")
			return
		}

		// 设置正确的MIME类型
		c.Data(http.StatusOK, "image/jpeg", photo.Data)
	})
	r.Run(":8081")
}
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
