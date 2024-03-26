package main

import (
	"BlogCMS/model"
	"BlogCMS/utils"
	"database/sql"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20
	// 处理文件上传
	r.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("photo")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 读取文件内容
		fileData, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer fileData.Close()
		bytes, err := ioutil.ReadAll(fileData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 创建图片记录
		image := model.UploadImg{
			ImgUuid: utils.GenerateFileUUid(),
			Data:    bytes,
		}

		// 保存到数据库
		if err := model.InsertImage(utils.GetMysqlDB(), &image); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "image_id": image.ImgUuid})

	})
	r.GET("/photo/:id", func(c *gin.Context) {
		// 从请求中获取图片ID
		photoID := c.Param("id")

		// 从数据库中检索图片

		photo, err := model.FindImageByName(utils.GetMysqlDB(), photoID)
		if err != nil {
			if err == sql.ErrNoRows {
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

	r.Run(":8080")
}
