package google

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/linhoi/gopay/common/httpx"
	"github.com/linhoi/kit/log"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/androidpublisher/v3"
	"google.golang.org/api/option"
)

type Client struct {
	credentials     []byte
	googlePublisher *androidpublisher.Service
	client          *http.Client
}

func NewClient(credentialsJSON []byte) (*Client, error) {
	service, err := androidpublisher.NewService(context.Background(), option.WithCredentialsJSON(credentialsJSON))
	if err != nil {
		return nil, err
	}
	return &Client{
		credentials:     credentialsJSON,
		googlePublisher: service,
		client:          httpx.NewClient(),
	}, nil
}

type API interface {
	GetToken(ctx context.Context) (*oauth2.Token, error)

	GetPurchase(ctx context.Context, packageName string, productID string, purchaseToken string) (*androidpublisher.ProductPurchase, error)
	GetPurchaseByReceipt(ctx context.Context, receipt Receipt) (*androidpublisher.ProductPurchase, error)

	GetSubscription(ctx context.Context, packageName, subscriptionID, purchaseToken string) (*androidpublisher.SubscriptionPurchase, error)
	GetSubscriptionByReceipt(ctx context.Context, receipt Receipt) (*androidpublisher.SubscriptionPurchase, error)

	Acknowledge(ctx context.Context, packageName string, productID string, token string) error
	AcknowledgeByReceipt(ctx context.Context, receipt Receipt) error

	GetVoidedPurchase(ctx context.Context, v VoidedPurchase) (voidedPurchases []*androidpublisher.VoidedPurchase, nextPageToken string, err error)
	DealVoidPurchase(ctx context.Context, d DealVoidedPurchase) (dealFailed []*androidpublisher.VoidedPurchase, err error)

	// ReceiveRealTimeDeveloperNotification ...
	// RTDN means RealTimeDeveloperNotification
	// https://developer.android.google.cn/google/play/billing/rtdn-reference?hl=zh-cn
	ReceiveRealTimeDeveloperNotification(ctx context.Context, req *http.Request) error
}

var _ API = (*Client)(nil)

func (c *Client) GetToken(ctx context.Context) (tokenInfo *oauth2.Token, err error) {
	conf, err := google.JWTConfigFromJSON(c.credentials, androidpublisher.AndroidpublisherScope)
	if err != nil {
		log.L(ctx).Warn("google get jwt config failed", zap.Error(err))
		return nil, err
	}

	authCtx := context.WithValue(context.Background(), oauth2.HTTPClient, c.client)
	authTransport := conf.Client(authCtx).Transport.(*oauth2.Transport)
	token, err := authTransport.Source.Token()
	if err != nil {
		log.L(ctx).Warn("auth transport get token failed", zap.Error(err))
		return nil, err
	}

	return token, nil
}

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
	err := c.googlePublisher.Purchases.Products.Acknowledge(packageName, productID, token, nil).Do()
	if err != nil {
		log.L(ctx).Warn("google play get subscription failed", zap.Error(err))
		return err
	}

	return nil
}

func (c *Client) GetPurchaseByReceipt(ctx context.Context, receipt Receipt) (*androidpublisher.ProductPurchase, error) {
	return c.GetPurchase(ctx, receipt.PackageName, receipt.ProductID, receipt.Token)
}

func (c *Client) GetSubscriptionByReceipt(ctx context.Context, receipt Receipt) (*androidpublisher.SubscriptionPurchase, error) {
	return c.GetSubscription(ctx, receipt.PackageName, receipt.SubscriptionID, receipt.Token)
}

func (c *Client) AcknowledgeByReceipt(ctx context.Context, receipt Receipt) error {
	return c.Acknowledge(ctx, receipt.PackageName, receipt.ProductID, receipt.Token)
}

