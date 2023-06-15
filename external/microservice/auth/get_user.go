package auth

import (
	"fmt"

	"github.com/vesicash/mor-api/external/external_models"
	"github.com/vesicash/mor-api/internal/config"
)

func (r *RequestObj) GetUser() (external_models.User, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.GetUserModel
		logger           = r.Logger
		idata            = r.RequestData
	)

	data, ok := idata.(external_models.GetUserRequestModel)
	if !ok {
		logger.Error("get user", idata, "request data format error")
		return outBoundResponse.Data.User, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("get user", data)
	err := r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("get user", outBoundResponse, err.Error())
		return outBoundResponse.Data.User, err
	}
	logger.Info("get user", outBoundResponse)

	return outBoundResponse.Data.User, nil
}

func (r *RequestObj) GetUsers() ([]external_models.User, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.GetUsersResponse
		logger           = r.Logger
		idata            = r.RequestData
	)

	data, ok := idata.(external_models.GetUsersRequest)
	if !ok {
		logger.Error("get users", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("get user", data)
	queryParams := "?is_mor_enabled=false"
	if data.IsMorEnabled {
		queryParams = "?is_mor_enabled=true"
	}

	if data.Search != "" {
		queryParams += fmt.Sprintf("&search=%s", data.Search)
	}

	err := r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("get users", outBoundResponse, err.Error())
		return outBoundResponse.Data, err
	}
	logger.Info("get user", outBoundResponse)

	return outBoundResponse.Data, nil
}

func (r *RequestObj) GetUsersByBusinessID() ([]external_models.User, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.GetUsersByBusinessIDModel
		logger           = r.Logger
		idata            = r.RequestData
	)

	data, ok := idata.(string)
	if !ok {
		logger.Error("get users by business id", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("get users by business id", data)
	err := r.getNewSendRequestObject(data, headers, "/"+data).SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("get users by business id", outBoundResponse, err.Error())
		return outBoundResponse.Data, err
	}
	logger.Info("get users by business id", outBoundResponse)

	return outBoundResponse.Data, nil
}

func (r *RequestObj) SetUserAuthorizationRequiredStatus() (bool, error) {

	var (
		appKey           = config.GetConfig().App.Key
		outBoundResponse external_models.SetUserAuthorizationRequiredStatusResponse
		logger           = r.Logger
		idata            = r.RequestData
	)

	data, ok := idata.(external_models.SetUserAuthorizationRequiredStatusModel)
	if !ok {
		logger.Error("set user authorization required status", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"v-app":        appKey,
	}

	logger.Info("set user authorization required status", data)
	err := r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Error("set user authorization required status", outBoundResponse, err.Error())
		return outBoundResponse.Data, err
	}
	logger.Info("set user authorization required status", outBoundResponse)

	return outBoundResponse.Data, nil
}
