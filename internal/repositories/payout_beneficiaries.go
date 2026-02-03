package repositories

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/models"
	"github.com/levion-studio/paybazaar/pkg"
)

type Beneficiary interface {
	GetBeneficiaries(echo.Context) (*[]models.BeneficiaryModel, error)
	AddNewBeneficiary(echo.Context) error
	VerifyBeneficiary(echo.Context) (*models.VerifyBeneficiaryResponseModel, error)
	DeleteBeneficiary(echo.Context) error
}

type beneficiaryRepo struct {
	query *database.Database
	jwt   *pkg.JwtUtils
}

func NewBeneficiaryRepo(query *database.Database) *beneficiaryRepo {
	return &beneficiaryRepo{query: query}
}

func (r *beneficiaryRepo) GetBeneficiaries(e echo.Context) (*[]models.BeneficiaryModel, error) {
	var phone = e.Param("phone")
	if phone == "" {
		return nil, fmt.Errorf("phone number not found")
	}
	res, err := r.query.GetBeneficiaries(phone)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch benificiries")
	}
	return res, nil
}

func (r *beneficiaryRepo) AddNewBeneficiary(e echo.Context) error {
	req := &models.BeneficiaryModel{}
	if err := e.Bind(req); err != nil {
		return err
	}
	err := r.query.AddNewBeneficiary(req)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("failed to add benificary")
	}
	return nil
}

func (r *beneficiaryRepo) VerifyBeneficiary(c echo.Context) (*models.VerifyBeneficiaryResponseModel, error) {
	var req models.VerifyBeneficiaryRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return nil, err
	}
	referenceId := uuid.NewString()
	token, err := r.jwt.GenerateTokenForPayoutBeneVerification(referenceId)
	if err != nil {
		return nil, err
	}

	payload := map[string]string{
		"refid":          referenceId,
		"account_number": req.AccountNumber,
		"ifsc_code":      req.IFSCCode,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	apiReq, err := http.NewRequest(
		"POST",
		"https://api.verifya2z.com/api/v1/verification/penny_drop_v2",
		strings.NewReader(string(bodyBytes)),
	)
	if err != nil {
		return nil, err
	}

	apiReq.Header.Add("accept", "application/json")
	apiReq.Header.Add("Token", token)
	apiReq.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(apiReq)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	// decode API response directly and return
	var resp models.VerifyBeneficiaryResponseModel
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (r *beneficiaryRepo) DeleteBeneficiary(e echo.Context) error {
	beneficiaryId := e.Param("ben_id")
	if beneficiaryId == "" {
		return fmt.Errorf("beneficiary id not found")
	}
	if err := r.query.DeleteBeneficiary(beneficiaryId); err != nil {
		return fmt.Errorf("failed to delete beneficiary")
	}
	return nil
}