func (c *Client) ReceiveRealTimeDeveloperNotification(ctx context.Context, req *http.Request) error {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.L(ctx).Warn("read rtdn body failed", zap.Error(err))
		return err
	}

	var rtdnBody RTDNBody
	err = json.Unmarshal(body, &rtdnBody)
	if err != nil {
		log.L(ctx).Warn("read rtdn json unmarshal failed", zap.Error(err))
		return err
	}

	data, err := base64.StdEncoding.DecodeString(rtdnBody.Message.Data)
	if err != nil {
		log.L(ctx).Warn("read rtdn decode failed", zap.Error(err))
		return err
	}

	var rtdnData RTDNData
	err = json.Unmarshal(data, &rtdnData)
	if err != nil {
		log.L(ctx).Warn("read rtdn json unmarshal after decode failed", zap.Error(err))
		return err
	}

	if rtdnData.OneTimeProductNotification == nil {
		log.L(ctx).Info("rtdn is not one time product notification", zap.Error(err))
		return nil
	}

	if rtdnData.OneTimeProductNotification.NotificationType != ONE_TIME_PRODUCT_PURCHASED {
		log.L(ctx).Info("rtdn is not one time product purchases notification", zap.Error(err))
		return nil
	}

	res, err := c.GetPurchase(ctx, rtdnData.PackageName, rtdnData.OneTimeProductNotification.Sku, rtdnData.OneTimeProductNotification.PurchaseToken)
	if res != nil {
		return nil
	}
	if err != nil {
		log.L(ctx).Info("rtdn get purchase  failed", zap.Error(err), zap.Any("rtdn", rtdnData))
		return err
	}

	err = c.Acknowledge(ctx, rtdnData.PackageName, rtdnData.OneTimeProductNotification.Sku, rtdnData.OneTimeProductNotification.PurchaseToken)
	if err != nil {
		log.L(ctx).Info("rtdn acknowledge  failed", zap.Error(err), zap.Any("rtdn", rtdnData))
		return err
	}

	return nil
}

func (c *Client) GetVoidedPurchase(ctx context.Context, v VoidedPurchase) (voidedPurchases []*androidpublisher.VoidedPurchase, nextPageToken string, err error) {
	req := c.googlePublisher.Purchases.Voidedpurchases.List(v.PackageName).Fields().StartTime(v.StartTime.UnixNano() / int64(time.Millisecond)).EndTime(v.EndTime.UnixNano() / int64(time.Millisecond))
	if v.Token != "" {
		req = req.Token(v.Token)
	}

	res, err := req.Do()
	if err != nil {
		fmt.Println(res)
		log.L(ctx).Warn("get voided purchase failed", zap.Error(err), zap.Any("req", v))
		return nil, "", err
	}
	if res.TokenPagination != nil {
		nextPageToken = res.TokenPagination.NextPageToken
	}

	return res.VoidedPurchases, nextPageToken, nil
}

func (c *Client) DealVoidPurchase(ctx context.Context, d DealVoidedPurchase) (dealFailed []*androidpublisher.VoidedPurchase, err error) {
	var nextPageToken string
	var voidedPurchases []*androidpublisher.VoidedPurchase

	res, err := c.googlePublisher.Purchases.Voidedpurchases.List(d.PackageName).StartTime(d.StartTime.Unix()).EndTime(d.EndTime.Unix()).Do()
	if err != nil {
		log.L(ctx).Warn("list voided purchase failed", zap.Error(err), zap.Any("req", d))
		return nil, err
	}
	if res.TokenPagination != nil {
		nextPageToken = res.TokenPagination.NextPageToken
	}
	voidedPurchases = res.VoidedPurchases

	for {
		for _, voidedPurchase := range voidedPurchases {
			if d.DealFunc != nil {
				err := d.DealFunc(ctx, voidedPurchase)
				if err != nil {
					dealFailed = append(dealFailed, voidedPurchase)
				}
			}
		}

		if nextPageToken == "" {
			break
		}

		res, err := c.googlePublisher.Purchases.Voidedpurchases.List(d.PackageName).StartTime(d.StartTime.Unix()).EndTime(d.EndTime.Unix()).Token(nextPageToken).Do()
		if err != nil {
			log.L(ctx).Warn("list voided purchase failed", zap.Error(err), zap.Any("req", d))
			return dealFailed, err
		}
		if res.TokenPagination != nil {
			nextPageToken = res.TokenPagination.NextPageToken
		}
		voidedPurchases = res.VoidedPurchases
	}

	return dealFailed, nil
}
