package mor

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/vesicash/mor-api/external/external_models"
	"github.com/vesicash/mor-api/external/request"
	"github.com/vesicash/mor-api/internal/config"
	"github.com/vesicash/mor-api/internal/models"
	"github.com/vesicash/mor-api/pkg/repository/storage/postgresql"
	"github.com/vesicash/mor-api/services"
)

func RequestWithdrawalService(extReq request.ExternalRequest, db postgresql.Databases, user external_models.User, req models.RequestWithdrawalRequest) (int, error) {
	var (
		withdrawal            = models.Withdrawal{MerchantID: int64(user.AccountID), Status: models.TransactionPending}
		withdrawalSum float64 = 0
	)

	checkTime := int(time.Now().Add(8736 * time.Hour).Unix())
	if req.WithdrawalDate > checkTime {
		return http.StatusBadRequest, fmt.Errorf("invalid timestamp, time must not be more than 1 year after today")
	}

	req.Currency = strings.ToUpper(strings.ReplaceAll(strings.ToUpper(req.Currency), "MOR_", ""))
	req.Currency = strings.ToUpper(strings.ReplaceAll(strings.ToUpper(req.Currency), "ESCROW_", ""))
	morWallet := strings.ToUpper(fmt.Sprintf("MOR_%v", req.Currency))

	wallet, err := services.GetWalletBalanceByAccountIdAndCurrency(extReq, int(user.AccountID), morWallet)
	if err != nil {
		return http.StatusBadRequest, err
	}

	withdrawals, _, err := withdrawal.GetWithdrawals(db.MOR, nil, []int{}, 0, 0)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	for _, w := range withdrawals {
		withdrawalSum += w.Amount
	}

	remainingBalance := wallet.Available - withdrawalSum
	if remainingBalance < req.Amount {
		return http.StatusBadRequest, fmt.Errorf("insufficient wallet balance")
	}

	withdrawal.Amount = req.Amount
	withdrawal.Currency = req.Currency
	withdrawal.WithdrawalDate = time.Unix(int64(req.WithdrawalDate), 0)

	err = withdrawal.CreateWithdrawal(db.MOR)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = SlackNotify(extReq, config.GetConfig().Slack.WithdrawalChannelID, `
	MOR WITHDRAWAL REQUEST FROM (`+strconv.Itoa(int(user.AccountID))+`) `+fmt.Sprintf("%v %v", user.Lastname, user.Firstname)+`
	Currency: `+withdrawal.Currency+`
	Amount: `+fmt.Sprintf("%v", withdrawal.Amount)+`
	Withdrawal Date: `+fmt.Sprintf("%v", withdrawal.WithdrawalDate)+`
	`)
	if err != nil && !extReq.Test {
		extReq.Logger.Error("error sending notification to slack: ", err.Error())
	}

	return http.StatusOK, nil
}

func GetWithdrawalsService(extReq request.ExternalRequest, db postgresql.Databases, paginator postgresql.Pagination, req models.GetWithdrawalRequest) ([]models.Withdrawal, postgresql.PaginationResponse, int, error) {
	var (
		withdrawal = models.Withdrawal{Status: models.TransactionPending}
	)

	if req.CurrencyFilter != "" {
		req.CurrencyFilter = strings.ToUpper(strings.ReplaceAll(strings.ToUpper(req.CurrencyFilter), "MOR_", ""))
		req.CurrencyFilter = strings.ToUpper(strings.ReplaceAll(strings.ToUpper(req.CurrencyFilter), "ESCROW_", ""))
		withdrawal.Currency = req.CurrencyFilter
	}

	if req.Status != "" {
		withdrawal.Status = models.TransactionStatus(req.Status)
	}

	usersIDs := []int{}
	if req.Search != "" {
		users, _ := services.GetUsers(extReq, true, req.Search)
		for _, u := range users {
			usersIDs = append(usersIDs, int(u.AccountID))
		}
	}

	withdrawals, pagination, err := withdrawal.GetWithdrawals(db.MOR, &paginator, usersIDs, req.FromTime, req.ToTime)
	if err != nil {
		return withdrawals, pagination, http.StatusInternalServerError, err
	}

	withdrawals, err = GetMorWithdrawalsDetails(extReq, db, withdrawals)
	if err != nil {
		return withdrawals, pagination, http.StatusInternalServerError, err
	}

	return withdrawals, pagination, http.StatusOK, nil
}

func CompleteWithdrawalService(extReq request.ExternalRequest, db postgresql.Databases, withdrawalID int) (int, error) {
	var (
		withdrawal = models.Withdrawal{ID: uint(withdrawalID)}
	)

	code, err := withdrawal.GetWithdrawalByID(db.MOR)
	if err != nil {
		return code, err
	}

	_, err = services.DebitWallet(extReq, db, withdrawal.Amount, withdrawal.Currency, int(withdrawal.MerchantID), "no", "yes", "")
	if err != nil {
		return http.StatusInternalServerError, err
	}

	withdrawal.Status = models.TransactionSuccessful
	err = withdrawal.UpdateAllFields(db.MOR)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func GetMorWithdrawalsDetails(extReq request.ExternalRequest, db postgresql.Databases, withdrawals []models.Withdrawal) ([]models.Withdrawal, error) {

	type withdrawalAndError struct {
		Withdrawal models.Withdrawal
		Err        error
	}

	var newWithdrawals []models.Withdrawal
	var errs []string
	var wg sync.WaitGroup
	results := make(chan withdrawalAndError, len(withdrawals))

	// Loop through the data slice and spawn a goroutine for each item.
	for _, withdrawal := range withdrawals {
		wg.Add(1)
		go func(extReq request.ExternalRequest, db postgresql.Databases, withdrawal models.Withdrawal, wg *sync.WaitGroup, results chan withdrawalAndError) {
			defer wg.Done()
			user, err := services.GetUserWithAccountID(extReq, int(withdrawal.MerchantID))
			withdrawal.Merchant = user

			results <- withdrawalAndError{
				Withdrawal: withdrawal,
				Err:        err,
			}

		}(extReq, db, withdrawal, &wg, results)
	}

	wg.Wait()
	close(results)

	// Collect the results from the channel and append them to the processedData slice.
	for result := range results {
		if result.Err != nil {
			errs = append(errs, result.Err.Error())
		} else {
			newWithdrawals = append(newWithdrawals, result.Withdrawal)
		}
	}

	if len(errs) > 0 {
		extReq.Logger.Error(fmt.Sprintf("error getting user details: %v", strings.Join(errs, ", ")))
	}

	return newWithdrawals, nil
}
