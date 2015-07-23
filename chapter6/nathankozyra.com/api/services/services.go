package service

import (
	"code.google.com/p/goauth2/oauth"
	"fmt"
)

type OauthService struct {
	clientID     string
	clientSecret string
	scope        string
	redirectURL  string
	authURL      string
	tokenURL     string
	requestURL   string
	code         string
}

var OauthServices map[string]OauthService

func InitServices() {

	OauthServices = map[string]OauthService{}

	OauthServices["facebook"] = OauthService{
		clientID:     "-",
		clientSecret: "-",
		scope:        "",
		redirectURL:  "-",
		authURL:      "-",
		tokenURL:     "-",
		requestURL:   "-",
		code:         "",
	}
	OauthServices["google"] = OauthService{
		clientID:     "-",
		clientSecret: "-",
		scope:        "https://www.googleapis.com/auth/plus.login",
		redirectURL:  "http://www.mastergoco.com/connect/google",
		authURL:      "https://accounts.google.com/o/oauth2/auth",
		tokenURL:     "https://accounts.google.com/o/oauth2/token",
		requestURL:   "https://graph.facebook.com/me",
		code:         "",
	}
}

func PostMessage(service string, authCode string, scope string) bool {
	OauthServices[service].scope = scope
	token, err = transport.Exchange(*authCode)
	if err != nil {
		log.Fatal("Exchange:", err)
	}

}

func GetAccessTokenURL(service string, scope string) string {

	oauthConnection := &oauth.Config{
		ClientId:     OauthServices[service].clientID,
		ClientSecret: OauthServices[service].clientSecret,
		RedirectURL:  OauthServices[service].redirectURL,
		Scope:        OauthServices[service].scope,
		AuthURL:      OauthServices[service].authURL,
		TokenURL:     OauthServices[service].tokenURL,
	}

	fmt.Println(OauthServices[service])

	url := oauthConnection.AuthCodeURL("")
	return url

}
