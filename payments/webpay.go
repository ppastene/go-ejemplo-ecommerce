package payments

import (
	"encoding/json"

	"github.com/ppastene/go-shopping-lite/utils"
)

type Webpay struct {
	ApiKey, ApiSecret string
}

type payload struct {
	BuyOrder  string  `json:"buy_order"`
	SessionId string  `json:"session_id"`
	Amount    float64 `json:"amount"`
	ReturnUrl string  `json:"return_url"`
}

type transactionCreateResponse struct {
	Token string `json:"token"`
	Url   string `json:"url"`
}

type transactionStatusResponse struct {
	Vci                string      `json:"vci"`
	Amount             float64     `json:"amount"`
	Status             string      `json:"status"`
	BuyOrder           string      `json:"buy_order"`
	SessionId          string      `json:"session_id"`
	CardDetail         cardDetails `json:"card_detail"`
	AccountingDate     string      `json:"accounting_date"`
	TransactionDate    string      `json:"transaction_date"`
	AuthorizationCode  string      `json:"authorization_code"`
	PaymentTypeCode    string      `json:"payment_type_code"`
	ResponseCode       int         `json:"response_code"`
	InstallmentsAmount int         `json:"installments_amount"`
	InstallmentsNumber int         `json:"installments_number"`
	Balance            int         `json:"balance"`
}

type cardDetails struct {
	CardNumber string `json:"card_number"`
}

func NewWebpay(apiKey, apiSecret string) *Webpay {
	return &Webpay{
		ApiKey:    apiKey,
		ApiSecret: apiSecret,
	}
}

func (w *Webpay) Create(buyOrder string, sessionId string, amount float64, returnUrl string) (transactionCreateResponse, error) {
	payload := &payload{
		BuyOrder:  buyOrder,
		SessionId: sessionId,
		Amount:    amount,
		ReturnUrl: returnUrl,
	}
	jsonData, _ := json.Marshal(payload)
	client := utils.NewRequest("https://webpay3gint.transbank.cl/rswebpaytransaction/api/webpay/v1.3/transactions/")
	headers := map[string]string{
		"Tbk-Api-Key-Id":     w.ApiKey,
		"Tbk-Api-Key-Secret": w.ApiSecret,
		"Content-Type":       "application/json",
	}
	client.Headers(headers)
	client.Body(jsonData)
	res, err := client.Post("")
	if err != nil {
		return transactionCreateResponse{}, err
	}
	var result transactionCreateResponse
	err = json.Unmarshal(res.Body(), &result)
	if err != nil {
		return transactionCreateResponse{}, err
	}
	return result, nil
}

func (w *Webpay) Commit(token string) (transactionStatusResponse, error) {
	client := utils.NewRequest("https://webpay3gint.transbank.cl/rswebpaytransaction/api/webpay/v1.3/transactions/" + token)
	headers := map[string]string{
		"Tbk-Api-Key-Id":     w.ApiKey,
		"Tbk-Api-Key-Secret": w.ApiSecret,
		"Content-Type":       "application/json",
	}
	client.Headers(headers)
	res, err := client.Put("")
	if err != nil {
		return transactionStatusResponse{}, err
	}
	var response transactionStatusResponse
	err = json.Unmarshal(res.Body(), &response)
	if err != nil {
		return transactionStatusResponse{}, err
	}
	return response, nil
}
