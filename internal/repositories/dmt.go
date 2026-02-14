package repositories

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	GetDmtBeneficiary(echo.Context) (*models.DMTGetBeneficiaryResponseModel, error)
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
	apiUrl := `https://v2bapi.rechargkit.biz/rkitdmr/checkWalletExist`
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
	apiUrl := `https://v2bapi.rechargkit.biz/rkitdmr/createWalletRequest`
	reqBody, err := json.Marshal(map[string]any{
		"mobile_no":      req.MobileNumber,
		"lat":            req.Latitude,
		"long":           req.Longitude,
		"aadhaar_number": req.AadharNumber,
		"pid_data":       req.PidData,
		"is_iris":        req.IsIris,
	})
	if err != nil {
		return nil, err
	}

	fmt.Printf("%s", string(reqBody))

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
	apiUrl := `https://v2bapi.rechargkit.biz/rkitdmr/verifyWalletRequest`
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
	apiUrl := `https://v2bapi.rechargkit.biz/rkitdmr/addBeneficiaryRequest`
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

func (dr *dmtRepository) GetDmtBeneficiary(c echo.Context) (*models.DMTGetBeneficiaryResponseModel, error) {
	var req models.DMTGetBeneficiaryRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return nil, err
	}
	apiUrl := `https://v2bapi.rechargkit.biz/rkitdmr/getUserDetails`
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
	fmt.Println(string(respBytes))

	var apiResponse models.DMTGetBeneficiaryResponseModel
	if err := json.Unmarshal(respBytes, &apiResponse); err != nil {
		return nil, err
	}
	return &apiResponse, nil
}

func (dr *dmtRepository) GetDMTBankList(c echo.Context) (*models.DMTBankListResponseModel, error) {
	apiUrl := `https://v2bapi.rechargkit.biz/rkitdmr/getBankList`
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


