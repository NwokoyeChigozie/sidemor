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
		services.GetOrCreateWalletBalance(extReq, int(user.AccountID), fmt.Sprintf("MOR_%v", v))
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

func UpdateDocumentStatusService(extReq request.ExternalRequest, db postgresql.Databases, settingsID int, req models.UpdateDocumentStatusRequest) (models.Setting, int, error) {
	var (
		setting = models.Setting{}
	)

	if code, err := setting.GetSettingByID(db.MOR, settingsID); err != nil {
		if code == http.StatusInternalServerError {
			return models.Setting{}, code, err
		}
	}

	err := setting.UpdateVerificationSettings(db.MOR, req)

	if err != nil {
		return setting, http.StatusInternalServerError, err
	}
	// loop through verification and check if all of them have been approved, if yes, call toggle-mor-status endpoint
	var isAllVerified bool = true
	if len(setting.Verifications) < 1 {
		isAllVerified = false
	} else {
		for _, v := range setting.Verifications {
			if v.Status != models.Verified {
				isAllVerified = false
			}
		}
	}

	if isAllVerified {
		setting.IsVerified = true
		err = setting.UpdateAllFields(db.MOR)
		if err != nil {
			return setting, http.StatusInternalServerError, err
		}

		err = services.ToggleMORStatus(extReq, uint(setting.AccountID), true)
		if err != nil {
			return setting, http.StatusInternalServerError, err
		}
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

func AddRemoveOrGetWalletsService(extReq request.ExternalRequest, db postgresql.Databases, user external_models.User, action string, req models.AddRemoveOrGetWalletsRequest) (map[string]external_models.WalletBalance, int, error) {
	var (
		setting = models.Setting{AccountID: int64(user.AccountID)}
	)

	code, err := setting.GetSettingByAccountID(db.MOR)
	if err != nil {
		if code == http.StatusInternalServerError {
			return map[string]external_models.WalletBalance{}, code, err
		}

		err := setting.CreateSetting(db.MOR)
		if err != nil {
			return map[string]external_models.WalletBalance{}, http.StatusInternalServerError, err
		}
	}

	fmt.Println(action)
	currencies := setting.CurrencyCodes
	if strings.EqualFold(action, "add") {
		currencies, err = handleAddCurrencies(extReq, db, int(user.AccountID), currencies, req.CurrencyCodes)
	} else if strings.EqualFold(action, "delete") {
		currencies, err = handleRemoveCurrencies(extReq, db, int(user.AccountID), currencies, req.CurrencyCodes)
	} else if strings.EqualFold(action, "get") {
		if len(req.CurrencyCodes) > 0 {
			// morCurrencies := []string{}
			var wallets = map[string]external_models.WalletBalance{}
			for _, c := range req.CurrencyCodes {
				wC := strings.ToUpper(c)
				w, err := services.GetOrCreateWalletBalance(extReq, int(user.AccountID), wC)
				if err == nil {
					wallets[w.Currency] = w
				}
				// morCurrencies = append(morCurrencies, fmt.Sprintf("MOR_%s", wC))
			}
			// wallets, err := services.GetWalletBalancesByCurrencies(extReq, db, int(user.AccountID), morCurrencies)
			// if err != nil {
			// 	return wallets, http.StatusInternalServerError, err
			// }
			return wallets, http.StatusOK, nil
		}
	}

	if err != nil {
		return map[string]external_models.WalletBalance{}, http.StatusInternalServerError, err
	}

	setting.CurrencyCodes = currencies
	err = setting.UpdateAllFields(db.MOR)
	if err != nil {
		return map[string]external_models.WalletBalance{}, http.StatusInternalServerError, err
	}

	morCurrencies := []string{}
	for _, c := range currencies {
		morCurrencies = append(morCurrencies, fmt.Sprintf("MOR_%s", strings.ToUpper(c)))
	}
	wallets, err := services.GetWalletBalancesByCurrencies(extReq, db, int(user.AccountID), morCurrencies)
	if err != nil {
		return wallets, http.StatusInternalServerError, err
	}

	return wallets, http.StatusOK, nil
}

func handleAddCurrencies(extReq request.ExternalRequest, db postgresql.Databases, accountID int, availableCurrencies []string, newCurrencies []string) ([]string, error) {
	fmt.Println("started adding", newCurrencies)
	for _, c := range newCurrencies {
		c = strings.ToUpper(c)
		fmt.Println("started adding", c)
		morCurrency := fmt.Sprintf("MOR_%v", c)
		fmt.Println("mor currency", morCurrency)
		if !utility.InStringSlice(c, availableCurrencies) {
			fmt.Println("not in currency", availableCurrencies)
			_, err := services.CreateWalletBalance(extReq, accountID, morCurrency, 0)
			if err != nil {
				return availableCurrencies, err
			}

			availableCurrencies = append(availableCurrencies, c)
		} else {
			_, err := services.GetWalletBalanceByAccountIdAndCurrency(extReq, accountID, morCurrency)
			if err != nil {
				_, err := services.CreateWalletBalance(extReq, accountID, morCurrency, 0)
				if err != nil {
					return availableCurrencies, err
				}
			}
		}
	}
	return availableCurrencies, nil
}

func handleRemoveCurrencies(extReq request.ExternalRequest, db postgresql.Databases, accountID int, availableCurrencies []string, removeCurrencies []string) ([]string, error) {
	for _, c := range removeCurrencies {
		c = strings.ToUpper(c)
		creditWallet := "MOR_USD"
		morCurrency := fmt.Sprintf("MOR_%v", c)
		if utility.InStringSlice(c, availableCurrencies) {
			if c == "USD" {
				continue
			}

			walletBalance, _ := services.GetWalletBalanceByAccountIdAndCurrency(extReq, accountID, morCurrency)
			initialBalance := walletBalance.Available

			if initialBalance == 0 {
				availableCurrencies = utility.RemoveString(availableCurrencies, c)
				continue
			}

			rate, err := services.GetRateByCurrencies(extReq, c, "USD")
			if err != nil {
				extReq.Logger.Error(fmt.Sprintf("error getting rate for currencies %v -> %v: %v", c, "USD", err.Error()))
				continue
			}

			if rate.ID <= 0 {
				continue
			}

			var multiplier float64 = 0
			if rate.InitialAmount > 0 {
				multiplier = rate.Amount / rate.InitialAmount
			}

			convertedBalance := multiplier * initialBalance

			err = services.CreateExchangeTransaction(extReq, accountID, int(rate.ID), initialBalance, convertedBalance, "completed")
			if err != nil {
				return availableCurrencies, err
			}

			_, err = services.DebitWallet(extReq, db, initialBalance, morCurrency, accountID, "no", "no", "")
			if err != nil {
				return availableCurrencies, err
			}

			_, err = services.CreditWallet(extReq, db, convertedBalance, creditWallet, accountID, false, "no", "no", "")
			if err != nil {
				return availableCurrencies, err
			}

			availableCurrencies = utility.RemoveString(availableCurrencies, c)
		}
	}

	return availableCurrencies, nil
}
