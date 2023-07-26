package tokens

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
)

type TokenInfo struct {
	Exp    int64
	UserID int
}

// Ideally I would probably run a cleanup task every day or so to remove expired tokens. Technically noone would be able to abuse them, but they take some memory space for no reason
var authedUsers map[string]*TokenInfo = make(map[string]*TokenInfo)

func CreateAuthToken(userId int) (*http.Cookie, error) {
	tokenBytes := make([]byte, 16)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return nil, err
	}

	token := base64.RawURLEncoding.EncodeToString(tokenBytes)
	// This token is supposed to be long term, 30 days to be exact. But for now it is just 30 seconds
	maxAge := 30 // 2592000

	cookie := &http.Cookie{
		Name:     "_ltat_",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   maxAge,
		SameSite: http.SameSiteNoneMode,
	}

	tokenInfo := &TokenInfo{
		Exp:    time.Now().Add(time.Duration(maxAge) * time.Second).Unix(),
		UserID: userId,
	}

	authedUsers[token] = tokenInfo
	return cookie, nil
}

// Returns userId and token string when the auth token is valid. 0, "" and error otherwise
func CheckAuthToken(r *http.Request) (int, string, error) {
	cookie, err := r.Cookie("_ltat_")
	if err != nil {
		return 0, "", fmt.Errorf("could not get the auth cookie: %v", err)
	}

	err = cookie.Valid()
	if err != nil {
		return 0, "", fmt.Errorf("the cookie is invalid: %v", err)
	}

	token := cookie.Value
	tokenInfo, exists := authedUsers[token]
	if !exists {
		return 0, "", fmt.Errorf("the token does not exist")
	}

	if tokenInfo.Exp <= time.Now().Unix() {
		RemoveAuthToken(cookie.Value)
		return 0, "", fmt.Errorf("the token has expired")
	}

	return tokenInfo.UserID, token, nil
}

func UpdateAuthToken(userId int, token string) (*http.Cookie, error) {
	cookie, err := CreateAuthToken(userId)
	if err != nil {
		return nil, err
	} else {
		RemoveAuthToken(token)
		return cookie, nil
	}
}

func RemoveAuthToken(token string) {
	delete(authedUsers, token)
}
