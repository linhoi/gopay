package google

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/linhoi/kit/log"
	"go.uber.org/zap"
	"google.golang.org/api/androidpublisher/v3"
	"google.golang.org/api/option"
)

type Client struct {
	googlePublisher *androidpublisher.Service
}

func NewClient(credentialsJSON []byte) (*Client, error) {
	service, err := androidpublisher.NewService(context.Background(), option.WithCredentialsJSON(credentialsJSON))
	if err != nil {
		return nil, err
	}
	return &Client{
		googlePublisher: service,
	}, nil
}

type API interface {
	GetPurchase(ctx context.Context, packageName string, productID string, purchaseToken string) (*androidpublisher.ProductPurchase, error)
	GetSubscription(ctx context.Context, packageName, subscriptionID, purchaseToken string)  (*androidpublisher.SubscriptionPurchase, error)
	Acknowledge(ctx context.Context, packageName string, productID string, token string) error

	// ReceiveRealTimeDeveloperNotification ...
	// RTDN means RealTimeDeveloperNotification
	// https://developer.android.google.cn/google/play/billing/rtdn-reference?hl=zh-cn
	ReceiveRealTimeDeveloperNotification(ctx context.Context, req *http.Request) error
}

var _ API = (*Client)(nil)

func (c *Client) GetPurchase(ctx context.Context, packageName string, productID string, purchaseToken string) (*androidpublisher.ProductPurchase, error) {
	res, err := c.googlePublisher.Purchases.Products.Get(packageName, productID, purchaseToken).Do()
	if err != nil {
		log.L(ctx).Warn("google play get purchase failed", zap.Error(err))
		return nil, err
	}

	if res.PurchaseState != PurchaseStatePurchased {
		log.L(ctx).Warn("google play purchase status is not purchase", zap.Int64("purchase_status", res.PurchaseState))
		return res, err
	}

	return res, nil
}

func (c *Client) GetSubscription(ctx context.Context, packageName, subscriptionID, purchaseToken string) (*androidpublisher.SubscriptionPurchase, error) {
	res, err := c.googlePublisher.Purchases.Subscriptions.Get(packageName, subscriptionID, purchaseToken).Do()
	if err != nil {
		log.L(ctx).Warn("google play get subscription failed", zap.Error(err))
		return nil, err
	}

	return res, nil

}

func (c *Client) Acknowledge(ctx context.Context, packageName string, productID string, token string) error {
	err := c.googlePublisher.Purchases.Products.Acknowledge(packageName, productID, token,nil).Do()
	if err != nil {
		log.L(ctx).Warn("google play get subscription failed", zap.Error(err))
		return err
	}

	return nil
}

func (c *Client) ReceiveRealTimeDeveloperNotification(ctx context.Context, req *http.Request) error {
	body , err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.L(ctx).Warn("read rtdn body failed",zap.Error(err))
		return err
	}

	var rtdnBody RTDNBody
	err = json.Unmarshal(body,&rtdnBody)
	if err != nil {
		log.L(ctx).Warn("read rtdn json unmarshal failed",zap.Error(err))
		return err
	}

	data,err := base64.StdEncoding.DecodeString(rtdnBody.Message.Data)
	if err != nil {
		log.L(ctx).Warn("read rtdn decode failed",zap.Error(err))
		return err
	}

	var rtdnData RTDNData
	err = json.Unmarshal(data,&rtdnData)
	if err != nil {
		log.L(ctx).Warn("read rtdn json unmarshal after decode failed",zap.Error(err))
		return err
	}

	if rtdnData.OneTimeProductNotification == nil {
		log.L(ctx).Info("rtdn is not one time product notification",zap.Error(err))
		return nil
	}

	if rtdnData.OneTimeProductNotification.NotificationType != ONE_TIME_PRODUCT_PURCHASED {
		log.L(ctx).Info("rtdn is not one time product purchases notification",zap.Error(err))
		return nil
	}

	res, err := c.GetPurchase(ctx, rtdnData.PackageName, rtdnData.OneTimeProductNotification.Sku, rtdnData.OneTimeProductNotification.PurchaseToken)
	if res != nil {
		return nil
	}
	if err != nil {
		log.L(ctx).Info("rtdn get purchase  failed",zap.Error(err),zap.Any("rtdn",rtdnData))
		return err
	}

	err = c.Acknowledge(ctx, rtdnData.PackageName, rtdnData.OneTimeProductNotification.Sku, rtdnData.OneTimeProductNotification.PurchaseToken)
	if err != nil {
		log.L(ctx).Info("rtdn acknowledge  failed",zap.Error(err),zap.Any("rtdn",rtdnData))
		return err
	}

	return nil
}
