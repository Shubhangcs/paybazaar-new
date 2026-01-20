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

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/models"
	"github.com/levion-studio/paybazaar/pkg"
)

type PayoutInterface interface {
	CreatePayout(echo.Context) error
}

type payoutRepository struct {
	db       *database.Database
	jwtUtils *pkg.JwtUtils
}

func NewPayoutRepository(db *database.Database, jwtUtils *pkg.JwtUtils) *payoutRepository {
	return &payoutRepository{
		db,
		jwtUtils,
	}
}

func (pr *payoutRepository) CreatePayout(c echo.Context) error {
	var payoutRequest models.CreatePayoutRequestModel
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()
	commision, err := pr.db.GetPayoutCommisionQuery(ctx, payoutRequest.RetailerID, "PAYOUT")

	isValid, err := pr.db.ValidateRequestQuery(ctx, payoutRequest, commision.RetailerCommision)
	if err != nil {
		return err
	}

	if !isValid {
		return fmt.Errorf("incorrect mpin or kyc status or insufficient balance")
	}

	if payoutRequest.Amount < 1000 || payoutRequest.Amount > 25000 {
		return fmt.Errorf("invalid amount")
	}

	// ---------------- API CALL ----------------

	apiUrl := `https://v2bapi.rechargkit.biz/rkitpayout/payoutTransfer`

	reqBody, err := json.Marshal(map[string]any{
		"mobile_no":          payoutRequest.MobileNumber,
		"account_number":     payoutRequest.BeneficiaryAccountNumber,
		"ifsc":               payoutRequest.BeneficiaryIFSCCode,
		"bank_name":          payoutRequest.BeneficiaryBankName,
		"beneficiary_name":   payoutRequest.BeneficiaryName,
		"amount":             payoutRequest.Amount,
		"partner_request_id": payoutRequest.PartnerRequestID,
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

	var res models.PayoutAPIResponseModel
	if err := json.Unmarshal(respBytes, &res); err != nil {
		return err
	}

	// Basic response sanity check
	if res.Status == 0 {
		fmt.Println(res)
		return fmt.Errorf("invalid payout gateway response")
	}

	if res.Status == 1 {
		return pr.db.PayoutPendingOrSuccessQuery(ctx, payoutRequest, *commision, payoutRequest.AdminID, "SUCCESS")
	}

	if res.Status == 2 {
		return pr.db.PayoutPendingOrSuccessQuery(ctx, payoutRequest, *commision, payoutRequest.AdminID, "PENDING")
	}

	if res.Status == 3 {
		return pr.db.PayoutFailedQuery(ctx, payoutRequest)
	}

	// ---------------- DB TRANSACTION ----------------
	return fmt.Errorf("invalid payout status")
}
