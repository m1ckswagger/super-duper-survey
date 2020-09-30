package main

func initializeRoutes() {
	// Handle index route
	router.GET("/", showIndexPage)

	// Handle single article
	router.GET("/article/view/:article_id", getArticle)
}
