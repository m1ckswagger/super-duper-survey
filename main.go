package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	ldap "github.com/go-ldap/ldap/v3"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Joke struct {
	ID    int    `json:"id" binding:"required"`
	Likes int    `json:"likes"`
	Joke  string `json:"joke" binding:"required"`
}

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

var jokes = []Joke{
	{1, 0, "Did you hear about the restaurant on the moon? Great food, no atmosphere."},
	{2, 0, "What do you call a fake noodle? An Impasta."},
	{3, 0, "How many apples grow on a tree? All of them."},
	{4, 0, "Want to hear a joke about paper? Nevermind it's tearable."},
	{5, 0, "I just watched a program about beavers. It was the best dam program I've ever seen."},
	{6, 0, "Why did the coffee file a police report? It got mugged."},
	{7, 0, "How does a penguin build it's house? Igloos it together."},
}

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

var cache redis.Conn

func main() {
	l, err := ldap.DialURL(fmt.Sprintf("ldaps://ldap.sipwise.com"))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	err = l.UnauthenticatedBind("")
	if err != nil {
		fmt.Println(err)
		return
	}
	username := "mberger"
	searchRequest := ldap.NewSearchRequest(
		"dc=sipwise,dc=com",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(uid=%s))", username),
		[]string{"dn"},
		nil,
	)
	sr, err := l.Search(searchRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(sr.Entries) == 0 {
		fmt.Println("No results found")
		return
	}
	userdn := sr.Entries[0].DN
	fmt.Println(userdn)
	initCache()
	gin.SetMode(gin.ReleaseMode)
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	logFileJSON, err := os.OpenFile("log.json", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer logFileJSON.Close()

	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: false,
		},
	)
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Use(logger.SetLogger())
	// Custom logger
	subLog := zerolog.New(logFileJSON).With().Str("service", "shardik").Logger()
	router.Use(logger.SetLogger(logger.Config{
		Logger:   &subLog,
		UTC:      true,
		SkipPath: []string{"/api"},
	}))
	//router.Use(static.Serve("/", static.LocalFile("./views", true)))

	router.GET("/", func(c *gin.Context) {
		c.HTML(
			// Set the HTTP status to 200 (OK)
			http.StatusOK,
			// Use the index.html template
			"index.html",
			// pass the data that the page uses (in this case, 'title')
			gin.H{
				"title": "Shardik",
			},
		)
	})

	api := router.Group("/api")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
		api.GET("/jokes", JokeHandler)
		api.POST("/jokes/like/:jokeID", LikeJokeHandler)
		api.POST("/login", Login)
	}
	router.Run(":8080")
}

func initCache() {
	conn, err := redis.DialURL("redis://localhost")
	if err != nil {
		panic(err)
	}
	cache = conn
}
func JokeHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, jokes)
}

func LikeJokeHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	if jokeid, err := strconv.Atoi(c.Param("jokeID")); err == nil {
		for j := range jokes {
			if jokes[j].ID == jokeid {
				jokes[j].Likes++
			}
		}
		for i := 0; i < len(jokes); i++ {
			if jokes[i].ID == jokeid {
				jokes[i].Likes++
			}
		}
		c.JSON(http.StatusOK, &jokes)
	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}
}

func Login(c *gin.Context) {
	var creds Credentials
	err := json.NewDecoder(c.Request.Body).Decode(&creds)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	fmt.Println(creds)
	expectedPassword, ok := users[creds.Username]
	if !ok || expectedPassword != creds.Password {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	id, _ := uuid.NewRandom()
	sessionToken := id.String()
	_, err = cache.Do("SETEX", sessionToken, "120", creds.Username)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	c.SetCookie("gin_cookie", sessionToken, 120, "/", "localhost", false, true)
}