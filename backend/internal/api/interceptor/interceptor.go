package interceptor

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/MattBrs/OcelotMDM/internal/user"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Interceptor struct {
	userService *user.Service
}

func NewAuthInterceptor(userService *user.Service) *Interceptor {
	return &Interceptor{userService}
}

func (i *Interceptor) CheckAuth(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")

	if authHeader == "" {
		ctx.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "authorization header is missing"},
		)
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	authToken := strings.Split(authHeader, " ")
	if len(authToken) != 2 || authToken[0] != "Bearer" {
		ctx.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "authorization token is not in the correct format"},
		)
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	tokenString := authToken[1]
	tk, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil || !tk.Valid {
		ctx.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "invalid or expired token"},
		)
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	claims, ok := tk.Claims.(jwt.MapClaims)
	if !ok {
		ctx.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "invalid token"},
		)
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		ctx.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "expiredtoken"},
		)
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	user, err := i.userService.GetUserById(ctx, claims["id"].(string))
	if err != nil {
		ctx.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "user not found or not authorized"},
		)
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx.Set("currentUser", user)
	ctx.Next()
}
