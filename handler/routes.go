package handler

import (
	"github.com/labstack/echo/v4"
	"golang-starter-pack/router/middleware"
	"golang-starter-pack/utils"
)

func (h *Handler) Register(v1 *echo.Group) {
	jwtMiddleware := middleware.JWT(utils.JWTSecret)
	guestPlayers := v1.Group("/players")
	guestPlayers.POST("", h.SignUp)
	guestPlayers.POST("/login", h.Login)

	player := v1.Group("/player", jwtMiddleware)
	player.GET("", h.CurrentPlayer)
	player.PUT("", h.UpdatePlayer)

	profiles := v1.Group("/profiles", jwtMiddleware)
	profiles.GET("/:username", h.GetProfile)
	profiles.POST("/:username/follow", h.Follow)
	profiles.DELETE("/:username/follow", h.Unfollow)

	items := v1.Group("/items", middleware.JWTWithConfig(
		middleware.JWTConfig{
			Skipper: func(c echo.Context) bool {
				if c.Request().Method == "GET" && c.Path() != "/api/items/feed" {
					return true
				}
				return false
			},
			SigningKey: utils.JWTSecret,
		},
	))
	items.POST("", h.CreateItem)
	items.GET("/feed", h.Feed)
	items.PUT("/:slug", h.UpdateItem)
	items.DELETE("/:slug", h.DeleteItem)
	items.POST("/:slug/comments", h.AddComment)
	items.DELETE("/:slug/comments/:id", h.DeleteComment)
	items.POST("/:slug/favorite", h.Favorite)
	items.DELETE("/:slug/favorite", h.Unfavorite)
	items.GET("", h.Items)
	items.GET("/:slug", h.GetItem)
	items.GET("/:slug/comments", h.GetComments)

	tags := v1.Group("/tags")
	tags.GET("", h.Tags)
}
