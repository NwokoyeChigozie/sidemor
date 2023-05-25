package mor

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/mor-api/services/mor-api"
	"github.com/vesicash/mor-api/utility"
)

func (base *Controller) MerchantWebhooks(c *gin.Context) {
	requestBody, err := c.GetRawData()
	if err != nil {
		base.ExtReq.Logger.Error("monnify callback log error", "Failed to read request body", err.Error())
	}

	code, err := mor.MerchantWebhooksService(c, base.ExtReq, base.Db, requestBody)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successful", nil)
	c.JSON(http.StatusOK, rd)

}
