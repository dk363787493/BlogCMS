package model

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

type UploadImg struct {
	ID        uint      `gorm:"column:id;primary_key;autoIncrement"`
	ImgUuid   string    `gorm:"column:img_uuid"`
	Data      []byte    `gorm:"column:data"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// TableName overrides the table name used by UploadImg to `upload_img`
func (UploadImg) TableName() string {
	return "upload_imgs"
}

// FindImageByName 根据图片名称查询图片记录
func FindImageByName(db *gorm.DB, ImgUuid string) (*UploadImg, error) {
	var image UploadImg
	result := db.Where("img_uuid = ?", ImgUuid).First(&image)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // 没有找到记录不一定是一个错误
		}
		return nil, result.Error
	}
	return &image, nil
}

func InsertImage(db *gorm.DB, image *UploadImg) error {
	result := db.Create(&image)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
