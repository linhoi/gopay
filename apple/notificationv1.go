package apple

type NotificationV1 struct {
	AutoRenewAdamID              string           `json:"auto_renew_adam_id"`                // An identifier that App Store Connect generates and the App Store uses to uniquely identify the auto-renewable subscription that the user's subscription renews. Treat this value as a 64-bit integer.
	AutoRenewProductID           string           `json:"auto_renew_product_id"`             // The product identifier of the auto-renewable subscription that the user's subscription renews.
	AutoRenewStatus              string           `json:"auto_renew_status"`                 // The current renewal status for an auto-renewable subscription product. Note that these values are different from those of the auto_renew_status in the receipt.  Possible values: true, false
	AutoRenewStatusChangeDate    string           `json:"auto_renew_status_change_date"`     //The time at which the renewal status for an auto-renewable subscription was turned on or off, in a date-time format similar to the ISO 8601 standard.
	AutoRenewStatusChangeDateMS  string           `json:"auto_renew_status_change_date_ms"`  //The time at which the renewal status for an auto-renewable subscription was turned on or off, in UNIX epoch time format, in milliseconds. Use this time format for processing dates.
	AutoRenewStatusChangeDatePTS string           `json:"auto_renew_status_change_date_pst"` // The time at which the renewal status for an auto-renewable subscription was turned on or off, in the Pacific time zone.
	Environment                  string           `json:"environment"`                       // The environment for which the receipt was generated.  Possible values: Sandbox, PROD
	ExpirationIntent             int              `json:"expiration_intent"`                 // The reason a subscription expired. This field is only present for an expired auto-renewable subscription. See expiration_intent for more information.
	NotificationType             NotificationType `json:"notification_type"`                 // The subscription event that triggered the notification.
	Password                     string           `json:"password"`                          // The same value as the shared secret you submit in the password field of the requestBody when validating receipts.
	UnifiedReceipt               UnifiedReceipt   `json:"unified_receipt"`                   // An object that contains information about the most recent in-app purchase transactions for the app.
	Bid                          string           `json:"bid"`                               // A string that contains the app bundle ID.
	Bvrs                         string           `json:"bvrs"`                              // A string that contains the app bundle version.
}

type NotificationType string

const (
	// CANCEL Indicates that either Apple customer support canceled the subscription or the user upgraded their subscription. The cancellation_date key contains the date and time of the change.
	CANCEL NotificationType = "CANCEL"
	// DIDChangeRenewalPref Indicates the customer made a change in their subscription plan that takes effect at the next renewal. The currently active plan is not affected.
	DIDChangeRenewalPref NotificationType = "DID_CHANGE_RENEWAL_PREF"
	// DIDChangeRenewalStatus Indicates a change in the subscription renewal status. Check auto_renew_status_change_date_ms and auto_renew_status in the JSON response to know the date and time of the last status update and the current renewal status.
	DIDChangeRenewalStatus NotificationType = "DID_CHANGE_RENEWAL_STATUS"
	// DIDFailToRenew Indicates a subscription that failed to renew due to a billing issue. Check is_in_billing_retry_period to know the current retry status of the subscription, and grace_period_expires_date to know the new service expiration date if the subscription is in a billing grace period.
	DIDFailToRenew NotificationType = "DID_FAIL_TO_RENEW"
	// DIDRecover Indicates a successful automatic renewal of an expired subscription that failed to renew in the past. Check expires_date to determine the next renewal date and time.
	DIDRecover NotificationType = "DID_RECOVER"
	// InitialBuy Occurs at the user's initial purchase of the subscription. Store latest_receipt on your server as a token to verify the users subscription status at any time by validating it with the App Store.
	InitialBuy NotificationType = "INITIAL_BUY"
	// InteractiveRenewal Indicates the customer renewed a subscription interactively, either by using your apps interface, or on the App Store in the account's Subscriptions settings. Make service available immediately.
	InteractiveRenewal NotificationType = "INTERACTIVE_RENEWAL"
	// RENEWAL Indicates a successful automatic renewal of an expired subscription that failed to renew in the past. Check expires_date to determine the next renewal date and time.
	RENEWAL NotificationType = "RENEWAL"
	// REFUND Indicates that App Store successfully refunded a transaction. The cancellation_date_ms contains the timestamp of the refunded transaction; the original_transaction_id and product_id identify the original transaction and product, and cancellation_reason contains the reason.
	REFUND NotificationType = "REFUND"
)

