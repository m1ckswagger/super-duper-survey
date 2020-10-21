package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func showIndexPage(c *gin.Context) {
	catalogs := getAllCatalogs()
	session := sessions.Default(c)
	id := session.Get("ID")
	if id == nil {
		uid, _ := uuid.NewRandom()
		fmt.Println(uid.String())
		session.Set("ID", uid.String())
	}
	ocaq := session.Get("catalog0")
	if ocaq == nil {
		ocaq = false
	}
	sheff := session.Get("catalog1")
	if sheff == nil {
		sheff = false
	}
	ocaiComplete := session.Get("ocai")
	if ocaiComplete == nil {
		ocaiComplete = false
	}
	session.Save()
	render(c, gin.H{"title": "Home Page",
		"payload":      catalogs,
		"ocaq":         ocaq.(bool),
		"sheff":        sheff.(bool),
		"ocai":         ocai,
		"ocaiComplete": ocaiComplete.(bool)},
		"index.html")
}

func getCatalog(c *gin.Context) {
	if catalogID, err := strconv.Atoi(c.Param("catalog_id")); err == nil {
		if catalog, err := getCatalogByID(catalogID); err == nil {
			session := sessions.Default(c)
			completed := session.Get(fmt.Sprintf("catalog%d", catalogID))
			if completed != nil && completed.(bool) {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
			render(c, gin.H{
				"title":   catalog.Name,
				"payload": catalog}, "catalog.html")
		} else {
			c.AbortWithError(http.StatusNotFound, err)
		}
	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}
}

func submitCatalog(c *gin.Context) {
	if catalogID, err := strconv.Atoi(c.Param("catalog_id")); err == nil {
		if catalog, err := getCatalogByID(catalogID); err == nil {
			queryPrep := "INSERT INTO `answers` (`surveyID`, `session_id`, `num`, `answer`) " +
				"VALUES (%d, '%s', %d, '%s');"
			session := sessions.Default(c)
			id := session.Get("ID")
			if id == nil {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
			c.Request.ParseForm()
			for k, v := range c.Request.PostForm {
				num, _ := strconv.Atoi(k)
				query := fmt.Sprintf(queryPrep, catalog.ID, id, num, v[0])
				_, err := db.Exec(query)
				if err != nil {
					fmt.Println(err)
				}
			}
			session.Set(fmt.Sprintf("catalog%d", catalogID), true)
			session.Save()
			render(c, gin.H{"title": "Success!"}, "survey-submit-successful.html")
		} else {
			c.AbortWithError(http.StatusNotFound, err)
		}
	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}
}
func showOCAICatalog(c *gin.Context) {
	session := sessions.Default(c)
	completed := session.Get("ocai")
	if completed != nil && completed.(bool) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	render(c, gin.H{"title": ocai.Name, "payload": ocai}, "ocai.html")
}
func submitOCAICatalog(c *gin.Context) {
	answers := make([]*OCAICategory, 6)
	session := sessions.Default(c)
	id := session.Get("ID")
	if id == nil {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	c.Request.ParseForm()
	for k, v := range c.Request.PostForm {
		index, _ := strconv.Atoi(string(k[0]))
		index--
		a := answers[index]
		if a == nil {
			a = &OCAICategory{ID: index + 1}
		}
		switch suffix := string(k[1:]); suffix {
		case "an":
			a.A.Now, _ = strconv.Atoi(v[0])
		case "ap":
			a.A.Preferred, _ = strconv.Atoi(v[0])
		case "bn":
			a.B.Now, _ = strconv.Atoi(v[0])
		case "bp":
			a.B.Preferred, _ = strconv.Atoi(v[0])
		case "cn":
			a.C.Now, _ = strconv.Atoi(v[0])
		case "cp":
			a.C.Preferred, _ = strconv.Atoi(v[0])
		case "dn":
			a.D.Now, _ = strconv.Atoi(v[0])
		case "dp":
			a.D.Preferred, _ = strconv.Atoi(v[0])
		}
		answers[index] = a
	}
	queryPrep := "INSERT INTO `ocai` (`session_id`, `num`, `a_now`, `a_preferred`, `b_now`, `b_preferred`, `c_now`, `c_preferred`, `d_now`, `d_preferred`) " +
		"VALUES ('%s', %d, %d, %d, %d, %d, %d, %d, %d, %d);"
	for _, a := range answers {
		query := fmt.Sprintf(queryPrep, id, a.ID, a.A.Now, a.A.Preferred, a.B.Now, a.B.Preferred, a.C.Now, a.C.Preferred, a.D.Now, a.D.Preferred)
		_, err := db.Exec(query)
		if err != nil {
			fmt.Println(err)
		}
	}
	session.Set("ocai", true)
	session.Save()
	render(c, gin.H{"title": "Success!"}, "survey-submit-successful.html")
}
