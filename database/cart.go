package database

import (
	"context"
	"errors"
	"log"

	"githum.com/muhammadAslam/ecommerce/models"
	"gorm.io/gorm"
)

var (
	ErrCanNotFindProduct         = errors.New("can't find product")
	ErrCantDecodeProduct         = errors.New("can't decode product string")
	ErrUserIdIsNotValid          = errors.New("user id is not valid")
	ErrCantUpdateUser            = errors.New("can't update user")
	ErrCantRemoveItemCart        = errors.New("can't remove item cart")
	ErrCantGetCartItems          = errors.New("can't get cart items")
	ErrCantCheckoutCart          = errors.New("can't checkout cart")
	ErrCantUpdateProductQuantity = errors.New("can't update product quantity")
	ErrCantDeleteProduct         = errors.New("can't delete product")
	ErrCantGetInstantBuy         = errors.New("can't get instant buy")
	ErrCantBuyCartItem           = errors.New("can't buy cart item")
	ErrQuantityMustBePositive    = errors.New("quantity must be positive")
	ErrProductIdIsNotValid       = errors.New("product id is not valid")
	ErrCantFindProductInCart     = errors.New("can't find product in cart")
	ErrCantFindUserAddress       = errors.New("can't find user address")
)

func AddProductToCart(ctx context.Context, db *gorm.DB, productId int64, userId int64) error {
	// Validate userId
	if userId <= 0 {
		return ErrUserIdIsNotValid
	}

	// Fetch the product from the database to ensure it exists
	var product models.Product
	if err := db.WithContext(ctx).First(&product, "id = ?", productId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrCanNotFindProduct
		}
		return err
	}

	// Fetch the user by userId
	var user models.User
	if err := db.WithContext(ctx).First(&user, "id = ?", userId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrUserIdIsNotValid
		}
		return err
	}

	// Check if the product is already in the user's cart
	var existingCart models.UserProduct
	if err := db.WithContext(ctx).First(&existingCart, "user_id = ? AND product_id = ?", userId, productId).Error; err == nil {
		// If the product exists, update the quantity
		existingCart.Quantity += 1
		if err := db.WithContext(ctx).Save(&existingCart).Error; err != nil {
			log.Println("Failed to update product quantity in cart:", err)
			return err
		}
	} else if err == gorm.ErrRecordNotFound {
		// If the product does not exist in the cart, add it as a new entry
		newCart := models.UserProduct{
			UserID:      userId,
			ProductID:   productId,
			ProductName: product.Name,
			Price:       product.Price,
			Quantity:    1,
			Rating:      product.Rating,
			Image:       product.Image,
		}

		if err := db.WithContext(ctx).Create(&newCart).Error; err != nil {
			log.Println("Failed to add product to cart:", err)
			return err
		}
	} else {

		return err
	}

	return nil
}

func RemoveProductFromCart(ctx context.Context, db *gorm.DB, productId int64, userId int64) error {
	// Validate userId
	if userId <= 0 {
		return ErrUserIdIsNotValid
	}
	var userProduct models.UserProduct
	if err := db.WithContext(ctx).First(&userProduct, "user_id = ? AND product_id = ?", userId, productId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrCantRemoveItemCart
		}
		return err
	}
	if err := db.WithContext(ctx).Delete(&userProduct).Error; err != nil {
		log.Println("Failed to remove product from cart:", err)
		return err
	}
	return nil
}

func GetCartItems(ctx context.Context, db *gorm.DB, userId int64) ([]models.UserProduct, error) {
	// Validate userId
	if userId <= 0 {
		return nil, ErrUserIdIsNotValid
	}
	var userProducts []models.UserProduct
	if err := db.WithContext(ctx).Find(&userProducts, "user_id = ?", userId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrCantGetCartItems
		}
		return nil, err
	}
	return userProducts, nil
}

