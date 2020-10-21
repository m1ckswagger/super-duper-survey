package main

func initializeRoutes() {
	// Handle index route
	router.GET("/", showIndexPage)

	// Handle single article
	catalogRoutes := router.Group("/catalog")
	{
		catalogRoutes.GET("/view/:catalog_id", getCatalog)
		catalogRoutes.POST("/submit/:catalog_id", submitCatalog)
	}
	ocaiRoutes := router.Group("/ocai")
	{
		ocaiRoutes.GET("/", showOCAICatalog)
		ocaiRoutes.POST("/submit", submitOCAICatalog)
	}
}
