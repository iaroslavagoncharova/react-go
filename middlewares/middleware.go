package middlewares

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/iaroslavagoncharova/react-go/config"
)

func AuthMiddleware(c *fiber.Ctx) error {
    // jwt token from authorization header
    authHeader := c.Get("Authorization")
    if authHeader == "" {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid token"})
    }

    tokenString := authHeader[len("Bearer "):]

    // validate token and return user ID and role
    userID, role, err := ValidateJWT(tokenString)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid token"})
    }

    // set id and role in locals
    c.Locals("user_id", userID)
    c.Locals("role", role)

    return c.Next()
}

func GenerateJWT(userID string, role string) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
        "role":    role,
        "exp":     time.Now().Add(time.Hour * 72).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(config.JWTSecret))
    if err != nil {
        return "", err
    }
    return tokenString, nil
}

func ValidateJWT(tokenString string) (string, string, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(config.JWTSecret), nil
    })

    if err != nil {
        return "", "", err
    }

    // extract claims if the token is valid
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        userID := claims["user_id"].(string)
        role := claims["role"].(string)
        fmt.Printf("Token is valid, User ID: %s, Role: %s\n", userID, role)
        return userID, role, nil
    }

    return "", "", fmt.Errorf("invalid token")
}


func Authorize(roles ...string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        userRole := c.Locals("role")
        for _, role := range roles {
            if userRole == role {
                return c.Next()
            }
        }
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You are not authorized to access this resource"})
    }
}