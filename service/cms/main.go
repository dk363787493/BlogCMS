package main

import (
	"BlogCMS/model"
	"BlogCMS/utils"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func main() {
	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20
	r.Use(CORSMiddleware())
	// 处理文件上传
	r.POST("/upload_img", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			fmt.Println("err:", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 读取文件内容
		fileData, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			fmt.Println("err:", err.Error())
			return
		}
		defer fileData.Close()
		bytes, err := ioutil.ReadAll(fileData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			fmt.Println("err:", err.Error())
			return
		}

		// 创建图片记录
		image := model.UploadImg{
			ImgUuid: utils.GenerateFileUUid(),
			Data:    bytes,
		}
		// 保存到数据库
		if err := model.InsertImage(&image); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			fmt.Println("err:", err.Error())
			return
		}
		fmt.Println("ImgUuid:", image.ImgUuid)
		// TODO 更换生成域名
		c.JSON(http.StatusOK, gin.H{"success": true, "location": fmt.Sprintf("http://localhost:8081/photo/%s", image.ImgUuid)})
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
	r.POST("/upload_article", func(c *gin.Context) {
		content := c.PostForm("content")
		article := new(model.Article)
		article.Content = content
		article.CoverPath = c.PostForm("cover-path")
		article.Title = c.PostForm("title")
		article.Tags = c.PostForm("tags")
		article.Description = c.PostForm("description")
		article.CategoryLevelOne, _ = strconv.Atoi(c.PostForm("category-level1"))
		article.CategoryLevelTwo, _ = strconv.Atoi(c.PostForm("category-level2"))
		insertArticle, err := model.InsertArticle(article)
		if err != nil {
			fmt.Println("err:", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("insertArticle:", insertArticle)
		c.JSON(http.StatusOK, gin.H{"success": true})
	})
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
		id, _ := strconv.Atoi(c.Param("id"))
		articles, err := model.GetArticlesById(uint(id))
		if err != nil {
			fmt.Println("err:", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "data": articles})
	})
	r.Run(":8080")
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
