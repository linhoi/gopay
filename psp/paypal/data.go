package paypal

import (
	"net/url"
	"time"
)

type TransactionStatus string

const (
	timeFmt = "2006-01-02T15:04:05.000Z"

	// https://developer.paypal.com/docs/api/transaction-search/v1/
	TransactionDeny    TransactionStatus = "D"
	TransactionPending TransactionStatus = "P"
	TransactionSuccess TransactionStatus = "S"
	TransactionRefund  TransactionStatus = "V"
)

type AccessTokeResp struct {
	Scope       string `json:"scope"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	AppID       string `json:"app_id"`
	ExpiresIn   int    `json:"expires_in"`
	Nonce       string `json:"nonce"`
}

type ListDisputesReq struct {
	DisputeState     string `url:"dispute_state"`
	PageSize         int    `url:"page_size"`
	NextPageToken    string `url:"next_page_token,omitempty"`
	UpdateTimeBefore string `url:"update_time_before,omitempty"`
	UpdateTimeAfter  string `url:"update_time_after,omitempty"`
}

type ListDisputesResp struct {
	Items []DisputeItem `json:"items"`
	Links []HATEOASLink `json:"links"`
}

type DisputeItem struct {
	DisputeID     string    `json:"dispute_id"`
	CreateTime    time.Time `json:"create_time"`
	UpdateTime    time.Time `json:"update_time"`
	Status        string    `json:"status"`
	Reason        string    `json:"reason"`
	DisputeState  string    `json:"dispute_state"`
	DisputeAmount struct {
		CurrencyCode string `json:"currency_code"`
		Value        string `json:"value"`
	} `json:"dispute_amount"`
	Links []HATEOASLink `json:"links"`
}

type HATEOASLink struct {
	Href   string `json:"href"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}

type ShowDisputeDetailsResp struct {
	DisputeID            string                `json:"dispute_id"`
	CreateTime           time.Time             `json:"create_time"`
	UpdateTime           time.Time             `json:"update_time"`
	DisputedTransactions []DisputedTransaction `json:"disputed_transactions"`
	Reason               string                `json:"reason"`
	Status               string                `json:"status"`
	DisputeAmount        struct {
		CurrencyCode string `json:"currency_code"`
		Value        string `json:"value"`
	} `json:"dispute_amount"`
	DisputeOutcome struct {
		OutcomeCode    string `json:"outcome_code"`
		AmountRefunded struct {
			CurrencyCode string `json:"currency_code"`
			Value        string `json:"value"`
		} `json:"amount_refunded"`
	} `json:"dispute_outcome"`

	// not necessary now

	//DisputeLifeCycleStage string `json:"dispute_life_cycle_stage"`
	//DisputeChannel        string `json:"dispute_channel"`
	//Messages              []struct {
	//	PostedBy   string    `json:"posted_by"`
	//	TimePosted time.Time `json:"time_posted"`
	//	Content    string    `json:"content"`
	//	Documents  []struct {
	//		Name string `json:"name"`
	//		Url  string `json:"url"`
	//	} `json:"documents"`
	//} `json:"messages"`
	//Extensions struct {
	//	MerchandizeDisputeProperties struct {
	//		IssueType      string `json:"issue_type"`
	//		ServiceDetails struct {
	//			SubReasons  []string `json:"sub_reasons"`
	//			PurchaseUrl string   `json:"purchase_url"`
	//		} `json:"service_details"`
	//	} `json:"merchandize_dispute_properties"`
	//} `json:"extensions"`
	//Offer struct {
	//	BuyerRequestedAmount struct {
	//		CurrencyCode string `json:"currency_code"`
	//		Value        string `json:"value"`
	//	} `json:"buyer_requested_amount"`
	//} `json:"offer"`
	//Links []struct {
	//	Href   string `json:"href"`
	//	Rel    string `json:"rel"`
	//	Method string `json:"method"`
	//} `json:"links"`
}

func (s ShowDisputeDetailsResp) GetDisputeIDAndTransactionID() string {
	disputeIDAndTransactionID := ""
	for i, transaction := range s.DisputedTransactions {
		if i == 0 {
			disputeIDAndTransactionID = transaction.SellerTransactionID
		} else {
			disputeIDAndTransactionID = disputeIDAndTransactionID + "_" + transaction.SellerTransactionID
		}
	}

	disputeIDAndTransactionID = disputeIDAndTransactionID + "_" + s.DisputeID
	return disputeIDAndTransactionID
}

type DisputedTransaction struct {
	SellerTransactionID string    `json:"seller_transaction_id"`
	CreateTime          time.Time `json:"create_time"`
	TransactionStatus   string    `json:"transaction_status"`
	GrossAmount         struct {
		CurrencyCode string `json:"currency_code"`
		Value        string `json:"value"`
	} `json:"gross_amount"`
	Buyer struct {
		Name string `json:"name"`
	} `json:"buyer"`
	Seller struct {
		Email      string `json:"email"`
		MerchantID string `json:"merchant_id"`
		Name       string `json:"name"`
	} `json:"seller"`
}

// https://developer.paypal.com/docs/integration/direct/transaction-search/#list-transactions
type TransactionSearchReq struct {
	TransactionID     string            `url:"transaction_id,omitempty"`
	Fields            string            `url:"fields"`
	PageSize          int               `url:"page_size"`
	Page              int               `url:"page"`
	StartDate         string            `url:"start_date"`
	EndDate           string            `url:"end_date"`
	TransactionStatus TransactionStatus `url:"transaction_status,omitempty"`
}

