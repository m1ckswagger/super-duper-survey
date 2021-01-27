package main

// This is a very important package

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"database/sql"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var router *gin.Engine
var ocai *OCAICatalog
var catalogs []*RegularCatalog
var db *sql.DB

func main() {
	var err error
	fmt.Println(err)
	db, err = sql.Open("mysql", "root:survey@tcp(db:3306)/survey")
	if err != nil {
		log.Panic(err)
	}
	initCatalogs()
	ocai = generateOCAICatalog()

	ocai.WriteJSON("data/ocai.json")

	// gin.SetMode(gin.ReleaseMode)

	router = gin.Default()
	store := cookie.NewStore([]byte("secret"))

	router.Use(sessions.Sessions("SESSIONID", store))
	router.LoadHTMLGlob("templates/*")

	initializeRoutes()
	log.Println("Starting router")
	//router.RunTLS(":443", "cert/fullchain.pem", "cert/privkey.pem")
	router.Run(":80")
}

func render(c *gin.Context, data gin.H, templateName string) {
	data["catalogs"] = getAllCatalogs()
	data["ocai"] = ocai
	headers := c.Request.Header.Get("Accept")
	log.Println(headers)
	log.Println(c.Request.Header.Get("auth_token"))
	switch {
	case strings.Contains(headers, "application/json"):
		c.JSON(http.StatusOK, data["payload"])
	case headers == "application/xml":
		c.XML(http.StatusOK, data["payload"])
	default:
		c.HTML(http.StatusOK, templateName, data)
	}
}