func CheckoutCart(ctx context.Context, db *gorm.DB, userId int64) error {
	// Validate userId
	if userId <= 0 {
		return ErrUserIdIsNotValid
	}
	var userProducts []models.UserProduct
	if err := db.WithContext(ctx).Find(&userProducts, "user_id = ?", userId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrCantCheckoutCart
		}
		return err

	}
	// calculate the total amount
	var totalAmount float64
	for _, product := range userProducts {
		totalAmount += product.Price * float64(product.Quantity)
	}
	// get latest user address id from Address table
	var userAddress models.Address
	if err := db.WithContext(ctx).Last(&userAddress, "user_id =?", userId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrCantCheckoutCart
		}
		return err
	}
	// make order now
	var order models.Order
	for _, product := range userProducts {
		order = models.Order{
			UserID:        userId,
			AddressID:     userAddress.ID,
			TotalPrice:    totalAmount,
			OrderStatus:   "ordered",
			PaymentMethod: "cod",
		}
		if err := db.WithContext(ctx).Create(&order).Error; err != nil {
			log.Println("Failed to make order:", err)
			return err
		}
		orderItem := models.OrderItem{
			OrderID:   order.ID,
			ProductID: product.ProductID,
			Quantity:  product.Quantity,
			Price:     product.Price,
		}
		if err := db.WithContext(ctx).Create(&orderItem).Error; err != nil {
			log.Println("Failed to make order item:", err)
			return err
		}
	}
	payment := models.Payment{
		OrderID:     order.ID,
		Amount:      totalAmount,
		PaymentType: "cod",
	}
	if err := db.WithContext(ctx).Create(&payment).Error; err != nil {
		log.Println("Failed to add payment record:", err)
		return err
	}
	// remove cartitems from userProducts list
	for _, product := range userProducts {
		db.WithContext(ctx).Delete(&models.UserProduct{}, "user_id =? AND product_id =?", userId, product.ProductID)
	}

	return nil
}

func UpdateProductQuantity(ctx context.Context, db *gorm.DB, userId int64, productId int64, qty int) error {
	if userId <= 0 {
		return ErrUserIdIsNotValid
	}
	if qty <= 0 {
		return ErrQuantityMustBePositive
	}
	if productId <= 0 {
		return ErrProductIdIsNotValid
	}
	var userProduct models.UserProduct
	if err := db.WithContext(ctx).First(&userProduct, "user_id = ? AND product_id = ?", userId, productId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrCantUpdateProductQuantity
		}
		return err
	}
	userProduct.Quantity = qty
	if err := db.WithContext(ctx).Save(&userProduct).Error; err != nil {
		log.Println("Failed to update product quantity:", err)
		return err
	}
	return nil

}

func DeleteProductFromCart(ctx context.Context, db *gorm.DB, userId int64) error {
	if userId <= 0 {
		return ErrUserIdIsNotValid
	}
	var userProducts []models.UserProduct
	if err := db.WithContext(ctx).Find(&userProducts, "user_id =?", userId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrCantDeleteProduct
		}
		return err
	}
	for _, product := range userProducts {
		db.WithContext(ctx).Delete(&models.UserProduct{}, "user_id =? AND product_id =?", userId, product.ProductID)
	}
	return nil
}
func GetInstantBuyProduct(ctx context.Context, db *gorm.DB, productId int64, uerId int64) error {
	if uerId <= 0 {
		return ErrUserIdIsNotValid
	}
	if productId <= 0 {
		return ErrProductIdIsNotValid
	}
	var product models.Product
	if err := db.WithContext(ctx).First(&product, "id = ?", productId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrCanNotFindProduct
		}
		return err
	}
	var userProduct models.UserProduct
	if err := db.WithContext(ctx).First(&userProduct, "user_id = ? AND product_id = ?", uerId, productId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrCantFindProductInCart
		}
		return err
	}
	var userAddress models.Address
	if err := db.WithContext(ctx).Last(&userAddress, "user_id =?", uerId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrCantFindUserAddress
		}
		return err
	}
	order := models.Order{
		UserID:        uerId,
		AddressID:     userAddress.ID,
		TotalPrice:    product.Price * float64(userProduct.Quantity),
		OrderStatus:   "ordered",
		PaymentMethod: "cod",
	}
	if err := db.WithContext(ctx).Create(&order).Error; err != nil {
		log.Println("Failed to make order", err)
		return err
	}
	orderItem := models.OrderItem{
		OrderID:   order.ID,
		ProductID: product.ID,
		Quantity:  userProduct.Quantity,
		Price:     product.Price,
	}
	if err := db.WithContext(ctx).Create(&orderItem).Error; err != nil {
		log.Println("Failed to make order item", err)
		return err

	}
	payment := models.Payment{
		OrderID:     order.ID,
		Amount:      product.Price * float64(userProduct.Quantity),
		PaymentType: "cod",
	}
	if err := db.WithContext(ctx).Create(&payment).Error; err != nil {
		log.Println("Failed to add payment record", err)
		return err

	}
	if err := db.WithContext(ctx).Delete(&userProduct).Error; err != nil {
		log.Println("Failed to remove product from cart", err)
		return err
	}
	return nil
}
