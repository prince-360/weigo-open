package main

import (
	"os"
	"weigo/api/config"
	"weigo/api/database"
	MediaRoutes "weigo/api/routes/media"
	StreamRoutes "weigo/api/routes/stream"
	UserRoutes "weigo/api/routes/user"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/urfave/cli"
)

func server() {
	r := echo.New()
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())
	r.Use(middleware.Recover())

	r.Static("/dist", "../ui/dist")
	r.Static("/uploaded", "../ui/uploaded")
	r.File("/", "../ui/index.html")

	UserRoutes.Init(r.Group("/user"))
	StreamRoutes.Init(r.Group("/stream"))
	MediaRoutes.Init(r.Group("/media"))
	r.Start(":1314")
}

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "db",
			Value:  "postgres://weigo@127.0.0.1/weigo",
			EnvVar: "DB",
		},
		cli.StringFlag{
			Name:   "ui-domain",
			Value:  "http://localhost:8080",
			EnvVar: "UI_DOMAIN",
		},
		// Google oauth settings
		cli.StringFlag{
			Name:   "google-client-key",
			EnvVar: "GOOGLE_CLIENT_KEY",
		},
		cli.StringFlag{
			Name:   "google-client-secret",
			EnvVar: "GOOGLE_CLIENT_SECRET",
		},
		cli.StringFlag{
			Name:   "google-redirect",
			EnvVar: "GOOGLE_REDIRECT",
		},
		// Github oauth settings
		cli.StringFlag{
			Name:   "github-client-key",
			EnvVar: "GITHUB_CLIENT_KEY",
		},
		cli.StringFlag{
			Name:   "github-client-secret",
			EnvVar: "GITHUB_CLIENT_SECRET",
		},
		cli.StringFlag{
			Name:   "github-redirect",
			EnvVar: "GITHUB_REDIRECT",
		},
	}
	app.Action = func(c *cli.Context) error {
		database.Connect(c.String("db"))
		config.Init(c)
		server()
		return nil
	}
	app.Run(os.Args)
}
