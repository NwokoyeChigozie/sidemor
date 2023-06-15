package mor

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/vesicash/mor-api/external/request"
	"github.com/vesicash/mor-api/internal/models"
	"github.com/vesicash/mor-api/pkg/repository/storage/postgresql"
	"github.com/vesicash/mor-api/services"
)

func GetVerificationSettingsService(extReq request.ExternalRequest, db postgresql.Databases, paginator postgresql.Pagination, req models.GetSettingsRequest) ([]models.Setting, postgresql.PaginationResponse, int, error) {
	var (
		setting       = models.Setting{}
		isVerified    *bool
		trueD, falseD = true, false
	)

	if strings.EqualFold(req.Status, "true") {
		isVerified = &trueD
	} else if strings.EqualFold(req.Status, "false") {
		isVerified = &falseD
	}

	usersIDs := []int{}
	if req.Search != "" {
		users, _ := services.GetUsers(extReq, true, req.Search)
		for _, u := range users {
			usersIDs = append(usersIDs, int(u.AccountID))
		}
	}

	settings, pagination, err := setting.GetSettings(db.MOR, paginator, usersIDs, req.FromTime, req.ToTime, isVerified)
	if err != nil {
		return settings, pagination, http.StatusInternalServerError, err
	}

	settings, err = GetMorSettingsDetails(extReq, db, settings)
	if err != nil {
		return settings, pagination, http.StatusInternalServerError, err
	}

	return settings, pagination, http.StatusOK, nil
}

func GetMorSettingDetails(extReq request.ExternalRequest, db postgresql.Databases, setting models.Setting) (models.Setting, error) {
	user, err := services.GetUserWithAccountID(extReq, int(setting.AccountID))
	if err != nil {
		return setting, err
	}

	setting.Email = user.EmailAddress
	setting.FullName = user.Lastname + " " + user.Firstname
	setting.AccountType = user.AccountType

	return setting, nil
}

func GetMorSettingsDetails(extReq request.ExternalRequest, db postgresql.Databases, settings []models.Setting) ([]models.Setting, error) {

	type settingAndError struct {
		Setting models.Setting
		Err     error
	}

	var newSettings []models.Setting
	var errs []string
	var wg sync.WaitGroup
	results := make(chan settingAndError, len(settings))

	// Loop through the data slice and spawn a goroutine for each item.
	for _, setting := range settings {
		wg.Add(1)
		go func(extReq request.ExternalRequest, db postgresql.Databases, setting models.Setting, wg *sync.WaitGroup, results chan settingAndError) {
			defer wg.Done()
			setting, err := GetMorSettingDetails(extReq, db, setting)
			results <- settingAndError{
				Setting: setting,
				Err:     err,
			}

		}(extReq, db, setting, &wg, results)
	}

	wg.Wait()
	close(results)

	// Collect the results from the channel and append them to the processedData slice.
	for result := range results {
		if result.Err != nil {
			errs = append(errs, result.Err.Error())
		} else {
			newSettings = append(newSettings, result.Setting)
		}
	}

	if len(errs) > 0 {
		extReq.Logger.Error(fmt.Sprintf("error getting mor transaction details: %v", strings.Join(errs, ", ")))
	}

	return newSettings, nil
}