type ReceiptInfo struct {
	CancellationDate            string `json:"cancellation_date"`             // The time Apple customer support canceled a transaction, in a date-time format similar to the ISO 8601. This field is only present for refunded transactions.
	CancellationDateMs          string `json:"cancellation_date_ms"`          // The time Apple customer support canceled a transaction, or the time an auto-renewable subscription plan was upgraded, in UNIX epoch time format, in milliseconds. This field is only present for refunded transactions. Use this time format for processing dates. See cancellation_date_ms for more information.
	CancellationDatePst         string `json:"cancellation_date_pst"`         // The time Apple customer support canceled a transaction, in the Pacific Time zone. This field is only present for refunded transactions.
	CancellationReason          string `json:"cancellation_reason"`           // The reason for a refunded transaction. When a customer cancels a transaction, the App Store gives them a refund and provides a value for this key. A value of 1 indicates that the customer canceled their transaction due to an actual or perceived issue within your app. A value of 0 indicates that the transaction was canceled for another reason; for example, if the customer made the purchase accidentally.  Possible values: 1, 0
	ExpiresDate                 string `json:"expires_date"`                  // The time a subscription expires or when it will renew, in a date-time format similar to the ISO 8601.
	ExpiresDateMs               string `json:"expires_date_ms"`               // The time a subscription expires or when it will renew, in UNIX epoch time format, in milliseconds. Use this time format for processing dates. See expires_date_ms for more information.
	ExpiresDatePst              string `json:"expires_date_pst"`              // The time a subscription expires or when it will renew, in the Pacific Time zone.
	IsInIntroOfferPeriod        string `json:"is_in_intro_offer_period"`      // An indicator of whether an auto-renewable subscription is in the introductory price period. See is_in_intro_offer_period for more information.  Possible values: true, false
	IsTrialPeriod               string `json:"is_trial_period"`               // An indicator of whether a subscription is in the free trial period. See is_trial_period for more information.
	IsUpgraded                  string `json:"is_upgraded"`                   // An indicator that a subscription has been canceled due to an upgrade. This field is only present for upgrade transactions.  Value: true
	OriginalPurchaseDate        string `json:"original_purchase_date"`        // The time of the original app purchase, in a date-time format similar to ISO 8601.
	OriginalPurchaseDateMs      string `json:"original_purchase_date_ms"`     // The time of the original app purchase, in UNIX epoch time format, in milliseconds. Use this time format for processing dates. For an auto-renewable subscription, this value indicates the date of the subscription's initial purchase. The original purchase date applies to all product types and remains the same in all transactions for the same product ID. This value corresponds to the original transactions transactionDate property in StoreKit.
	OriginalPurchaseDatePst     string `json:"original_purchase_date_pst"`    // The time of the original app purchase, in the Pacific Time zone.
	OriginalTransactionID       string `json:"original_transaction_id"`       // The transaction identifier of the original purchase. See original_transaction_id for more information.
	ProductID                   string `json:"product_id"`                    // The unique identifier of the product purchased. You provide this value when creating the product in App Store Connect, and it corresponds to the productIdentifier property of the SKPayment object stored in the transaction's payment property.
	PromotionalOfferID          string `json:"promotional_offer_id"`          // The identifier of the subscription offer redeemed by the user. See promotional_offer_id for more information.
	PurchaseDate                string `json:"purchase_date"`                 // The time the App Store charged the user's account for a purchased or restored product, or the time the App Store charged the users account for a subscription purchase or renewal after a lapse, in a date-time format similar to ISO 8601.
	PurchaseDateMs              string `json:"purchase_date_ms"`              // For consumable, non-consumable, and non-renewing subscription products, the time the App Store charged the user's account for a purchased or restored product, in the UNIX epoch time format, in milliseconds. For auto-renewable subscriptions, the time the App Store charged the users account for a subscription purchase or renewal after a lapse, in the UNIX epoch time format, in milliseconds. Use this time format for processing dates.
	PurchaseDatePst             string `json:"purchase_date_pst"`             // The time the App Store charged the user's account for a purchased or restored product, or the time the App Store charged the users account for a subscription purchase or renewal after a lapse, in the Pacific Time zone.
	Quantity                    string `json:"quantity"`                      // The number of consumable products purchased. This value corresponds to the quantity property of the SKPayment object stored in the transaction's payment property. The value is usually 1 unless modified with a mutable payment. The maximum value is 10.
	SubscriptionGroupIdentifier string `json:"subscription_group_identifier"` // The identifier of the subscription group to which the subscription belongs. The value for this field is identical to the subscriptionGroupIdentifier property in SKProduct.
	TransactionID               string `json:"transaction_id"`                // A unique identifier for a transaction such as a purchase, restore, or renewal. See transaction_id for more information.
	WebOrderLineItemID          string `json:"web_order_line_item_id"`        // A unique identifier for purchase events across devices, including subscription-renewal events. This value is the primary key for identifying subscription purchases.
}

type UnifiedReceipt struct {
	Environment        string         `json:"environment"`          //The environment for which the receipt was generated.  Possible values: Sandbox, Production
	LatestReceipt      string         `json:"latest_receipt"`       //The latest Base64-encoded app receipt.
	LatestReceiptInfo  []*ReceiptInfo `json:"latest_receipt_info"`  //An array that contains the latest 100 in-app purchase transactions of the decoded value in latest_receipt. This array excludes transactions for consumable products that your app has marked as finished. The contents of this array are identical to those in responseBody.Latest_receipt_info in the verifyReceipt endpoint response for receipt validation.
	PendingRenewalInfo []string       `json:"pending_renewal_info"` //An array where each element contains the pending renewal information for each auto-renewable subscription identified in product_id. The contents of this array are identical to those in responseBody.Pending_renewal_info in the verifyReciept endpoint response for receipt validation.
	Status             int            `json:"status"`               //The status code, where 0 indicates that the notification is valid.  Value: 0
}
