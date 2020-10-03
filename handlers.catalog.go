package main

import (
	"github.com/gin-gonic/gin"
)

func showCatalogs(c *gin.Context) {
	render(c, gin.H{"title": "Catalog", "payload": sheff}, "catalog.html")
}
