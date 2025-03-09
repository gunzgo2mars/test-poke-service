package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type IJwt interface {
	// Jwt
	GenerateToken(payload jwt.MapClaims) (string, error)

	// Middleware
	ValidateMiddleware(next echo.HandlerFunc) echo.HandlerFunc
}

type JwtEntities struct {
	secretKey string
}

func NewJwt(secret string) IJwt {
	return &JwtEntities{
		secretKey: secret,
	}
}

func (j *JwtEntities) GenerateToken(payload jwt.MapClaims) (string, error) {
	// Add issued at and expiration time to the payload
	payload["iat"] = time.Now().Unix()
	payload["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expires in 1 day

	// Create a new token object, specifying signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	// Sign the token with the secret key and get the complete, signed token
	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *JwtEntities) ValidateMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing JWT token")
		}

		// Parse the token
		token, err := jwt.Parse(
			strings.TrimPrefix(tokenString, "Bearer "),
			func(token *jwt.Token) (interface{}, error) {
				return []byte(j.secretKey), nil
			},
		)
		if err != nil || !token.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid JWT token")
		}

		// Set the token claims in the context
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid JWT claims")
		}

		c.Set("user", claims)

		return next(c)
	}
}
