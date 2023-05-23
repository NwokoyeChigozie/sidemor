package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/vesicash/mor-api/external/request"
	"github.com/vesicash/mor-api/pkg/controller/mor"
	"github.com/vesicash/mor-api/pkg/repository/storage/postgresql"
	"github.com/vesicash/mor-api/utility"
)

func Mor(r *gin.Engine, ApiVersion string, validator *validator.Validate, db postgresql.Databases, logger *utility.Logger) *gin.Engine {
	extReq := request.ExternalRequest{Logger: logger, Test: false}
	mor := mor.Controller{Db: db, Validator: validator, Logger: logger, ExtReq: extReq}

	morjobsUrl := r.Group(fmt.Sprintf("%v/jobs", ApiVersion))
	{
		morjobsUrl.POST("/start", mor.StartCronJob)
		morjobsUrl.POST("/start-bulk", mor.StartCronJobsBulk)
		morjobsUrl.POST("/stop", mor.StopCronJob)
		morjobsUrl.PATCH("/update_interval", mor.UpdateCronJobInterval)
	}
	return r
}
