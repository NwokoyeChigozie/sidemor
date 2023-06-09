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

func GetPayoutService(extReq request.ExternalRequest, db postgresql.Databases, payoutID int) (models.Payout, int, error) {
	var (
		payout = models.Payout{ID: uint(payoutID)}
	)

	code, err := payout.GetPayoutByID(db.MOR)
	if err != nil {
		return payout, code, err
	}

	payout, err = GetMorPayoutDetails(extReq, db, payout)
	if err != nil {
		return payout, http.StatusInternalServerError, err
	}

	return payout, http.StatusOK, nil
}

func GetPayoutsService(extReq request.ExternalRequest, db postgresql.Databases, paginator postgresql.Pagination, req models.GetPayoutRequest) ([]models.Payout, postgresql.PaginationResponse, int, error) {
	var (
		payout = models.Payout{}
	)

	if req.CurrencyFilter != "" {
		country, _ := services.GetCountryByCurrency(extReq, extReq.Logger, strings.ToUpper(req.CurrencyFilter))
		payout.CountryID = int64(country.ID)
	}

	if req.Status != "" {
		payout.Status = models.TransactionStatus(req.Status)
	}

	payouts, pagination, err := payout.GetPayouts(db.MOR, paginator, req.Search, req.FromTime, req.ToTime)
	if err != nil {
		return payouts, pagination, http.StatusInternalServerError, err
	}

	payouts, err = GetMorPayoutsDetails(extReq, db, payouts)
	if err != nil {
		return payouts, pagination, http.StatusInternalServerError, err
	}

	return payouts, pagination, http.StatusOK, nil
}

func GetMorPayoutDetails(extReq request.ExternalRequest, db postgresql.Databases, payout models.Payout) (models.Payout, error) {
	user, err := services.GetUserWithAccountID(extReq, int(payout.MerchantID))
	if err != nil {
		return payout, err
	}

	country, err := services.GetCountryByID(extReq, extReq.Logger, int(payout.CountryID))
	if err != nil {
		return payout, err
	}

	payout.MerchantName = user.Lastname + " " + user.Firstname
	payout.MerchantEmail = user.EmailAddress
	payout.Currency = country.CurrencyCode

	return payout, nil
}

func GetMorPayoutsDetails(extReq request.ExternalRequest, db postgresql.Databases, payouts []models.Payout) ([]models.Payout, error) {

	type payoutAndError struct {
		Payout models.Payout
		Err    error
	}

	var newPayouts []models.Payout
	var errs []string
	var wg sync.WaitGroup
	results := make(chan payoutAndError, len(payouts))

	// Loop through the data slice and spawn a goroutine for each item.
	for _, payout := range payouts {
		wg.Add(1)
		go func(extReq request.ExternalRequest, db postgresql.Databases, payout models.Payout, wg *sync.WaitGroup, results chan payoutAndError) {
			defer wg.Done()
			payout, err := GetMorPayoutDetails(extReq, db, payout)
			results <- payoutAndError{
				Payout: payout,
				Err:    err,
			}

		}(extReq, db, payout, &wg, results)
	}

	wg.Wait()
	close(results)

	// Collect the results from the channel and append them to the processedData slice.
	for result := range results {
		if result.Err != nil {
			errs = append(errs, result.Err.Error())
		} else {
			newPayouts = append(newPayouts, result.Payout)
		}
	}

	if len(errs) > 0 {
		extReq.Logger.Error(fmt.Sprintf("error getting mor payouts details: %v", strings.Join(errs, ", ")))
	}

	return newPayouts, nil
}
