package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/models"
)

type DMTInterface interface {
	CheckDMTWalletExists(echo.Context) (*models.DMTWalletCheckResponseModel, error)
	CreateDMTWallet(echo.Context) (*models.DMTCreateWalletResponseModel, error)
	VerifyDMTWallet(echo.Context) (*models.DMTWalletVerificationResponseModel, error)
	AddDMTBeneficiary(echo.Context) (*models.DMTAddBeneficiaryResponseModel, error)
	GetDMTBankList(echo.Context) (*models.DMTBankListResponseModel, error)
}

type dmtRepository struct {
	db *database.Database
}

func NewDMTRepository(db *database.Database) *dmtRepository {
	return &dmtRepository{
		db,
	}
}

func (dr *dmtRepository) CheckDMTWalletExists(c echo.Context) (*models.DMTWalletCheckResponseModel, error) {
	var req models.DMTWalletCheckRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return nil, err
	}
	apiUrl := `https://v2bapi.rechargkit.biz/moneytransfer/checkWalletExist`
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	apiRequest, err := http.NewRequest(
		http.MethodPost,
		apiUrl,
		bytes.NewReader(reqBody),
	)
	if err != nil {
		return nil, err
	}

	apiRequest.Header.Set("Content-Type", "application/json")
	apiRequest.Header.Set("Authorization", "Bearer "+os.Getenv("RKIT_API_TOKEN"))

	client := &http.Client{Timeout: 20 * time.Second}

	resp, err := client.Do(apiRequest)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResponse models.DMTWalletCheckResponseModel
	if err := json.Unmarshal(respBytes, &apiResponse); err != nil {
		return nil, err
	}
	return &apiResponse, nil
}

func (dr *dmtRepository) CreateDMTWallet(c echo.Context) (*models.DMTCreateWalletResponseModel, error) {
	var req models.DMTCreateWalletRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return nil, err
	}
	apiUrl := `https://v2bapi.rechargkit.biz/moneytransfer/createWalletRequest`
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()
	aadharNumber, err := dr.db.GetRetailerAadharNumberForDMTQuery(ctx, req.RetailerID)
	if err != nil {
		return nil, err
	}
	req.AadharNumber = aadharNumber
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	apiRequest, err := http.NewRequest(
		http.MethodPost,
		apiUrl,
		bytes.NewReader(reqBody),
	)
	if err != nil {
		return nil, err
	}

	apiRequest.Header.Set("Content-Type", "application/json")
	apiRequest.Header.Set("Authorization", "Bearer "+os.Getenv("RKIT_API_TOKEN"))

	client := &http.Client{Timeout: 20 * time.Second}

	resp, err := client.Do(apiRequest)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Println(string(respBytes))

	var apiResponse models.DMTCreateWalletResponseModel
	if err := json.Unmarshal(respBytes, &apiResponse); err != nil {
		return nil, err
	}
	return &apiResponse, nil
}

func (dr *dmtRepository) VerifyDMTWallet(c echo.Context) (*models.DMTWalletVerificationResponseModel, error) {
	var req models.DMTWalletVerificationRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return nil, err
	}
	apiUrl := `https://v2bapi.rechargkit.biz/moneytransfer/verifyWalletRequest`
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()
	aadharNumber, err := dr.db.GetRetailerAadharNumberForDMTQuery(ctx, req.RetailerID)
	if err != nil {
		return nil, err
	}
	req.StateResp = aadharNumber
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	apiRequest, err := http.NewRequest(
		http.MethodPost,
		apiUrl,
		bytes.NewReader(reqBody),
	)
	if err != nil {
		return nil, err
	}

	apiRequest.Header.Set("Content-Type", "application/json")
	apiRequest.Header.Set("Authorization", "Bearer "+os.Getenv("RKIT_API_TOKEN"))

	client := &http.Client{Timeout: 20 * time.Second}

	resp, err := client.Do(apiRequest)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResponse models.DMTWalletVerificationResponseModel
	if err := json.Unmarshal(respBytes, &apiResponse); err != nil {
		return nil, err
	}
	return &apiResponse, nil
}

func (dr *dmtRepository) AddDMTBeneficiary(c echo.Context) (*models.DMTAddBeneficiaryResponseModel, error) {
	var req models.DMTAddBeneficiaryRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return nil, err
	}
	apiUrl := `https://v2bapi.rechargkit.biz/moneytransfer/addBeneficiaryRequest`
	req.PartnerRequestID = uuid.NewString()
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	apiRequest, err := http.NewRequest(
		http.MethodPost,
		apiUrl,
		bytes.NewReader(reqBody),
	)
	if err != nil {
		return nil, err
	}

	apiRequest.Header.Set("Content-Type", "application/json")
	apiRequest.Header.Set("Authorization", "Bearer "+os.Getenv("RKIT_API_TOKEN"))

	client := &http.Client{Timeout: 20 * time.Second}

	resp, err := client.Do(apiRequest)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResponse models.DMTAddBeneficiaryResponseModel
	if err := json.Unmarshal(respBytes, &apiResponse); err != nil {
		return nil, err
	}
	return &apiResponse, nil
}

func (dr *dmtRepository) GetDMTBankList(c echo.Context) (*models.DMTBankListResponseModel, error) {
	apiUrl := `https://v2a.rechargkit.biz/moneytransfer/getBankList`

	apiRequest, err := http.NewRequest(
		http.MethodGet,
		apiUrl,
		nil,
	)
	if err != nil {
		return nil, err
	}

	apiRequest.Header.Set("Content-Type", "application/json")
	apiRequest.Header.Set("Authorization", "Bearer "+os.Getenv("RKIT_API_TOKEN"))

	client := &http.Client{Timeout: 20 * time.Second}

	resp, err := client.Do(apiRequest)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResponse models.DMTBankListResponseModel
	if err := json.Unmarshal(respBytes, &apiResponse); err != nil {
		return nil, err
	}
	return &apiResponse, nil
}
