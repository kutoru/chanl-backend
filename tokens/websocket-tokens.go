package tokens

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func CreateWebsocketToken(channelId int, userId int) (string, error) {
	claims := jwt.MapClaims{}
	claims["channelId"] = channelId
	claims["userId"] = userId
	claims["exp"] = time.Now().Add(5 * time.Second).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("WS_TOKEN_KEY")))

	if err != nil {
		return "", err
	} else {
		return tokenString, nil
	}
}

func ParseWebsocketToken(tokenString string) (int, int, error) {
	token, err := jwt.Parse(tokenString, parseHelper)
	if err != nil {
		return 0, 0, err
	}

	if !token.Valid {
		return 0, 0, fmt.Errorf("token is invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, 0, fmt.Errorf("claims are invalid")
	}

	channelId, exists := claims["channelId"].(float64)
	if !exists {
		return 0, 0, fmt.Errorf("channelId is invalid")
	}

	userId, exists := claims["userId"].(float64)
	if !exists {
		return 0, 0, fmt.Errorf("userId is invalid")
	}

	return int(channelId), int(userId), nil
}

func parseHelper(token *jwt.Token) (interface{}, error) {
	_, ok := token.Method.(*jwt.SigningMethodHMAC)
	if !ok {
		return nil, fmt.Errorf("invalid signing method: %v", token.Header["alg"])
	}

	return []byte(os.Getenv("WS_TOKEN_KEY")), nil
}
