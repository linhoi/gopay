package apple

import "time"

type Config struct {
	BundleID string `yaml:"bundle_id"`
	PSW      string `yaml:"psw"`
	RootCA   string `yaml:"root_ca"`
}

// Receipt is the receipt for an in-app purchase.
type Receipt struct {
	Quantity              int       `json:"quantity,omitempty"`
	ProductID             string    `json:"product_id,omitempty"`
	TransactionID         string    `json:"transaction_id,omitempty"`
	PurchaseDate          time.Time `json:"purchase_date"`
	OriginalTransactionID string    `json:"original_transaction_id,omitempty"`
	OriginalPurchaseDate  time.Time `json:"original_purchase_date"`
	ExpiresDate           time.Time `json:"expires_date"`
	WebOrderLineItemID    int       `json:"web_order_line_item_id,omitempty"`
	CancellationDate      time.Time `json:"cancellation_date"`
}

// Valid ...
func (r *Receipt) Valid() bool {
	if r.ExpiresDate.Unix() > 0 && r.ExpiresDate.Before(time.Now()) {
		return false
	}

	if r.CancellationDate.Unix() > 0 {
		return false
	}

	return true
}

// Receipts is the app receipt.
type Receipts struct {
	BundleID                   string    `json:"bundle_id,omitempty"`
	ApplicationVersion         string    `json:"application_version,omitempty"`
	OpaqueValue                []byte    `json:"opaque_value,omitempty"`
	SHA1Hash                   []byte    `json:"sha_1_hash,omitempty"`
	ReceiptCreationDate        time.Time `json:"receipt_creation_date"`
	InApp                      []Receipt `json:"in_app,omitempty"`
	OriginalApplicationVersion string    `json:"original_application_version,omitempty"`
	ExpirationDate             time.Time `json:"expiration_date"`

	rawBundleID []byte
}

// Valid ...
func (r *Receipts) Valid(bundleID string) bool {
	if r.BundleID != bundleID {
		return false
	}

	return true
}
