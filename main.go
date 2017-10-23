package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday"
	_ "github.com/heroku/x/hmetrics/onload"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

var ctx
var conf


func main() {

	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     "L7xxd11b2cd54f5341b5ad7acdc90e3c6807",
		ClientSecret: "4169a027773442f6b09d307a6a851715",
		Scopes:       []string{"oob"},
		RedirectURL:  "https://www.getpostman.com/oauth2/callback",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://apis.hootsuite.com/auth/oauth/v2/authorize",
			TokenURL: "https://apis.hootsuite.com/auth/oauth/v2/token",
		},
	}



	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")


	router.GET("/login", func(c *gin.Context) {
		url := conf.AuthCodeURL("state", oauth2.AccessTypeOnline)
		c.Redirect(0, url)
	})





	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/mark", func(c *gin.Context) {
		c.String(http.StatusOK, string(blackfriday.MarkdownBasic([]byte("**hi!**"))))
	})

	router.Run(":" + port)
}
