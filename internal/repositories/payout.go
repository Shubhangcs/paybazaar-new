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

	if err := bindAndValidate(c, &payoutRequest); err != nil {
		return err
	}

	// 🔐 Retailer ID MUST come from JWT
	claims, ok := c.Get("user").(*models.AccessTokenClaims)
	if !ok || claims.UserID == "" {
		return fmt.Errorf("unauthorized")
	}
	payoutRequest.RetailerID = claims.UserID

	// Amount validation
	if payoutRequest.Amount < 1000 || payoutRequest.Amount > 25000 {
		return fmt.Errorf("failed to payout invalid amount")
	}

	ctx, cancel := context.WithTimeout(
		c.Request().Context(),
		30*time.Second,
	)
	defer cancel()

	// MPIN + KYC check
	isValid, err := pr.db.VerifyMPINAndKycQuery(
		ctx,
		payoutRequest.RetailerID,
		payoutRequest.MPIN,
	)
	if err != nil {
		return err
	}

	fmt.Println(isValid)
	if !isValid {
		return fmt.Errorf("invalid mpin or incomplete kyc")
	}

	// Commission split
	commision, err := pr.db.GetPayoutCommisionSplit(
		ctx,
		payoutRequest.RetailerID,
		payoutRequest.Amount,
	)
	if err != nil {
		return err
	}

	// Wallet balance check
	hasBalance, err := pr.db.CheckRetailerWalletBalance(
		ctx,
		payoutRequest.RetailerID,
		payoutRequest.Amount,
		commision.TotalCommision,
	)
	if err != nil {
		return err
	}

	if !hasBalance {
		return fmt.Errorf("insufficient balance")
	}

	// Generate partner request ID
	partnerReqID := uuid.New()
	payoutRequest.PartnerRequestID = partnerReqID.String()

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

	// ---------------- DB TRANSACTION ----------------
	return pr.db.CreatePayoutQuery(ctx, payoutRequest, res, *commision)
}
