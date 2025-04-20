//Handle panic berguna untuk menghandle error panic.
//error panic itu akan menshutdown app secara paksa
//ketika ada error panic app akan terhenenti secara otomatis
//disini digunakan recover, jadi nanti keika ada error panic ketika menggunakan recover app tidak akan tershutdown dan teteap berjalan

package middlewares

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"user-service/common/response"
	"user-service/config"
	"user-service/constants"
	errConstant "user-service/constants/error"
	services "user-service/services/user"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func HandlePanic() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorf("Recovered from panic: %v", r)
				c.JSON(http.StatusInternalServerError, response.Response{
					Status:  constants.Error,
					Message: errConstant.ErrInternalServerError,
				})

				c.Abort()
			}
		}()
		c.Next()
	}
}

// rate limiter berfungsi untuk memberi batasan req yang masuk ke session
func RateLimiter(lmt *limiter.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)
		if err != nil {
			c.JSON(http.StatusTooManyRequests, response.Response{
				Status:  constants.Error,
				Message: errConstant.ErrTooManyRequests.Error(),
			})
			c.Abort()
			return
		}
		//jika semuanya normal
		c.Next()
	}
}

func extractBearerToken(token string) string {
	arrayToken := strings.Split(token, " ")
	if len(arrayToken) == 2 {
		return arrayToken[1]
	}
	return ""
}

func responseUnauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, response.Response{
		Status:  constants.Error,
		Message: message,
	})
	c.Abort()
}

func validateAPIKey(c *gin.Context) error {
	apiKey := c.GetHeader(constants.XApiKey)
	requestAt := c.GetHeader(constants.XRequestAt)
	serviceName := c.GetHeader(constants.XserviceName)
	signatureKey := config.Config.SignatureKey

	raw := fmt.Sprintf("%s:%s:%s", serviceName, signatureKey, requestAt)
	hash := sha256.New()
	hash.Write([]byte(raw))
	resultHash := hex.EncodeToString(hash.Sum(nil))

	fmt.Println("====== DEBUG API KEY VALIDATION ======")
	fmt.Println("x-api-key        :", apiKey)
	fmt.Println("x-request-at     :", requestAt)
	fmt.Println("x-service-name   :", serviceName)
	fmt.Println("signatureKey     :", signatureKey)
	fmt.Println("Raw string       :", raw)
	fmt.Println("Expected hash    :", resultHash)
	fmt.Println("Match?           :", apiKey == resultHash)
	fmt.Println("======================================")
	fmt.Println("🔍 [DEBUG] generated hash :", resultHash)

	if apiKey != resultHash {
		fmt.Println("❌ [ERROR] API Key tidak valid")
		return errConstant.ErrUnauthorized
	}
	fmt.Println("✅ [INFO] API Key valid")
	return nil
}

func validateBearerToken(c *gin.Context, token string) error {
	if !strings.Contains(token, "Bearer") {
		fmt.Println("❌ [ERROR] Authorization header tidak mengandung Bearer")
		return errConstant.ErrUnauthorized
	}

	tokenString := extractBearerToken(token)
	if tokenString == "" {
		fmt.Println("❌ [ERROR] Bearer token kosong")
		return errConstant.ErrUnauthorized
	}

	claims := &services.Claims{}
	tokenJwt, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		_, oke := token.Method.(*jwt.SigningMethodHMAC)
		if !oke {
			fmt.Println("❌ [ERROR] JWT method tidak sesuai")
			return nil, errConstant.ErrInvalidToken
		}

		jwtSecret := []byte(config.Config.JwtSecretKey)
		return jwtSecret, nil
	})

	if err != nil {
		fmt.Println("❌ [ERROR] JWT parse gagal:", err)
		return errConstant.ErrUnauthorized
	}
	if !tokenJwt.Valid {
		fmt.Println("❌ [ERROR] JWT tidak valid")
		return errConstant.ErrUnauthorized
	}

	fmt.Println("✅ [INFO] JWT valid")
	userLogin := c.Request.WithContext(context.WithValue(c.Request.Context(), constants.UserLogin, claims.User))
	c.Request = userLogin
	c.Set(constants.Token, token)
	return nil
}

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("[DEBUG] Authenticate Middleware Kepanggil")
		var err error
		token := c.GetHeader("Authorization")
		if token == "" {
			fmt.Println("❌ [ERROR] Authorization header tidak ditemukan")
			fmt.Println("Headers:", c.Request.Header)
			responseUnauthorized(c, errConstant.ErrUnauthorized.Error())
			return
		}

		err = validateBearerToken(c, token)
		if err != nil {
			fmt.Println("❌ [ERROR] Validasi Bearer token gagal:", err)
			responseUnauthorized(c, err.Error())
			return
		}

		err = validateAPIKey(c)
		if err != nil {
			fmt.Println("❌ [ERROR] Validasi API Key gagal:", err)
			responseUnauthorized(c, err.Error())
			return
		}
		fmt.Println("✅ [INFO] Bearer token valid")
		fmt.Println("✅ [INFO] API Key valid")
		c.Next()
	}
}
