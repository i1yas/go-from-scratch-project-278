package httpapi

import "github.com/gin-gonic/gin"

// RegisterRoutes register routes with handlers
func RegisterRoutes(
	router gin.IRouter,
	links *LinksHandler,
) {
	rootAPI := router.Group("/api")

	linksAPI := rootAPI.Group("/links")
	linksAPI.GET("", links.GetLinks)
	linksAPI.GET("/:id", links.GetLinkByID)
}
