package paypal

import (
	"context"
	"time"
)

// API is an interface for paypal REST APIs
// https://developer.paypal.com/docs/integration/direct/customer-disputes/
type API interface {
	// GetAccessToken AccessToken ...
	GetAccessToken(ctx context.Context) (string, error)

	// 争议查询
	ListDisputes(ctx context.Context) (disputeItems []DisputeItem, error error)
	ShowDisputeDetails(ctx context.Context, disputeID string) (detail ShowDisputeDetailsResp, error error)

	// 交易查询
	//https://developer.paypal.com/docs/api/transaction-search/v1/#transactions_get
	GetTransaction(ctx context.Context, transactionID string, startTime, endTime time.Time) (transaction TransactionInfo, err error)
	GetRefundTransaction(ctx context.Context, startTime, endTime time.Time) (transactions []TransactionInfo, err error)

	//NVP
	SetExpressCheckout(ctx context.Context, req SetExpressCheckoutReq) (SetExpressCheckoutResp , error)
	GetExpressCheckoutDetails(ctx context.Context, req GetExpressCheckoutDetailsReq) (GetExpressCheckoutDetailsResp , error)
	DoExpressCheckoutPayment(ctx context.Context, req GetExpressCheckoutDetailsReq) (GetExpressCheckoutDetailsResp , error)
}


