package config

import (
	"log"
	"os"
)

var JWTSecret []byte

func LoadConfig() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is not set")
	}
	JWTSecret = []byte(secret)
}