package main

import (
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
	db, err = sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/surveys")
	print(db, err)
	initCatalogs()
	ocai = generateOCAICatalog()
	ocai.WriteJSON("data/ocai.json")
	router = gin.Default()
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("SESSIONID", store))
	router.LoadHTMLGlob("templates/*")
	initializeRoutes()
	router.Run(":443")
}

func render(c *gin.Context, data gin.H, templateName string) {
	loggedInInterface, _ := c.Get("is_logged_in")
	data["is_logged_in"] = loggedInInterface.(bool)
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

func initCatalogs() {
	ocaq := NewRegularCatalog("Organizational Culture Assessment Questionaire", "data/ocaq.txt", "data/ocaq.json", []string{"Yes", "No"})
	catalogs = append(catalogs, ocaq)
	sheff := NewRegularCatalog("Sheffield Culture Survey", "data/sheffield.txt", "data/sheffield.json", []string{
		"Strongly disagree",
		"Disagree",
		"Neutral",
		"Agree",
		"Strongly Agree",
	})
	catalogs = append(catalogs, sheff)
}
