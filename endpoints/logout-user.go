package endpoints

import (
	"log"
	"net/http"

	"github.com/kutoru/chanl-backend/tokens"
)

func logoutUser(w http.ResponseWriter, r *http.Request) {
	log.Println("logoutUser called")

	_, token, err := tokens.CheckAuthToken(r)
	if err == nil {
		tokens.RemoveAuthToken(token)
	}

	cookie := &http.Cookie{
		Name:     "_ltat_",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
		SameSite: http.SameSiteNoneMode,
	}

	http.SetCookie(w, cookie)
	sendJSONResponse(w, true)
}
