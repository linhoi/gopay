package paypal

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/dghubble/sling"
	"github.com/linhoi/kit/log"
	"go.uber.org/zap"
	"gopay/common/httpx"
)

const (
	// App const name
	App string = "paypal_app"

	//const number
	clientTimeout                   = 10000 * time.Millisecond
	getTokenBeforeExpire            = 30 * time.Second
	maxPageSize                     = 50
	maxPageSizeForTranscationSearch = 500

	//papyal api path
	getAccessTokenPath     = "/v1/oauth2/token"       // https://developer.paypal.com/docs/api/reference/get-an-access-token/
	listDisputesPath       = "/v1/customer/disputes"  // https://developer.paypal.com/docs/api/customer-disputes/v1/#disputes_list
	showDisputeDetailsPath = "/v1/customer/disputes/" //https://developer.paypal.com/docs/api/customer-disputes/v1/#disputes_get
	getTransactionPath     = "/v1/reporting/transactions"
	nvp                    = "/nvp" // https://developer.paypal.com/docs/nvp-soap-api/nvp/

	// enum
	disputeStateRESOLVED = "RESOLVED"
	linkReasonTypeSelf   = "self" //https://www.iana.org/assignments/link-relations/link-relations.xhtml
	linkReasonTypeNext   = "next"
	refundByUser         = "RESOLVED_BUYER_FAVOUR" // https://developer.paypal.com/docs/api/customer-disputes/v1/#disputes_get

)

type Client struct {
	client        *sling.Sling
	config        Config
	token         string
	tokenExpireAt time.Time
}

type Config struct {
	Host                     string         `yaml:"host"`
	ClientID                 string         `yaml:"client_id"`
	Secret                   string         `yaml:"secret"`
	NVPAndSOAPAPICredentials APICredentials `yaml:"api_credentials"`
}

// APICredentials https://developer.paypal.com/docs/nvp-soap-api/apiCredentials/#api-signatures
type APICredentials struct {
	Username  string
	Password  string
	Signature string
}

func NewClient(c Config) *Client {
	hc := httpx.NewClient()
	client := sling.New().Client(hc)
	return &Client{
		client: client,
		config: c,
	}
}

// make sure Client implement API interface
var _ API = (*Client)(nil)

// GetAccessToken always get a token without expire.
func (c *Client) GetAccessToken(ctx context.Context) (string, error) {
	if c.token != "" && c.tokenExpireAt.After(time.Now().Add(-getTokenBeforeExpire)) {
		return c.token, nil
	}

	token, expireAt, err := c.getAccessToken(ctx)
	if err != nil {
		return "", err
	}

	c.token = token
	c.tokenExpireAt = expireAt

	log.L(ctx).Info("get token success", zap.String("expireAt", expireAt.Format(time.RFC3339)))
	return token, nil
}

