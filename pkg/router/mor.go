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

	morAuthUrl := r.Group(fmt.Sprintf("%v", ApiVersion), middleware.Authorize(db, extReq, middleware.AuthType))
	{
		morAuthUrl.GET("/customers", mor.GetCustomers)
		morAuthUrl.GET("/transactions/get", mor.GetMerchantTransactions)
		morAuthUrl.GET("/transactions/summary/:account_id", mor.GetMerchantTransactionsSummary)
		morAuthUrl.GET("/payouts/get", mor.GetMerchantPayouts)
		morAuthUrl.POST("/withdrawal/request", mor.RequestWithdrawal)
	}

	morSettingsAuthUrl := r.Group(fmt.Sprintf("%v/settings", ApiVersion), middleware.Authorize(db, extReq, middleware.AuthType))
	{
		morSettingsAuthUrl.GET("/get", mor.GetSettings)
		morSettingsAuthUrl.POST("/save", mor.SaveSettings)
		morSettingsAuthUrl.POST("/payment-methods/:action", mor.EnableOrDisablePaymentMethods)
		morSettingsAuthUrl.POST("/wallets/:action", mor.AddRemoveOrGetWallets)
	}

	paymentBusinessAdminUrl := r.Group(fmt.Sprintf("%v/admin", ApiVersion), middleware.Authorize(db, extReq, middleware.BusinessAdmin))
	{
		paymentBusinessAdminUrl.POST("/transaction/record", mor.RecordTransaction)
		paymentBusinessAdminUrl.GET("/transaction/get/:id", mor.GetTransaction)
		paymentBusinessAdminUrl.GET("/transactions/get", mor.GetTransactions)
		paymentBusinessAdminUrl.GET("/transactions/summary", mor.GetTransactionsSummary)
		paymentBusinessAdminUrl.GET("/transactions/summary/:account_id", mor.GetMerchantTransactionsSummary)

		paymentBusinessAdminUrl.GET("/settings/get", mor.GetVerificationSettings)
		paymentBusinessAdminUrl.POST("/settings/:id/document", mor.UpdateDocumentStatus)

		paymentBusinessAdminUrl.GET("/payout/get/:id", mor.GetPayout)
		paymentBusinessAdminUrl.GET("/payouts/get", mor.GetPayouts)
		paymentBusinessAdminUrl.POST("/payout/to-wallet", mor.PayOutToWallets)

		paymentBusinessAdminUrl.GET("/withdrawal/get-all", mor.GetWithdrawals)
		paymentBusinessAdminUrl.PATCH("/withdrawal/complete/:withdrawal_id", mor.CompleteWithdrawal)
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
