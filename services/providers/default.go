package providers

import (
	"github.com/gin-gonic/gin"
	"github.com/vesicash/mor-api/external/request"
	"github.com/vesicash/mor-api/pkg/repository/storage/postgresql"
)

var (
	paymentHistorySuccessful = "successful"
)

func HandleDefaultMerchantWebhook(c *gin.Context, extReq request.ExternalRequest, db postgresql.Databases, requestBody []byte) error {
	return nil
}
