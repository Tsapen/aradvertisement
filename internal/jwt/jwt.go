package jwt

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Tsapen/aradvertisement/internal/auth"
	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
)

// Secrets struct contains environment variables for creating tokens.
type Secrets struct {
	SetManually   bool
	AccessSecret  string
	RefreshSecret string
}

// PrepareAuthEnvironment set secret
func PrepareAuthEnvironment(s Secrets) error {
	if err := os.Setenv("ACCESS_SECRET", s.AccessSecret); err != nil {
		return err
	}

	if err := os.Setenv("REFRESH_SECRET", s.RefreshSecret); err != nil {
		return err
	}

	return nil
}

// CreateToken creates token.
func CreateToken(username string) (*auth.TokenDetails, error) {
	var td = &auth.TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUUID = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUUID = uuid.NewV4().String()

	var err error

	var atClaims = jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUUID
	atClaims["username"] = username
	atClaims["exp"] = td.AtExpires

	var at = jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	var rtClaims = jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUUID
	rtClaims["username"] = username
	rtClaims["exp"] = td.RtExpires

	var rt = jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}

	return td, nil
}

// TokenValid checks token validity.
func TokenValid(tokenString string) error {
	var token, err = verifyToken(tokenString)
	if err != nil {
		return err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}

	return nil
}

func verifyToken(tokenString string) (*jwt.Token, error) {
	var token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

// ExtractTokenMetadata extracts metadata from request.
func ExtractTokenMetadata(tokenString string) (*auth.AccessDetails, error) {
	var token, err = verifyToken(tokenString)
	if err != nil {
		return nil, err
	}

	var claims, ok = token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		var accessUUID, ok = claims["access_uuid"].(string)
		if !ok {
			return nil, errors.New("uuid is not string")
		}

		var username string
		username, ok = claims["username"].(string)
		if !ok {
			return nil, errors.New("username is not string")
		}

		return &auth.AccessDetails{
			AccessUUID: accessUUID,
			Username:   username,
		}, nil
	}

	return nil, errors.New("can't authorize")
}

// Parse parses, validates, and returns a token.
func Parse(refreshToken string) (map[string]interface{}, error) {
	var token, err = jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok {
		return nil, errors.New("bad token format")
	}

	var mapClaims, ok = token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("bad token format")
	}

	return mapClaims, nil
}
