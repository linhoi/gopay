package apple

import (
	"context"
	"crypto/x509"
	"encoding/base64"

	"github.com/awa/go-iap/appstore"
	"github.com/linhoi/gopay/common/httpx"
	"github.com/linhoi/kit/log"
	"go.uber.org/zap"
)

type IAP interface {
	// LocalValidateReceipt ...
	LocalValidateReceipt(ctx context.Context, receipt string) (*Receipts, error)
	// Verify ...
	Verify(ctx context.Context, receipt string) (*appstore.IAPResponse, error)
}

type Client struct {
	config   *Config
	appstore *appstore.Client
}

// NewClient ...
func NewClient(config *Config) *Client {
	hc := appstore.NewWithClient(httpx.NewClient())

	return &Client{
		config:   config,
		appstore: hc,
	}
}

// Verify : reference https://developer.apple.com/documentation/storekit/original_api_for_in-app_purchase/validating_receipts_with_the_app_store#//apple_ref/doc/uid/TP40010573-CH104-SW1
func (c *Client) Verify(ctx context.Context, receipt string) (*appstore.IAPResponse, error) {
	resp := &appstore.IAPResponse{}
	err := c.appstore.Verify(ctx, appstore.IAPRequest{ReceiptData: receipt}, &resp)
	if err != nil {
		log.L(ctx).Warn("appstore verify failed", zap.Error(err))
		return nil, err
	}

	iapStatus := Status(resp.Status)
	if iapStatus == StatusSuccess {
		return resp, nil
	}

	return resp, iapStatus.Error()
}

// LocalValidateReceipt ...
// reference : https://developer.apple.com/library/archive/releasenotes/General/ValidateAppStoreReceipt/Chapters/ValidateLocally.html#//apple_ref/doc/uid/TP40010573-CH1-SW2
func (c Client) LocalValidateReceipt(ctx context.Context, receipt string) (*Receipts, error) {
	data, err := base64.StdEncoding.DecodeString(receipt)
	if err != nil {
		log.L(ctx).Warn("decode receipt failed", zap.Error(err))
		return nil, err
	}

	caBytes, err := base64.StdEncoding.DecodeString(c.config.RootCA)
	if err != nil {
		log.L(ctx).Warn("base64 decode ca failed", zap.Error(err))
		return nil, err
	}

	ca, err := x509.ParseCertificate(caBytes)
	if err != nil {
		log.L(ctx).Warn("parse ca failed", zap.Error(err))
		return nil, err
	}

	receipts, err := ParseReceipt(ca, data)
	if err != nil {
		log.L(ctx).Warn("parse receipt failed", zap.Error(err))
		return nil, err
	}

	return &receipts, nil
}

var _ IAP = (*Client)(nil)
