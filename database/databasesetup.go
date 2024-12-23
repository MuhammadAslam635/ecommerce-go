package database

import (
	"fmt"
	"log"

	"githum.com/muhammadAslam/ecommerce/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DBSet() *gorm.DB {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "user=postgres password=azzan310 dbname=ecommerce port=5432 sslmode=disable TimeZone=Asia/Karachi",
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", err)
	}
	err = db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Product{},
		&models.UserProduct{},
		&models.Address{},
		&models.Order{},
		&models.OrderItem{},
		&models.Payment{},
		&models.Review{},
	)
	if err != nil {
		log.Fatal("failed to migrate models: " + err.Error())
	}
	fmt.Println("successfully migrated")
	return db

}

var Client = DBSet()

type ProductData struct {
	DB        *gorm.DB
	TableName string
}

func NewProductData(db *gorm.DB, tableName string) *ProductData {
	return &ProductData{
		DB:        db,
		TableName: tableName,
	}
}

func (pd *ProductData) GetById(id int64) (*models.Product, error) {
	var product models.Product
	if err := pd.DB.First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

// UserData retrieves users from the database
type UserData struct {
	DB        *gorm.DB
	TableName string
}

func NewUserData(db *gorm.DB, tableName string) *UserData {
	return &UserData{
		DB:        db,
		TableName: tableName,
	}
}