type TransactionSearchResp struct {
	TransactionDetails []struct {
		TransactionInfo TransactionInfo `json:"transaction_info"`

		// Not necessary Body Now

		//PayerInfo       struct {
		//	AccountId     string `json:"account_id"`
		//	EmailAddress  string `json:"email_address"`
		//	AddressStatus string `json:"address_status"`
		//	PayerStatus   string `json:"payer_status"`
		//	PayerName     struct {
		//		GivenName         string `json:"given_name"`
		//		Surname           string `json:"surname"`
		//		AlternateFullName string `json:"alternate_full_name"`
		//	} `json:"payer_name"`
		//	CountryCode string `json:"country_code"`
		//} `json:"payer_info"`
		//ShippingInfo struct {
		//	Name    string `json:"name"`
		//	Address struct {
		//		Line1       string `json:"line1"`
		//		Line2       string `json:"line2"`
		//		City        string `json:"city"`
		//		CountryCode string `json:"country_code"`
		//		PostalCode  string `json:"postal_code"`
		//	} `json:"address"`
		//} `json:"shipping_info"`
		//CartInfo struct {
		//	ItemDetails []struct {
		//		ItemCode        string `json:"item_code,omitempty"`
		//		ItemName        string `json:"item_name"`
		//		ItemDescription string `json:"item_description,omitempty"`
		//		ItemQuantity    string `json:"item_quantity"`
		//		ItemUnitPrice   struct {
		//			CurrencyCode string `json:"currency_code"`
		//			Value        string `json:"value"`
		//		} `json:"item_unit_price"`
		//		ItemAmount struct {
		//			CurrencyCode string `json:"currency_code"`
		//			Value        string `json:"value"`
		//		} `json:"item_amount"`
		//		TaxAmounts []struct {
		//			TaxAmount struct {
		//				CurrencyCode string `json:"currency_code"`
		//				Value        string `json:"value"`
		//			} `json:"tax_amount"`
		//		} `json:"tax_amounts,omitempty"`
		//		TotalItemAmount struct {
		//			CurrencyCode string `json:"currency_code"`
		//			Value        string `json:"value"`
		//		} `json:"total_item_amount"`
		//		InvoiceNumber string `json:"invoice_number"`
		//	} `json:"item_details"`
		//} `json:"cart_info"`
		//StoreInfo struct {
		//} `json:"store_info"`
		//AuctionInfo struct {
		//} `json:"auction_info"`
		//IncentiveInfo struct {
		//} `json:"incentive_info"`
	} `json:"transaction_details"`
	AccountNumber         string `json:"account_number"`
	LastRefreshedDatetime string `json:"last_refreshed_datetime"`
	Page                  int    `json:"page"`
	TotalItems            int    `json:"total_items"`
	TotalPages            int    `json:"total_pages"`
	Links                 []struct {
		Href   string `json:"href"`
		Rel    string `json:"rel"`
		Method string `json:"method"`
	} `json:"links"`
}

type TransactionInfo struct {
	PaypalAccountId           string `json:"paypal_account_id"`
	TransactionId             string `json:"transaction_id"`
	TransactionEventCode      string `json:"transaction_event_code"`
	TransactionInitiationDate string `json:"transaction_initiation_date"`
	TransactionUpdatedDate    string `json:"transaction_updated_date"`
	TransactionAmount         struct {
		CurrencyCode string `json:"currency_code"`
		Value        string `json:"value"`
	} `json:"transaction_amount"`
	FeeAmount struct {
		CurrencyCode string `json:"currency_code"`
		Value        string `json:"value"`
	} `json:"fee_amount"`
	InsuranceAmount struct {
		CurrencyCode string `json:"currency_code"`
		Value        string `json:"value"`
	} `json:"insurance_amount"`
	ShippingAmount struct {
		CurrencyCode string `json:"currency_code"`
		Value        string `json:"value"`
	} `json:"shipping_amount"`
	ShippingDiscountAmount struct {
		CurrencyCode string `json:"currency_code"`
		Value        string `json:"value"`
	} `json:"shipping_discount_amount"`
	TransactionStatus     string `json:"transaction_status"`
	TransactionSubject    string `json:"transaction_subject"`
	TransactionNote       string `json:"transaction_note"`
	InvoiceId             string `json:"invoice_id"`
	CustomField           string `json:"custom_field"`
	ProtectionEligibility string `json:"protection_eligibility"`
}

// GetInvoiceID get invoiceID sent to PapPal by method SetExpressCheckout
func (t *TransactionInfo) GetInvoiceID() string {
	if t == nil {
		return ""
	}

	if t.InvoiceId != "" {
		return t.InvoiceId
	}

	return t.CustomField
}

type SetExpressCheckoutReq struct {
	ReturnURL    string `url:"RETURNURL"`
	Amount       string `url:"PAYMENTREQUEST_0_AMT"`
	NoShipping   string `url:"NOSHIPPING"`
	CurrencyCode string `url:"PAYMENTREQUEST_0_CURRENCYCODE"`
	Desc         string `url:"PAYMENTREQUEST_0_DESC"`
	Customer     string `url:"PAYMENTREQUEST_0_CUSTOM"`
	Invoice      string `url:"PAYMENTREQUEST_0_INVNUM"`
}

type setExpressCheckoutReq struct {
	NVPBase
	SetExpressCheckoutReq
}

type SetExpressCheckoutResp struct {
	Body url.Values
}

type NVPBase struct {
	Method    string `url:"METHOD"`
	Version   string `url:"VERSION"`
	User      string `url:"USER"`
	Pwd       string `url:"PWD"`
	Signature string `url:"SIGNATURE"`
}

type GetExpressCheckoutDetailsReq struct {
	Token   string `url:"TOKEN"`
}

type getExpressCheckoutDetailsReq struct {
	NVPBase
	GetExpressCheckoutDetailsReq
}

type GetExpressCheckoutDetailsResp struct {
	Body url.Values
}