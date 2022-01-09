package paypal

type NVPMethod string
type NVPVersion string

const (
	SetExpressCheckout NVPMethod = "SetExpressCheckout"
	GetExpressCheckoutDetails NVPMethod = "GetExpressCheckoutDetails"
	DoExpressCheckoutPayment NVPMethod = "DoExpressCheckoutPayment"

	Version124 NVPVersion = "124.0"
)