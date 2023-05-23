package router

import (
	"net/http"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/vesicash/mor-api/internal/config"
	"github.com/vesicash/mor-api/pkg/middleware"
	"github.com/vesicash/mor-api/pkg/repository/storage/postgresql"
	"github.com/vesicash/mor-api/utility"
)

func Setup(logger *utility.Logger, validator *validator.Validate, db postgresql.Databases, appConfiguration *config.App) *gin.Engine {
	if appConfiguration.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()

	// Middlewares
	// r.Use(gin.Logger())
	r.ForwardedByClientIP = true
	r.SetTrustedProxies(config.GetConfig().Server.TrustedProxies)
	r.Use(middleware.PrometheusMiddleware())
	r.Use(middleware.Security())
	r.Use(middleware.Throttle())
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.MaxMultipartMemory = 1 << 20 // 1MB

	ApiVersion := "v2"
	Health(r, ApiVersion, validator, db, logger)
	Mor(r, ApiVersion, validator, db, logger)

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "Welcome to mor-api",
			"status":  http.StatusOK,
		})
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"name":    "Not Found",
			"message": "Page not found.",
			"code":    404,
			"status":  http.StatusNotFound,
		})
	})

	return r
}
