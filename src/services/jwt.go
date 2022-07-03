package services

import (
	"context"
	"encoding/json"
	"fmt"
	c "github.com/pbalan/ontrack/src/config"
	"time"

	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type JwtCustomClaim struct {
	ID string `json:"id"`
	jwt.StandardClaims
}

var jwtSecret = []byte(getJwtSecret())
var configuration c.Configurations

func getJwtSecret() string {
	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}
	secret, _ := json.Marshal(configuration.Jwt.Secret)
	if string(secret) == "" {
		log.Fatal("JWT_SECRET not configured")
	}
	return string(secret)
}

func JwtGenerate(ctx context.Context, userID string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, &JwtCustomClaim{
		ID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	})

	token, err := t.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}

func JwtValidate(ctx context.Context, token string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, &JwtCustomClaim{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there's a problem with the signing method")
		}
		return jwtSecret, nil
	})
}
