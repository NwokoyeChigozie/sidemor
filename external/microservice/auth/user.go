package auth

import (
	"fmt"
	"github.com/vesicash/mor-api/external/external_models"
	"github.com/vesicash/mor-api/internal/config"
)

func (r *RequestObj) ToggleMORStatus() (interface{}, error) {
	var (
		outBoundResponse external_models.ToggleMORStatusResponse
		logger           = r.Logger
		idata            = r.RequestData
		appKey           = config.GetConfig().App.Key
	)
	data, ok := idata.(external_models.ToggleMORStatusReq)
	if !ok {
		logger.Error("toggle mor status", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("toggle mor status", data)
	err := r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("toggle mor status", outBoundResponse, err.Error())
		return outBoundResponse.Data, err
	}
	logger.Info("toggle mor status", outBoundResponse)

	return outBoundResponse.Data, nil
}
