package notification_mocks

import (
	"fmt"

	"github.com/vesicash/mor-api/external/external_models"
	"github.com/vesicash/mor-api/utility"
)

func SendAuthorizedNotification(logger *utility.Logger, idata interface{}) (interface{}, error) {
	data, ok := idata.(external_models.AuthorizeNotificationRequest)
	if !ok {
		logger.Error("authorized notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}

	logger.Info("authorized notification", data)

	return nil, nil
}

func SendAuthorizationNotification(logger *utility.Logger, idata interface{}) (interface{}, error) {

	data, ok := idata.(external_models.AuthorizeNotificationRequest)
	if !ok {
		logger.Error("authorization notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}

	logger.Info("verification failed notification", data)

	return nil, nil
}
