package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/models"
)

type PayoutInterface interface {
	CreatePayoutTransaction(echo.Context) error
	GetAllPayoutTransactions(echo.Context) ([]models.GetAllPayoutTransactionsResponseModel, error)
	GetPayoutTransactionsByRetailerId(echo.Context) ([]models.GetRetailerPayoutTransactionsResponseModel, error)
}

type payoutRepository struct {
	db *database.Database
}

func NewPayoutRepository(db *database.Database) *payoutRepository {
	return &payoutRepository{
		db,
	}
}

func (pr *payoutRepository) CreatePayoutTransaction(c echo.Context) error {
	var req models.CreatePayoutRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()

	commision, err := pr.db.GetPayoutCommisionQuery(ctx, req.RetailerId, req.Amount)
	if err != nil {
		return err
	}

	if err := pr.db.VerifyRetailerForPayoutTransactionQuery(ctx, req.RetailerId, req.Amount+commision.TotalCommision); err != nil {
		return err
	}
	req.PartnerRequestId = uuid.NewString()

	apiUrl := `https://v2bapi.rechargkit.biz/rkitpayout/payoutTransfer`
	reqBody, err := json.Marshal(map[string]any{
		"mobile_no":          req.MobileNumber,
		"account_no":         req.AccountNumber,
		"ifsc":               req.IFSCCode,
		"bank_name":          req.BankName,
		"beneficiary_name":   req.BeneficiaryName,
		"amount":             req.Amount,
		"transfer_type":      req.TransferType,
		"partner_request_id": req.PartnerRequestId,
	})
	if err != nil {
		return err
	}

	apiRequest, err := http.NewRequest(
		http.MethodPost,
		apiUrl,
		bytes.NewReader(reqBody),
	)
	if err != nil {
		return err
	}

	apiRequest.Header.Set("Content-Type", "application/json")
	apiRequest.Header.Set("Authorization", "Bearer "+os.Getenv("RKIT_API_TOKEN"))

	client := &http.Client{Timeout: 20 * time.Second}

	resp, err := client.Do(apiRequest)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var apiResponse struct {
		Error                 int    `json:"error"`
		Message               string `json:"msg"`
		Status                int    `json:"status"`
		OrderId               string `json:"orderid"`
		OperatorTransactionId string `json:"optransid"`
		PartnerRequestId      string `json:"partnerreqid"`
	}

	if err := json.Unmarshal(respBytes, &apiResponse); err != nil {
		return err
	}
	req.OrderId = apiResponse.OrderId
	req.OperatorTransactionId = apiResponse.OperatorTransactionId

	if apiResponse.Status == 1 {
		req.TransactionStatus = "SUCCESS"
		if err := pr.db.CreatePayoutSuccessOrPendingQuery(ctx, req, *commision); err != nil {
			return err
		}
		return nil
	}

	if apiResponse.Status == 2 {
		req.TransactionStatus = "PENDING"
		if err := pr.db.CreatePayoutSuccessOrPendingQuery(ctx, req, *commision); err != nil {
			return err
		}
		return nil
	}

	if apiResponse.Status == 3 {
		req.TransactionStatus = "FAILED"
		if err := pr.db.CreatePayoutFailureQuery(ctx, req, *commision); err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("invalid status from recharge kit")
}

func (pr *payoutRepository) GetAllPayoutTransactions(c echo.Context) ([]models.GetAllPayoutTransactionsResponseModel, error) {
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()
	limit, offset := parsePagination(c)
	return pr.db.GetAllPayoutTransactionsQuery(ctx, limit, offset)
}

func (pr *payoutRepository) GetPayoutTransactionsByRetailerId(c echo.Context) ([]models.GetRetailerPayoutTransactionsResponseModel, error) {
	var retailerId = c.Param("retailer_id")
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()
	limit, offset := parsePagination(c)
	return pr.db.GetPayoutTransactionsByRetailerIdQuery(ctx, retailerId, limit, offset)
}
