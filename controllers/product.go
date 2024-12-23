package controllers

import (
	"context"
	"errors"
	"log"

	"githum.com/muhammadAslam/ecommerce/models"
	"gorm.io/gorm"
)

func GetById(ctx context.Context, db *gorm.DB, id int) error {
	if id == 0 {
		log.Printf("product id %d not found", id)
		return errors.New("product not found")
	}
	var product models.Product
	err := db.WithContext(ctx).First(&product, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("product id %d not found", id)
			return errors.New("product not found")
		}
		return err
	}
	return nil

}
