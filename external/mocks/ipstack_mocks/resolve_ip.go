package ipstack_mocks

import (
	"fmt"

	"github.com/vesicash/mor-api/external/external_models"
	"github.com/vesicash/mor-api/internal/config"
	"github.com/vesicash/mor-api/utility"
)

func IpstackResolveIp(logger *utility.Logger, idata interface{}) (external_models.IPStackResolveIPResponse, error) {

	var (
		key              = config.GetConfig().IPStack.Key
		outBoundResponse external_models.IPStackResolveIPResponse
	)

	ip, ok := idata.(string)
	if !ok {
		logger.Error("ipstack resolve ip", idata, "request data format error")
		return outBoundResponse, fmt.Errorf("request data format error")
	}
	outBoundResponse.Ip = ip
	outBoundResponse.City = "city"
	outBoundResponse.CountryName = "name"

	path := "/" + ip + "?access_key=" + key

	logger.Info("ipstack resolve ip", ip, path)

	return outBoundResponse, nil
}
