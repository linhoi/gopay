package google

const (
	PurchaseStatePurchased = 0
	PurchaseStateCanceled  = 1
	PurchaseStatePending   = 2
)

type NotificationType int

const (
	SUBSCRIPTION_RECOVERED              NotificationType = 1  //- 从帐号保留状态恢复了订阅。
	SUBSCRIPTION_RENEWED                NotificationType = 2  //- 续订了处于活动状态的订阅。
	SUBSCRIPTION_CANCELED               NotificationType = 3  //- 自愿或非自愿地取消了订阅。如果是自愿取消，在用户取消时发送。
	SUBSCRIPTION_PURCHASED              NotificationType = 4  //- 购买了新的订阅。
	SUBSCRIPTION_ON_HOLD                NotificationType = 5  //- 订阅已进入帐号保留状态（如果已启用）。
	SUBSCRIPTION_IN_GRACE_PERIOD        NotificationType = 6  //- 订阅已进入宽限期（如果已启用）。
	SUBSCRIPTION_RESTARTED              NotificationType = 7  // - 用户已通过 Play > 帐号 > 订阅恢复了订阅。订阅已取消，但在用户恢复时尚未到期。如需了解详情，请参阅 [恢复](/google/play/billing/subscriptions#restore)。
	SUBSCRIPTION_PRICE_CHANGE_CONFIRMED NotificationType = 8  //- 用户已成功确认订阅价格变动。
	SUBSCRIPTION_DEFERRED               NotificationType = 9  //- 订阅的续订时间点已延期。
	SUBSCRIPTION_PAUSED                 NotificationType = 10 //- 订阅已暂停。
	SUBSCRIPTION_PAUSE_SCHEDULE_CHANGED NotificationType = 11 // - 订阅暂停计划已更改。
	SUBSCRIPTION_REVOKED                NotificationType = 12 //- 用户在到期时间之前已撤消订阅。
	SUBSCRIPTION_EXPIRED                NotificationType = 13 //- 订阅已到期。

	// 通知的类型。它可以具有以下值：
	ONE_TIME_PRODUCT_PURCHASED NotificationType = 1 //- 用户成功购买了一次性商品。
	ONE_TIME_PRODUCT_CANCELED  NotificationType = 2 // - 用户已取消待处理的一次性商品购买交易。
)

type RTDNBody struct {
	Message struct {
		Attributes struct {
			Key string `json:"key"`
		} `json:"attributes"`
		Data      string `json:"data"`
		MessageId string `json:"messageId"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

type RTDNData struct {
	Version                    string `json:"version"`
	PackageName                string `json:"packageName"`
	EventTimeMillis            string `json:"eventTimeMillis"`
	OneTimeProductNotification *struct {
		Version          string           `json:"version"`
		NotificationType NotificationType `json:"notificationType"`
		PurchaseToken    string           `json:"purchaseToken"`
		Sku              string           `json:"sku"`
	} `json:"oneTimeProductNotification"`
	SubscriptionNotification *struct {
		Version          string           `json:"version"`
		Notificationtype NotificationType `json:"notificationtype"`
		Purchasetoken    string           `json:"purchasetoken"`
		Subscriptionid   string           `json:"subscriptionid"`
	} `json:"subscriptionnotification"`
	TestNotification *struct {
		Version string `json:"version"`
	} `json:"testNotification"`
}
