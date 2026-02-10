package models

type GetPostpaidMobileRechargeBillFetchAPIRequestModel struct {
	MobileNumber string `json:"mobile_no" validate:"required"`
	OperatorCode int    `json:"operator_code" validate:"required"`
}

type GetPostpaidMobileRechargeBillFetchAPIResponseModel struct {
	Error      int                                `json:"error"`
	Message    string                             `json:"msg"`
	Status     int                                `json:"status"`
	BillAmount GetPostpaidBillAmountAPIResponseModel `json:"billAmount"`
}

type GetPostpaidBillAmountAPIResponseModel struct {
	BillAmount    float64 `json:"billAmount"`
	BillNetAmount float64 `json:"billnetamount"`
	BillDate      string  `json:"billdate"`
	BillDueDate   string  `json:"dueDate"`
	AcceptPayment bool    `json:"acceptPayment"`
	AcceptPartPay bool    `json:"acceptPartPay"`
	CellNumber    string  `json:"cellNumber"`
	UserName      string  `json:"userName"`
}

type CreatePostpaidMobileRechargeAPIRequestModel struct {
	RetailerID       string  `json:"retailer_id" validate:"required"`
	MobileNumber     string  `json:"mobile_no" validate:"required"`
	OperatorCode     int     `json:"operator_code" validate:"required"`
	Amount           float64 `json:"amount" validate:"required"`
	OperatorCircle   int     `json:"circle" validate:"required"`
	PartnerRequestID string  `json:"partner_request_id,omitempty"`
	OperatorName     string  `json:"operator_name" validate:"required"`
	CircleName       string  `json:"circle_name" validate:"required"`
}

type GetPostpaidMobileRechargeAPIResponseModel struct {
	Error                 int    `json:"error"`
	Message               string `json:"msg"`
	Status                int    `json:"status"`
	OrderID               string `json:"orderid"`
	OperatorTransactionID string `json:"optransid"`
}
