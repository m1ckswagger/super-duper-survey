package main

import (
	"fmt"
	"net/http"

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
	db, err = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/surveys")
	if err != nil {
		panic(err)
	}
	initCatalogs()
	ocai = generateOCAICatalog()
	ocai.WriteJSON("data/ocai.json")
	gin.SetMode(gin.ReleaseMode)
	router = gin.Default()
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("SESSIONID", store))
	router.LoadHTMLGlob("templates/*")
	initializeRoutes()
	router.RunTLS(":443", "cert/fullchain.pem", "cert/privkey.pem")
}

func render(c *gin.Context, data gin.H, templateName string) {
	data["catalogs"] = getAllCatalogs()
	data["ocai"] = ocai

	switch c.Request.Header.Get("Accept") {
	case "application/json":
		c.JSON(http.StatusOK, data["payload"])
	case "application/xml":
		c.XML(http.StatusOK, data["payload"])
	default:
		c.HTML(http.StatusOK, templateName, data)
	}
}