// getAccessToken: In general, access tokens have a life of 15 minutes or eight hours depending on the scopes associated.
func (c *Client) getAccessToken(ctx context.Context) (token string, expireAt time.Time, err error) {
	path := c.config.Host + getAccessTokenPath
	req, err := http.NewRequest(http.MethodPost, path, bytes.NewBuffer([]byte("grant_type=client_credentials")))
	if err != nil {
		log.L(ctx).Error("GetAccessToken", zap.String("err", err.Error()))
		return "", time.Now(), err
	}
	req.SetBasicAuth(c.config.ClientID, c.config.Secret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var resBody AccessTokeResp
	var errResp ErrorResp
	_, err = c.client.Do(req, &resBody, &errResp)
	if err = errResp.Err(err); err != nil {
		log.L(ctx).Error("GetAccessToken", zap.String("err", err.Error()))
		return "", time.Now(), err
	}

	if resBody.AccessToken == "" {
		return "", time.Now(), errors.New("access token not found")
	}

	return resBody.AccessToken, time.Now().Add(time.Duration(resBody.ExpiresIn) * time.Second), nil
}

// ListDisputes list all resolve dispute.
func (c *Client) ListDisputes(ctx context.Context) (disputeItems []DisputeItem, error error) {
	accessToken, err := c.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	path := c.config.Host + listDisputesPath
	var firstPageDisputes ListDisputesResp
	var errResp ErrorResp
	_, err = c.client.Set("Authorization", "Bearer "+accessToken).
		QueryStruct(ListDisputesReq{DisputeState: disputeStateRESOLVED, PageSize: maxPageSize}).
		Path(path).Receive(&firstPageDisputes, &errResp)

	if err = errResp.Err(err); err != nil {
		log.S(ctx).Errorw("ListDisputes", "err", err)
		return nil, err
	}

	return c.listAllDisputes(ctx, accessToken, firstPageDisputes)
}

// listFirstPageDisputes ...
func (c *Client) listFirstPageDisputes(ctx context.Context) (disputeItems []DisputeItem, nextPageUrl string, error error) {
	accessToken, err := c.GetAccessToken(ctx)
	if err != nil {
		return nil, "", err
	}

	path := c.config.Host + listDisputesPath
	var firstPageDisputes ListDisputesResp
	var errResp ErrorResp
	_, err = c.client.Set("Authorization", "Bearer "+accessToken).
		QueryStruct(ListDisputesReq{DisputeState: disputeStateRESOLVED, PageSize: maxPageSize}).
		Path(path).Receive(&firstPageDisputes, &errResp)

	if err = errResp.Err(err); err != nil {
		log.S(ctx).Errorw("listFirstPageDisputes", "err", err)
		return nil, "", err
	}

	return firstPageDisputes.Items, findNextPageToken(firstPageDisputes), nil
}

// listFirstPageDisputes ...
func (c *Client) listFirstPageDisputesV2(ctx context.Context, startTime, endTime time.Time) (disputeItems []DisputeItem, nextPageUrl string, error error) {
	accessToken, err := c.GetAccessToken(ctx)
	if err != nil {
		return nil, "", err
	}

	path := c.config.Host + listDisputesPath
	var firstPageDisputes ListDisputesResp
	var errResp ErrorResp
	_, err = c.client.Set("Authorization", "Bearer "+accessToken).
		QueryStruct(ListDisputesReq{DisputeState: disputeStateRESOLVED, PageSize: maxPageSize,
			UpdateTimeAfter: startTime.Format(timeFmt), UpdateTimeBefore: endTime.Format(timeFmt)}).
		Path(path).Receive(&firstPageDisputes, &errResp)

	if err = errResp.Err(err); err != nil {
		log.S(ctx).Errorw("listFirstPageDisputesV2", "err", err)

		return nil, "", err
	}

	return firstPageDisputes.Items, findNextPageToken(firstPageDisputes), nil
}

// ListDisputes list all resolve dispute.
func (c *Client) ListOnePageDisputes(ctx context.Context, pagePageToken string) (disputeItems []DisputeItem, nextPageToken string, error error) {
	if pagePageToken == "" {
		return c.listFirstPageDisputes(ctx)
	}

	accessToken, err := c.GetAccessToken(ctx)
	if err != nil {
		return nil, "", err
	}

	var nextPageDisputes ListDisputesResp
	var errResp ErrorResp

	path := c.config.Host + listDisputesPath
	_, err = c.client.Set("Authorization", "Bearer "+accessToken).
		QueryStruct(ListDisputesReq{DisputeState: disputeStateRESOLVED, PageSize: maxPageSize, NextPageToken: pagePageToken}).
		Path(path).Receive(&nextPageDisputes, &errResp)
	if err = errResp.Err(err); err != nil {
		log.S(ctx).Errorw("ListOnePageDisputes", "err", err)

		return nil, "", err
	}

	return nextPageDisputes.Items, findNextPageToken(nextPageDisputes), nil
}

func (c *Client) ListOnePageDisputesV2(ctx context.Context, pagePageToken string, startTime, endTime time.Time) (disputeItems []DisputeItem, nextPageToken string, error error) {
	if pagePageToken == "" {
		return c.listFirstPageDisputesV2(ctx, startTime, endTime)
	}

	accessToken, err := c.GetAccessToken(ctx)
	if err != nil {
		return nil, "", err
	}

	var nextPageDisputes ListDisputesResp
	var errResp ErrorResp

	path := c.config.Host + listDisputesPath
	_, err = c.client.Set("Authorization", "Bearer "+accessToken).
		QueryStruct(
			ListDisputesReq{DisputeState: disputeStateRESOLVED, PageSize: maxPageSize, NextPageToken: pagePageToken,
				UpdateTimeAfter: startTime.Format(timeFmt), UpdateTimeBefore: endTime.Format(timeFmt)}).
		Path(path).Receive(&nextPageDisputes, &errResp)
	if err = errResp.Err(err); err != nil {
		log.S(ctx).Errorw("ListOnePageDisputesV2", "err", err)

		return nil, "", err
	}

	return nextPageDisputes.Items, findNextPageToken(nextPageDisputes), nil
}

func (c *Client) listAllDisputes(ctx context.Context, accessToken string, firstPageDisputes ListDisputesResp) (disputeItems []DisputeItem, err error) {
	allDisputeItems := make([]DisputeItem, 0, len(firstPageDisputes.Items))
	thisPage := firstPageDisputes

	for {
		path := findNextPageToken(thisPage)
		if path == "" {
			allDisputeItems = append(allDisputeItems, thisPage.Items...)
			break
		}

		var nextPageDisputes ListDisputesResp
		var errResp ErrorResp
		_, err = c.client.Set("Authorization", "Bearer "+accessToken).
			Path(path).Receive(&nextPageDisputes, &errResp)
		if err = errResp.Err(err); err != nil {
			log.S(ctx).Errorw("listAllDisputes", "err", err)

			break
		}

		thisPage = nextPageDisputes
	}

	return allDisputeItems, nil
}

func findNextPageToken(disputes ListDisputesResp) string {
	if len(disputes.Items) < maxPageSize {
		return ""
	}

	for _, link := range disputes.Links {
		if link.Rel == linkReasonTypeNext {
			u, err := url.Parse(link.Href)
			if err != nil {
				return ""
			}
			q := u.Query()
			return q.Get("next_page_token")
		}
	}
	return ""
}

// ShowDisputeDetails get dispute detail by disputeID.
func (c *Client) ShowDisputeDetails(ctx context.Context, disputeID string) (detail ShowDisputeDetailsResp, error error) {
	accessToken, err := c.GetAccessToken(ctx)
	if err != nil {
		return detail, err
	}

	path := c.config.Host + showDisputeDetailsPath + disputeID
	var errResp ErrorResp
	_, err = c.client.Set("Authorization", "Bearer "+accessToken).
		Path(path).Receive(&detail, &errResp)
	if err = errResp.Err(err); err != nil {
		log.S(ctx).Errorw("ShowDisputeDetails", "err", err)

		return detail, err
	}

	return detail, nil
}

// GetUserRefundDispute get user refund dispute if this dispute refunded by user.
func (c *Client) GetUserRefundDispute(ctx context.Context, disputeID string) (detail ShowDisputeDetailsResp, isUserRefund bool) {
	if disputeID == "" {
		return detail, false
	}

	detail, err := c.ShowDisputeDetails(ctx, disputeID)
	if err != nil {
		log.S(ctx).Errorw("GetUserRefundDispute ShowDisputeDetails", "err", err)

		return detail, false
	}

	if detail.DisputeOutcome.OutcomeCode == refundByUser {
		return detail, true
	}

	return detail, false
}

// GetUserRefundDisputedDetails get user refund disputes transaction details.
// userRefundDisputes is a map where key is invoiceID and value is disputeDetails
// invoiceID is SDK orderNumber sent to PayPal, details: https://developer.paypal.com/docs/nvp-soap-api/set-express-checkout-nvp
func (c *Client) GetUserRefundDisputedDetails(ctx context.Context, disputes []DisputeItem) (userRefundDisputeDetailByInvoiceID map[string]ShowDisputeDetailsResp, noInvoiceIDUserRefundDisputeDetails []ShowDisputeDetailsResp) {
	userRefundDisputeDetailByInvoiceID = make(map[string]ShowDisputeDetailsResp, len(disputes)/2)
	noInvoiceIDUserRefundDisputeDetails = make([]ShowDisputeDetailsResp, 0, len(disputes)/2)
	for _, dispute := range disputes {
		detail, isUserRefund := c.GetUserRefundDispute(ctx, dispute.DisputeID)
		if !isUserRefund {
			continue
		}
		log.S(ctx).Infow("a order refund by user in paypal", "dispute", dispute)

		for _, transaction := range detail.DisputedTransactions {
			if _, exist := userRefundDisputeDetailByInvoiceID[transaction.SellerTransactionID]; exist {
				continue
			}

			beginTime, endTime := parseCreateTimeRange(transaction.CreateTime)
			transactionInfo, err := c.GetTransaction(ctx, transaction.SellerTransactionID, beginTime, endTime)
			if err != nil {
				log.S(ctx).Errorw("GetTransaction", "err", err)

				continue
			}

			log.S(ctx).Infow("Get A Transaction", "transaction", transaction)

			if invoiceID := transactionInfo.GetInvoiceID(); invoiceID != "" {
				userRefundDisputeDetailByInvoiceID[invoiceID] = detail
			} else {
				noInvoiceIDUserRefundDisputeDetails = append(noInvoiceIDUserRefundDisputeDetails, detail)
			}

		}

	}

	return userRefundDisputeDetailByInvoiceID, noInvoiceIDUserRefundDisputeDetails
}

func parseCreateTimeRange(createTime time.Time) (time.Time, time.Time) {
	return createTime.Add(-time.Hour), createTime.Add(time.Hour)
}

// FiltrateUserRefundDispute filters user refund disputes.
func (c *Client) FiltrateUserRefundDispute(ctx context.Context, disputes []DisputeItem) (userRefundDisputes []DisputeItem) {
	userRefundDisputes = make([]DisputeItem, 0, len(disputes)/2)
	for _, dispute := range disputes {
		if _, isUserRefund := c.GetUserRefundDispute(ctx, dispute.DisputeID); isUserRefund {
			userRefundDisputes = append(userRefundDisputes, dispute)
		}
	}

	return userRefundDisputes
}

// GetTransaction get first transaction by transactionID
// time range between startTime and endTime must less than 31 days.
//https://developer.paypal.com/docs/integration/direct/transaction-search/#list-transactions
/*
curl example
curl -v -X GET https://api-m.sandbox.paypal.com/v1/reporting/transactions?transaction_id=5TY05013RG002845M&fields=all&page_size=100&page=1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <Access-Token>"
*/
func (c *Client) GetTransaction(ctx context.Context, transactionID string, startTime, endTime time.Time) (transaction TransactionInfo, err error) {
	if transactionID == "" {
		return transaction, errors.New("transactionID is zero")
	}

	accessToken, err := c.GetAccessToken(ctx)
	if err != nil {
		return transaction, err
	}

	var resp TransactionSearchResp
	var errResp ErrorResp
	path := c.config.Host + getTransactionPath
	_, err = c.client.Set("Authorization", "Bearer "+accessToken).
		Path(path).QueryStruct(TransactionSearchReq{TransactionID: transactionID, StartDate: startTime.Format(timeFmt), EndDate: endTime.Format(timeFmt), Fields: "all", PageSize: 100, Page: 1}).
		Receive(&resp, &errResp)

	if err = errResp.Err(err); err != nil {
		log.S(ctx).Errorw("GetTransaction", "path", path,"err",err)

		return transaction, err
	}

	if len(resp.TransactionDetails) > 0 {
		return resp.TransactionDetails[0].TransactionInfo, nil
	}

	return transaction, errors.New("transaction not fund")
}

func (c *Client) GetRefundTransaction(ctx context.Context, startTime, endTime time.Time) (transactions []TransactionInfo, err error) {
	transaction, totalPage, err := c.getRefundTransactionByPage(ctx, startTime, endTime, 1)
	if err != nil {
		return transactions, err
	}

	transactions = make([]TransactionInfo, 0, maxPageSizeForTranscationSearch)
	transactions = append(transactions, transaction...)
	for page := 2; page <= totalPage; page++ {
		transaction, _, err := c.getRefundTransactionByPage(ctx, startTime, endTime, page)
		if err != nil {
			log.S(ctx).Errorw("getRefundTransactionByPage", "err",err)

			return transactions, err
		}

		transactions = append(transactions, transaction...)
	}

	return transactions, nil
}

func (c *Client) getRefundTransactionByPage(ctx context.Context, startTime, endTime time.Time, page int) (transactions []TransactionInfo, totalPage int, err error) {
	accessToken, err := c.GetAccessToken(ctx)
	if err != nil {
		return nil, 0, err
	}

	var resp TransactionSearchResp
	var errResp ErrorResp
	path := c.config.Host + getTransactionPath
	_, err = c.client.Set("Authorization", "Bearer "+accessToken).
		Path(path).QueryStruct(
		TransactionSearchReq{
			StartDate: startTime.Format(timeFmt), EndDate: endTime.Format(timeFmt), Fields: "all", PageSize: maxPageSizeForTranscationSearch, Page: page,
			TransactionStatus: TransactionRefund}).
		Receive(&resp, &errResp)

	if err = errResp.Err(err); err != nil {
		return nil, 0, err
	}

	transactions = make([]TransactionInfo, 0, len(resp.TransactionDetails))
	for _, detail := range resp.TransactionDetails {
		transactions = append(transactions, detail.TransactionInfo)
	}

	return transactions, resp.TotalPages, nil
}
