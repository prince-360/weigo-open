package user

import (
	"encoding/json"
	"log"
	"net/http"
	"weigo/api/config"
	"weigo/api/database"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/labstack/echo"
)

func newGoogleAuthClient() *oauth2.Config {
	c := config.GetGoogleConfig()
	return &oauth2.Config{
		ClientID:     c["key"],
		ClientSecret: c["secret"],
		RedirectURL:  c["redirect"],
		Endpoint:     google.Endpoint,
	}
}

type googleUser struct {
	Username string `json:"nickname"`
	ID       string `json:"id"`
	Emails   []struct {
		Value string `json:"value"`
		Type  string `json:"type"`
	} `json:"emails"`
}

func googleUserData(tok *oauth2.Token) (*database.OauthProfile, error) {
	ggAuth := newGoogleAuthClient()
	client := ggAuth.Client(oauth2.NoContext, tok)
	resp, err := client.Get("https://www.googleapis.com/plus/v1/people/me")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	g := googleUser{}
	decoder.Decode(&g)
	op := &database.OauthProfile{}
	op.Type = "google"
	op.AccountID = g.ID
	op.Username = g.Username
	if g.Emails == nil {
		return op, nil
	}
	for k := range g.Emails {
		if g.Emails[k].Type == "account" {
			op.Email = g.Emails[0].Value
			break
		}
	}
	return op, nil
}

func googleOauth(c echo.Context) error {
	ggAuth := newGoogleAuthClient()
	code := c.QueryParam("code")
	tok, err := ggAuth.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Println(err)
		return c.Redirect(http.StatusTemporaryRedirect, config.GetUIDomain()+"#/login")
	}
	op, err := googleUserData(tok)
	if err != nil {
		log.Println(err)
		return c.Redirect(http.StatusTemporaryRedirect, config.GetUIDomain()+"#/login")
	}
	err = database.OauthProfile{}.Insert(op)
	if err != nil {
		log.Println(err)
		return c.Redirect(http.StatusTemporaryRedirect, config.GetUIDomain()+"#/login")
	}
	profile, err := database.ProfileGetByGoogleID(op.AccountID)
	if err != nil {
		log.Println(err)
		return c.Redirect(http.StatusTemporaryRedirect, config.GetUIDomain()+"#/login")
	}
	if profile != nil {
		return c.Redirect(http.StatusTemporaryRedirect, config.GetUIDomain()+"#/login/auth/"+op.Key)
	}
	p, err := database.ProfileGetByEmail(op.Email)
	if err != nil {
		log.Println(err)
		return c.Redirect(http.StatusTemporaryRedirect, config.GetUIDomain()+"#/login")
	}
	if p != nil {
		err := p.ChangeGoogleID(op.AccountID)
		if err != nil {
			log.Println(err)
			return c.Redirect(http.StatusTemporaryRedirect, config.GetUIDomain()+"#/login")
		}
		return c.Redirect(http.StatusTemporaryRedirect, config.GetUIDomain()+"#/login/auth/"+op.Key)
	}
	return c.Redirect(http.StatusTemporaryRedirect, config.GetUIDomain()+"#/login/setup/"+op.Key)
}
