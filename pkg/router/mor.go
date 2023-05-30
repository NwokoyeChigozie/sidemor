package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/vesicash/mor-api/external/request"
	"github.com/vesicash/mor-api/pkg/controller/mor"
	"github.com/vesicash/mor-api/pkg/middleware"
	"github.com/vesicash/mor-api/pkg/repository/storage/postgresql"
	"github.com/vesicash/mor-api/utility"
)

func Mor(r *gin.Engine, ApiVersion string, validator *validator.Validate, db postgresql.Databases, logger *utility.Logger) *gin.Engine {
	extReq := request.ExternalRequest{Logger: logger, Test: false}
	mor := mor.Controller{Db: db, Validator: validator, Logger: logger, ExtReq: extReq}

	morUrl := r.Group(fmt.Sprintf("%v", ApiVersion))
	{
		morUrl.POST("/webhook/:account_id", mor.MerchantWebhooks)
	}

	paymentApiUrl := r.Group(fmt.Sprintf("%v", ApiVersion), middleware.Authorize(db, extReq, middleware.AuthType))
	{
		paymentApiUrl.GET("/customers", mor.GetCustomers)
		paymentApiUrl.POST("/save-settings", mor.SaveSettings)
		paymentApiUrl.POST("/get-settings", mor.GetSettings)
	}

	morjobsUrl := r.Group(fmt.Sprintf("%v/jobs", ApiVersion))
	{
		morjobsUrl.POST("/start", mor.StartCronJob)
		morjobsUrl.POST("/start-bulk", mor.StartCronJobsBulk)
		morjobsUrl.POST("/stop", mor.StopCronJob)
		morjobsUrl.PATCH("/update_interval", mor.UpdateCronJobInterval)
	}

	return r
}
