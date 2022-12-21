package jwt

import (
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
    "github.com/kataras/iris"
	"github.com/zhashkevych/auth/pkg/parser"
	"..\models"
)


//Signin - authorization method
func (a *Authorizer) SignIn(username, password string) (string, error) {
	//check if user exists and password is correct
	//connect to db
	dsn := "host=localhost user=postgres password=mypas dbname=pears port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	//get user from db
	var user models.User
	db.Where("login = ?", username).First(&user)
	if user.password != password {
		return "", fmt.Errorf("invalid username or password")
	}

	//create token
	claims := models.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(a.expireDuration)),
			IssueAt:   jwt.At(time.Now()),
		},
		Username: username,
	}

	//return token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.signingkey)
}

//ParseToken - parse token method
func ParseToken(accessToken string, signingKey []byte) (string, error){
	token, err := jwt.ParseWithClaims(accessToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return signingKey, nil
	}

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.Username, nil
	}

	return "", fmt.Errorf("invalid token")
}

//Middleware - method for checking token
func Middleware (c iris.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err := parser.ParseToken(headerParts[1], SIGNING_KEY)
	if err != nil {
		status := http.StatusBadRequest
		if err == jwt.ErrInvalidToken {
			status = http.StatusUnauthorized
		}
		c.AbortWithStatus(status)
		return
	}
}











