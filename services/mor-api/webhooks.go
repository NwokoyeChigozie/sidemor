package mor

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/mor-api/external/request"
	"github.com/vesicash/mor-api/internal/models"
	"github.com/vesicash/mor-api/pkg/repository/storage/postgresql"
	"github.com/vesicash/mor-api/services/providers"
	"github.com/vesicash/mor-api/utility"
)

func MerchantWebhooksService(c *gin.Context, extReq request.ExternalRequest, db postgresql.Databases, requestBody []byte) (int, error) {
	var (
		err error
	)

	provider, err := GetProvider(c, requestBody)
	if err != nil {
		extReq.Logger.Error(fmt.Sprintf("webhook log error for %v %v", provider, err.Error()))
		return http.StatusInternalServerError, err
	}

	logWebhookData(extReq, db, provider, requestBody)

	switch strings.ToLower(provider) {
	case "rave":
		err = providers.HandleRaveMerchantWebhook(c, extReq, db, requestBody)
	case "e-transact":
		err = providers.HandleETransactMerchantWebhook(c, extReq, db, requestBody)
	case "paystack":
		err = providers.HandlePaystackMerchantWebhook(c, extReq, db, requestBody)
	default:
		err = providers.HandleDefaultMerchantWebhook(c, extReq, db, requestBody)
	}

	if err != nil {
		extReq.Logger.Error(fmt.Sprintf("webhook log error for %v %v", provider, err.Error()))
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func GetProvider(c *gin.Context, requestBody []byte) (string, error) {
	var (
		provider string
	)

	if utility.GetHeader(c, "verif-hash") != "" {
		provider = "rave"
	}

	// rave, e-transact, paystack
	return provider, nil
}

func logWebhookData(extReq request.ExternalRequest, db postgresql.Databases, provider string, requestBody []byte) error {
	extReq.Logger.Info(fmt.Sprintf("webhook log info for %v %v", provider, string(requestBody)))
	webhookLog := models.WebhookLog{
		Log:      string(requestBody),
		Provider: provider,
	}
	err := webhookLog.CreateWebhookLog(db.Payment)
	if err != nil {
		return err
	}
	return nil
}
