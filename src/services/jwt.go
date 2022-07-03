package services

import (
	"context"
	"encoding/json"
	"fmt"
	c "github.com/pbalan/ontrack/src/config"
	"github.com/pbalan/ontrack/src/graph/model"
	"gorm.io/gorm"
	"time"

	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/twinj/uuid"
)

type JwtCustomClaim struct {
	ID string `json:"id"`
	jwt.StandardClaims
}

var jwtSecret = []byte(getJwtSecret())
var jwtRefreshSecret = []byte(getRefreshSecret())
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

func getRefreshSecret() string {
	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}
	refreshSecret, _ := json.Marshal(configuration.Jwt.RefreshSecret)
	if string(refreshSecret) == "" {
		log.Fatal("JWT_REFRESH_SECRET not configured")
	}
	return string(refreshSecret)
}

func JwtGenerate(ctx context.Context, userID string) (*model.TokenDetail, error) {
	td := &model.TokenDetail{}
	td.AtValidUpto = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUUID = uuid.NewV4().String()

	td.RtValidUpto = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUUID = uuid.NewV4().String()

	var err error
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUUID
	atClaims["user_id"] = userID
	atClaims["exp"] = td.AtValidUpto
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}

	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUUID
	rtClaims["user_id"] = userID
	rtClaims["exp"] = td.RtValidUpto
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString(jwtRefreshSecret)
	if err != nil {
		return nil, err
	}

	return td, nil
}

func SaveAuth(db *gorm.DB, td *model.TokenDetail) (Token *model.TokenDetail, err error) {
	err = db.Create(&td).Error

	if err != nil {
		return td, err
	}

	return td, nil
}

func VerifyAuth(db *gorm.DB, userId uint64, td *model.TokenDetail) (exists bool, err error) {
	err = db.Where(
		"refresh_token = ? AND refresh_uuid = ? AND rt_valid_upto >= ?",
		td.RefreshToken,
		userId,
		time.Now().Unix()).First(&td).Error

	if err != nil {
		return true, err
	}

	return false, nil
}

func JwtValidate(ctx context.Context, token string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, &JwtCustomClaim{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there's a problem with the signing method")
		}
		return jwtSecret, nil
	})
}
