package main

import (
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday"
	_ "github.com/heroku/x/hmetrics/onload"
	//"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"io/ioutil"
	"fmt"
	"time"
)

//var client http.Client

func main() {

	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     "L7xxd11b2cd54f5341b5ad7acdc90e3c6807",
		ClientSecret: "4169a027773442f6b09d307a6a851715",
		Scopes:       []string{"oob"},
		RedirectURL:  "https://intense-plateau-14305.herokuapp.com/authCallback",
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
		authURL := conf.AuthCodeURL("state", oauth2.AccessTypeOnline)
		c.Redirect(301, authURL)
	})


	router.GET("/authCallback", func(c *gin.Context) {
		code := c.Query("code")
		tok, err := conf.Exchange(ctx, code)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Token:", tok)
		client := conf.Client(ctx, tok)
		fmt.Println(client)

		currentTime := time.Now()
		endTime := currentTime.AddDate(0,0,14)// . Add(time.Hour * time.Duration(336)) //add 336 hours (14 days * 24hrs)

		fmt.Println("CurrentTime:", currentTime.Format(time.RFC3339))
		fmt.Println("EndTime:", endTime.Format(time.RFC3339))

		msgURL := fmt.Sprintf("https://apis.hootsuite.com/v1/messages?startTime=%s&endTime=%s&limit=100", url.QueryEscape(currentTime.Format(time.RFC3339)), url.QueryEscape(endTime.Format(time.RFC3339)))

		fmt.Println("Request URL:", msgURL)
		resp, err := client.Get(msgURL)
		fmt.Println("GET request completed: ", resp)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body[:]))
		c.String(http.StatusOK, string(body[:]))
	})



	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/mark", func(c *gin.Context) {
		c.String(http.StatusOK, string(blackfriday.MarkdownBasic([]byte("**hi!**"))))
	})

	router.Run(":" + port)
}
