package main

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/aerosystems/nix-trainee-5-6-7-8/docs" // docs are generated by Swag CLI, you have to import it
)

func (app *Config) NewRouter() *echo.Echo {
	e := echo.New()

	e.GET("/docs/*", echoSwagger.WrapHandler)

	e.POST("/v1/users/registration", app.BaseHandler.Registration)
	e.POST("/v1/users/confirmation", app.BaseHandler.Confirmation)

	e.POST("/v1/users/login", app.BaseHandler.Login)
	e.POST("/v1/users/logout", app.BaseHandler.Logout, app.AuthorizationMiddleware())

	e.POST("/v1/tokens/refresh", app.BaseHandler.RefreshToken, app.AuthorizationMiddleware())

	e.GET("/v1/users/login/google", app.BaseHandler.Oauth2GoogleLogin)
	e.GET("/v1/callback/google", app.BaseHandler.Oauth2GoogleCallback)

	e.GET("/v1/comments", app.BaseHandler.ReadComments, app.AuthorizationMiddleware())
	e.GET("/v1/comments/:id", app.BaseHandler.ReadComment, app.AuthorizationMiddleware())
	e.POST("/v1/comments", app.BaseHandler.CreateComment, app.AuthorizationMiddleware())
	e.PATCH("/v1/comments/:id", app.BaseHandler.UpdateComment, app.AuthorizationMiddleware())
	e.DELETE("/v1/comments/:id", app.BaseHandler.DeleteComment, app.AuthorizationMiddleware())

	e.GET("/v1/posts", app.BaseHandler.ReadPosts, app.AuthorizationMiddleware())
	e.GET("/v1/posts/:id", app.BaseHandler.ReadPost, app.AuthorizationMiddleware())
	e.POST("/v1/posts", app.BaseHandler.CreatePost, app.AuthorizationMiddleware())
	e.PATCH("/v1/posts/:id", app.BaseHandler.UpdatePost, app.AuthorizationMiddleware())
	e.DELETE("/v1/posts/:id", app.BaseHandler.DeletePost, app.AuthorizationMiddleware())

	return e
}
