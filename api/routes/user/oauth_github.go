package user

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"weigo/api/config"
	"weigo/api/database"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"

	"github.com/labstack/echo"
)

func newGithubAuthClient() *oauth2.Config {
	c := config.GetGithubConfig()
	return &oauth2.Config{
		ClientID:     c["key"],
		ClientSecret: c["secret"],
		RedirectURL:  c["redirect"],
		Endpoint:     github.Endpoint,
	}
}

type githubUser struct {
	ID       int    `json:"id"`
	Username string `json:"login"`
}

func githubUserData(tok *oauth2.Token) (*database.OauthProfile, error) {
	ghAuth := newGithubAuthClient()
	client := ghAuth.Client(oauth2.NoContext, tok)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	var g = githubUser{}
	err = decoder.Decode(&g)
	if err != nil {
		return nil, err
	}
	if g.ID == 0 {
		return nil, errors.New("no id")
	}
	return &database.OauthProfile{
		AccountID: strconv.Itoa(g.ID),
		Username:  g.Username,
		Type:      "github",
	}, nil
}

func githubOauth(c echo.Context) error {
	code := c.QueryParam("code")
	ghAuth := newGithubAuthClient()
	tok, err := ghAuth.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Println(err)
		return c.Redirect(http.StatusTemporaryRedirect, config.GetUIDomain()+"#/login")
	}
	op, err := githubUserData(tok)
	if err != nil {
		log.Println(err)
		return c.Redirect(http.StatusTemporaryRedirect, config.GetUIDomain()+"#/login")
	}
	err = database.OauthProfile{}.Insert(op)
	if err != nil {
		log.Println(err)
		return c.Redirect(http.StatusTemporaryRedirect, config.GetUIDomain()+"#/login")
	}
	profile, err := database.ProfileGetByGithubID(op.AccountID)
	if err != nil {
		log.Println(err)
		return c.Redirect(http.StatusTemporaryRedirect, config.GetUIDomain()+"#/login")
	}
	if profile != nil {
		return c.Redirect(http.StatusTemporaryRedirect, config.GetUIDomain()+"#/login/auth/"+op.Key)
	}
	return c.Redirect(http.StatusTemporaryRedirect, config.GetUIDomain()+"#/login/setup/"+op.Key)
}
