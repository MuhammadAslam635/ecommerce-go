package tokens

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"githum.com/muhammadAslam/ecommerce/models"
	"gorm.io/gorm"
)

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(db *gorm.DB, email string, name string) (signedToken string, signedRefreshToken string, err error) {
	claims := &models.SignedDetails{
		Email: email,
		Name:  name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	refreshClaims := &models.SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(168 * time.Hour).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}

func UpdateAllTokens(db *gorm.DB, signedToken string, signedRefreshToken string, userId int64) error {
	user := models.User{}
	result := db.First(&user, userId)
	if result.Error != nil {
		return result.Error
	}

	user.Token = signedToken
	user.RefreshToken = signedRefreshToken
	user.UpdatedAt = time.Now()

	return db.Save(&user).Error
}
func ValidateToken(signedToken string) (*models.SignedDetails, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(
		signedToken,
		&models.SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			// Ensure the signing method is HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		return nil, err
	}

	// Validate the token claims
	claims, ok := token.Claims.(*models.SignedDetails)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Check token expiration
	if claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("token expired")
	}

	return claims, nil
}
