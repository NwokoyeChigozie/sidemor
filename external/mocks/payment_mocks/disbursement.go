package payment_mocks

import (
	"fmt"

	"github.com/vesicash/mor-api/external/external_models"
	"github.com/vesicash/mor-api/utility"
)

func RequestManualRefund(logger *utility.Logger, idata interface{}) (map[string]interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
	)
	data, ok := idata.(string)
	if !ok {
		logger.Error("request manual refund", idata, "request data format error")
		return outBoundResponse, fmt.Errorf("request data format error")
	}

	reqData := external_models.TransactionIDRequestModel{
		TransactionId: data,
	}

	logger.Info("request manual refund", reqData)

	return outBoundResponse, nil
}
