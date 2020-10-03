package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine
var ocaq *Catalog
var sheff *Catalog
var ocai *CatalogOCAI

func main() {
	ocaq = NewCatalog("data/ocaq.txt", "data/ocaq.json", []string{"Yes", "No"})
	sheff = NewCatalog("data/sheffield.txt", "data/sheffield.json", []string{
		"Strongly disagree",
		"Disagree",
		"Neutral",
		"Agree",
		"Strongly Agree",
	})
	ocai = generateOCAICatalog()
	ocai.WriteJSON("data/ocai.json")
	// os.Exit(1)
	router = gin.Default()
	router.LoadHTMLGlob("templates/*")
	initializeRoutes()
	router.Run(":8080")
}

func render(c *gin.Context, data gin.H, templateName string) {
	loggedInInterface, _ := c.Get("is_logged_in")
	data["is_logged_in"] = loggedInInterface.(bool)

	switch c.Request.Header.Get("Accept") {
	case "application/json":
		c.JSON(http.StatusOK, data["payload"])
	case "application/xml":
		c.XML(http.StatusOK, data["payload"])
	default:
		c.HTML(http.StatusOK, templateName, data)
	}
}
