package apple

import "errors"

type Status int

func (s Status) Error() error {
	message, ok := statusMessage[int(s)]
	if ok {
		return errors.New(message)
	}
	return errors.New("unknown status")
}

var (
	// reference: // https://developer.apple.com/documentation/appstorereceipts/status
	statusMessage = map[int]string{
		21000: "The request to the App Store was not made using the HTTP POST request method",
		21001: "This status code is no longer sent by the App Store",
		21002: "The data in the receipt-data property was malformed or the service experienced a temporary issue. Try again",
		21003: "The receipt could not be authenticated",
		21004: "The shared secret you provided does not match the shared secret on file for your account",
		21005: "The receipt server was temporarily unable to provide the receipt. Try again",
		21006: "This receipt is valid but the subscription has expired. When this status code is returned to your server, the receipt data is also decoded and returned as part of the response. Only returned for iOS 6-style transaction receipts for auto-renewable subscriptions",
		21007: "This receipt is from the test environment, but it was sent to the production environment for verification",
		21008: "This receipt is from the production environment, but it was sent to the test environment for verification",
		21009: "Internal data access error. Try again later",
		21010: "The user account cannot be found or has been deleted",
	}
)

var (
	StatusSuccess Status = 0
	StatusMethodError Status = 21000
	StatusReceiptMalformedOrServiceError Status = 21002
	StatusReceiptUnauthenticated Status = 21003
	StatusSharedSecretUnMatch    Status = 21004
	StatusReceiptProvideError    Status = 21005
	StatusReceiptSubscriptionExpired Status = 21006
	StatusEnvironmentDisMatchFromTest Status = 21007
	StatusEnvironmentDisMatchFromProduction Status = 21008
	StatusDataAccessError Status = 21009
	StatusAccountError Status = 21010
)
