package mor

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/vesicash/mor-api/external/external_models"
	"github.com/vesicash/mor-api/external/request"
	"github.com/vesicash/mor-api/internal/models"
	"github.com/vesicash/mor-api/pkg/repository/storage/postgresql"
	"github.com/vesicash/mor-api/services"
	"github.com/vesicash/mor-api/utility"
)

func SaveSettingsService(extReq request.ExternalRequest, db postgresql.Databases, user external_models.User, req models.SaveSettingsRequest) (models.Setting, int, error) {
	var (
		setting          = models.Setting{AccountID: int64(user.AccountID)}
		prevCountries    = []int{}
		verificationsMap = map[int]models.SettingsVerification{}
		verifications    = []models.SettingsVerification{}
	)

	code, err := setting.GetSettingByAccountID(db.MOR)
	if err != nil {
		if code == http.StatusInternalServerError {
			return models.Setting{}, code, err
		}

		err := setting.CreateSetting(db.MOR)
		if err != nil {
			return models.Setting{}, http.StatusInternalServerError, err
		}
	}

	for _, v := range setting.Countries {
		prevCountries = append(prevCountries, int(v.ID))
	}

	for _, v := range req.Countries {
		if !utility.InIntSlice(v, prevCountries) {
			country, err := services.GetCountryByID(extReq, extReq.Logger, v)
			if err != nil {
				return models.Setting{}, http.StatusBadRequest, fmt.Errorf("country with id:%v not found, %v", v, err.Error())
			}
			setting.Countries = append(setting.Countries, models.SettingsCountries{
				ID:           uint(v),
				Name:         country.Name,
				CurrencyCode: country.CurrencyCode,
			})
		}
	}

	for _, v := range req.WalletCurrencyCodes {
		v = strings.ToUpper(v)
		if !utility.InStringSlice(v, setting.CurrencyCodes) {
			setting.CurrencyCodes = append(setting.CurrencyCodes, v)
		}
	}

	setting.BusinessTypeID = req.BusinessTypeID
	setting.UsageType = req.UsageType

	for _, v := range setting.Verifications {
		verificationsMap[int(v.CountryID)] = v
	}

	for _, d := range req.Documents {
		if existingDocument, ok := verificationsMap[int(d.CountryID)]; ok {
			existingDocument.DocumentUrl = d.DocumentUrl
			existingDocument.Status = models.VerificationPending
		} else {
			verificationsMap[int(d.CountryID)] = models.SettingsVerification{
				CountryID:   d.CountryID,
				DocumentUrl: d.DocumentUrl,
				Status:      models.VerificationPending,
			}
		}
	}

	for _, v := range verificationsMap {
		verifications = append(verifications, v)
	}

	setting.Verifications = verifications

	err = setting.UpdateAllFields(db.MOR)
	if err != nil {
		return models.Setting{}, http.StatusInternalServerError, err
	}

	return setting, http.StatusOK, nil
}

func GetSettingsService(extReq request.ExternalRequest, db postgresql.Databases, user external_models.User) (*models.Setting, int, error) {
	var (
		setting = models.Setting{AccountID: int64(user.AccountID)}
	)

	code, err := setting.GetSettingByAccountID(db.MOR)
	if err != nil {
		if code == http.StatusInternalServerError {
			return &models.Setting{}, code, err
		}

		return nil, http.StatusOK, nil
	}

	return &setting, http.StatusOK, nil
}

func EnableOrDisablePaymentMethodsService(extReq request.ExternalRequest, db postgresql.Databases, user external_models.User, action string, req models.EnableOrDisablePaymentMethodsRequest) (models.Setting, int, error) {
	var (
		setting = models.Setting{AccountID: int64(user.AccountID)}
	)

	code, err := setting.GetSettingByAccountID(db.MOR)
	if err != nil {
		if code == http.StatusInternalServerError {
			return models.Setting{}, code, err
		}

		err := setting.CreateSetting(db.MOR)
		if err != nil {
			return models.Setting{}, http.StatusInternalServerError, err
		}
	}

	paymentMethods := setting.PaymentMethods
	if strings.EqualFold(action, "enable") {
		for _, p := range req.Methods {
			if !p.In(paymentMethods) {
				paymentMethods = append(paymentMethods, p)
			}
		}
	} else if strings.EqualFold(action, "disable") {
		newPaymentMethods := []models.PaymentMethod{}
		for _, v := range setting.PaymentMethods {
			if !v.In(req.Methods) {
				newPaymentMethods = append(newPaymentMethods, v)
			}
		}
		paymentMethods = newPaymentMethods
	}

	setting.PaymentMethods = paymentMethods
	err = setting.UpdateAllFields(db.MOR)
	if err != nil {
		return setting, http.StatusInternalServerError, err
	}

	return setting, http.StatusOK, nil
}
