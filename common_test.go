package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

var tmpArticleList []article

// This function is used for setup before executing the test functions
func TestMain(m *testing.M) {
	//Set gin to test mode
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

// Helper function to create a router during testing
func getRouter(withTemplates bool) *gin.Engine {
	r := gin.Default()
	if withTemplates {
		r.LoadHTMLGlob("templates/*")
	}
	return r
}

// Helper function to process a request and test it's response
func testHTTPResponse(t *testing.T, r *gin.Engine, req *http.Request, f func(w *httptest.ResponseRecorder) bool) {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if !f(w) {
		t.Fail()
	}
}

// This function is used to store the main lists into the temporary one for testing
func saveLists() {
	tmpArticleList = articleList
}

func restoreLists() {
	articleList = tmpArticleList
}
