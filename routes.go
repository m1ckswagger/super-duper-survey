package main

func initializeRoutes() {
	// Handle index route
	router.GET("/", showIndexPage)

	// Handle single article
	router.GET("/article/view/:article_id", getArticle)

	userRoutes := router.Group("/u")
	{
		userRoutes.GET("/register", showRegistrationPage)
		userRoutes.POST("/register", register)
		userRoutes.GET("/login", showLoginPage)
		userRoutes.POST("/login", performLogin)
		userRoutes.GET("/logout", logout)
	}
}
