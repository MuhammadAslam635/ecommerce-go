package models

import (
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID            int64         `gorm:"primary_key"`
	Name          string        `gorm:"not null"`
	Email         string        `gorm:"not null;uniqueIndex"`
	Phone         string        `gorm:"not null;uniqueIndex"`
	Password      string        `gorm:"not null"`
	Roles         string        `gorm:"not null"`
	Token         string        `gorm:"null"`
	RefreshToken  string        `gorm:"null"`
	AddressDetail []Address     `gorm:"foreignKey:UserID"`
	OrderStatus   []Order       `gorm:"foreignKey:UserID"`
	UserCart      []UserProduct `gorm:"foreignKey:UserID"`
	Reviews       []Review      `gorm:"foreignKey:UserID"`
	Orders        []Order       `gorm:"foreignKey:UserID"`
	OrderItems    []OrderItem   `gorm:"foreignKey:UserID"`
}

type OrderItem struct {
	gorm.Model
	ID        int64 `gorm:"primary_key"`
	UserID    int64 `gorm:"not null"`
	OrderID   int64 `gorm:"not null"`
	ProductID int64 `gorm:"not null"`
	Quantity  int
	Price     float64
}

type Category struct {
	gorm.Model
	ID       int64     `gorm:"primary_key"`
	Name     string    `gorm:"not null"`
	Slug     string    `gorm:"not null;unique"`
	Products []Product `gorm:"foreignKey:CategoryID"`
}

type Product struct {
	gorm.Model
	ID          int64       `gorm:"primary_key"`
	CategoryID  int64       `gorm:"not null"`
	Category    Category    `gorm:"foreignKey:CategoryID"`
	Name        string      `gorm:"not null"`
	Description string      `gorm:"not null"`
	Price       float64     `gorm:"not null"`
	Quantity    int         `gorm:"not null"`
	Image       string      `gorm:"null"`
	Rating      int         `gorm:"null"`
	OrderItems  []OrderItem `gorm:"foreignKey:ProductID"`
}

type UserProduct struct {
	gorm.Model
	ID          int64   `gorm:"primary_key"`
	UserID      int64   `gorm:"not null"`
	ProductID   int64   `gorm:"not null"`
	ProductName string  `gorm:"not null"`
	Price       float64 `gorm:"not null"`
	Quantity    int     `gorm:"not null"`
	Rating      int     `gorm:"null"`
	Image       string  `gorm:"null"`
}

type Address struct {
	gorm.Model
	ID      int64  `gorm:"primary_key"`
	UserID  int64  `gorm:"not null"`
	User    User   `gorm:"foreignKey:UserID"`
	Street  string `gorm:"not null"`
	City    string `gorm:"not null"`
	State   string `gorm:"not null"`
	Country string `gorm:"not null"`
}

type Order struct {
	gorm.Model
	ID            int64   `gorm:"primary_key"`
	UserID        int64   `gorm:"not null"`
	User          User    `gorm:"foreignKey:UserID"`
	AddressID     int64   `gorm:"not null"`
	Address       Address `gorm:"foreignKey:AddressID"`
	TotalPrice    float64 `gorm:"not null"`
	OrderStatus   string  `gorm:"not null"`
	PaymentMethod string  `gorm:"not null"`
}

type Payment struct {
	gorm.Model
	ID          int64   `gorm:"primary_key"`
	OrderID     int64   `gorm:"not null"`
	Order       Order   `gorm:"foreignKey:OrderID"`
	PaymentType string  `gorm:"not null"`
	Amount      float64 `gorm:"not null"`
}

type Review struct {
	gorm.Model
	ID        int64   `gorm:"primary_key"`
	UserID    int64   `gorm:"not null"`
	User      User    `gorm:"foreignKey:UserID"`
	ProductID int64   `gorm:"not null"`
	Product   Product `gorm:"foreignKey:ProductID"`
	Rating    int     `gorm:"not null"`
	Comment   string  `gorm:"not null"`
}
type SignedDetails struct {
	Email string
	Name  string
	jwt.StandardClaims
}
