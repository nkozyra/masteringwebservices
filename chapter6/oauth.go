package main

import (
	"code.google.com/p/goauth2/oauth"
	"fmt"
)

var (
	clientID     = "--"
	clientSecret = "--"
	scope        = "https://api.twitter.com/1.1/statuses/home_timeline.json"
	redirectURL  = "http://www.mastergoco.com"
	authURL      = "https://www.facebook.com/dialog/oauth"
	tokenURL     = "https://graph.facebook.com/oauth/access_token"
	requestURL   = "https://graph.facebook.com/me"
	code         = ""
)

func main() {

	oauthConnection := &oauth.Config{
		ClientId:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scope:        scope,
		AuthURL:      authURL,
		TokenURL:     tokenURL,
	}
	//transport := &oauth.Transport{Config: oauthConnection}

	url := oauthConnection.AuthCodeURL("")
	fmt.Print("Visit this URL to get a code, then run again with -code=YOUR_CODE\n\n")
	fmt.Println(url)

}
