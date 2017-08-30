package config

import "github.com/urfave/cli"

var (
	googleConfig = map[string]string{}
	githubConfig = map[string]string{}
	uiDomain     = ""
)

// Init .
func Init(c *cli.Context) {
	googleConfig["key"] = c.String("google-client-key")
	googleConfig["secret"] = c.String("google-client-secret")
	googleConfig["redirect"] = c.String("google-redirect")

	githubConfig["key"] = c.String("github-client-key")
	githubConfig["secret"] = c.String("github-client-secret")
	githubConfig["redirect"] = c.String("github-redirect")

	uiDomain = c.String("ui-domain")
}

// GetUIDomain .
func GetUIDomain() string {
	return uiDomain
}

// GetGoogleConfig .
func GetGoogleConfig() map[string]string {
	return googleConfig
}

// GetGithubConfig .
func GetGithubConfig() map[string]string {
	return githubConfig
}
